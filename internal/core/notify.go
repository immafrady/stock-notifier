package core

import (
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
)

type MsgRecord struct {
	Tag     string
	Title   string
	Message string
}

var m map[string][]utils.Msg
var msgChan chan MsgRecord

func init() {
	m = make(map[string][]utils.Msg)
	msgChan = make(chan MsgRecord)
	go func() {
		for {
			select {
			case msg := <-msgChan:
				m[msg.Tag] = append(m[msg.Tag], utils.Msg{Title: msg.Title, Message: msg.Message})
			}
		}
	}()
}

func SendToPending(tag, title, message string) {
	msgChan <- MsgRecord{Tag: tag, Title: title, Message: message}
}

func ShowNotifications() {
	for tag, msgs := range m {
		if len(msgs) > 1 {
			title := fmt.Sprintf("%s (%v条通知)", tag, len(msgs))
			var messages string
			for _, msg := range msgs {
				messages += fmt.Sprintf("* %s\n", msg.Title)
				if msg.Message != "" {
					messages += fmt.Sprintf("%s\n", msg.Message)
				}
				messages += "--------------\n"
			}
			utils.Notify(title, messages)
		} else if len(msgs) == 1 {
			utils.Notify(tag, fmt.Sprintf("* %s\n%s", msgs[0].Title, msgs[0].Message))
		}
	}
	m = make(map[string][]utils.Msg) // 重置
}
