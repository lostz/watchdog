package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/lostz/watchdog/server"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return
	}
	s.Run()

}
