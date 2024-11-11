package core

import (
	"log"
	"runtime"
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
	core.ticker = time.NewTicker(1 * time.Second)
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
		log.Println("项目已运行，如无消息弹出，请检查消息是否被屏蔽")
		if runtime.GOOS == "darwin" {
			log.Println("mac系统请将【脚本编辑器】的消息权限打开，并将其添加到勿扰模式的白名单中")
		}
		// 设置播报
		core.SetupBroadcast(c)
		// 开始循环
		count := 0
		for {
			select {
			case t := <-core.ticker.C:
				count++
				if count%5 == 0 {
					core.i++
					core.Updates(t, false)
				}
				if count%5 == 1 {
					// 下一秒展示消息
					ShowNotifications()
				}
			}
		}
	}

}
