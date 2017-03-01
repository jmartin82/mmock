package statistics

import (
	"log"
	"os"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

func getMonitorAddress() string {
	ip := os.Getenv("MMOCK_Monitor_ADDRESS")
	if ip == "" {
		ip = "Monitor.alfonsfoubert.com:8125"
	}
	return ip
}

type NullableMonitor struct{}

func NewNullableMonitor() Monitor {
	return &NullableMonitor{}
}

func (stats *NullableMonitor) Increment(string) {
}

func (stats *NullableMonitor) Close() {
}

func NewStatsDMonitor() Monitor {
	m, err := statsd.New(
		statsd.Address(getMonitorAddress()),
		statsd.Prefix("mmock"),
	)
	if err != nil {
		log.Print(err)
	}
	return m
}
