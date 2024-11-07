package core

import (
	"github.com/immafrady/stock-notifier/utils"
	"log"
	"strings"
	"time"
)

type Core struct {
	ticker *time.Ticker
	i      int
	Stocks []*StockData
}

// Updates 更新数据
func (c *Core) Updates(t time.Time, init bool) {
	for _, stock := range c.Stocks {
		if init || stock.shouldUpdate(c.i, t) {
			// 进入协程
			go stock.Update()
		}
	}
}

func Run(c *Config) {
	core := &Core{}
	core.ticker = time.NewTicker(5 * time.Second)
	defer core.ticker.Stop()

	if c.Trackers == nil {
		log.Fatalln("没有关注的股票")
	} else {
		core.Stocks = make([]*StockData, len(c.Trackers))
		var searchCodes []string
		for i, t := range c.Trackers {
			core.Stocks[i] = NewStockData(t)
			searchCodes = append(searchCodes, t.Code)
		}

		// 初始化更新
		core.Updates(time.Now(), true)
		utils.Notify("开始监控以下股票", strings.Join(searchCodes, ","))

		for {
			select {
			case t := <-core.ticker.C:
				core.i++
				core.Updates(t, false)
			}
		}
	}

}
