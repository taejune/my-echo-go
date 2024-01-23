package server

import (
	"github.com/taejune/echo-server-go/pkg/base"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	MaxSize = 1024 * 1024 * 5 / 4
)

func Echo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	msg := parseHttpRequest(r)

	if base.Option.Payload {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read body")
		}

		var contents []byte
		if r.Header.Get("Content-Type") == "application/json" || r.Header.Get("Content-Type") == "application/x-ndjson" {
			contents, _ = base.MarshalJson(string(body), base.Option.Size, base.Option.Pretty)
		}
		msg["Payload"] = string(contents)
	} else {
		msg["Payload"] = ""
	}

	out, err := base.MarshalJson(msg, MaxSize, false)
	if err != nil {
		log.Println("failed to marshal json")
	}

	if _, err := w.Write(out); err != nil {
		log.Println("failed to response with" + err.Error())
	}
}

func parseHttpRequest(r *http.Request) map[string]any {
	output := make(map[string]any)
	output["Remote"] = r.RemoteAddr
	output["Host"] = r.Host
	output["Method"] = r.Method
	output["Proto"] = r.Proto
	output["URI"] = r.RequestURI
	output["Encoding"] = r.TransferEncoding
	output["Cookies"] = make([]*http.Cookie, len(r.Cookies()))
	output["Cookies"] = r.Cookies()
	output["ContentLength"] = r.ContentLength
	output["UserAgent"] = r.UserAgent()
	output["Payload"] = ""

	output["Headers"] = make(map[string]interface{})
	headers := output["Headers"].(map[string]interface{})
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ",")
	}

	return output
}
