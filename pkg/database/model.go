package database

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

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

func (m *Model) Update(tx *gorm.DB, a any) error {
	res := tx.Updates(a)
	return res.Error
}

func NewID() string {
	return uuid.New().String()
}

func GenNestedPreloads(t any) []string {
	return genNestedGenPreloads(reflect.TypeOf(t).Elem())
}

func genNestedGenPreloads(t reflect.Type) []string {
	var res []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}

		fieldName := field.Name
		var fieldType reflect.Type
		switch field.Type.Kind() {
		case reflect.Struct:
			fieldType = field.Type
		case reflect.Slice:
			fieldType = field.Type.Elem()
		default:
			continue
		}

		idField, ok := fieldType.FieldByName("ID")
		if !ok {
			continue
		}
		gormTag := idField.Tag.Get("gorm")
		if !strings.Contains(gormTag, "primarykey") {
			continue
		}

		v := genNestedGenPreloads(fieldType)
		if len(v) == 0 {
			res = append(res, fieldName)
		} else {
			for _, vv := range v {
				res = append(res, fmt.Sprintf("%s.%s", fieldName, vv))
			}
		}
	}
	return res
}
