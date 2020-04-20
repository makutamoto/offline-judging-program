package main

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
)

func compile(language, file string) (bool, string) {
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "./languages/"+language+"/compile.sh", file, tempPrefix+"a.out")
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: childUID, Gid: childGID}
	cmd.Run()
	return cmd.ProcessState.ExitCode() == 0, string(stderr.Bytes())
}

func compileString(language, code string) (bool, string) {
	temp := tempPrefix + "code"
	file, err := os.Create(temp)
	file.Chown(childUID, childGID)
	if err != nil {
		return false, "Compile System Error"
	}
	defer file.Close()
	if _, err := file.WriteString(code); err != nil {
		return false, "Compile System Error"
	}
	return compile(language, temp)
}
