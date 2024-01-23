package server

import (
	"github.com/taejune/echo-server-go/pkg/base"
	"io"
	"log"
	"net/http"
	"strings"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	msg := parseHttpRequest(r)

	defer r.Body.Close()

	if base.Option.Payload {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read body")
		}
		msg["Payload"] = string(body)
	} else {
		msg["Payload"] = ""
	}

	out, err := base.MarshalJson(msg, base.Option.Size, false)
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
