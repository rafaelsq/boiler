package resolver

import (
	"context"
	"fmt"
	"testing"

	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	assert.Equal(t, iface.ErrNotFound, Wrap(context.TODO(), iface.ErrNotFound))

	assert.Equal(t, "service failed", Wrap(context.TODO(), fmt.Errorf("opz"), "fail").Error())

	assert.Equal(t, "opz", Wrap(context.WithValue(
		context.TODO(), iface.ContextKeyDebug{}, true), fmt.Errorf("opz")).Error())
}
