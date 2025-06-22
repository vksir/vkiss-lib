package database

import (
	"github.com/google/uuid"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
	"time"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

type Model struct {
	ID        string    `gorm:"primarykey;size:36" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

func (m *Model) Preload(tx *gorm.DB) *gorm.DB {
	return tx
}

func (m *Model) Update(tx *gorm.DB, a any) error {
	res := tx.Updates(a)
	return res.Error
}

func Init(dbPath string, migrate []any) {
	sql := sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbPath,
	}

	var err error
	db, err = gorm.Open(sql, &gorm.Config{})
	errutil.Check(err)

	res := db.Exec("PRAGMA foreign_keys = ON", nil)
	errutil.Check(res.Error)

	err = db.AutoMigrate(migrate...)
	errutil.Check(err)
}
