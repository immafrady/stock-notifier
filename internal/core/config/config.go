package config

import (
	"github.com/immafrady/stock-notifier/utils"
	"github.com/marcozac/go-jsonc"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

type Config struct {
	Updates   []*Update    `json:"updates" yaml:"updates"`
	Broadcast []*Broadcast `json:"broadcast" yaml:"broadcast"`
	Trackers  []*Tracker   `json:"trackers" yaml:"trackers"`
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
