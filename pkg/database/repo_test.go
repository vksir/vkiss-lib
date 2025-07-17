package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"github.com/vksir/vkiss-lib/pkg/util/fileutil"
	"gorm.io/gorm"
	"path/filepath"
	"testing"
)

type M struct {
	ID string `gorm:"primarykey" json:"id"`
}

func (m *M) Update(tx *gorm.DB, a any) error {
	res := tx.Updates(a)
	return res.Error
}

type A struct {
	M
	HasOneB B      `json:"b" gorm:"constraint:OnDelete:CASCADE"`
	V       string `json:"v"`
}

func (m *A) Update(tx *gorm.DB, a any) error {
	v := a.(*A)
	cc := *v
	err := tx.Model(&cc).
		Association("HasOneB").
		Unscoped().
		Clear()
	if err != nil {
		return errutil.Wrap(err)
	}

	res := tx.Updates(v)
	if res.Error != nil {
		return errutil.Wrap(res.Error)
	}
	return nil
}

type B struct {
	M
	AID      string `json:"-" gorm:"not null"`
	HasManyC []C    `json:"c" gorm:"constraint:OnDelete:CASCADE"`
	V        string `json:"v"`
}

type C struct {
	M
	BID string `json:"-" gorm:"not null"`
	V   string `json:"v"`
}

type D struct {
	M
	BelongsToAID string `json:"-" gorm:"not null"`
	BelongsToA   A      `json:"a" gorm:"constraint:OnDelete:CASCADE"`
	EID          string `json:"-" gorm:"not null"`
	V            string `json:"v"`
}

type E struct {
	M
	HasManyD []D    `json:"d" gorm:"constraint:OnDelete:CASCADE"`
	V        string `json:"v"`
}

func (m *E) Update(tx *gorm.DB, a any) error {
	v := a.(*E)

	cc := *v
	err := tx.Model(&cc).
		Association("HasManyD").
		Unscoped().
		Clear()
	if err != nil {
		return errutil.Wrap(err)
	}

	res := tx.Updates(v)
	if res.Error != nil {
		return errutil.Wrap(res.Error)
	}
	return nil
}

type F struct {
	M
	HasManyG []G    `json:"g" gorm:"constraint:OnDelete:CASCADE"`
	V        string `json:"v"`
}

func (m *F) Update(tx *gorm.DB, a any) error {
	v := a.(*F)

	cc := *v
	err := tx.Model(&cc).
		Association("HasManyG").
		Unscoped().
		Clear()
	if err != nil {
		return errutil.Wrap(err)
	}

	res := tx.Updates(v)
	if res.Error != nil {
		return errutil.Wrap(res.Error)
	}
	return nil
}

type G struct {
	M
	FID      string `json:"-" gorm:"not null"`
	HasManyH []H    `json:"h" gorm:"constraint:OnDelete:CASCADE"`
	V        string `json:"v"`
}

type H struct {
	M
	GID string `json:"-" gorm:"not null"`
	V   string `json:"v"`
}

func initModRepo(t *testing.T) []A {
	if db == nil {
		dbPath := filepath.Join(fileutil.Home, "aurora-admin", "aurora-admin.db")
		fmt.Println("dbPath:", dbPath)
		Init(dbPath, []any{})
	}

	db.Exec("DROP TABLE IF EXISTS ds;")
	db.Exec("DROP TABLE IF EXISTS es;")
	db.Exec("DROP TABLE IF EXISTS cs;")
	db.Exec("DROP TABLE IF EXISTS bs;")
	db.Exec("DROP TABLE IF EXISTS 'as';")
	db.Exec("DROP TABLE IF EXISTS hs;")
	db.Exec("DROP TABLE IF EXISTS gs;")
	db.Exec("DROP TABLE IF EXISTS fs;")
	err := db.AutoMigrate(&A{}, &B{}, &C{}, &D{}, &E{}, &F{}, &G{}, &H{})
	errutil.Check(err)

	aRepo := NewRepo[A](db)

	a0 := []*A{
		{
			M: M{ID: "a1"},
			V: "a1",
			HasOneB: B{
				M: M{ID: "b1"},
				V: "b1",
				HasManyC: []C{
					{
						M: M{ID: "c1"},
						V: "c1",
					},
					{
						M: M{ID: "c2"},
						V: "c2",
					},
				},
			},
		},
		{
			M: M{ID: "aa1"},
			V: "aa1",
			HasOneB: B{
				M: M{ID: "bb1"},
				V: "bb1",
				HasManyC: []C{
					{
						M: M{ID: "cc1"},
						V: "cc1",
					},
					{
						M: M{ID: "cc2"},
						V: "cc2",
					},
				},
			},
		},
	}

	err = aRepo.Create(context.Background(), a0...)
	assert.Nil(t, err)
	a1, err := aRepo.GetAll(context.Background())
	assert.Nil(t, err)
	return a1
}

