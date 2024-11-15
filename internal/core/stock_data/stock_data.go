package stock_data

import (
	"fmt"
	"github.com/immafrady/stock-notifier/internal/core/config"
	"log"
	"strconv"
	"sync"
	"time"
)

// StockData 完整数据
type StockData struct {
	mutex     sync.Mutex
	Disable   bool                  // 如果遇到异常，禁用
	Frequency int                   // 格式化后的更新频率
	MaxLogs   int                   // 最多的log数
	ApiData   *ApiData              // 最新数据
	Config    *config.ConfigTracker // 配置
	Trackers  Trackers
	PriceLogs []*PriceLog
}

// PriceLog 价格日志
type PriceLog struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
}

// NewStockData 新建一个
func NewStockData(c *config.ConfigTracker) *StockData {
	data := &StockData{
		Config: c,
	}
	{
		// 处理缓存数
		if c.Continuous > 2 {
			data.MaxLogs = c.Continuous
		} else {
			data.MaxLogs = 2 // 默认两个才能发生比较
		}
	}
	{
		var setDefaultFrequency = func() {
			data.Frequency = 1
			log.Printf("[error]【%s】转换错误,将取默认值:1", c.Frequency)
		}

		// 处理更新频率
		if len(c.Frequency) < 2 {
			setDefaultFrequency()
		} else {
			l := len(c.Frequency)
			unit := c.Frequency[l-1]
			valueStr := c.Frequency[:l-1]
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				setDefaultFrequency()
			} else {
				durationInSeconds := 0
				switch unit {
				case 's':
					durationInSeconds = value
				case 'm':
					durationInSeconds = value * 60
				case 'h':
					durationInSeconds = value * 60 * 60
				}
				if durationInSeconds == 0 {
					setDefaultFrequency()
				} else {
					data.Frequency = durationInSeconds / 5
				}
			}
		}
	}
	go func() {
		data.update()
		// 延迟挂载监听
		data.Trackers = NewTrackers(data)
	}()
	return data
}

// shouldUpdate 判断是否可以更新
func (s *StockData) shouldUpdate(i int, t time.Time) bool {
	if s.Disable {
		return false
	}
	if i%s.Frequency == 0 {
		// 根据频率判断，再进入下层判断
		if s.Config.Updates != nil {
			if s.ApiData != nil {
				// 如果有获取过数据，以获取数据时为准
				if t.Sub(s.ApiData.UpdateAt).Abs() < time.Minute {
					t = s.ApiData.UpdateAt
				}
			}
			for _, update := range s.Config.Updates {
				from, to := update.Range()
				if t.After(from) && t.Before(to) {
					return true
				}
			}
		} else {
			// 没有的时候也给允许吧
			return true
		}
	}

	return false
}

// update 实际的update方法
func (s *StockData) update() {
	if s.mutex.TryLock() {
		apiData, done := NewApiData(s.Config.Code)
		if apiData == nil {
			if done {
				// 没有数据&终止，remove
				s.Disable = true
				log.Printf("代码：%s查询失败，停止查询", s.Config.Code)
			}
		} else {
			if s.ApiData == nil || apiData.UpdateAt.After(s.ApiData.UpdateAt) {
				// 只有当更新时间大于最新时间时，才会取更新
				s.ApiData = apiData
				// 更新日志
				s.PriceLogs = append(s.PriceLogs, &PriceLog{
					Timestamp: apiData.UpdateAt,
					Price:     apiData.Current,
				})
				if len(s.PriceLogs) > s.MaxLogs {
					s.PriceLogs = s.PriceLogs[:s.MaxLogs]
				}
				// 开始触发监控
				s.Trackers.Evaluate(s)
			}
		}
		s.mutex.Unlock()
	}
}

// TryUpdate 核心更新逻辑
func (s *StockData) TryUpdate(i int, t time.Time) {
	if s.shouldUpdate(i, t) {
		s.update()
	}
}

// Shout 发送通知
func (s *StockData) Shout(title, message string) {
	base := fmt.Sprintf("%s %s %s(%.2f%%)",
		s.ApiData.UpdateAt.Format("15:04:05"),
		s.ApiData.Name,
		s.ApiData.ParsePrice(s.ApiData.Current),
		s.ApiData.Percentage,
	)
	SendToPending(base, title, message)
}
