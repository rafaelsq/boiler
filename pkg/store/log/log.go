package log

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"runtime/debug"
)

var (
	// ProjectFolder is the main folder of the project
	ProjectFolder = []byte("boiler")
	// IgnoreList are all the files ignored from the log
	IgnoreList = [][]byte{
		[]byte("/middleware.go"),
		[]byte("/errors/errors.go"),
	}
	bl    = []byte("\n")
	end   = []byte("\033[0m")
	red   = []byte("\033[38;5;1m")
	dark  = []byte("\033[38;5;237m")
	gray  = []byte("\033[38;5;240m")
	light = []byte("\033[38;5;244m")
	white = []byte("\033[38;5;250m")
	rXP   = regexp.MustCompile(fmt.Sprintf(`(.+/)(%s/.*/)(([^/]+\..{2,3}):\d+)(.*)`, ProjectFolder))
)

// WriteStack write the stack to the writer
func WriteStack(w io.Writer) {
	for i, l := range bytes.Split(debug.Stack(), []byte("\n")) {
		if i == 0 {
			_, _ = w.Write(bl)
			_, _ = w.Write(red)
			_, _ = w.Write(l)
			_, _ = w.Write(end)
			_, _ = w.Write(bl)
			continue
		}

		ms := rXP.FindSubmatch(l)
		ignore := false
		if len(ms) != 0 {
			for _, ig := range IgnoreList {
				if bytes.Contains(ms[0], ig) {
					ignore = true
					break
				}
			}
		}

		if len(ms) != 0 && !ignore {
			// fmt.Println(string(l))
			// for i := 0; i < len(ms); i++ {
			// 	fmt.Println(" -", string(ms[i]))

			// }

			_, _ = w.Write(light)
			_, _ = w.Write(ms[1])
			_, _ = w.Write(white)
			_, _ = w.Write(ms[2])
			_, _ = w.Write(red)
			_, _ = w.Write(ms[3])
			_, _ = w.Write(gray)
			for i := 5; i < len(ms); i++ {
				_, _ = w.Write(ms[i])
			}
			_, _ = w.Write(end)
			_, _ = w.Write(bl)
			continue
		}

		if bytes.Contains(l, ProjectFolder) {
			_, _ = w.Write(gray)
			_, _ = w.Write(l)
			_, _ = w.Write(end)
			_, _ = w.Write(bl)
			continue
		}

		_, _ = w.Write(dark)
		_, _ = w.Write(l)
		_, _ = w.Write(end)
		_, _ = w.Write(bl)
	}
}
