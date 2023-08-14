package main

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	WorkDurationMins       int  `json:"work_duration_mins"`
	ShortBreakDurationMins int  `json:"short_break_duration_mins"`
	LongBreakDurationMins  int  `json:"long_break_duration_mins"`
	ShortBreaks            int  `json:"short_breaks"`
	SendNotification       bool `json:"send_notification"`
}

func (c Config) WorkDuration() time.Duration {
	return time.Minute * time.Duration(c.WorkDurationMins)
}

func (c Config) ShortBreakDuration() time.Duration {
	return time.Minute * time.Duration(c.ShortBreakDurationMins)
}

func (c Config) LongBreakDuration() time.Duration {
	return time.Minute * time.Duration(c.LongBreakDurationMins)
}

func ReadConfig() (Config, error) {
	var config Config

	f, err := os.ReadFile("./config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal([]byte(f), &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
