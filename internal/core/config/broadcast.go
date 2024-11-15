package config

// Broadcast 播报时间点
type Broadcast struct {
	Time  TimeString `json:"time" yaml:"time"`
	Label string     `json:"label" yaml:"label"`
}
