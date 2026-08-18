package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BaseUnScopedRepository[T any] struct {
	udp  *gorm.DB
	uctx context.Context
}

func (repo *BaseUnScopedRepository[T]) UExists(query string, args ...any) (bool, error) {
	count, err := repo.UCount(query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *BaseUnScopedRepository[T]) UFirst(order, query string, args ...any) (*T, error) {
	whereChain, err := repo.newUChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}
	result, err := whereChain.First(repo.uctx)
	return &result, err
}

func (repo *BaseUnScopedRepository[T]) UWhere(order, query string, args ...any) ([]T, error) {
	whereChain, err := repo.newUChain(order, query, args...)
	if err != nil {
		return nil, err
	}
	if whereChain == nil {
		return nil, nil
	}
	return whereChain.Find(repo.uctx)
}

func (repo *BaseUnScopedRepository[T]) newUChain(order, query string, args ...any) (gorm.ChainInterface[T], error) {
	exists, err := repo.UExists(query, args...)
	if err != nil {
		return nil, fmt.Errorf("レコードの参照で問題が発生しました: %w", err)
	}
	if !exists {
		return nil, nil
	}
	return gorm.G[T](repo.udp).Scopes(func(stmt *gorm.Statement) {
		stmt.Unscoped = true
	}).Where(query, args...).Order(order), nil
}

func (repo *BaseUnScopedRepository[T]) UAll(order string) ([]T, error) {
	return gorm.G[T](repo.udp).Scopes(func(stmt *gorm.Statement) {
		stmt.Unscoped = true
	}).Order(order).Find(repo.uctx)
}

func (repo *BaseUnScopedRepository[T]) UCount(query string, args ...any) (int64, error) {
	count, err := gorm.G[T](repo.udp).Scopes(func(stmt *gorm.Statement) {
		stmt.Unscoped = true
	}).Where(query, args...).Count(repo.uctx, "*")
	if err != nil {
		return 0, fmt.Errorf("レコードのカウントで問題が発生しました: %w", err)
	}
	return count, nil
}
