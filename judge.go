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
	"time"
)

// #include<unistd.h>
// long getClock(void) {
//	 return sysconf(_SC_CLK_TCK);
// }
import "C"

var clockTck = int(C.getClock())

func compareValue(correct string, answer string) bool {
	integerC, err := strconv.Atoi(correct)
	if err == nil {
		integerA, err := strconv.Atoi(answer)
		return err == nil && integerC == integerA
	}
	numberC, err := strconv.ParseFloat(correct, 64)
	if err == nil {
		numberA, err := strconv.ParseFloat(correct, 64)
		return err == nil && math.Abs((numberC-numberA)/numberC) < 10e-3
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

func testCode(code string, limit int, test string, solution string) (resultType, time.Duration) {
	var stdout bytes.Buffer
	var tleFlag = false
	var waFlag = false
	var correctScan, answerScan bool
	testFile, err := os.Open(test)
	if err != nil {
		fmt.Printf("Couldn't load test case '%s': skipped.\n", test)
		return resultAccepted, 0
	}
	defer testFile.Close()
	cmd := exec.Command(code)
	cmd.Stdin = bufio.NewReader(testFile)
	cmd.Stdout = &stdout
	cmd.Start()
	for {
		execTime := getTime(cmd.Process)
		if execTime > limit {
			cmd.Process.Kill()
			tleFlag = true
			break
		} else if execTime == -1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	cmd.Wait()
	execTime := cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime()
	if tleFlag {
		return resultTimeLimitExceeded, execTime
	}
	scannerOut := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	scannerOut.Split(bufio.ScanWords)
	solutionFile, err := os.Open(solution)
	if err != nil {
		fmt.Printf("Couldn't load solution '%s': skipped.\n", solution)
		return resultAccepted, 0
	}
	scannerCorrect := bufio.NewScanner(solutionFile)
	for {
		answerScan = scannerOut.Scan()
		correctScan = scannerCorrect.Scan()
		if !answerScan || !correctScan {
			break
		}
		if !compareValue(scannerCorrect.Text(), scannerOut.Text()) {
			waFlag = true
			break
		}
	}
	if answerScan == correctScan {
		if waFlag {
			return resultWrongAnswer, execTime
		}
		return resultAccepted, execTime
	}
	return resultWrongAnswer, execTime
}
