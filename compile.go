package main

import (
	"bytes"
	"os"
	"os/exec"
)

func compile(language, file string) (bool, string) {
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "./languages/"+language+"/compile.sh", file, "./temp/a.out")
	cmd.Stderr = &stderr
	cmd.Run()
	return cmd.ProcessState.ExitCode() == 0, string(stderr.Bytes())
}

func compileString(language, code string) (bool, string) {
	file, err := os.Create("./temp/code")
	if err != nil {
		return false, "Compile System Error"
	}
	defer file.Close()
	if _, err := file.WriteString(code); err != nil {
		return false, "Compile System Error"
	}
	return compile(language, "./temp/code")
}
