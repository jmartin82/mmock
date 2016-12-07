package logging

import "fmt"

type ChannelLogger struct {
	ChannelLog chan string
}

func (logger ChannelLogger) Printf(format string, v ...interface{}) {
	logger.ChannelLog <- fmt.Sprintf(format, v...)
}

func (logger ChannelLogger) Print(v ...interface{}) {
	logger.ChannelLog <- fmt.Sprint(v...)
}

func (logger ChannelLogger) Println(v ...interface{}) {
	logger.ChannelLog <- fmt.Sprintln(v...)
}

func (logger ChannelLogger) Fatalf(format string, v ...interface{}) {
	logger.ChannelLog <- fmt.Sprintf(format, v...)
}

func (logger ChannelLogger) Fatal(v ...interface{}) {
	logger.ChannelLog <- fmt.Sprint(v...)
}

func (logger ChannelLogger) Fatalln(v ...interface{}) {
	logger.ChannelLog <- fmt.Sprintln(v...)
}
