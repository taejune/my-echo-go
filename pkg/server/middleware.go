package server

import (
	"bytes"
	"fmt"
	"github.com/taejune/echo-server-go/pkg/base"
	"io"
	"log"
	"net/http"
	"time"
)

var hostname string = "localhsot"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		output := fmt.Sprintf("[%s] %s %s %s %s%s %s %s %d",
			time.Now().Format(time.StampMicro), hostname, r.RemoteAddr, r.Host, r.Method, r.RequestURI, r.Proto, r.UserAgent(), r.ContentLength)

		var body []byte
		if base.Option.Payload {
			body, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			var contents []byte
			if r.Header.Get("Content-Type") == "application/json" || r.Header.Get("Content-Type") == "application/x-ndjson" {
				contents, _ = base.MarshalJson(string(body), base.Option.Size, base.Option.Pretty)
			}
			output = fmt.Sprintf("%s\n%s", output, contents)
		}

		log.Printf(output)
		next.ServeHTTP(w, r)
	})
}
