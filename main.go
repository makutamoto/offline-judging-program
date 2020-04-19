package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	var result resultType
	if len(os.Args) != 3 {
		log.Fatal("Please specify a limit and accuracy.")
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
	info := parseStdin()
	compiled, output := compileString(info.Code)
	fmt.Println(output)
	if compiled {
		for _, test := range info.Tests {
			res, execTime := testCode("./temp/a.out", limit, accuracy, test.In, test.Out)
			result.update(res)
			fmt.Println(res, execTime)
		}
	} else {
		result.update(resultCompileError)
	}
	fmt.Println(result)
	os.Exit(int(result))
}
