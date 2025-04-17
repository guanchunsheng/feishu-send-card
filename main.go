// 通过群组机器人向群里边发卡片
// 配置文件conf.json中，webhook和secret是群组机器人的配置，cardID和version是卡片的配置
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type sendCard struct {
	webhook   string
	timestamp int64
	secret    string
	signature string
	cardID    string
	version   string
}

// 群组机器人api：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
// 卡片发送api https://open.feishu.cn/document/feishu-cards/quick-start/send-message-cards-with-custom-bot
type msg struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Card      struct {
		Type string `json:"type"`
		Data struct {
			TemplateID          string `json:"template_id"`
			TemplateVersionName string `json:"template_version_name"`
		} `json:"data"`
	} `json:"card"`
}

func main() {
	sd := sendCard{}
	sd.readConf()

	var err error
	sd.timestamp = time.Now().Unix()
	sd.signature, err = genSign(sd.secret, sd.timestamp)
	if err != nil {
		panic(fmt.Errorf("签名生成失败：%w", err))
	}
	sd.send()
}

func (s *sendCard) readConf() {
	conf := viper.New()
	conf.AddConfigPath("./")
	conf.SetConfigName("conf")
	conf.SetConfigType("json")
	if err := conf.ReadInConfig(); err != nil {
		panic(fmt.Errorf("conf.ini文件读取失败，请检查配置文件是否正常:\n%w", err))
	}

	s.webhook = conf.GetString("webhook")
	s.secret = conf.GetString("secret")
	s.cardID = conf.GetString("cardid")
	s.version = conf.GetString("version")

	if len(s.webhook) == 0 {
		panic(fmt.Errorf("webhook为空，请检查配置文件conf.json"))
	}
	if len(s.secret) == 0 {
		panic(fmt.Errorf("secret为空，请检查配置文件conf.json"))
	}
	if len(s.cardID) == 0 {
		panic(fmt.Errorf("cardid为空，请检查配置文件conf.json"))
	}
	if len(s.version) == 0 {
		panic(fmt.Errorf("version为空，请检查配置文件conf.json"))
	}
}

func (s *sendCard) send() {
	request := msg{}
	request.Timestamp = strconv.FormatInt(s.timestamp, 10)
	request.Sign = s.signature
	request.MsgType = "interactive"
	request.Card.Type = "template"
	request.Card.Data.TemplateID = s.cardID
	request.Card.Data.TemplateVersionName = s.version

	buf, err := json.Marshal(request)
	if err != nil {
		panic(fmt.Errorf("请求body构造失败：%w", err))
	}

	fmt.Println(bytes.NewBuffer(buf))

	resp, err := http.Post(s.webhook, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		panic(fmt.Errorf("构建HTTP请求失败：%w", err))
	}
	response, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(fmt.Errorf("解析返回结果失败：%w", err))
	}
	fmt.Println(string(response))

}

// 飞书提供的signature库
func genSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
