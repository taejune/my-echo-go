package main

import (
	"context"
	"os"
	//"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	// b64 "encoding/base64"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + string(port)

	mux := http.NewServeMux()
	mux.Handle("/", logging_middleware(http.HandlerFunc(echo)))

	certPath, isCertExist := os.LookupEnv("CERT_PATH")
	privateKeyPath, isPrivateKeyExist := os.LookupEnv("PRIVATE_KEY_PATH")

	if isCertExist && isPrivateKeyExist {
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			log.Fatal(err)
			return
		}
		if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
			log.Fatal(err)
			return
		}
		log.Println("Listening TLS enabled on" + port)
		log.Fatal(http.ListenAndServeTLS(port, certPath, privateKeyPath, mux))
	} else {
		log.Println("Listening on " + port)
		log.Fatal(http.ListenAndServe(port, mux))
	}
}

func logging_middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, port, _ := net.SplitHostPort(r.Host)
		data := make(map[string]interface{})

		data["method"] = r.Method
		data["host"], data["port"], _ = net.SplitHostPort(r.Host)
		//data["scheme"] = r.URL.Scheme
		//data["path"] = r.URL.Path
		//data["content-length"] = r.ContentLength
		data["headers"] = make(map[string]interface{})
		headers := data["headers"].(map[string]interface{})
		for k, v := range r.Header {
			headers[k] = strings.Join(v, ",")
		}

		// body
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic("Couldn't get body")
		}
		data["body"] = body

		// query params
		//data["query"] = make(map[string]interface{})
		//q := data["query"].(map[string]interface{})
		//for k, v := range r.URL.Query() {
		//	q[k] = v
		//}

		info, _ := json.MarshalIndent(data, "", " ")
		ctxWithReqInfo := context.WithValue(r.Context(), "req-info", data)
		rWithInfo := r.WithContext(ctxWithReqInfo)

		log.Printf("[%s] %s -> http://%s:%s%s\n%v\n", r.Method, r.RemoteAddr, host, port, r.RequestURI, string(info))

		next.ServeHTTP(w, rWithInfo)
	})
}

func echo(w http.ResponseWriter, r *http.Request) {
	info := r.Context().Value("req-info")

	payload, _ := json.MarshalIndent(info, "", " ")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		log.Println("failed to response with" + err.Error())
	}
}
