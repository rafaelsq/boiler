package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
)

func main() {
	c := make(chan notify.EventInfo)

	if err := notify.Watch("./...", c, notify.Create, notify.Write, notify.Remove); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	gqlctx, gqlcancel := context.WithCancel(context.Background())
	gctx, gcancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	tctx, tcancel := context.WithCancel(context.Background())
	pctx, pcancel := context.WithCancel(context.Background())
	go buildNRun(ctx)
	for {
		select {
		case ei := <-c:
			path := ei.Path()
			if strings.HasSuffix(path, "schema.graphql") {
				gqlcancel()
				gqlctx, gqlcancel = context.WithCancel(context.Background())
				go func() {
					gqlctx := gqlctx
					if err := run(gqlctx, "go", "run", "github.com/99designs/gqlgen"); err != nil &&
						gqlctx.Err() != context.Canceled {

						fmt.Println("graph fail;", err)
					}
				}()
			} else if strings.Contains(path, "pkg/iface/") {
				gcancel()
				gctx, gcancel = context.WithCancel(context.Background())
				go func() {
					gctx := gctx
					if err := run(gctx, "go", "generate", "./..."); err != nil &&
						gctx.Err() != context.Canceled {

						fmt.Println("gen fail;", err)
					}
				}()
			} else if strings.HasSuffix(path, "_test.go") {
				tcancel()
				tctx, tcancel = context.WithCancel(context.Background())

				pieces := strings.Split(path, "/")
				pkg := strings.Join(pieces[:len(pieces)-1], "/")

				go run(tctx, "go", "test", "-mod=vendor", "-cover", pkg)
			} else if strings.HasSuffix(path, ".go") {
				cancel()
				ctx, cancel = context.WithCancel(context.Background())
				go buildNRun(ctx)
			} else if strings.HasSuffix(ei.Path(), ".proto") {
				pcancel()
				pctx, pcancel = context.WithCancel(context.Background())
				go func() {
					pctx := pctx
					if err := run(pctx, "make", "proto"); err != nil &&
						pctx.Err() != context.Canceled {

						fmt.Println("proto fail;", err)
					}
				}()
			}
		}
	}
}

func buildNRun(ctx context.Context) {
	err := run(ctx, "go", "build", "-mod=vendor", "cmd/server/server.go")
	if err != nil && err != context.Canceled && ctx.Err() != context.Canceled {
		fmt.Println("build failed;", err)
	} else {
		if err := run(ctx, "./server"); err != nil && err != context.Canceled && ctx.Err() != context.Canceled {
			fmt.Println("run failed;", err)
		}
	}
}

func run(ctx context.Context, command string, args ...string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Millisecond * 200):
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
		case <-done:
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
