package convutil

import (
	"encoding/json"
	"fmt"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
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

func Json(a any) ([]byte, error) {
	return json.MarshalIndent(a, "", "    ")
}

func MustJsonString(a any) string {
	content, err := Json(a)
	errutil.Check(err)
	return string(content)
}
