package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime/debug"

	"github.com/rafaelsq/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Log write an error in the log
func Log(err error) {
	errs := errors.List(err)
	fmt.Printf("\x0b[38;5;1mERROR\x1b[0m(%d):\n", len(errs))
	for _, err := range errs {
		fmt.Printf(" %v\n", err)
		if er, is := err.(*errors.Error); is {
			if er.Args != nil {
				fmt.Printf("  \x1b[38;5;1margs\x1b[0m: ")
				_ = json.NewEncoder(os.Stderr).Encode(er.Args)
			}
			fmt.Printf("  \x1b[38;5;1mfile\x1b[0m: %s\n", er.Caller)
		}
	}
}

// Zerolog log all the related errors of an error
func Zerolog(err error) {
	errs := errors.List(err)
	lg := log.Error().Timestamp()
	for i, err := range errs {
		if er, is := err.(*errors.Error); is {
			d := zerolog.Dict().
				Str("error", er.Error()).
				Str("caller", er.Caller)
			if er.Args != nil {
				d.Interface("args", er.Args)
			}
			lg.Dict(fmt.Sprintf("err%d", i), d)
		}
	}

	lg.Msg(err.Error())
}

var (
	// ProjectFolder is the main folder of the project
	ProjectFolder = []byte("boiler")
	// IgnoreList are all the files ignored from the log
	IgnoreList = [][]byte{
		[]byte("/middleware.go"),
		[]byte("/errors/errors.go"),
	}
	bl    = []byte("\n")
	end   = []byte("\x1b[0m")
	red   = []byte("\x1b[38;5;1m")
	dark  = []byte("\x1b[38;5;237m")
	gray  = []byte("\x1b[38;5;240m")
	light = []byte("\x1b[38;5;244m")
	white = []byte("\x1b[38;5;250m")
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
