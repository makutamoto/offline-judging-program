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

func testCode(code string, limit int, test string, correct string) (resultType, time.Duration) {
	var stdout bytes.Buffer
	var result resultType
	var correctScan, answerScan bool
	testFile, err := os.Open(test)
	if err != nil {
		fmt.Printf("Couldn't load test case '%s': skipped.\n", test)
		return result, 0
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
	correctFile, err := os.Open(correct)
	if err != nil {
		fmt.Printf("Couldn't load correct '%s': skipped.\n", correct)
		return result, 0
	}
	scannerCorrect := bufio.NewScanner(correctFile)
	for {
		answerScan = scannerOut.Scan()
		correctScan = scannerCorrect.Scan()
		if !answerScan || !correctScan {
			break
		} else if !compareValue(scannerCorrect.Text(), scannerOut.Text()) {
			result.update(resultWrongAnswer)
			break
		}
	}
	if answerScan != correctScan {
		result.update(resultWrongAnswer)
	}
	return result, execTime
}
