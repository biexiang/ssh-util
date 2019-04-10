package entry

import (
	"fmt"
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
	fmt.Println("==========  ssh  ==========")
	app.PrintServer()
	app.PrintMenu()
	fmt.Println("Please Input server number to connect or action name")

	ret1,ret2,isValid := app.CheckInput()
	if !isValid {
		ret1,ret2,isValid = app.CheckInput()
	}

	if ret1 != "" {
		fmt.Println("Start Action " + ret1)
	}else if ret2 != -1 {
		fmt.Println("Choose Number " + strconv.Itoa(ret2),app.serverMap[ret2])
		app.serverMap[ret2-1].Connect()
	}
}

func (app *App) PrintServer() {
	fmt.Println("========== server ==========")
	for index,server := range app.serverMap {
		s := strconv.Itoa(index + 1) + SEP + server.Name + SEP + server.User + "@" + server.Host + ":" + server.Port
		fmt.Println(s)
	}
}

func (app *App) PrintMenu() {
	//fmt.Println("========== Action ==========")
	//fmt.Println(strings.Join(action,"\n"))
}

func (app *App) CheckInput() (string,int,bool) {
	var input = ""
	for{
		fmt.Scanln(&input)
		isActionExists := func(inputAction string) (string,bool) {
			for _,sepAction := range action {
				if inputAction == sepAction {
					return sepAction,true
				}
			}
			return "",false
		}
		if ret,isExists := isActionExists(input); isExists {
			return ret,-1,isExists
		}

		i,_ := strconv.Atoi(input)
		if _,isExists := app.serverMap[i]; !isExists {
			return "",-1,false
		}else {
			return "",i,true
		}
	}
}