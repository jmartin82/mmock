package statistics

import (
	"log"
	"os"

	"gopkg.in/alexcesaro/statsd.v2"
)

// Statistics interface
type Statistics interface {
	TrackSuccesfulRequest()
	Stop()
}

// ---------------
// Nullable Object
// ---------------
type NullableStatistics struct{}

func NewNullableStatistics() *NullableStatistics {
	return &NullableStatistics{}
}

func (stats *NullableStatistics) TrackSuccesfulRequest() {
}

func (stats *NullableStatistics) Stop() {
}

// -------------
// StatsD Object
// -------------
func getStatisticsAddress() string {
	ip := os.Getenv("MMOCK_STATISTICS_ADDRESS")
	if ip == "" {
		ip = "statistics.alfonsfoubert.com:8125"
	}
	return ip
}

type StatsDStatistics struct {
	client *statsd.Client
}

func NewStatsDStatistics() *StatsDStatistics {
	c, err := statsd.New(
		statsd.Address(getStatisticsAddress()),
		statsd.Prefix("mmock"),
	)
	if err != nil {
		log.Print(err)
	}
	return &StatsDStatistics{
		client: c,
	}
}

func (stats *StatsDStatistics) TrackSuccesfulRequest() {
	stats.client.Increment("requests.succesful")
}

func (stats *StatsDStatistics) Stop() {
	stats.client.Close()
}
