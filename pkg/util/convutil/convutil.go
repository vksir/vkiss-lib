package convutil

import (
	"fmt"
	"strconv"
)

func String(a any) string {
	switch a.(type) {
	case int:
		return strconv.Itoa(a.(int))
	case bool:
		if a.(bool) {
			return "true"
		} else {
			return "false"
		}
	default:
		return fmt.Sprintf("%v", a)
	}
}
