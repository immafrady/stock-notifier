package core

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type BroadcastTrigger struct {
	*ConfigBroadcast
	Next time.Time
}

type Broadcast struct {
	Stocks   []*StockData
	Triggers []*BroadcastTrigger
}

// NewBroadcast 新增播报
func NewBroadcast(c *Core, config *Config) *Broadcast {
	b := &Broadcast{
		Stocks:   c.Stocks,
		Triggers: make([]*BroadcastTrigger, len(config.Broadcast)),
	}
	t := time.Now()
	for i, broadcast := range config.Broadcast {
		hour, minute := broadcast.Time.Get()
		next := time.Date(t.Year(), t.Month(), t.Day(), hour, minute, 0, 0, t.Location())
		if t.After(next) {
			next = next.Add(24 * time.Hour)
		}
		b.Triggers[i] = &BroadcastTrigger{
			ConfigBroadcast: broadcast,
			Next:            next,
		}
		log.Printf("播报【%s】时间: %v", broadcast.Label, next)
	}
	return b
}

// Broadcast 轮询并播报
func (b *Broadcast) Broadcast(t time.Time) {
	for _, trigger := range b.Triggers {
		if t.After(trigger.Next) {
			var msgs []string
			for _, s := range b.Stocks {
				msgs = append(msgs, fmt.Sprintf("%s %s(%.2f%% %s): \n开盘价: %s 昨收价: %s\n最高价: %s 最低价: %s\n",
					s.ApiData.Name,
					s.ApiData.ParsePrice(s.ApiData.Current),
					s.ApiData.Percentage,
					s.ApiData.ParsePrice(s.ApiData.Diff),
					s.ApiData.ParsePrice(s.ApiData.Opening),
					s.ApiData.ParsePrice(s.ApiData.YesterdayClosed),
					s.ApiData.ParsePrice(s.ApiData.High),
					s.ApiData.ParsePrice(s.ApiData.Low),
				))
			}
			SendToPending(
				trigger.Label,
				fmt.Sprintf("监控%v个股票", len(b.Stocks)),
				strings.Join(msgs, "------------\n"),
			)

			// 增加时间
			trigger.Next = trigger.Next.Add(24 * time.Hour)
			log.Printf("本次播报完成，下次播报【%s】时间: %v", trigger.Label, trigger.Next)
		}
	}
}
