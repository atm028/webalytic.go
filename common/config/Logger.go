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

func (o *Logger) GetLogger(modue string) bunyan.Logger {
	fmt.Println("GetLogget: name: ", o.Cfg.Name())
	config := bunyan.Config{
		Name: o.Cfg.Name(),
		Streams: []bunyan.Stream{
			{
				Name:   modue,
				Level:  o.Cfg.Level(),
				Stream: os.Stdout,
			},
			{
				Name:  modue,
				Level: o.Cfg.Level(),
				Path:  strings.Join([]string{o.Cfg.Path(), "/", o.Cfg.Name(), ".log"}, ""),
			},
		},
	}
	logger, err := bunyan.CreateLogger(config)
	if err != nil {
		fmt.Println("Cannot create logger: ", err)
	}
	return logger
}
