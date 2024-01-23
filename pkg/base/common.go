package base

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type PayloadPrintOpt struct {
	Size    int64
	Payload bool
	Pretty  bool
	Format  string
}

var Option PayloadPrintOpt

func InitByEnv() {
	switch strings.ToLower(os.Getenv("include_payload")) {
	case "true":
		Option.Payload = true
	case "false":
		Option.Payload = false
	default:
		Option.Payload = false
	}

	var err error
	Option.Size, err = strconv.ParseInt(os.Getenv("max_payload_len"), 10, 64)
	if err != nil {
		Option.Size = 1000
	}

	switch strings.ToLower(os.Getenv("pretty")) {
	case "true":
		Option.Pretty = true
	case "false":
		Option.Pretty = false
	default:
		Option.Pretty = false
	}
}

func MarshalJson(v interface{}, maxSize int64, pretty bool) ([]byte, error) {
	var dat []byte
	var err error
	if pretty {
		dat, err = json.MarshalIndent(v, "", "  ")
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
	r = strings.ReplaceAll(r, " ", "")

	return []byte(r), nil
}
