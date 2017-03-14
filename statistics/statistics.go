package statistics

type Monitor interface {
	Increment(metric string)
	Close()
}

type Statistics struct {
	monitor Monitor
}

func (s *Statistics) Increment(metric string) {
	s.monitor.Increment(metric)
}

func (s *Statistics) Stop() {
	s.monitor.Close()
}

func (s *Statistics) SetMonitor(monitor Monitor) {
	s.monitor = monitor
}

func NewStatistics() *Statistics {
	return &Statistics{
		monitor: NewStatsHatMonitor(),
	}
}

var statistics = NewStatistics()

func TrackMockRequest() {
	statistics.Increment("requests.mock")
}

func TrackConsoleRequest() {
	statistics.Increment("requests.console")
}

func TrackVerifyRequest() {
	statistics.Increment("requests.verify")
}

func TrackScenarioFeature() {
	statistics.Increment("feature.scenario")
}

func TrackProxyFeature() {
	statistics.Increment("feature.proxy")
}

func SetMonitor(monitor Monitor) {
	statistics.SetMonitor(monitor)
}

func Stop() {
	statistics.Stop()
}
