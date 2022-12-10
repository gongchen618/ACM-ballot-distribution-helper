package qq

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type qqSendResponse struct {
	Status string `json:"status"`
}

func QQSendAndFindWhetherSuccess(userId int64, message string) bool {
	sendMsgUrl := "http://127.0.0.1:5722/send_private_msg"
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(nil))
	if err != nil {
		logrus.Info(err, userId, "message send1")
		return false
	}

	q := req.URL.Query()
	q.Add("user_id", strconv.FormatInt(userId, 10))
	q.Add("message", message)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Info(err, userId, "message send2")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	qqSendResp := qqSendResponse{}
	if err = json.Unmarshal(body, &qqSendResp); err != nil || qqSendResp.Status == "failed" {
		return false
	}
	return true
}
