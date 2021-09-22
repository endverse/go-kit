package alert

import (
	"bytes"
	"html/template"

	"github.com/spf13/pflag"
)

type Content struct {
	Platform     string         `json:"platform"`
	Title        string         `json:"title"`
	TitleContent string         `json:"titleContent"`
	Template     TemplateRender `json:"template"`
}

func (c *Content) Execute() (string, error) {
	var w bytes.Buffer

	templateRender := template.Must(template.New(c.Template.Name()).Parse(ContentTitle + c.Template.Tempalte()))

	err := templateRender.Execute(&w, c)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}

type AlertType string

// AlertTypes
const (
	AlertArsenalPlatformAlertType  AlertType = "Arsenal告警助手"
	NoticeArsenalPlatformAlertType AlertType = "Arsenal通知助手"
)

// AlertType IDs
const (
	AlertArsenalPlatformAlertTypeID  = "ebcfc36a33b81b6eumc1Jd3X"
	NoticeArsenalPlatformAlertTypeID = "87664fbe6e6302e7umc1Jd3W"
)

func (at AlertType) String() string {
	return string(at)
}

func (at AlertType) ID() string {
	switch at {
	case AlertArsenalPlatformAlertType:
		return AlertArsenalPlatformAlertTypeID
	case NoticeArsenalPlatformAlertType:
		return NoticeArsenalPlatformAlertTypeID
	default:
	}

	return ""
}

type GetBossHiUserRequest struct {
	AppId  string `json:"appId"`
	Ts     string `json:"ts"`
	Emails string `json:"emails"`
	S      string `json:"s"`
}

type GetUserResponse struct {
	Code int32               `json:"code"`
	Msg  string              `json:"msg"`
	Data GetUserResponseData `json:"data"`
}

type GetUserResponseData struct {
	Result []BossHiUser `json:"result"`
}

type BossHiUser struct {
	Email  string `json:"email"`
	UserId int64  `json:"userId"`
}

type SendMsgRequest struct {
	AppId    string `json:"appId"`
	Ts       string `json:"ts"`
	RbotId   string `json:"robotId"`
	UserIds  string `json:"userIds"`
	GroupIds string `json:"groupIds"`
	Text     string `json:"text"`
	S        string `json:"s"`
}

type SendMsgResponse struct {
	Code int32               `json:"code"`
	Msg  string              `json:"msg"`
	Data SendMsgResponseData `json:"data"`
}

type SendMsgResponseData struct {
	UserIdMsgs  []UserIdMsg  `json:"userIdMsgs"`
	GroupIdMsgs []GroupIdMsg `json:"groupIdMsgs"`
}

type UserIdMsg struct {
	MsgId  int64 `json:"msgId"`
	UserId int64 `json:"userId"`
}

type GroupIdMsg struct {
	MsgId   int64 `json:"msgId"`
	GroupId int64 `json:"groupId"`
}

type AlerterConfiguration struct {
	Ak          string `json:"ak"`
	Sk          string `json:"sk"`
	Host        string `json:"host"`
	MessagePath string `json:"messagePath"`
	UserIdPath  string `json:"useridPath"`
}

func (o *AlerterConfiguration) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.Ak, "ak", o.Ak, "Alerter access key.")
	fs.StringVar(&o.Sk, "sk", o.Sk, "Alerter secret key.")
	fs.StringVar(&o.Host, "host", o.Host, "Alerter remote server host.")
	fs.StringVar(&o.MessagePath, "msg-path", o.MessagePath, "Alerter send alert message URI.")
	fs.StringVar(&o.UserIdPath, "user-path", o.UserIdPath, "Alerter get users id URI.")
}

func (o *AlerterConfiguration) Validate() []error {
	if o == nil {
		return nil
	}

	var errs []error
	return errs
}

func ConvertWithDefault(cfg *AlerterConfiguration) {
	cfg.Ak = "7Awt41Z1"
	cfg.Host = "https://inner-hi.zhipin.com"
	cfg.MessagePath = "/api/open/message/sendText"
	cfg.UserIdPath = "/api/open/user/getIdByEmails"
}
