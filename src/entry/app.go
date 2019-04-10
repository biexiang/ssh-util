package entry

import (
	"fmt"
	"os"
	"strconv"
)

const SEP  = "\t => "
var action = []string{"add","edit","remove","exit"}

type App struct {
	serverMap map[int]Server
}

func GetApp() *App {

	var serverMap = make(map[int]Server)
	for index,server := range C.Servers {

		if server.Host == "" {
			continue
		}

		if server.User == "" {
			server.User = C.Default.User
		}

		if server.Pass == "" {
			server.Pass = C.Default.Pass
		}

		if server.Port == "" {
			server.Port = C.Default.Port
		}

		if server.Method == "" {
			server.Method = C.Default.Method
		}

		if server.Key == "" {
			server.Key = C.Default.Key
		}

		serverMap[index] = server
	}

	return &App{
		serverMap: serverMap,
	}
}

func (app *App) PrintAndServe() {
	Clear()
	app.PrintServer()
	fmt.Println("Please Input server number to connect or action name")

	ret,isValid := app.CheckInput()
	if !isValid {
		ret,isValid = app.CheckInput()
	}

	if !isValid {
		fmt.Println("Input Not Valid")
		os.Exit(1)
	}
	app.serverMap[ret-1].Connect()
}

func (app *App) PrintServer() {
	fmt.Println("========== server ==========")
	for index,server := range app.serverMap {
		s := strconv.Itoa(index + 1) + SEP + server.Name + SEP + server.User + "@" + server.Host + ":" + server.Port
		fmt.Println(s)
	}
}

func (app *App) CheckInput() (int,bool) {
	var input int
	for{
		fmt.Scanln(&input)
		if _,isExists := app.serverMap[input - 1]; !isExists {
			return -1,false
		}else {
			return input,true
		}
	}
}