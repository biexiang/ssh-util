package entry

import (
	"os"
	"os/exec"
	"runtime"
)

func Clear() {
	var cmd *exec.Cmd
	if "windows" == runtime.GOOS {
		cmd = exec.Command("cmd","/c","cls")
	}else {
		cmd = exec.Command("/bin/bash","-c","clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
