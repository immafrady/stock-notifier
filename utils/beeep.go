package utils

import (
	"github.com/gen2brain/beeep"
	"log"
)

type Msg struct {
	Title   string
	Message string
}

func Notify(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		log.Println("[error]发送消息失败", err)
	}
}