func initModRepoE(t *testing.T) ([]A, E) {
	a1 := initModRepo(t)

	eRepo := NewRepo[E](db)
	e0 := E{
		M: M{ID: "e1"},
		HasManyD: []D{
			{
				M:            M{ID: "d1"},
				BelongsToAID: "a1",
				V:            "d1",
			},
			{
				M:            M{ID: "d2"},
				BelongsToAID: "aa1",
				V:            "d2",
			},
		},
		V: "e1",
	}
	err := eRepo.Create(context.Background(), &e0)
	assert.Nil(t, err)
	e1, err := eRepo.Get(context.Background(), e0.ID)
	assert.Nil(t, err)
	return a1, e1
}

func initModRepoF(t *testing.T) F {
	initModRepo(t)

	fRepo := NewRepo[F](db)
	f0 := F{
		M: M{ID: "f1"},
		HasManyG: []G{
			{
				M: M{ID: "g1"},
				HasManyH: []H{
					{
						M: M{ID: "h1"},
						V: "h1",
					},
					{
						M: M{ID: "h2"},
						V: "h2",
					},
				},
				V: "g1",
			},
			{
				M: M{ID: "g2"},
				HasManyH: []H{
					{
						M: M{ID: "h3"},
						V: "h3",
					},
					{
						M: M{ID: "h4"},
						V: "h4",
					},
				},
				V: "g2",
			},
		},
		V: "f1",
	}
	err := fRepo.Create(context.Background(), &f0)
	assert.Nil(t, err)
	f1, err := fRepo.Get(context.Background(), f0.ID)
	assert.Nil(t, err)
	return f1
}

func assertModelEqMap(t *testing.T, v1, v2 any) {
	j1, err := json.Marshal(&v1)
	assert.Nil(t, err)
	j2, err := json.Marshal(&v2)
	assert.Nil(t, err)
	assert.JSONEq(t, string(j1), string(j2))
}

func TestCascadeCreate(t *testing.T) {
	a1 := initModRepo(t)

	assertModelEqMap(t, []map[string]any{
		{
			"id": "a1",
			"v":  "a1",
			"b": map[string]any{
				"id": "b1",
				"v":  "b1",
				"c": []map[string]any{
					{
						"id": "c1",
						"v":  "c1",
					},
					{
						"id": "c2",
						"v":  "c2",
					},
				},
			},
		},
		{
			"id": "aa1",
			"v":  "aa1",
			"b": map[string]any{
				"id": "bb1",
				"v":  "bb1",
				"c": []map[string]any{
					{
						"id": "cc1",
						"v":  "cc1",
					},
					{
						"id": "cc2",
						"v":  "cc2",
					},
				},
			},
		},
	}, a1)
}

func TestCascadeDelete(t *testing.T) {
	a1 := initModRepo(t)

	aRepo := NewRepo[A](db)
	bRepo := NewRepo[B](db)
	cRepo := NewRepo[C](db)
	err := aRepo.Delete(context.Background(), a1[0].ID, a1[1].ID)
	assert.Nil(t, err)

	aAll, err := aRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, aAll, 0)
	bAll, err := bRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, bAll, 0)
	cAll, err := cRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, cAll, 0)
}

func TestCascadeUpdate(t *testing.T) {
	initModRepo(t)

	aRepo := NewRepo[A](db)
	bRepo := NewRepo[B](db)
	cRepo := NewRepo[C](db)

	a2 := []*A{
		{
			M: M{ID: "a1"},
			V: "a2",
			HasOneB: B{
				M: M{ID: "b1"},
				V: "b2",
				HasManyC: []C{
					{
						M: M{ID: "c1"},
						V: "c4",
					},
					{
						M: M{ID: "c3"},
						V: "c3",
					},
				},
			},
		},
		{
			M: M{ID: "aa1"},
			V: "aa2",
			HasOneB: B{
				M:        M{ID: "bb1"},
				V:        "bb2",
				HasManyC: []C{},
			},
		},
	}

	err := aRepo.Update(context.Background(), a2...)
	assert.Nil(t, err)

	a3, err := aRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assertModelEqMap(t, []map[string]any{
		{
			"id": "a1",
			"v":  "a2",
			"b": map[string]any{
				"id": "b1",
				"v":  "b2",
				"c": []map[string]any{
					{
						"id": "c1",
						"v":  "c4",
					},
					{
						"id": "c3",
						"v":  "c3",
					},
				},
			},
		},
		{
			"id": "aa1",
			"v":  "aa2",
			"b": map[string]any{
				"id": "bb1",
				"v":  "bb2",
				"c":  []map[string]any{},
			},
		},
	}, a3)

	aAll, err := aRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, aAll, 2)
	bAll, err := bRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, bAll, 2)
	cAll, err := cRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, cAll, 2)
}

