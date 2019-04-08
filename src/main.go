package main

import "github.com/biexiang/ssh-util/src/entry"

func main()  {
	app := entry.GetApp()
	app.PrintAndServe()
}