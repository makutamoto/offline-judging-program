package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

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
