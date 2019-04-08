package entry

import (
	"flag"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Version string `yaml:"version"`
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name 	 string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Method	 string `yaml:"method"`
	Pass     string `yaml:"pass"`
	Key		 string `yaml:"key"`
}

var C = Config{}

func init() {
	var configPath string
	var help string
	flag.StringVar(&configPath,"config","./config.yaml","Config Path")
	flag.StringVar(&help,"h","","Show Help")

	ret, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
	}
	yaml.Unmarshal(ret, &C)
}

func (server Server) Connect() {
	server.parseAuth()
}

func (server Server) parseAuth() {
	sshList := []ssh.AuthMethod{}
	log.Println(sshList)
}
