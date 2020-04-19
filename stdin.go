package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type testCase struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type inputType struct {
	Code  string     `json:"code"`
	Tests []testCase `json:"tests"`
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
