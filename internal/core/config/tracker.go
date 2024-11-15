package config

// Tracker 跟踪股票配置
type Tracker struct {
	Code                string    `json:"code" yaml:"code"`                               // 股市编码
	Frequency           string    `json:"frequency" yaml:"frequency"`                     // 轮循频次，最低5s (可选单位 s:秒 m:分钟 h:小时)
	Continuous          int       `json:"continuous" yaml:"continuous"`                   // 轮询频次内连续上涨/下跌x次时弹出提示
	PriceDiff           float64   `json:"priceDiff" yaml:"priceDiff"`                     // 两次轮询【差价】超过该数值时弹出提示
	PercentDiff         string    `json:"percentDiff" yaml:"percentDiff"`                 // 两次轮询【相差百分比】超过该数值时弹出提示
	RealTimePercentDiff string    `json:"realTimePercentDiff" yaml:"realTimePercentDiff"` // 溢价/折价率超过x%时弹出提示（仅ETF生效）
	TargetHighPrice     float64   `json:"targetHighPrice" yaml:"targetHighPrice"`         // 多于越目标价位时提醒(0表示不监控)
	TargetLowPrice      float64   `json:"targetLowPrice" yaml:"targetLowPrice"`           // 少于目标价位时提醒(0表示不监控)
	Updates             []*Update `json:"updates,omitempty" yaml:"updates"`               // 更新时间段（不传时取默认值）
}
