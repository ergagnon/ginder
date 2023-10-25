package internal

import "fmt"

type app struct {
	Config AppConfig
}

func NewApp(cfg AppConfig) *app {
	return &app{
		Config: cfg,
	}
}

func (a *app) Run() {
	fmt.Printf(a.Config.Message)
}
