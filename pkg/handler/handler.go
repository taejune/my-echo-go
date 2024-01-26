package handler

import (
	"encoding/json"
	"fmt"
	"github.com/taejune/echo-server-go/pkg/base"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DefaultLen = 1024
	MaxLen     = DefaultLen * 1024 * 4
)

type Option struct {
	Quiet  bool   // whether to include request into response (do not echo)
	Size   int64  // max size of payload
	Pretty bool   // whether to pretty print json
	Format string // json or text
}

var option Option
var hostname string

const (
	FormatJson = "json"
	FormatText = "text"
)

func init() {
	option = Option{
		Quiet:  false,
		Size:   DefaultLen,
		Pretty: false,
		Format: FormatText,
	}
	hostname, _ = os.Hostname()
}

func InitOption() {
	if base.StringToBool(os.Getenv("ECHO_JSON")) {
		option.Format = FormatJson
	}
	option.Pretty = base.StringToBool(os.Getenv("ECHO_PRETTY"))
	switch strings.ToUpper(os.Getenv("ECHO_LEVEL")) {
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
		option.Size = DefaultLen * 100
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

func Echo(w http.ResponseWriter, r *http.Request) {
	msg := parseHttpRequest(r, option.Format)
	if _, err := w.Write(msg); err != nil {
		log.Println("failed to response with" + err.Error())
	}
}

func parseHttpRequest(r *http.Request, format string) []byte {
	defer r.Body.Close()

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read body")
	}
	body, err := base.MarshalJson(string(dat), option.Size, option.Pretty)
	if err != nil {
		log.Println("failed to marshal body")
	}

	output := []byte{}
	if format == FormatJson {
		req := make(map[string]any)
		req["Time"] = time.Now().Format(time.StampMicro)
		req["Remote"] = r.RemoteAddr
		req["Host"] = r.Host
		req["Method"] = r.Method
		req["Proto"] = r.Proto
		req["URI"] = r.RequestURI
		if len(r.TransferEncoding) > 0 {
			req["Encoding"] = r.TransferEncoding
		}
		if len(r.Cookies()) > 0 {
			req["Cookies"] = make([]*http.Cookie, len(r.Cookies()))
			req["Cookies"] = r.Cookies()

		}
		req["ContentLength"] = r.ContentLength
		req["UserAgent"] = r.UserAgent()

		req["Headers"] = make(map[string]interface{})
		headers := req["Headers"].(map[string]interface{})
		for k, v := range r.Header {
			headers[k] = strings.Join(v, ",")
		}

		if !option.Quiet {
			req["Body"] = string(body)
		}

		output, _ = json.Marshal(req)
	} else {
		req := fmt.Sprintf("[%s] %s %s %s %s%s %s %s %d",
			time.Now().Format(time.StampMicro), hostname, r.RemoteAddr, r.Host,
			r.Method, r.RequestURI, r.Proto, r.UserAgent(), r.ContentLength)

		if option.Quiet {
			output = []byte(fmt.Sprintf("%s", req))
		} else {
			output = []byte(fmt.Sprintf("%s %s", req, body))
		}
	}

	return output
}
