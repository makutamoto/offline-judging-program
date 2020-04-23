package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type statusType struct {
	WholeResult int    `json:"whole_result"`
	Result      int    `json:"result"`
	Time        int64  `json:"time"`
	Memory      int64  `json:"memory"`
	CurrentCase int    `json:"current_case"`
	WholeCase   int    `json:"whole_case"`
	Description string `json:"description"`
}

type testType struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type problemType struct {
	Limit    int        `json:"limit"`
	Accuracy int        `json:"accuracy"`
	Tests    []testType `json:"tests"`
}

type inputType struct {
	Language string      `json:"language"`
	Code     string      `json:"code"`
	Problem  problemType `json:"problem"`
}

func sendStatus(wholeResult resultType, result resultType, time int64, memory int64, currentCase int, wholeCase int, description string) {
	status := statusType{WholeResult: int(wholeResult), Result: int(result), Time: time, Memory: memory, CurrentCase: currentCase, WholeCase: wholeCase, Description: description}
	bytes, err := json.Marshal(status)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(bytes))
}

func parseStdin() inputType {
	var input inputType
	scanner := bufio.NewScanner(os.Stdin)
	bytes := make([]byte, 0)
	for scanner.Scan() && scanner.Text() != "" {
		bytes = append(bytes, scanner.Bytes()...)
	}
	if err := json.Unmarshal(bytes, &input); err != nil {
		log.Fatal(err)
	}
	return input
}
