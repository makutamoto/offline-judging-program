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
	"time"
)

// #include<unistd.h>
// long getClock(void) {
//	 return sysconf(_SC_CLK_TCK);
// }
import "C"

var clockTck = int(C.getClock())

func compareValue(correct string, answer string, accuracy float64) bool {
	integerC, err := strconv.Atoi(correct)
	if err == nil {
		integerA, err := strconv.Atoi(answer)
		return err == nil && integerC == integerA
	}
	numberC, err := strconv.ParseFloat(correct, 64)
	if err == nil {
		numberA, err := strconv.ParseFloat(correct, 64)
		return err == nil && math.Abs((numberC-numberA)/numberC) <= accuracy
	}
	return correct == answer
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

func testCode(code string, limit int, accuracy float64, testIn string, testOut string) (resultType, time.Duration) {
	var stdout bytes.Buffer
	var result resultType
	var correctScan, answerScan bool
	cmd := exec.Command("bash", "./data/run.sh", code)
	cmd.Stdin = strings.NewReader(testIn)
	cmd.Stdout = &stdout
	if err := cmd.Start(); err != nil {
		result.update(resultSystemError)
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
	execTime := cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime()
	if result == resultTimeLimitExceeded {
		return result, execTime
	}
	scannerOut := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	scannerOut.Split(bufio.ScanWords)
	scannerCorrect := bufio.NewScanner(strings.NewReader(testOut))
	scannerCorrect.Split(bufio.ScanWords)
	for {
		answerScan = scannerOut.Scan()
		correctScan = scannerCorrect.Scan()
		if !answerScan || !correctScan {
			break
		} else if !compareValue(scannerCorrect.Text(), scannerOut.Text(), accuracy) {
			result.update(resultWrongAnswer)
			break
		}
	}
	if answerScan != correctScan {
		result.update(resultWrongAnswer)
	}
	return result, execTime
}
