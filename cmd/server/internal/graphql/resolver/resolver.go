package resolver

import (
	"context"
	"errors"
	"fmt"

	"boiler/pkg/store"
	"boiler/pkg/store/config"

	"github.com/rs/zerolog/log"
)

// Wrap wrap error
func Wrap(ctx context.Context, err error, args ...string) error {
	if errors.Is(err, store.ErrNotFound) {
		return err
	}

	if debug := ctx.Value(config.ContextKeyDebug{}); debug != nil {
		return err
	}

	msg := err.Error()
	if len(msg) != 0 {
		msg = args[0]
	}

	log.Error().Err(err).Msg(msg)

	return fmt.Errorf("service failed")
}
