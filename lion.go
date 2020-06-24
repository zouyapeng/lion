package lion

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LoginURI = "/getToken.do"
	QueryURI = "/configQuery.do"
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

type Lion struct {
	url        string
	username   string
	password   string
	token      string
	expireTime int64
}

func (l Lion) String() string {
	return fmt.Sprintf("<Lion: %s User: %s>", l.url, l.username)
}

func (l *Lion) QueryConfigByProject(env, projectName string) (config map[string]interface{}, err error) {
	if l.expireTime < time.Now().Unix() {
		logrus.Printf("token %s is expired", l.token)
		if err = login(l); err != nil {
			return
		}
		logrus.Printf("new token %s", l.token)
	}

	r, err := doRequest(l.url+QueryURI, map[string]string{
		"app":         l.username,
		"token":       l.token,
		"env":         env,
		"projectName": projectName,
	})

	if err != nil {
		logrus.Error(err)
		return
	}
	var resp Response
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return
	}

	if resp.Code != 0 {
		return config, errors.New(resp.Msg)
	}

	return resp.Result.(map[string]interface{}), nil
}

func (l *Lion) QueryConfigByKey(env string, key interface{}) (config map[string]interface{}, err error) {
	if l.expireTime < time.Now().Unix() {
		logrus.Printf("token %s is expired", l.token)
		if err = login(l); err != nil {
			return
		}
		logrus.Printf("new token %s", l.token)
	}

	var keyString string
	switch key.(type) {
	case string:
		keyString = key.(string)
	case []string:
		keyString = strings.Join(key.([]string), ",")
	}

	r, err := doRequest(l.url+QueryURI, map[string]string{
		"app":   l.username,
		"token": l.token,
		"env":   env,
		"key":   keyString,
	})

	if err != nil {
		logrus.Error(err)
		return
	}
	var resp Response
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return
	}

	if resp.Code != 0 {
		return config, errors.New(resp.Msg)
	}

	return resp.Result.(map[string]interface{}), nil
}

func login(l *Lion) (err error) {
	r, err := doRequest(l.url+LoginURI, map[string]string{
		"app": l.username,
		"pwd": l.password,
	})
	if err != nil {
		return
	}

	var resp Response
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return
	}

	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}

	l.token = resp.Result.(map[string]interface{})["token"].(string)
	l.expireTime = int64(resp.Result.(map[string]interface{})["expireTime"].(float64) / 1000)

	return
}

// Login 初始化lion结构体，获取token
func Init(url, username, password string) (l *Lion, err error) {
	tmpLion := &Lion{
		url:      url,
		username: username,
		password: password,
	}

	if err = login(tmpLion); err != nil {
		return
	}
	return tmpLion, nil
}
