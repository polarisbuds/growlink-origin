package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	dp   *gorm.DB
	bctx context.Context
}

func NewBaseRepository[T any](ctx context.Context, dp *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		dp:   dp,
		bctx: ctx,
	}
}

func (repo *BaseRepository[T]) Exists(query string, args ...any) (bool, error) {
	count, err := repo.Count(query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *BaseRepository[T]) Insert(model *T) error {
	return gorm.G[T](repo.dp).Create(repo.bctx, model)
}

func (repo *BaseRepository[T]) InsertInBatch(models []T, batchSize int) error {
	return gorm.G[T](repo.dp).CreateInBatches(repo.bctx, &models, batchSize)
}

func (repo *BaseRepository[T]) First(order, query string, args ...any) (*T, error) {
	whereChain, err := repo.WhereChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}
	result, err := whereChain.First(repo.bctx)
	return &result, err
}

func (repo *BaseRepository[T]) Where(order, query string, args ...any) ([]T, error) {
	whereChain, err := repo.WhereChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}
	return whereChain.Find(repo.bctx)
}

func (repo *BaseRepository[T]) WhereInBatch(batchSize int, order, query string, args ...any) ([]T, error) {
	whereChain, err := repo.WhereChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}

	result := []T{}
	if err := whereChain.FindInBatches(repo.bctx, batchSize, func(data []T, batch int) error {
		result = append(result, data...)
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *BaseRepository[T]) WhereChain(order, query string, args ...any) (gorm.ChainInterface[T], error) {
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
	return gorm.G[T](repo.dp).Order(order).Find(repo.bctx)
}

func (repo *BaseRepository[T]) Count(query string, args ...any) (int64, error) {
	count, err := gorm.G[T](repo.dp).Where(query, args...).Count(repo.bctx, "*")
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
	return gorm.G[T](repo.dp).Where(query, args...).Limit(limit).Delete(repo.bctx)
}

func (repo *BaseRepository[T]) Update(name string, value any, query string, args ...any) (int, error) {
	exists, err := repo.Exists(query, args...)
	if err != nil {
		return 0, fmt.Errorf("レコードの参照で問題が発生しました: %w", err)
	}
	if !exists {
		return 0, nil
	}
	return gorm.G[T](repo.dp).Where(query, args...).Update(repo.bctx, name, value)
}

func (repo *BaseRepository[T]) Updates(value T, query string, args ...any) (int, error) {
	exists, err := repo.Exists(query, args...)
	if err != nil {
		return 0, fmt.Errorf("レコードの参照で問題が発生しました: %w", err)
	}
	if !exists {
		return 0, nil
	}
	return gorm.G[T](repo.dp).Where(query, args...).Updates(repo.bctx, value)
}
