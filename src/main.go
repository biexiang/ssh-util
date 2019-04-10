package main

import (
	"github.com/biexiang/ssh-util/src/entry"
)

func main()  {
	if entry.Host != "" {
		defaultConfig := entry.C.Default
		directServer := entry.Server{
			Host:	entry.Host,
			User:	defaultConfig.User,
			Pass:	defaultConfig.Pass,
			Port:	defaultConfig.Port,
			Key:	defaultConfig.Key,
			Method: defaultConfig.Method,
		}
		directServer.Connect()
	}else {
		app := entry.GetApp()
		app.PrintAndServe()
	}
}