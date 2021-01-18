package resolver

import (
	"context"
	"fmt"
	"testing"

	"boiler/pkg/errors"
	"boiler/pkg/store/config"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	assert.Equal(t, errors.ErrNotFound, Wrap(context.TODO(), errors.ErrNotFound))

	assert.Equal(t, "service failed", Wrap(context.TODO(), fmt.Errorf("opz"), "fail").Error())

	assert.Equal(t, "opz", Wrap(context.WithValue(
		context.TODO(), config.ContextKeyDebug{}, true), fmt.Errorf("opz")).Error())
}
