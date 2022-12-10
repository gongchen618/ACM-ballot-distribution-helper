package controller

import (
	"ballot/data"
	"ballot/qq"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func MessageReverseHandler(c echo.Context) error {
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := struct {
		PostType string `json:"post_type"`
	}{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		logrus.WithFields(logrus.Fields{"err": err.Error()})
		return c.JSON(http.StatusBadRequest, nil)
	}

	if p.PostType == "message" {
		p2 := struct {
			MessageType string `json:"message_type"`
		}{}
		if err := json.Unmarshal(bodyBytes, &p2); err != nil {
			logrus.WithFields(logrus.Fields{"err": err.Error()})
			return c.JSON(http.StatusBadRequest, nil)
		}
		if p2.MessageType == "private" {
			return privateMessageHandler(c)
		}
	}

	return c.JSON(http.StatusOK, nil)
}

type requestPrivateMessage struct {
	SubType    string `json:"sub_type"`
	UserId     int64  `json:"user_id"`
	RawMessage string `json:"raw_message"`
}

var (
	QQAdmin = int64(123456) // 填写用于更新信息的 QQ 账号
)

func privateMessageHandler(c echo.Context) error {
	req := requestPrivateMessage{}
	if err := c.Bind(&req); err != nil {
		logrus.WithFields(logrus.Fields{"err": err.Error()})
		return c.JSON(http.StatusBadRequest, nil)
	}

	if req.UserId != QQAdmin {
		return c.JSON(http.StatusOK, nil)
	}
	go func(message string) {
		workBallotMessage(message)
	}(req.RawMessage)

	return c.JSON(http.StatusOK, nil)
}

func workBallotMessage(message string) {
	word := strings.Split(message, " ")
	if len(word) == 2 {
		member, vis := data.Member[word[0]]
		if vis == false {
			return
		}
		_, vis = data.Problem[data.Accept{
			CF:      word[0],
			Problem: word[1],
		}]
		if vis == true {
			logrus.WithFields(logrus.Fields{"info": "重复的AC", "member": member})
			return
		}
		data.Problem[data.Accept{
			CF:      word[0],
			Problem: word[1],
		}] = true
		volunteers, vis := data.Volunteer[member.Room]
		volunteers.It = volunteers.It + 1
		if volunteers.It == len(volunteers.Info) {
			volunteers.It = 0
		}
		volunteer := volunteers.Info[volunteers.It]
		if vis == false {
			qq.QQSendAndFindWhetherSuccess(QQAdmin, fmt.Sprintf("%s找不到志愿者！", member.Room))
		}
		ok := qq.QQSendAndFindWhetherSuccess(volunteer.QQ, fmt.Sprintf("姓名：%s，房间：%s，座位：%s\nAccpet：%s", member.Name, member.Room, member.Seat, word[1]))
		if ok == false {
			logrus.WithFields(logrus.Fields{"info": "信息发送失败", "member": member, "message": message})
		} else {
			qq.QQSendAndFindWhetherSuccess(QQAdmin, fmt.Sprintf("\"%s\"已发送给%s", message, volunteer.Name))
		}
	}
}
