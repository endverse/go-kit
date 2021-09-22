package alert

import (
	"fmt"
	"html/template"
	"testing"
)

func TestErrorTemplateRender(t *testing.T) {

	errorTemplateRender := ErrorTemplateRender{
		ExpName:     template.HTML("lab-3"),
		NickName:    template.HTML("测试Template"),
		TaskName:    "testing-task-name",
		FailedNodes: "t2,t3",
		Progress:    "2/3",
		Result:      template.HTML("http://arsenal-gray.weizhipin.com/lab/detail?expName=lab-3&aflowOnceID=testing-task-name&tab=taskRecord"),
		StartTime:   "2021-09-22 00:00:00",
		EndTime:     "2021-09-22 23:09:09",
	}

	content := Content{
		Platform:     "Arsenal算法平台",
		Title:        "ERROR",
		TitleContent: "发现你有一个作业运行失败，请关注",
		Template:     &errorTemplateRender,
	}

	text, err := content.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("\n%s \n", text)
}
