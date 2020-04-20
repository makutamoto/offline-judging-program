package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// #include<unistd.h>
// long getClock(void) {
//	 return sysconf(_SC_CLK_TCK);
// }
import "C"

var clockTck = int(C.getClock())

func compareValue(test string, answer string, accuracy float64) bool {
	integerC, err := strconv.Atoi(test)
	if err == nil {
		integerA, err := strconv.Atoi(answer)
		return err == nil && integerC == integerA
	}
	numberC, err := strconv.ParseFloat(test, 64)
	if err == nil {
		numberA, err := strconv.ParseFloat(test, 64)
		return err == nil && math.Abs((numberC-numberA)/numberC) <= accuracy
	}
	return test == answer
}

func getTime(process *os.Process) int {
	var state byte
	path := fmt.Sprintf("/proc/%d/stat", process.Pid)
	fp, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	for i := 0; i < 13; i++ {
		scanner.Scan()
		if i == 2 {
			state = scanner.Text()[0]
		}
	}
	scanner.Scan()
	user, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}
	scanner.Scan()
	sys, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}
	if state == 'Z' {
		return -1
	}
	return 1000 * (user + sys) / clockTck
}

func testCode(language, code string, limit int, accuracy float64, testIn, testOut string) (resultType, time.Duration) {
	var stdout bytes.Buffer
	var result resultType
	var testScan, answerScan bool
	cmd := exec.Command("bash", "./languages/"+language+"/run.sh", code)
	cmd.Stdin = strings.NewReader(testIn)
	cmd.Stdout = &stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: childUID, Gid: childGID}
	if err := cmd.Start(); err != nil {
		result.update(resultInternalError)
		return result, 0
	}
	for {
		execTime := getTime(cmd.Process)
		if execTime > limit {
			cmd.Process.Kill()
			result.update(resultTimeLimitExceeded)
			break
		} else if execTime == -1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		result.update(resultReferenceError)
	}
	execTime := cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime()
	if result == resultTimeLimitExceeded {
		return result, execTime
	}
	scannerOut := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	scannerOut.Split(bufio.ScanWords)
	scannerTest := bufio.NewScanner(strings.NewReader(testOut))
	scannerTest.Split(bufio.ScanWords)
	for {
		answerScan = scannerOut.Scan()
		testScan = scannerTest.Scan()
		if !answerScan || !testScan {
			break
		} else if !compareValue(scannerTest.Text(), scannerOut.Text(), accuracy) {
			result.update(resultWrongAnswer)
			break
		}
	}
	if answerScan != testScan {
		result.update(resultWrongAnswer)
	}
	return result, execTime
}
