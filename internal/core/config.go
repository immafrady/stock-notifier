package core

import (
	"github.com/immafrady/stock-notifier/utils"
	"github.com/marcozac/go-jsonc"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

// ConfigUpdate 更新时间段
type ConfigUpdate struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

// ConfigTracker 跟踪股票配置
type ConfigTracker struct {
	Code        string          `json:"code" yaml:"code"`                 // 股市编码
	Frequency   string          `json:"frequency" yaml:"frequency"`       // 轮循频次，最低5s (可选单位 s:秒 m:分钟 h:小时)
	Continuous  int             `json:"continuous" yaml:"continuous"`     // 轮询频次内连续上涨/下跌x次时弹出提示
	PriceDiff   int             `json:"priceDiff" yaml:"priceDiff"`       // 两次轮询【差价】超过该数值时弹出提示
	PercentDiff int             `json:"percentDiff" yaml:"percentDiff"`   // 两次轮询【相差百分比】超过该数值时弹出提示
	Updates     []*ConfigUpdate `json:"updates,omitempty" yaml:"updates"` // 更新时间段（不传时取默认值）
}

type Config struct {
	Updates  []*ConfigUpdate  `json:"updates" yaml:"updates"`
	Trackers []*ConfigTracker `json:"trackers" yaml:"trackers"`
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
func NewConfig(p string) *Config {
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
		return newJsoncConfig(data)
	} else if strings.HasSuffix(p, "yaml") || strings.HasSuffix(p, "yml") {
		return newYamlConfig(data)
	} else {
		return nil
	}
}
