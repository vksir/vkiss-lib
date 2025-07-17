package database

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAutoNestedPreload(t *testing.T) {
	res := genNestedGenPreloads(reflect.TypeOf(&A{}).Elem())
	fmt.Println(res)
}
