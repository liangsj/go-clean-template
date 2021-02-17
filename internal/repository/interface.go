// Package repository implements databases for services. For each service in its own file.
package repository

import (
	"context"

	"github.com/evrone/go-service-template/internal/domain"
)

type Translation interface {
	Store(ctx context.Context, entity domain.Translation) error
	GetHistory(context.Context) ([]domain.Translation, error)
}
