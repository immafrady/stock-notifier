package core

import (
	"fmt"
	"math"
)

type Tracker struct {
	welcome         bool
	priceDiff       bool
	percentDiff     bool
	continuous      bool
	targetHighPrice bool
	targetLowPrice  bool
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
	if s.Config.PercentDiff > 0 {
		l := len(s.PriceLogs)
		if l > 2 {
			prev := s.PriceLogs[l-2].Price
			curr := s.PriceLogs[l-1].Price
			diff := (curr - prev) / prev
			if math.Abs(diff) > s.Config.PercentDiff {
				if !s.Tracker.percentDiff {
					s.Tracker.percentDiff = true
					s.Shout("当前异动已超监控值", fmt.Sprintf("监控时间段内价格变化为： %0.2f%%", diff))
				}
			} else {
				s.Tracker.percentDiff = false
			}

		}
	}
}

// TrackContinuous 监控连续
func (s *StockData) TrackContinuous() {
	l := len(s.PriceLogs)
	if s.Config.Continuous > 2 && l == s.Config.Continuous {
		prev := s.PriceLogs[0].Price
		curr := s.PriceLogs[1].Price
		if prev != curr {
			up := curr-prev > 0

			for i := 2; i < l; i++ {
				prev = s.PriceLogs[i-1].Price
				curr = s.PriceLogs[i].Price
				newUp := curr-prev > 0
				if up != newUp {
					// 重置
					s.Tracker.continuous = false
					return
				}
			}
		}

		if !s.Tracker.continuous {
			s.Tracker.continuous = true
			first := s.PriceLogs[0].Price
			last := s.PriceLogs[l-1].Price
			diff := last - first
			s.Shout(
				fmt.Sprintf("连续%v次记录呈现单边走势", s.Config.Continuous),
				fmt.Sprintf("监控时间段内价格变化为： %s(%0.2f%%)", s.ApiData.ParsePrice(diff), diff/first),
			)
		}
	}
}

// TrackTargetHighPrice 超越目标价位
func (s *StockData) TrackTargetHighPrice() {
	if s.Config.TargetHighPrice > 0 {
		if s.ApiData.Current > s.Config.TargetHighPrice {
			if !s.Tracker.targetHighPrice {
				s.Tracker.targetHighPrice = true
				s.Shout(
					fmt.Sprintf("超越目标价位%s", s.ApiData.ParsePrice(s.Config.TargetHighPrice)),
					"",
				)
			}
		} else {
			s.Tracker.targetHighPrice = false
		}
	}
}

// TrackTargetLowPrice 跌破目标价位
func (s *StockData) TrackTargetLowPrice() {
	if s.Config.TargetLowPrice > 0 {
		if s.ApiData.Current < s.Config.TargetLowPrice {
			if !s.Tracker.targetLowPrice {
				s.Tracker.targetLowPrice = true
				s.Shout(
					fmt.Sprintf("跌破目标价位%s", s.ApiData.ParsePrice(s.Config.TargetLowPrice)),
					"",
				)
			}
		} else {
			s.Tracker.targetLowPrice = false
		}
	}
}
