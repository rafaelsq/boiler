package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

type Error struct {
	msg    string
	args   map[string]interface{}
	caller string
}

func (e *Error) Error() string {
	return e.msg
}

func New(msg string) error {
	return &Error{msg: msg, caller: getFrame(1)}
}

func WithArg(msg, key string, value interface{}) error {
	return &Error{msg: msg, args: map[string]interface{}{key: value}, caller: getFrame(1)}
}

func WithArgs(msg string, args map[string]interface{}) error {
	return &Error{msg: msg, args: args, caller: getFrame(1)}
}

func Cause(err error) error {
	errs := multierr.Errors(err)
	if len(errs) != 0 {
		return errs[0]
	}
	return nil
}

func getFrame(skipFrames int) string {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	path := fmt.Sprintf("%s:%d", frame.File, frame.Line)
	ms := rXP.FindStringSubmatch(path)
	if len(ms) != 0 {
		return fmt.Sprintf("%s%s", ms[2], ms[4])
	}

	return path
}

func Log(err error) {
	errs := multierr.Errors(err)
	fmt.Printf("\x1b[38;5;1mERROR\x1b[0m(%d):\n", len(errs))
	for _, err := range errs {
		fmt.Printf(" %v\n", err)
		if er, is := err.(*Error); is {
			if er.args != nil {
				fmt.Printf("  \x1b[38;5;1margs\x1b[0m: ")
				_ = json.NewEncoder(os.Stderr).Encode(er.args)
				fmt.Printf("\n")
			}
			fmt.Printf("  \x1b[38;5;1mfile\x1b[0m: %s\n", er.caller)
		}
	}
}
func Zerolog(err error) {
	errs := multierr.Errors(err)
	lg := log.Error().Timestamp()
	for i, err := range errs {
		if er, is := err.(*Error); is {
			d := zerolog.Dict().
				Str("error", er.Error()).
				Str("caller", er.caller)
			if er.args != nil {
				d.Interface("args", er.args)
			}
			lg.Dict(fmt.Sprintf("err%d", i), d)
		}
	}

	lg.Msg(err.Error())
}

var (
	ProjectFolder = []byte("boiler")
	IgnoreList    = [][]byte{
		[]byte("/vendor/"),
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
	rXP   = regexp.MustCompile(fmt.Sprintf(`(.+/)(%s/(vendor/)?.*/)(([^/]+\..{2,3}):\d+)(.*)`, ProjectFolder))
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
