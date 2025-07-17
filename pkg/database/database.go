package database

import (
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
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

func ClearAll(db *gorm.DB, dst ...any) error {
	for _, m := range dst {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m)
	}
	return nil
}
