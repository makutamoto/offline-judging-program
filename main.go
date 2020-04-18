package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	var result resultType
	if len(os.Args) != 5 {
		log.Fatal("Please specify a code, limit, tests and corrects.")
	}
	code := os.Args[1]
	limit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	testDir := os.Args[3]
	tests, err := ioutil.ReadDir(testDir)
	if err != nil {
		log.Fatal(err)
	}
	correctDir := os.Args[4]
	corrects, err := ioutil.ReadDir(correctDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Judging '%s'...\n", code)
	for i, test := range tests {
		correct := corrects[i]
		if !test.IsDir() && !correct.IsDir() {
			res, execTime := testCode(code, limit, filepath.Join(testDir, test.Name()), filepath.Join(correctDir, correct.Name()))
			result.update(res)
			fmt.Println(res, execTime)
		}
	}
	fmt.Println(result)
	os.Exit(int(result))
}
