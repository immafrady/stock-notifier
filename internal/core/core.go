package core

import (
	"log"
	"time"
)

type Core struct {
	ticker *time.Ticker
	i      int
	Stocks []*StockData
}

// Updates 更新数据
func (c *Core) Updates(t time.Time, force bool) {
	for _, stock := range c.Stocks {
		if force || stock.shouldUpdate(c.i, t) {
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
		for i, t := range c.Trackers {
			core.Stocks[i] = NewStockData(t)
		}

		// 初始化更新
		core.Updates(time.Now(), true)

		for {
			select {
			case t := <-core.ticker.C:
				core.i++
				core.Updates(t, false)
			}
		}
	}

}
