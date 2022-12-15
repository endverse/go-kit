package mysql

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/endverse/go-kit/storage/mysql/config"

	"github.com/endverse/go-kit/storage/mysql/tracing"

	"github.com/endverse/go-kit/retry"

	"github.com/hex-techs/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Core struct {
	db *gorm.DB

	cfg *config.Configuration
}

func New(cfg *config.Configuration) *Core {
	return &Core{
		cfg: cfg,
	}
}

func (c *Core) Conn() error {
	var db *gorm.DB
	var err error

	run := func() error {
		if db, err = connect(c.cfg); err != nil {
			return err
		}
		return nil
	}

	err = retry.RetryFunc(c.cfg.MaxRetry, 1*time.Second, run)
	if err != nil {
		return err
	}

	c.db = db
	return nil
}

// return a new core with GORM DB with context
// WithContext change current instance db's context to ctx
func (c *Core) GetDB(ctx context.Context) (*gorm.DB, error) {
	iface := ctx.Value(CtxTransactionKey{})

	if iface != nil {
		tx, ok := iface.(*gorm.DB)
		if !ok {
			return nil, fmt.Errorf("unexpect context value type: %s", reflect.TypeOf(tx))
		}

		return tx, nil
	}

	return c.db.WithContext(ctx), nil
}

// connect to the mysql server
func connect(cfg *config.Configuration) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Url()), &gorm.Config{
		SkipDefaultTransaction: cfg.SkipDefaultTransaction,
		PrepareStmt:            cfg.PrepareStmt,
		QueryFields:            cfg.QueryFields,
		NamingStrategy:         schema.NamingStrategy{SingularTable: cfg.SingularTable},
		Logger: logger.New(&klog.GormLoggerWriter{}, logger.Config{
			SlowThreshold:             cfg.SlowThreshold,
			LogLevel:                  logger.LogLevel(cfg.LogLevel),
			IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
			Colorful:                  cfg.Colorful,
		}),
	})
	if err != nil {
		return nil, err
	}

	if cfg.EnableTracing {
		_ = db.Use(&tracing.OpentracingPlugin{})
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
