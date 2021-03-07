package errors

import (
	"bytes"
	"regexp"
	"runtime"
	"runtime/debug"
)

var (
	ProjectPrefix       = regexp.MustCompile(`(?i)^(Boiler/|main.main)`)
	ProjectPathPrefix   = regexp.MustCompile(`(?i)(Boiler/)`)
	IgnoreCallerPrefixs = regexp.MustCompile(
		`(?i)(generated|/middleware\.go|_gen.go|pkg/errors/|graphql/handle\.go|router\.go|rest/resp.go)`,
	)
)

func CallerByLevel(lvls ...int) string {
	var lvl int
	if len(lvls) > 0 {
		lvl = lvls[0]
	}

	lvl += 2

	_, stack := GetStack()
	if len(stack) > lvl {
		return stack[lvl]
	}

	return ""
}

func Caller() string {
	_, stack := GetStack()
	for _, s := range stack {
		if !IgnoreCallerPrefixs.Match([]byte(s)) {
			return s
		}
	}

	return ""
}

func GetStack() (string, []string) {

	var buf [2 << 13]byte
	runtime.Stack(buf[:], false)

	lines := bytes.Split(debug.Stack(), []byte("\n"))
	stack := make([]string, 0, len(lines)/2)
	for i, line := range lines {
		if ProjectPrefix.Match(line) {
			args := ProjectPathPrefix.Split(
				string(
					bytes.Split(
						bytes.TrimSpace(lines[i+1]),
						[]byte(" "),
					)[0],
				),
				2,
			)

			if len(args) > 1 {
				stack = append(stack, args[1])
			} else {
				stack = append(stack, args[0])
			}
		}
	}

	return string(lines[0]), stack
}
