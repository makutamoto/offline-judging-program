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
	if len(os.Args) != 5 {
		log.Fatal("Please specify a limit, accuracy, tests and corrects.")
	}
	limit, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	accuracyExp, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	accuracy := math.Pow10(accuracyExp)
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
	compiled, output := compileStdin()
	fmt.Println(output)
	if compiled {
		for i, test := range tests {
			correct := corrects[i]
			if !test.IsDir() && !correct.IsDir() {
				res, execTime := testCode("./temp/a.out", limit, accuracy, filepath.Join(testDir, test.Name()), filepath.Join(correctDir, correct.Name()))
				result.update(res)
				fmt.Println(res, execTime)
			}
		}
	} else {
		result.update(resultCompileError)
	}
	fmt.Println(result)
	os.Exit(int(result))
}
