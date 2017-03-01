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
		monitor: NewStatsDMonitor(),
	}
}

var statistics = NewStatistics()

func TrackSuccesfulRequest() {
	statistics.Increment("requests.succesful")
}

func Stop() {
	statistics.Stop()
}

func Disable() {
	statistics.SetMonitor(NewNullableMonitor())
}
