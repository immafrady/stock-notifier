package core

import (
	"fmt"
	"math"
	"strings"
)

type Tracker struct {
	welcome         bool
	priceDiff       bool
	percentDiff     bool
	continuous      bool // deprecated
	targetHighPrice int
	targetLowPrice  int
}

// TrackWelcome 首次弹出提醒
func (s *StockData) TrackWelcome() {
	if !s.Tracker.welcome {
		s.Tracker.welcome = true
		s.Shout("开始监控", "")
	}
}

// TrackPriceDiff 监控价格差
func (s *StockData) TrackPriceDiff() {
	if s.Config.PriceDiff > 0 {
		l := len(s.PriceLogs)
		if l > 2 {
			prev := s.PriceLogs[l-2].Price
			curr := s.PriceLogs[l-1].Price
			diff := curr - prev
			if math.Abs(diff) > s.Config.PriceDiff {
				if !s.Tracker.priceDiff {
					s.Tracker.priceDiff = true
					s.Shout("当前差价已超监控值", fmt.Sprintf("监控时间段内价格变化为： %s", s.ApiData.ParsePrice(diff)))
				}
			} else {
				s.Tracker.priceDiff = false
			}

		}
	}
}

// TrackPercentDiff 监控百分位差
func (s *StockData) TrackPercentDiff() {
	if s.PercentDiff != 0 {
		l := len(s.PriceLogs)
		if l > 2 {
			prev := s.PriceLogs[l-2].Price
			curr := s.PriceLogs[l-1].Price
			diff := (curr - prev) / prev
			if math.Abs(diff) > s.PercentDiff {
				s.Shout("当前异动已超监控值", fmt.Sprintf("监控时间段内价格变化为： %0.2f%%\n%s➡️%s", diff*100, s.ApiData.ParsePrice(prev), s.ApiData.ParsePrice(curr)))
			}

		}
	}
}

// TrackContinuous 监控连续
func (s *StockData) TrackContinuous() {
	l := len(s.PriceLogs)
	if s.Config.Continuous > 2 && l == s.Config.Continuous {

		isIncreasing := true
		isDecreasing := true

		for i := 1; i < len(s.PriceLogs); i++ {
			if s.PriceLogs[i].Price > s.PriceLogs[i-1].Price {
				isDecreasing = false // 如果出现上涨，则不是下跌
			} else if s.PriceLogs[i].Price < s.PriceLogs[i-1].Price {
				isIncreasing = false // 如果出现下跌，则不是上涨
			} else {
				return // 平价直接再见
			}
		}

		var trend string
		if isIncreasing {
			trend = "上涨"
		} else if isDecreasing {
			trend = "下跌"
		} else {
			return
		}
		first := s.PriceLogs[0].Price
		last := s.PriceLogs[len(s.PriceLogs)-1].Price
		diff := last - first
		rate := diff / first
		list := make([]string, len(s.PriceLogs))
		for i, p := range s.PriceLogs {
			list[i] = s.ApiData.ParsePrice(p.Price)
		}
		s.Shout(
			fmt.Sprintf("连续%v次记录呈现单边%s走势", s.Config.Continuous, trend),
			fmt.Sprintf("监控时间段内价格变化为： %s(%0.2f%%)\n%v", s.ApiData.ParsePrice(diff), rate, strings.Join(list, "➡️")),
		)
	}
}

// TrackTargetHighPrice 超越目标价位
func (s *StockData) TrackTargetHighPrice() {
	if s.Config.TargetHighPrice > 0 {
		if s.ApiData.Current > s.Config.TargetHighPrice && s.Tracker.targetHighPrice != 2 {
			s.Tracker.targetHighPrice = 2
			s.Shout(
				fmt.Sprintf("⬆️超越目标高位%s :)", s.ApiData.ParsePrice(s.Config.TargetHighPrice)),
				"",
			)
		} else if s.ApiData.Current == s.Config.TargetHighPrice && s.Tracker.targetHighPrice != 1 {
			s.Tracker.targetHighPrice = 1
			s.Shout(
				fmt.Sprintf("⏺️达到目标高位%s", s.ApiData.ParsePrice(s.Config.TargetHighPrice)),
				"",
			)
		} else if s.ApiData.Current < s.Config.TargetHighPrice && s.Tracker.targetHighPrice != 0 {
			s.Tracker.targetHighPrice = 0
			s.Shout(
				fmt.Sprintf("⬇️跌破目标高位%s :(", s.ApiData.ParsePrice(s.Config.TargetHighPrice)),
				"",
			)
		}
	}
}

// TrackTargetLowPrice 跌破目标价位
func (s *StockData) TrackTargetLowPrice() {
	if s.Config.TargetLowPrice > 0 {
		if s.ApiData.Current < s.Config.TargetLowPrice && s.Tracker.targetLowPrice != 2 {
			s.Tracker.targetLowPrice = 2
			s.Shout(
				fmt.Sprintf("⬇️跌破目标低位%s :(", s.ApiData.ParsePrice(s.Config.TargetLowPrice)),
				"",
			)
		} else if s.ApiData.Current == s.Config.TargetLowPrice && s.Tracker.targetLowPrice != 1 {
			s.Tracker.targetLowPrice = 1
			s.Shout(
				fmt.Sprintf("⏺️️达到目标低位%s", s.ApiData.ParsePrice(s.Config.TargetLowPrice)),
				"",
			)
		} else if s.ApiData.Current > s.Config.TargetLowPrice && s.Tracker.targetLowPrice != 0 {
			s.Tracker.targetLowPrice = 0
			s.Shout(
				fmt.Sprintf("⬆️超越目标低位%s :)", s.ApiData.ParsePrice(s.Config.TargetLowPrice)),
				"",
			)
		}
	}
}
