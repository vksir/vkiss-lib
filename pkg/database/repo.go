package database

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

type Constraint[T any] interface {
	Update(tx *gorm.DB, a any) error
	*T
}

type Repo[T any, PT Constraint[T]] struct {
	db       *gorm.DB
	pt       PT
	preloads []string
}

func NewRepo[T any, PT Constraint[T]](db *gorm.DB) *Repo[T, PT] {
	r := &Repo[T, PT]{db: db, pt: PT(new(T))}
	r.preloads = GenNestedPreloads(r.pt)
	return r
}

func (r *Repo[T, PT]) DB() *gorm.DB {
	return r.db
}

func (r *Repo[T, PT]) Transaction(fc func(tx *Repo[T, PT]) error, opts ...*sql.TxOptions) error {
	return r.db.Transaction(func(txDB *gorm.DB) error {
		txRepo := NewRepo[T, PT](txDB)
		return fc(txRepo)
	}, opts...)
}

// Get 级联查询。查询单个数据
func (r *Repo[T, PT]) Get(ctx context.Context, id string) (T, error) {
	var t T
	err := r.Preload(r.db.WithContext(ctx)).
		Where("id = ?", id).Take(&t).Error
	return t, err
}

// GetMany 级联查询。查询多个数据
func (r *Repo[T, PT]) GetMany(ctx context.Context, ids []string) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	var ts []T
	err := r.Preload(r.db.WithContext(ctx)).
		Where("id IN ?", ids).Find(&ts).Error
	return ts, err
}

// GetAll 级联查询。查询所有数据
func (r *Repo[T, PT]) GetAll(ctx context.Context) ([]T, error) {
	var ts []T
	err := r.Preload(r.db.WithContext(ctx)).
		Find(&ts).Error
	return ts, err
}

// Create 级联创建
func (r *Repo[T, PT]) Create(ctx context.Context, ts ...*T) error {
	if len(ts) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&ts).Error
}

// Delete 级联删除
func (r *Repo[T, PT]) Delete(ctx context.Context, ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	var t T
	return r.db.WithContext(ctx).Unscoped().
		Where("id IN ?", ids).Delete(&t).Error
}

// Update 级联更新。空值也会更新
func (r *Repo[T, PT]) Update(ctx context.Context, ts ...*T) error {
	return r.db.WithContext(ctx).Unscoped().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Select("*").Transaction(func(tx *gorm.DB) error {
		for _, t := range ts {
			err := r.pt.Update(tx, t)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo[T, PT]) GetOneByKey(ctx context.Context, k string, v any) (T, error) {
	var t T
	err := r.Preload(r.db.WithContext(ctx)).
		Where(fmt.Sprintf("%s = ?", k), v).Take(&t).Error
	return t, err
}

func (r *Repo[T, PT]) Exist(ctx context.Context, id string) bool {
	var c int64
	r.db.WithContext(ctx).
		Where("id = ?", id).Count(&c)
	return c != 0
}

func (r *Repo[T, PT]) ExistByKey(ctx context.Context, k string, v any) bool {
	var c int64
	r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", k), v).Count(&c)
	return c != 0
}

func (r *Repo[T, PT]) Preload(tx *gorm.DB) *gorm.DB {
	for _, preload := range r.preloads {
		tx = tx.Preload(preload)
	}
	return tx
}
