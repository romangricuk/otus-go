package app

import "context"

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
