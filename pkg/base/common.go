package base

import (
	"encoding/json"
	"strings"
)

func MarshalJson(v interface{}, maxSize int64, pretty bool) ([]byte, error) {
	var dat []byte
	var err error
	if pretty {
		dat, err = json.MarshalIndent(v, " ", " ")
		if err != nil {
			return nil, err
		}
	} else {
		dat, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}

	var printLen int64
	if int64(len(dat)) > maxSize {
		printLen = maxSize + 3
		dat[printLen-3] = '.'
		dat[printLen-2] = '.'
		dat[printLen-1] = '.'
	} else {
		printLen = int64(len(dat))
	}
	dat = dat[:printLen]

	r := strings.ReplaceAll(string(dat), "\\n", "")
	r = strings.ReplaceAll(r, "\\", "")
	r = strings.ReplaceAll(r, "\"", "")
	//r = strings.ReplaceAll(r, " ", "")

	return []byte(r), nil
}

func StringToBool(str string) bool {
	switch strings.ToLower(str) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no", "":
		return false
	default:
		return false
	}
}