func TestCascadeCreateE(t *testing.T) {
	_, e1 := initModRepoE(t)
	assertModelEqMap(t, map[string]any{
		"id": "e1",
		"v":  "e1",
		"d": []map[string]any{
			{
				"id": "d1",
				"v":  "d1",
				"a": map[string]any{
					"id": "a1",
					"v":  "a1",
					"b": map[string]any{
						"id": "b1",
						"v":  "b1",
						"c": []map[string]any{
							{
								"id": "c1",
								"v":  "c1",
							},
							{
								"id": "c2",
								"v":  "c2",
							},
						},
					},
				},
			},
			{
				"id": "d2",
				"v":  "d2",
				"a": map[string]any{
					"id": "aa1",
					"v":  "aa1",
					"b": map[string]any{
						"id": "bb1",
						"v":  "bb1",
						"c": []map[string]any{
							{
								"id": "cc1",
								"v":  "cc1",
							},
							{
								"id": "cc2",
								"v":  "cc2",
							},
						},
					},
				},
			},
		},
	}, e1)
}

func TestCascadeDeleteE(t *testing.T) {
	_, e1 := initModRepoE(t)
	aRepo := NewRepo[A](db)
	bRepo := NewRepo[B](db)
	cRepo := NewRepo[C](db)
	dRepo := NewRepo[D](db)
	eRepo := NewRepo[E](db)
	err := eRepo.Delete(context.Background(), e1.ID)
	assert.Nil(t, err)

	aAll, err := aRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, aAll, 2)
	bAll, err := bRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, bAll, 2)
	cAll, err := cRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, cAll, 4)
	dAll, err := dRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, dAll, 0)
	eAll, err := eRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, eAll, 0)
}

func TestDeleteAAffectE(t *testing.T) {
	a1, e1 := initModRepoE(t)
	aRepo := NewRepo[A](db)
	dRepo := NewRepo[D](db)
	eRepo := NewRepo[E](db)
	err := aRepo.Delete(context.Background(), a1[0].ID, a1[1].ID)
	assert.Nil(t, err)

	dAll, err := dRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, dAll, 0)

	e2, err := eRepo.Get(context.Background(), e1.ID)
	assert.Nil(t, err)
	assertModelEqMap(t, map[string]any{
		"id": "e1",
		"v":  "e1",
		"d":  []map[string]any{},
	}, e2)
}

func TestCascadeUpdateE(t *testing.T) {
	_, e1 := initModRepoE(t)
	dRepo := NewRepo[D](db)
	eRepo := NewRepo[E](db)

	e2 := E{
		M: M{ID: "e1"},
		HasManyD: []D{
			{
				M:            M{ID: "d3"},
				BelongsToAID: "a1",
				V:            "d33",
			},
		},
		V: "e1",
	}

	err := eRepo.Update(context.Background(), &e2)
	assert.Nil(t, err)

	dAll, err := dRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, dAll, 1)

	e3, err := eRepo.Get(context.Background(), e1.ID)
	assert.Nil(t, err)
	assertModelEqMap(t, map[string]any{
		"id": "e1",
		"v":  "e1",
		"d": []map[string]any{
			{
				"id": "d3",
				"v":  "d33",
				"a": map[string]any{
					"id": "a1",
					"v":  "a1",
					"b": map[string]any{
						"id": "b1",
						"v":  "b1",
						"c": []map[string]any{
							{
								"id": "c1",
								"v":  "c1",
							},
							{
								"id": "c2",
								"v":  "c2",
							},
						},
					},
				},
			},
		},
	}, e3)
}

func TestCascadeUpdateF(t *testing.T) {
	initModRepoF(t)
	fRepo := NewRepo[F](db)
	gRepo := NewRepo[G](db)
	hRepo := NewRepo[H](db)

	f2 := F{
		M: M{ID: "f1"},
		HasManyG: []G{
			{
				M: M{ID: "g1"},
				HasManyH: []H{
					{
						M: M{ID: "h1"},
						V: "h5",
					},
					{
						M: M{ID: "h6"},
						V: "h6",
					},
				},
				V: "g2",
			},
		},
		V: "f2",
	}

	err := fRepo.Update(context.Background(), &f2)
	assert.Nil(t, err)

	gAll, err := gRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, gAll, 1)
	hAll, err := hRepo.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, hAll, 2)

	f3, err := fRepo.Get(context.Background(), f2.ID)
	assert.Nil(t, err)
	assertModelEqMap(t, map[string]any{
		"id": "f1",
		"v":  "f2",
		"g": []map[string]any{
			{
				"id": "g1",
				"v":  "g2",
				"h": []map[string]any{
					{
						"id": "h1",
						"v":  "h5",
					},
					{
						"id": "h6",
						"v":  "h6",
					},
				},
			},
		},
	}, f3)
}
