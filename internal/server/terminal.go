package server

import "github.com/fatih/color"

type terminalColors struct {
	Info    func(a ...interface{}) string
	Success func(a ...interface{}) string
	Error   func(a ...interface{}) string
}

var terminal = terminalColors{
	Info:    color.New(color.FgBlue).SprintFunc(),
	Success: color.New(color.Bold, color.FgGreen).SprintFunc(),
	Error:   color.New(color.Bold, color.FgRed).SprintFunc(),
}
