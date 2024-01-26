package middleware

import (
	"bytes"
	"fmt"
	"github.com/taejune/echo-server-go/pkg/base"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Option struct {
	Quiet  bool   // whether to include request content into response
	Size   int64  // max size of payload
	Pretty bool   // whether to pretty print json
	Format string // json or text
}

var option Option
var hostname string = "localhost"

const (
	DefaultLen = 1024
	MaxLen     = DefaultLen * 1024 * 4
)

func init() {
	option = Option{
		Quiet:  true,
		Size:   MaxLen,
		Pretty: false,
		Format: "TEXT",
	}
}

func InitOption() {
	option.Pretty = base.StringToBool(os.Getenv("LOG_PRETTY"))
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "FATAL":
	case "ERROR":
	case "WARN":
		break
	case "INFO":
		option.Quiet = false
		option.Size = DefaultLen
		break
	case "DEBUG":
		option.Quiet = false
		option.Size = DefaultLen * 10
		break
	case "TRACE":
		option.Quiet = false
		option.Size = MaxLen
		break
	default:
		option.Quiet = true
		option.Size = DefaultLen
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		output := fmt.Sprintf("[%s] %s %s %s %s%s %s %s %d",
			time.Now().Format(time.StampMicro), hostname, r.RemoteAddr, r.Host,
			r.Method, r.RequestURI, r.Proto, r.UserAgent(), r.ContentLength)

		if !option.Quiet {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			body, _ = base.MarshalJson(string(body), option.Size, option.Pretty)
			output = fmt.Sprintf("%s\n%s", output, body)
		}

		log.Printf(output)
		next.ServeHTTP(w, r)
	})
}
