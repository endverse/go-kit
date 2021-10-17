package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

const Driver = "mysql"

type Configuration struct {
	// Url: user:passowrd@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	MaxRetry int

	LogLevel                  int
	MaxIdle                   int
	MaxOpenConns              int
	PrepareStmt               bool
	SingularTable             bool
	ConnMaxLifetime           time.Duration
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	SkipDefaultTransaction    bool
	QueryFields               bool
	Colorful                  bool

	EnableTracing bool
}

func (o *Configuration) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.Host, "mysql-host", o.Host, "Mysql database host.")
	fs.IntVar(&o.Port, "mysql-port", o.Port, "Mysql database port.")
	fs.StringVar(&o.User, "mysql-user", o.User, "Mysql database user.")
	fs.StringVar(&o.Password, "mysql-password", o.Password, "Mysql database password.")
	fs.StringVar(&o.DBName, "mysql-db", o.DBName, "Mysql database name.")
	fs.IntVar(&o.MaxIdle, "mysql-max-idle", o.MaxIdle, "Mysql max idle.")
	fs.IntVar(&o.MaxOpenConns, "mysql-max-open-conns", o.MaxOpenConns, "Mysql max active connection.")
	fs.IntVar(&o.LogLevel, "mysql-log-level", o.LogLevel, "LogLevel set log level, '1 - Silent' '2 - Error' '3 - Warn' '4 - Info'.")
	fs.BoolVar(&o.PrepareStmt, "prepare-stmt", o.PrepareStmt, "PrepareStmt executes the given query in cached statement.")
	fs.BoolVar(&o.QueryFields, "query-fields", o.QueryFields, "QueryFields executes the SQL query with all fields of the table")
	fs.BoolVar(&o.SingularTable, "singular-table", o.SingularTable, "LogMode set log mode, `true` for detailed logs, `false` for no log, default, will only print error logs.")
	fs.BoolVar(&o.IgnoreRecordNotFoundError, "ignore-record-not-found-error", o.IgnoreRecordNotFoundError, "Ignore record not found error.")
	fs.BoolVar(&o.SkipDefaultTransaction, "skip-default-tx", o.SkipDefaultTransaction, "GORM perform single create, update, delete operations in transactions by default to ensure database data integrity. You can disable it by setting `SkipDefaultTransaction` to true")
	fs.BoolVar(&o.Colorful, "color-log", o.Colorful, "Enable log color")
	fs.DurationVar(&o.ConnMaxLifetime, "conn-max-life-time", o.ConnMaxLifetime, "Mysql conn max life time.")
	fs.DurationVar(&o.SlowThreshold, "slow-threshold", o.SlowThreshold, "Mysql slow threshold.")
	fs.IntVar(&o.MaxRetry, "max-retry", o.MaxRetry, "Max connect to the mysql server retry count.")
	fs.BoolVar(&o.EnableTracing, "enable-tracing", o.EnableTracing, "Enable open tracing.")
}

func (o *Configuration) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	if o.Host == "" {
		errs = append(errs, errors.New("Mysql.Host is empty"))
	}
	if o.Port == 0 {
		errs = append(errs, errors.New("Mysql.Port is empty"))
	}
	if o.User == "" {
		errs = append(errs, errors.New("Mysql.User is empty"))
	}
	if o.Password == "" {
		errs = append(errs, errors.New("Mysql.Password is empty"))
	}
	if o.DBName == "" {
		errs = append(errs, errors.New("Mysql.DBName is empty"))
	}

	switch o.LogLevel {
	case 1, 2, 3, 4:
	default:
		errs = append(errs, errors.New("Mysql.LogLevel only support '1, 2, 3, 4'"))
	}

	return errs
}

func (o *Configuration) Url() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", o.User, o.Password, o.Host, o.Port, o.DBName)
}

func (o *Configuration) Hosts() int {
	return o.MaxIdle
}

func (o *Configuration) MaxIdles() int {
	return o.MaxIdle
}

func (o *Configuration) MaxActives() int {
	return o.MaxOpenConns
}

func (o *Configuration) MaxRetrys() int {
	return o.MaxRetry
}

func (o *Configuration) ConnMaxLifetimes() time.Duration {
	return o.ConnMaxLifetime
}

func (o *Configuration) SingularTables() bool {
	return o.SingularTable
}

func (o *Configuration) LogLevels() int {
	return o.LogLevel
}

var DefaultConfiguratuin = &Configuration{
	Host:                      "",
	Port:                      3306,
	User:                      "",
	Password:                  "",
	DBName:                    "",
	MaxRetry:                  15,
	LogLevel:                  2,
	MaxIdle:                   20,
	MaxOpenConns:              30,
	QueryFields:               true,
	PrepareStmt:               true,
	SingularTable:             true,
	ConnMaxLifetime:           10 * time.Second,
	SlowThreshold:             200 * time.Millisecond,
	IgnoreRecordNotFoundError: false,
	SkipDefaultTransaction:    false,
	Colorful:                  true,
	EnableTracing:             false,
}
