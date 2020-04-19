package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	var result resultType
	if len(os.Args) != 6 {
		log.Fatal("Please specify a code, limit, accuracy, tests and corrects.")
	}
	code := os.Args[1]
	limit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	accuracyExp, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	accuracy := math.Pow10(accuracyExp)
	testDir := os.Args[4]
	tests, err := ioutil.ReadDir(testDir)
	if err != nil {
		log.Fatal(err)
	}
	correctDir := os.Args[5]
	corrects, err := ioutil.ReadDir(correctDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Judging '%s'...\n", code)
	for i, test := range tests {
		correct := corrects[i]
		if !test.IsDir() && !correct.IsDir() {
			res, execTime := testCode(code, limit, accuracy, filepath.Join(testDir, test.Name()), filepath.Join(correctDir, correct.Name()))
			result.update(res)
			fmt.Println(res, execTime)
		}
	}
	fmt.Println(result)
	os.Exit(int(result))
}
