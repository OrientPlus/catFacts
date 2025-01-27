package repo

import (
	"caxfaxService1/internal/entity"
	"context"
)

type Repo interface {
	// Добавляет сущность entity.Fact в базу и возвращает его ID
	AddFact(ctx context.Context, fact entity.Fact) (int32, error)

	// Достает из базы сущность entity.Fact по его ID
	GetFact(ctx context.Context, factID int32) (entity.Fact, error)

	// Обновляет сущность entity.Fact в базе
	UpdateFact(ctx context.Context, factID int32, fact entity.Fact) error

	// Удаляет сущность entity.Fact по ее ID
	DeleteFact(ctx context.Context, factID int32) error
}
