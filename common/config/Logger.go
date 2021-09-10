package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/bhoriuchi/go-bunyan/bunyan"
)

type ILogger interface {
	GetLogger(module string) bunyan.Logger
}

type Logger struct {
	Cfg ILoggerConfig
}

func (o *Logger) GetLogger(module string) bunyan.Logger {
	path := o.Cfg.Path()
	name := o.Cfg.Name()
	level := o.Cfg.Level()
	host, err := os.Hostname()
	if err != nil {
		fmt.Println("Cannot get hostname")
		host = "localhost"
	}
	fullPath := strings.Join([]string{
		o.Cfg.Path(),
		"/", o.Cfg.Name(),
		"_", module,
		"_", host,
		".log"}, "")
	fmt.Println("GetLogget: hostName: ", host)
	fmt.Println("GetLogget: name: ", name)
	fmt.Println("GetLogget: path: ", path)
	fmt.Println("GetLogget: fullPath: ", fullPath)
	fmt.Println("GetLogget: level: ", level)
	config := bunyan.Config{
		Name: o.Cfg.Name(),
		Streams: []bunyan.Stream{
			{
				Name:   module,
				Level:  o.Cfg.Level(),
				Stream: os.Stdout,
			},
			{
				Name:  module,
				Level: o.Cfg.Level(),
				Path:  fullPath,
			},
		},
	}
	logger, err := bunyan.CreateLogger(config)
	if err != nil {
		fmt.Println("Cannot create logger: ", err)
	}
	return logger
}
