package main

import (
	"fmt"
	"math"
)

const tempPrefix = "/var/tmp/makutamoto-offline-judging-program-"

func main() {
	var result resultType
	info := parseStdin()
	accuracy := math.Pow10(info.Problem.Accuracy)
	compiled, compilerOutput := compileString(info.Language, info.Code)
	fmt.Println(len(info.Problem.Tests))
	if compiled {
		for _, test := range info.Problem.Tests {
			res, execTime := testCode(info.Language, tempPrefix+"a.out", info.Problem.Limit, accuracy, test.In, test.Out)
			result.update(res)
			fmt.Println(res, execTime.Milliseconds())
		}
	} else {
		result.update(resultCompileError)
	}
	fmt.Println(result)
	fmt.Println(compilerOutput)
}
