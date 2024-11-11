package core

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"strings"
	"time"
)

// SetupBroadcast 设置播报
func (c *Core) SetupBroadcast(config *Config) {
	if config.Broadcast != nil {
		for _, broadcast := range config.Broadcast {
			scheduler := gocron.NewScheduler(time.Local)
			var job *gocron.Job
			job, _ = scheduler.Every(1).Day().At(broadcast.Time).Do(func() {
				var msgs []string
				for _, s := range c.Stocks {
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
					broadcast.Label,
					fmt.Sprintf("监控%v个股票", len(c.Stocks)),
					strings.Join(msgs, "------------\n"),
				)
				log.Printf("本次播报完成，下次播报【%s】时间: %v", broadcast.Label, job.NextRun())
			})
			scheduler.StartAsync()
			log.Printf("播报【%s】时间: %v", broadcast.Label, job.NextRun())
		}
	}
}
