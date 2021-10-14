package mysql

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go-arsenal.kanzhun.tech/arsenal/go-kit/retry"
	"go-arsenal.kanzhun.tech/arsenal/go-kit/storage/mysql/config"
	"go-arsenal.kanzhun.tech/arsenal/go-kit/storage/mysql/tracing"

	"github.com/hex-techs/klog"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type core struct {
	db *gorm.DB

	cfg *config.Configuration
}

func New(cfg *config.Configuration) *core {
	return &core{
		cfg: cfg,
	}
}

func (c *core) Conn() error {
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
func (c *core) ctxDB(ctx context.Context) (*gorm.DB, error) {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		tx, ok := iface.(*gorm.DB)
		if !ok {
			return nil, fmt.Errorf("unexpect context value type: %s", reflect.TypeOf(tx))
		}

		return tx, nil
	}

	return c.db.WithContext(ctx), nil
}

// DB return gorm.DB, spanFinish func, some error.
// If EnableTracing is false, spanFinish is a empty func.
// param titles must has only one item.
func (c *core) DB(ctx context.Context, titles ...string) (*gorm.DB, func(), error) {
	var spanFinish = func() {}
	var title = "gromTracing"

	if c.cfg.EnableTracing {
		if len(titles) == 1 {
			title = titles[0]
		}
		// generate a new span.
		// NOTE: must need finish the span, if not, it can not send data.
		span := opentracing.StartSpan(title)
		spanFinish = func() { span.Finish() }

		// set the root span into the contxt, and generate a new context
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	db, err := c.ctxDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	return db, spanFinish, nil
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
		tracing.CreateGlobalJager(&jaegerconfig.Configuration{
			ServiceName: "GORM-Tracing",
			Disabled:    false,
			Sampler: &jaegerconfig.SamplerConfig{
				Type: jaeger.SamplerTypeConst,
				// The param's value is between 0 and 1,
				// if set to 1, it will output all operations to the Reporter.
				Param: 1,
			},
			Reporter: &jaegerconfig.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: cfg.TracingHostPort,
				User:               cfg.TracingUser,
				Password:           cfg.TracingPassword,
			},
		})

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
