package main

import (
	"encoding/json"
	"io/ioutil"
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
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &input); err != nil {
		log.Fatal(err)
	}
	return input
}
