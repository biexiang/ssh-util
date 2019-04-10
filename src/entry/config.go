package entry

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Config struct {
	Version string `yaml:"version"`
	Servers []Server `yaml:"servers"`
	Default Default `yaml:"default"`
}

type Server struct {
	Name 	 string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User 	 string `yaml:"user"`
	Method	 string `yaml:"method"`
	Pass     string `yaml:"pass"`
	Key		 string `yaml:"key"`

	termWidth int
	termHeight int
}


type Default struct {
	Port     string `yaml:"port"`
	User 	 string `yaml:"user"`
	Method	 string `yaml:"method"`
	Pass     string `yaml:"pass"`
	Key		 string `yaml:"key"`
}

var C = Config{}
var ConfigPath string
var Host string

func init() {

	flag.StringVar(&ConfigPath,"c","","Config Path")
	flag.StringVar(&Host,"h","","Specific Which Host To Connect")
	flag.Parse()

	if ConfigPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	ret, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Println(err)
	}
	yaml.Unmarshal(ret, &C)
}

func (server Server) Connect() {
	auth,err := server.parseAuth()
	if err != nil {
		fmt.Println("Error In ParseAuth")
		return
	}

	config := &ssh.ClientConfig{
		User:server.User,
		Auth:auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	address := server.Host + ":" + server.Port
	client,err := ssh.Dial("tcp",address,config)
	if err != nil {
		fmt.Println("Error In Build Connection with ", server.Host)
		return
	}
	defer client.Close()
	session,err := client.NewSession()
	if err != nil {
		fmt.Println("Error In Create Session")
		return
	}
	defer session.Close()
	//获取当前终端的文件描述符
	fd := int(os.Stdin.Fd())
	oldState,err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println("Raw Fail")
		return
	}
	//心跳检查
	keepAliveCh := server.keepAlive(session)
	defer close(keepAliveCh)
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin  = os.Stdin
	defer terminal.Restore(fd,oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	server.termWidth, server.termHeight, _ = terminal.GetSize(fd)
	if err := session.RequestPty("xterm-256color", server.termHeight, server.termWidth, modes); err != nil {
		fmt.Println("Create Term Fail")
		return
	}
	winChangeCh := server.listenWinChange(session,fd)
	defer close(winChangeCh)

	err = session.Shell()
	if err != nil {
		fmt.Println("Get Remote Shell Fail")
		return
	}

	err = session.Wait()
	if err != nil {
		fmt.Println("Remote Exit Error")
		return
	}
}

func (server Server) listenWinChange(session *ssh.Session,fd int) chan struct{} {
	term := make(chan struct{})
	go func() {
		for{
			select{
			case <-term:
				return
			default:
				width,height,_ := terminal.GetSize(fd)

				if width != server.termWidth || height != server.termHeight {
					server.termWidth = width
					server.termHeight = height
					session.WindowChange(height,width)
				}
				time.Sleep(time.Second * 1)
			}
		}
	}()

	return term
}

func (server Server) keepAlive(session *ssh.Session) chan struct{} {
	term := make(chan struct{})
	go func() {
		for{
			select {
			case <- term:
				return;
			default:
				_,err := session.SendRequest("What's Up",true,nil)
				if err != nil {
					fmt.Println("KeepAlive Heartbeat Lose")
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()
	return term
}

func (server Server) parseAuth() ([]ssh.AuthMethod,error){
	sshList := []ssh.AuthMethod{}
	method := strings.ToLower(server.Method)
	if method == "password" {
		sshList = append(sshList,ssh.Password(server.Pass))
	}

	if method == "key" {
		method,err := fetchPublicKey(server)
		if err != nil {
			return nil,err
		}
		sshList = append(sshList,method)
	}
	return sshList,nil
}

func fetchPublicKey(server Server) (ssh.AuthMethod, error) {

	var signer ssh.Signer

	if server.Key == "" {
		server.Key = "~/.ssh/id_rsa"
	}

	pemBytes,err := ioutil.ReadFile(server.Key)
	if err != nil {
		return nil,err
	}

	if server.Pass == "" {
		signer,err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer,err = ssh.ParsePrivateKeyWithPassphrase(pemBytes,[]byte(server.Pass))
	}

	if err != nil {
		return nil,err
	}

	return ssh.PublicKeys(signer),nil
}
