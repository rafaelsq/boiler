package graphql

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func NewMutation(db entity.DB) *Mutation {
	return &Mutation{db}
}

type Mutation struct {
	db entity.DB
}

func (r *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	panic("not implemented")
}
