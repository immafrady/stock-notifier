package config

import (
	"github.com/immafrady/stock-notifier/utils"
	"github.com/marcozac/go-jsonc"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TimeString 短的时间戳
type TimeString string

func (t TimeString) Get() (h, m int) {
	strs := strings.Split(string(t), ":")
	if len(strs) != 2 {
		return 0, 0
	} else {
		h, _ = strconv.Atoi(strs[0])
		m, _ = strconv.Atoi(strs[1])
		return
	}
}

// ConfigUpdate 更新时间段
type ConfigUpdate struct {
	From TimeString `json:"from" yaml:"from"`
	To   TimeString `json:"to" yaml:"to"`
}

// ConfigBroadcast 播报时间点
type ConfigBroadcast struct {
	Time  TimeString `json:"time" yaml:"time"`
	Label string     `json:"label" yaml:"label"`
}

func (u ConfigUpdate) Range() (from, to time.Time) {
	now := time.Now()
	year, month, day := now.Date()
	fh, fm := u.From.Get()
	th, tm := u.To.Get()
	return time.Date(year, month, day, fh, fm, 0, 0, now.Location()),
		time.Date(year, month, day, th, tm, 0, 0, now.Location())
}

// ConfigTracker 跟踪股票配置
type ConfigTracker struct {
	Code                string          `json:"code" yaml:"code"`                               // 股市编码
	Frequency           string          `json:"frequency" yaml:"frequency"`                     // 轮循频次，最低5s (可选单位 s:秒 m:分钟 h:小时)
	Continuous          int             `json:"continuous" yaml:"continuous"`                   // 轮询频次内连续上涨/下跌x次时弹出提示
	PriceDiff           float64         `json:"priceDiff" yaml:"priceDiff"`                     // 两次轮询【差价】超过该数值时弹出提示
	PercentDiff         string          `json:"percentDiff" yaml:"percentDiff"`                 // 两次轮询【相差百分比】超过该数值时弹出提示
	RealTimePercentDiff string          `json:"realTimePercentDiff" yaml:"realTimePercentDiff"` // 溢价/折价率超过x%时弹出提示（仅ETF生效）
	TargetHighPrice     float64         `json:"targetHighPrice" yaml:"targetHighPrice"`         // 多于越目标价位时提醒(0表示不监控)
	TargetLowPrice      float64         `json:"targetLowPrice" yaml:"targetLowPrice"`           // 少于目标价位时提醒(0表示不监控)
	Updates             []*ConfigUpdate `json:"updates,omitempty" yaml:"updates"`               // 更新时间段（不传时取默认值）
}

type Config struct {
	Updates   []*ConfigUpdate    `json:"updates" yaml:"updates"`
	Broadcast []*ConfigBroadcast `json:"broadcast" yaml:"broadcast"`
	Trackers  []*ConfigTracker   `json:"trackers" yaml:"trackers"`
}

// newYamlConfig yaml的配置
func newYamlConfig(bytes []byte) (c *Config) {
	err := yaml.Unmarshal(bytes, &c)
	if err == nil {
		return c
	}

	utils.PanicOnError(err, "配置yaml文件格式错误")
	return
}

// newJsoncConfig jsonc的配置
func newJsoncConfig(bytes []byte) (c *Config) {
	err := jsonc.Unmarshal(bytes, &c)
	if err == nil {
		return c
	}

	utils.PanicOnError(err, "配置json文件格式错误")
	return
}

// NewConfig 新配置
func NewConfig(p string) (c *Config) {
	var (
		err  error
		data []byte
	)
	p, err = filepath.Abs(p)
	utils.PanicOnError(err, "获取绝对路径失败")

	data, err = utils.ReadFile(p)
	utils.PanicOnError(err, "读取文件失败")

	p = strings.ToLower(p)

	if strings.HasSuffix(p, "jsonc") || strings.HasSuffix(p, "json") {
		c = newJsoncConfig(data)
	} else if strings.HasSuffix(p, "yaml") || strings.HasSuffix(p, "yml") {
		c = newYamlConfig(data)
	} else {
		return nil
	}

	// 将更新时间附上
	for _, tracker := range c.Trackers {
		if tracker.Updates == nil {
			tracker.Updates = c.Updates
		}
	}
	return
}
