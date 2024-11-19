package core

import (
	"github.com/immafrady/stock-notifier/internal/core/broadcast"
	"github.com/immafrady/stock-notifier/internal/core/config"
	"github.com/immafrady/stock-notifier/internal/core/stock_data"
	"log"
	"runtime"
	"time"
)

type Core struct {
	ticker    *time.Ticker
	i         int
	Stocks    map[string]*stock_data.StockData
	Broadcast *broadcast.Broadcast
}

// Updates 更新数据
func (c *Core) Updates(t time.Time) {
	for _, stock := range c.Stocks {
		// 进入协程
		go stock.TryUpdate(c.i, t)
	}
}

func Run(cfgFile string) {
	c := config.NewConfig(cfgFile)
	core := &Core{}
	core.ticker = time.NewTicker(1 * time.Second)
	defer core.ticker.Stop()

	if c.Trackers == nil {
		log.Fatalln("没有关注的股票")
	} else {
		core.Stocks = make(map[string]*stock_data.StockData)
		for _, t := range c.Trackers {
			if stock := core.Stocks[t.Code]; stock != nil {
				log.Printf("查询码[%s]已存在，仅生效第一个配置", t.Code)
			} else {
				core.Stocks[t.Code] = stock_data.NewStockData(t)
			}
		}

		// 初始化更新
		log.Println("项目已运行，如无消息弹出，请检查消息是否被屏蔽")
		if runtime.GOOS == "darwin" {
			log.Println("mac系统请将【脚本编辑器】的消息权限打开，并将其添加到勿扰模式的白名单中")
		}

		// 设置播报
		core.Broadcast = broadcast.NewBroadcast(core.Stocks, c)
		// 开始循环
		count := 0
		for {
			select {
			case t := <-core.ticker.C:
				count++
				if count%5 == 0 {
					core.i++
					core.Updates(t)
				}
				if count%5 == 1 {
					// 下一秒展示消息
					stock_data.ShowNotifications()
				}
				core.Broadcast.Broadcast(t)
			}
		}
	}

}
