package statistics

import (
	"github.com/stathat/go"
)

type NullableMonitor struct{}

func NewNullableMonitor() Monitor {
	return &NullableMonitor{}
}

func (stats *NullableMonitor) Increment(string) {
}

func (stats *NullableMonitor) Close() {
}

type StatsHatMonitor struct{}

func NewStatsHatMonitor() Monitor {
	return &StatsHatMonitor{}
}

func (stats *StatsHatMonitor) Increment(metric string) {
	stathat.PostEZCount(metric, "0uzDCBeE2Ni9cCF5", 1)
}

func (stats *StatsHatMonitor) Close() {
}
