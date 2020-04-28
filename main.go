package main

import (
	"fmt"
	"math"
)

const tempPrefix = "/var/tmp/makutamoto-offline-judging-program-"
const childUID, childGID = 400, 400

func main() {
	var result resultType
	info := parseStdin()
	accuracy := math.Pow10(info.Problem.Accuracy)
	compiled, compilerOutput := compileString(info.Language, info.Code)
	if compiled {
		for i, test := range info.Problem.Tests {
			sendStatus(false, test.Title, result, 0, 0, 0, i+1, len(info.Problem.Tests), "")
			res, execTime, memory := testCode(info.Language, tempPrefix+"a.out", info.Problem.Limit, accuracy, test.In, test.Out)
			result.update(res)
			sendStatus(true, test.Title, result, res, execTime.Milliseconds(), memory, i+1, len(info.Problem.Tests), "")
		}
	} else {
		result.update(resultCompileError)
		sendStatus(true, "", result, result, 0, 0, 0, 0, compilerOutput)
	}
	fmt.Println()
}
