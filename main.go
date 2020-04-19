package main

import (
	"fmt"
	"math"
	"os"
)

const tempPrefix = "/var/tmp/makutamoto-offline-judging-program-"

func main() {
	var result resultType
	info := parseStdin()
	accuracy := math.Pow10(info.Problem.Accuracy)
	compiled, _ := compileString(info.Language, info.Code)
	if compiled {
		for _, test := range info.Problem.Tests {
			res, execTime := testCode(info.Language, tempPrefix+"a.out", info.Problem.Limit, accuracy, test.In, test.Out)
			result.update(res)
			fmt.Println(res, execTime)
		}
	} else {
		result.update(resultCompileError)
	}
	fmt.Println(result)
	os.Exit(int(result))
}
