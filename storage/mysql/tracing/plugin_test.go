package tracing

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func initJaeger() (closer io.Closer, err error) {
	// 根据配置初始化Tracer 返回Closer
	tracer, closer, err := (&config.Configuration{
		ServiceName: "gormTracing",
		Disabled:    false,
		Sampler: &config.SamplerConfig{
			Type: jaeger.SamplerTypeConst,
			// param的值在0到1之间，设置为1则将所有的Operation输出到Reporter
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}).NewTracer()
	if err != nil {
		return
	}

	// 设置全局Tracer - 如果不设置将会导致上下文无法生成正确的Span
	opentracing.SetGlobalTracer(tracer)
	return
}

type Aflow struct {
	ID             uint       `json:"id" gorm:"primary_key;type:bigint(20) auto_increment;comment:'主键自增id'"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"comment:'创建时间';default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"comment:'更新时间';default:CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"comment:'删除时间，字段不为空的时候视为标记删除'" sql:"index"`
	Name           string     `json:"name" gorm:"unique_index:af_name;comment:'aflow名字';not null;default:''"`
	NickName       string     `json:"nickName" gorm:"comment:'算法作业别名';not null;default:''"`
	UserId         int        `json:"userId" gorm:"comment:'用户id';not null"`
	UserName       string     `json:"userName" gorm:"comment:'用户名';not null;default:''"`
	UserNameZh     string     `json:"userNameZh" gorm:"comment:'中文用户名';not null'default:''"`
	UserEmail      string     `json:"userEmail" gorm:"comment:'用户邮箱';not null;default:''"`
	TeamId         string     `json:"teamId" gorm:"comment:'团队id'"`
	TeamName       string     `json:"teamName" gorm:"comment:'团队名字';not null;default:''"`
	ExpId          uint       `json:"expId" gorm:"comment:'实验id';type:bigint(20);not null"`
	ExpName        string     `json:"expName" gorm:"comment:'实验名';not null;default:''"`
	ExpNameZh      string     `json:"expNameZh" gorm:"comment:'实验名';not null;default:''"`
	StepsCount     int        `json:"stepCount" gorm:"comment:'总步数';not null"`
	SourcePath     string     `json:"sourcePath" gorm:"comment:'源文件路径';not null;default:''"`
	ReleasePath    string     `json:"releasePath" gorm:"comment:'release文件路径';not null;default:''"`
	ScheduleStatus *int       `json:"scheduleStatus" gorm:"-;column:status"`
	Trigger        string     `json:"trigger" gorm:"type:longtext;comment:'触发作业'"`
	Version        string     `json:"version" gorm:"comment:'版本号';not null;default:''"` // format: 20210101111111
	Description    string     `json:"description" gorm:"comment:'描述信息';not null;default:''"`
}

const dsn = "root:michaelson@tcp(localhost:3306)/arsenal?charset=utf8mb4&parseTime=True&loc=Local"

func TestOpentracing(t *testing.T) {
	closer, err := initJaeger()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		// Logger: logger.New(&log.GormLoggerWriter{}, logger.Config{
		// 	SlowThreshold:             200 * time.Microsecond,
		// 	LogLevel:                  logger.LogLevel(4),
		// 	IgnoreRecordNotFoundError: false,
		// 	Colorful:                  true,
		// }),
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = db.Use(&OpentracingPlugin{})

	// 生成新的Span - 注意将span结束掉，不然无法发送对应的结果
	span := opentracing.StartSpan("gormTracing unit test")
	defer span.Finish()

	// 把生成的Root Span写入到Context上下文，获取一个子Context
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	session := db.WithContext(ctx)

	// Create
	session.Create(&Aflow{Name: "opentracing", NickName: "opentracing-nick"})
	// Read
	var aflow Aflow
	session.First(&aflow, 108)                       // 根据整形主键查找
	session.First(&aflow, "name = ?", "opentracing") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	session.Model(&aflow).Update("nick_name", "opentracing-nick")
	// Update - 更新多个字段
	session.Model(&aflow).Updates(Aflow{UserName: "michaelson", ExpName: "opentracing"}) // 仅更新非零值字段
	session.Model(&aflow).Updates(map[string]interface{}{"user_name": "michaelson1", "exp_name": "opentracing2"})

	// Delete - 删除 product
	session.Delete(&aflow, 108)
}
