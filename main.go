package main

import (
    "os"
    //"fmt"
    "strings"
    "net"
    "net/http"
    "log"
    "encoding/json"
    "io/ioutil"
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

    log.Println("Listening on " + port)
    http.ListenAndServe(port, mux)
}

func logging_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        host, port, _ := net.SplitHostPort(r.Host)
        log.Printf("[%s] %s -> http://%s:%s%s\n", r.Method, r.RemoteAddr, host, port, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

func echo(w http.ResponseWriter, r *http.Request) {
    data := make(map[string]interface{})

    w.Header().Set("Content-Type", "application/json")

    data["method"] = r.Method
    data["scheme"] = r.URL.Scheme
    data["host"], data["port"], _ = net.SplitHostPort(r.Host)
    data["path"] = r.URL.Path
    data["content-length"] = r.ContentLength

    // headers
    data["headers"] = make(map[string]interface{})
    headers := data["headers"].(map[string]interface{})
    for k, v := range r.Header  {
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
    data["query"] = make(map[string]interface{})
    q := data["query"].(map[string]interface{})
    for k, v := range r.URL.Query() {
        q[k] = v
    }

    payload, _ := json.Marshal(data)
    w.Write(payload)
}


