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
	ProjectFolder = []byte("boiler")
	IgnoreList    = [][]byte{
		[]byte("/middleware.go"),
		[]byte("/errors/errors.go"),
	}
	bl    = []byte("\n")
	end   = []byte("\x1b[0m")
	red   = []byte("\x1b[38;5;1m")
	dark  = []byte("\x1b[38;5;236m")
	gray  = []byte("\x1b[38;5;239m")
	light = []byte("\x1b[38;5;243m")
	white = []byte("\x1b[38;5;249m")
	rXP   = regexp.MustCompile(fmt.Sprintf(`(.+/)(%s/.*/)(([^/]+\..{2,3}):\d+)(.*)`, ProjectFolder))
)

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
			_, _ = w.Write(light)
			_, _ = w.Write(ms[1])
			_, _ = w.Write(white)
			_, _ = w.Write(ms[2])
			_, _ = w.Write(red)
			_, _ = w.Write(ms[4])
			_, _ = w.Write(gray)
			_, _ = w.Write(ms[6])
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
