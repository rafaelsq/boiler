package resolver

import (
	"context"
	"fmt"
	"testing"

	"boiler/pkg/store"
	"boiler/pkg/store/config"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	assert.Equal(t, store.ErrNotFound, Wrap(context.TODO(), store.ErrNotFound))

	assert.Equal(t, "service failed", Wrap(context.TODO(), fmt.Errorf("opz"), "fail").Error())

	assert.Equal(t, "opz", Wrap(context.WithValue(
		context.TODO(), config.ContextKeyDebug{}, true), fmt.Errorf("opz")).Error())
}
