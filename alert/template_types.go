package alert

import (
	"html/template"
)

type TemplateRender interface {
	Template() string
	Name() string
}

const (
	ContentTitle = `【{{.Platform}}】
{{.Title}}: {{.TitleContent}}
`
)

var _ TemplateRender = &ErrorTemplateRender{}
var _ TemplateRender = &NoticeTemplateRender{}

type ErrorTemplateRender struct {
	ExpName     template.HTML
	NickName    template.HTML
	TaskName    string
	FailedNodes template.HTML
	Progress    string
	Result      template.HTML
	StartTime   string
	EndTime     string
}

func (t *ErrorTemplateRender) Name() string {
	return "errorContent"
}

func (t *ErrorTemplateRender) Template() string {
	return `
实验名称: {{.Template.ExpName}}
作业名字: {{.Template.NickName}}
作业ID: {{.Template.TaskName}}
失败节点: {{.Template.FailedNodes}}
作业进度: {{.Template.Progress}}
结果地址: {{.Template.Result}}
开始时间: {{.Template.StartTime}}
失败时间: {{.Template.EndTime}}`
}

type NoticeTemplateRender struct {
	ExpName   template.HTML
	NickName  template.HTML
	TaskName  string
	Progress  string
	Result    template.HTML
	StartTime string
	EndTime   string
}

func (t *NoticeTemplateRender) Template() string {
	return `
实验名称: {{.Template.ExpName}}
作业名字: {{.Template.NickName}}
作业ID: {{.Template.TaskName}}
作业进度: {{.Template.Progress}}
结果地址: {{.Template.Result}}
开始时间: {{.Template.StartTime}}
结束时间: {{.Template.EndTime}}`
}

func (t *NoticeTemplateRender) Name() string {
	return "noticeContent"
}
