package core

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

// StockData 完整数据
type StockData struct {
	mutex     sync.Mutex
	Frequency int            // 格式化后的更新频率
	MaxLogs   int            // 最多的log数
	ApiData   *ApiData       // 单词数据
	Tracker   *ConfigTracker // 配置
	PriceLogs []*PriceLog
}

// PriceLog 价格日志
type PriceLog struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float32   `json:"price"`
}

// NewStockData 新建一个
func NewStockData(t *ConfigTracker) *StockData {
	data := &StockData{
		Tracker: t,
	}
	{
		// 处理缓存数
		if t.Continuous > 2 {
			data.MaxLogs = t.Continuous
		} else {
			data.MaxLogs = 2 // 默认两个才能发生比较
		}
		data.PriceLogs = make([]*PriceLog, data.MaxLogs)
	}
	{
		var setDefaultFrequency = func() {
			data.Frequency = 1
			log.Printf("[error]【%s】转换错误,将取默认值:1", t.Frequency)
		}

		// 处理更新频率
		if len(t.Frequency) < 2 {
			setDefaultFrequency()
		} else {
			l := len(t.Frequency)
			unit := t.Frequency[l-1]
			valueStr := t.Frequency[:l-1]
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
	return data
}

// 判断是否可以更新
func (s *StockData) shouldUpdate(i int, t time.Time) bool {
	if s.Frequency%i == 0 {
		// 根据频率判断
		return true
	}
	if s.Tracker.Updates != nil {
		if s.ApiData != nil {
			// 如果有获取过数据，以获取数据时为准
			t = s.ApiData.UpdateAt
		}
		for _, update := range s.Tracker.Updates {
			from, to := update.Range()
			if t.After(from) && t.Before(to) {
				return true
			}
		}
	} else {
		// 没有的时候也给允许吧
		return true
	}

	return false
}

// Update 核心更新逻辑
func (s *StockData) Update() {
	fmt.Println(s.Tracker.Code)
	if s.mutex.TryLock() {
		apiData := NewApiData(s.Tracker.Code)
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
		}
		fmt.Println(apiData)
	}
	s.mutex.Unlock()
}
