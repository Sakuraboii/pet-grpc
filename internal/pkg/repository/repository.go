package repository

import (
	"context"
	"errors"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type UsersRepo interface {
	Add(ctx context.Context, user *User) (int64, error)
	GetById(ctx context.Context, id int64) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) (bool, error)
	Delete(ctx context.Context, id int64) error
}

type CarsRepo interface {
	Add(ctx context.Context, car *Car) (int64, error)
	GetById(ctx context.Context, id int64) (*Car, error)
	List(ctx context.Context) ([]*Car, error)
	Update(ctx context.Context, user *Car) (bool, error)
	Delete(ctx context.Context, id int64) error
}
