package config

import "time"

// Update 更新时间段
type Update struct {
	From TimeString `json:"from" yaml:"from"`
	To   TimeString `json:"to" yaml:"to"`
}

func (u Update) Range() (from, to time.Time) {
	now := time.Now()
	year, month, day := now.Date()
	fh, fm := u.From.Get()
	th, tm := u.To.Get()
	return time.Date(year, month, day, fh, fm, 0, 0, now.Location()),
		time.Date(year, month, day, th, tm, 0, 0, now.Location())
}
