package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	dp  *gorm.DB
	ctx context.Context
}

type RepositoryOption struct {
	Order         string
	WhereOption   *WhereOption
	PreloadOption []*PreloadOption
}

type WhereOption struct {
	Query string
	Args  []any
}

type PreloadOption struct {
	Association string
	Query       func(db gorm.PreloadBuilder) error
}

func (repo *BaseRepository[T]) Exists(query string, args ...any) (bool, error) {
	count, err := repo.Count(query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *BaseRepository[T]) Insert(model *T) error {
	return gorm.G[T](repo.dp).Create(repo.ctx, model)
}

func (repo *BaseRepository[T]) InsertInBatch(models []T, batchSize int) error {
	return gorm.G[T](repo.dp).CreateInBatches(repo.ctx, &models, batchSize)
}

func (repo *BaseRepository[T]) First(order, query string, args ...any) (*T, error) {
	whereChain, err := repo.newChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}
	result, err := whereChain.First(repo.ctx)
	return &result, err
}

func (repo *BaseRepository[T]) Where(order, query string, args ...any) ([]T, error) {
	whereChain, err := repo.newChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}
	return whereChain.Find(repo.ctx)
}

func (repo *BaseRepository[T]) newChain(order, query string, args ...any) (gorm.ChainInterface[T], error) {
	exists, err := repo.Exists(query, args...)
	if err != nil {
		return nil, fmt.Errorf("レコードの参照で問題が発生しました: %w", err)
	}
	if !exists {
		return nil, nil
	}
	return gorm.G[T](repo.dp).Where(query, args...).Order(order), nil
}

func (repo *BaseRepository[T]) All(order string) ([]T, error) {
	return gorm.G[T](repo.dp).Order(order).Find(repo.ctx)
}

func (repo *BaseRepository[T]) Count(query string, args ...any) (int64, error) {
	count, err := gorm.G[T](repo.dp).Where(query, args...).Count(repo.ctx, "*")
	if err != nil {
		return 0, fmt.Errorf("レコードのカウントで問題が発生しました: %w", err)
	}
	return count, nil
}

func (repo *BaseRepository[T]) Delete(limit int, query string, args ...any) (int, error) {
	exists, err := repo.Exists(query, args...)
	if err != nil {
		return 0, fmt.Errorf("レコードの参照で問題が発生しました: %w", err)
	}
	if !exists {
		return 0, nil
	}
	return gorm.G[T](repo.dp).Where(query, args...).Limit(limit).Delete(repo.ctx)
}

func (repo *BaseRepository[T]) Update(name string, value any, query string, args ...any) (int, error) {
	exists, err := repo.Exists(query, args...)
	if err != nil {
		return 0, fmt.Errorf("レコードの参照で問題が発生しました: %w", err)
	}
	if !exists {
		return 0, nil
	}
	return gorm.G[T](repo.dp).Where(query, args...).Update(repo.ctx, name, value)
}
