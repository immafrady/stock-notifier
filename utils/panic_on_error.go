package utils

import (
	"fmt"
	"log"
)

func PanicOnError(err error, label string) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("[error]%s: %s", label, err.Error()))
	}
}
