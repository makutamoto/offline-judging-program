package main

type resultType int

const (
	resultAccepted resultType = iota
	resultWrongAnswer
	resultTimeLimitExceeded
	resultReferenceError
	resultCompileError
	resultSystemError
)

func (result resultType) String() string {
	switch result {
	case resultAccepted:
		return "AC"
	case resultSystemError:
		return "SE"
	case resultCompileError:
		return "CE"
	case resultReferenceError:
		return "RE"
	case resultTimeLimitExceeded:
		return "TLE"
	case resultWrongAnswer:
		return "WA"
	}
	return "SE"
}

func (result *resultType) update(new resultType) {
	if *result < new {
		*result = new
	}
}
