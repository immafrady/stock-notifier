package core

import (
	"fmt"
	"math"
)

type Tracker struct {
	priceDiff   bool
	percentDiff bool
	continuous  bool
	targetPrice bool
}

// TrackPriceDiff 监控价格差
func (s *StockData) TrackPriceDiff() {
	if s.Config.PriceDiff != 0 {
		l := len(s.PriceLogs)
		if l > 2 {
			prev := s.PriceLogs[l-2].Price
			curr := s.PriceLogs[l-1].Price
			diff := curr - prev
			if math.Abs(diff) > s.Config.PriceDiff && !s.Tracker.priceDiff {
				s.Tracker.priceDiff = true
				s.Shout("当前差价已超监控值", fmt.Sprintf("监控时间段内价格变化为： %s", s.ApiData.ParsePrice(diff)))
			} else {
				s.Tracker.priceDiff = false
			}

		}
	}
}

// TrackPercentDiff 监控百分位差
func (s *StockData) TrackPercentDiff() {
	if s.Config.PercentDiff != 0 {
		l := len(s.PriceLogs)
		if l > 2 {
			prev := s.PriceLogs[l-2].Price
			curr := s.PriceLogs[l-1].Price
			diff := (curr - prev) / prev
			if math.Abs(diff) > s.Config.PercentDiff && !s.Tracker.percentDiff {
				s.Tracker.percentDiff = true
				s.Shout("当前异动已超监控值", fmt.Sprintf("监控时间段内价格变化为： %0.2f%%", diff))
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
		} else {
			s.Tracker.continuous = false
		}
	}
}

func (s *StockData) TrackTargetPrice() {}
