package logging

import "log"

var logger ChannelLogger

func SetLogger(channelLogger ChannelLogger) {
	logger = channelLogger
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
	if logger != (ChannelLogger{}) {
		logger.Printf(format, v...)
	}
}

func Print(v ...interface{}) {
	log.Print(v...)
	if logger != (ChannelLogger{}) {
		logger.Print(v...)
	}
}

func Println(v ...interface{}) {
	log.Println(v...)
	if logger != (ChannelLogger{}) {
		logger.Println(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
	if logger != (ChannelLogger{}) {
		logger.Fatalf(format, v...)
	}
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
	if logger != (ChannelLogger{}) {
		logger.Fatal(v...)
	}
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
	if logger != (ChannelLogger{}) {
		logger.Fatalln(v...)
	}
}
