package alert

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type Alerter interface {
	AlertToUser(users string, content *Content, alertType AlertType) error
	AlertToGroup(groupIds string, content *Content, alertType AlertType) error
	GetUserIDs(emails string) (string, error)
}

type alerter struct {
	httpclient *resty.Client
	config     *AlerterConfiguration
}

func NewAlerter(config *AlerterConfiguration) (Alerter, error) {
	if config.Sk == "" {
		return nil, errors.New("Alerter sk is empty.")
	}

	return &alerter{
		config: config,
	}, nil
}

func (c *alerter) AlertToGroup(groupIds string, content *Content, alertType AlertType) error {
	if !c.config.AlertAdmin {
		return nil
	}

	ts := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

	text, err := content.Execute()
	if err != nil {
		return err
	}

	request := SendMsgRequest{
		AppId:    c.config.Ak,
		Ts:       ts,
		RbotId:   alertType.ID(),
		GroupIds: groupIds,
		Text:     text,
		S:        generateHMAC(map[string]string{"appId": c.config.Ak, "ts": ts, "robotId": alertType.ID(), "groupIds": groupIds, "text": text}, c.config.Sk),
	}

	response, err := c.httpclient.R().SetBody(&request).SetDoNotParseResponse(true).Post(generatePath(c.config.Host, c.config.MessagePath))
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(response.RawBody())
	defer response.RawBody().Close()

	var sendMsgResponse SendMsgResponse
	err = json.Unmarshal(data, &sendMsgResponse)
	if err != nil {
		return err
	}

	if sendMsgResponse.Code != 0 {
		return fmt.Errorf("send bosshi message failed %s", sendMsgResponse.Msg)
	}

	return nil
}

func (c *alerter) AlertToUser(users string, content *Content, alertType AlertType) error {
	userIds, err := c.GetUserIDs(users)
	if err != nil {
		return err
	}

	ts := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

	text, err := content.Execute()
	if err != nil {
		return err
	}

	request := SendMsgRequest{
		AppId:   c.config.Ak,
		Ts:      ts,
		RbotId:  alertType.ID(),
		UserIds: userIds,
		Text:    text,
		S:       generateHMAC(map[string]string{"appId": c.config.Ak, "ts": ts, "robotId": alertType.ID(), "userIds": userIds, "text": text}, c.config.Sk),
	}

	// http.Header{"Content-Type": {"application/json"}}
	response, err := c.httpclient.R().SetBody(&request).SetDoNotParseResponse(true).Post(generatePath(c.config.Host, c.config.MessagePath))
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(response.RawBody())
	defer response.RawBody().Close()

	var sendMsgResponse SendMsgResponse
	err = json.Unmarshal(data, &sendMsgResponse)
	if err != nil {
		return err
	}

	if sendMsgResponse.Code != 0 {
		return fmt.Errorf("send bosshi message failed %s", sendMsgResponse.Msg)
	}

	return nil
}

func (c *alerter) GetUserIDs(emails string) (string, error) {
	ts := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

	request := GetBossHiUserRequest{
		AppId:  c.config.Ak,
		Ts:     ts,
		Emails: emails,
		S:      generateHMAC(map[string]string{"appId": c.config.Ak, "ts": ts, "emails": emails}, c.config.Sk),
	}

	response, err := c.httpclient.R().SetBody(&request).SetDoNotParseResponse(true).Post(generatePath(c.config.Host, c.config.UserIdPath))
	if err != nil {
		return "", err
	}

	data, _ := ioutil.ReadAll(response.RawBody())
	defer response.RawBody().Close()

	var getUserResponse GetUserResponse
	err = json.Unmarshal(data, &getUserResponse)
	if err != nil {
		return "", err
	}

	if getUserResponse.Code != 0 {
		return "", fmt.Errorf("get bosshi users failed %s", getUserResponse.Msg)
	}

	return generateUserIDs(getUserResponse.Data.Result), nil
}
