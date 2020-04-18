package main

type resultType int

const (
	resultAccepted resultType = iota
	resultWrongAnswer
	resultTimeLimitExceeded
	resultReferenceError
	resultCompileError
	resultApplicationError
)

func (result resultType) String() string {
	switch result {
	case resultAccepted:
		return "AC"
	case resultApplicationError:
		return "AE"
	case resultCompileError:
		return "CE"
	case resultReferenceError:
		return "RE"
	case resultTimeLimitExceeded:
		return "TLE"
	case resultWrongAnswer:
		return "WA"
	}
	return "AE"
}

func (result *resultType) update(new resultType) {
	if *result < new {
		*result = new
	}
}
