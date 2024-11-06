package tracker

// ConfigUpdate 更新时间段
type ConfigUpdate struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

// ConfigTracker 跟踪股票配置
type ConfigTracker struct {
	Code        string          `json:"code" yaml:"code"`               // 股市编码
	Frequency   string          `json:"frequency" yaml:"frequency"`     // 轮循频次，最低5s (可选单位 s:秒 m:分钟 h:小时)
	Continuous  string          `json:"continuous" yaml:"continuous"`   // 轮询频次内连续上涨/下跌x次时弹出提示
	PriceDiff   string          `json:"priceDiff" yaml:"priceDiff"`     // 两次轮询【差价】超过该数值时弹出提示
	PercentDiff string          `json:"percentDiff" yaml:"percentDiff"` // 两次轮询【相差百分比】超过该数值时弹出提示
	Updates     *[]ConfigUpdate `json:"updates" yaml:"updates"`         // 更新时间段（不传时取默认值）
}

type Config struct {
	Updates *[]ConfigUpdate `json:"updates" yaml:"updates"`
	Tracker *ConfigTracker  `json:"tracker" yaml:"tracker"`
}
