package resolver

import (
	"context"
	"fmt"

	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/log"
	"github.com/rafaelsq/errors"
)

// Wrap wrap error
func Wrap(ctx context.Context, err error, args ...string) error {
	if er := errors.Cause(err); err == iface.ErrNotFound {
		return er
	}

	if debug := ctx.Value(iface.ContextKeyDebug{}); debug != nil {
		return err
	}

	msg := err.Error()
	if len(msg) != 0 {
		msg = args[0]
	}

	lg := errors.New(msg).SetParent(err)
	lg.Caller = errors.Caller(1)
	log.Log(lg)

	return fmt.Errorf("service failed")
}
