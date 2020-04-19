package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	var result resultType
	info := parseStdin()
	accuracy := math.Pow10(info.Problem.Accuracy)
	compiled, _ := compileString(info.Code)
	if compiled {
		for _, test := range info.Problem.Tests {
			res, execTime := testCode("./temp/a.out", info.Problem.Limit, accuracy, test.In, test.Out)
			result.update(res)
			fmt.Println(res, execTime)
		}
	} else {
		result.update(resultCompileError)
	}
	fmt.Println(result)
	os.Exit(int(result))
}
