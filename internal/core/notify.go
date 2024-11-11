package core

import (
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
	"log"
	"strings"
)

type MsgRecord struct {
	Tag      string
	Title    string
	Messages []string
}

var m map[string][]MsgRecord
var msgChan chan MsgRecord

func init() {
	m = make(map[string][]MsgRecord)
	msgChan = make(chan MsgRecord)
	go func() {
		for {
			select {
			case msg := <-msgChan:
				m[msg.Tag] = append(m[msg.Tag], msg)
			}
		}
	}()
}

func SendToPending(tag, title, message string) {
	rawMessages := strings.Split(message, "\n")
	var messages []string
	for _, rm := range rawMessages {
		if rm != "" {
			messages = append(messages, rm)
		}
	}
	log.Printf("发送消息\n> tag: %s\n> title: %s\n> message: %s\n\n", tag, title, strings.Join(messages, "\n > "))
	msgChan <- MsgRecord{Tag: tag, Title: title, Messages: messages}
}

func ShowNotifications() {
	for tag, msgs := range m {
		title := tag
		if len(msgs) > 1 {
			title += fmt.Sprintf(" (%v条通知)", len(msgs))
		}

		var messages []string
		for _, msg := range msgs {
			message := fmt.Sprintf("# %s\n", msg.Title)
			if len(msg.Messages) > 0 {
				for _, content := range msg.Messages {
					message += fmt.Sprintf("> %s\n", content)
				}
			}
			messages = append(messages, message)
		}
		utils.Notify(title, strings.Join(messages, "------------\n"))
	}
	m = make(map[string][]MsgRecord) // 重置
}
