package main

import (
	"github.com/taejune/echo-server-go/pkg/base"
	"github.com/taejune/echo-server-go/pkg/server"
	"log"
	"net/http"
	"os"
)

func main() {
	base.InitByEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.Handle("/", server.LoggingMiddleware(http.HandlerFunc(server.Echo)))

	certPath, isCertExist := os.LookupEnv("CERT_PATH")
	if certPath == "" {
		certPath = "tls/server.crt"
	}
	privateKeyPath, isPrivateKeyExist := os.LookupEnv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "tls/server.key"
	}

	if isCertExist && isPrivateKeyExist {
		go func() {
			if _, err := os.Stat(certPath); os.IsNotExist(err) {
				log.Println(err)
				return
			}
			if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
				log.Println(err)
				return
			}
			log.Println("Listening TLS enabled on" + port)
			log.Fatal(http.ListenAndServeTLS(":443", certPath, privateKeyPath, mux))
		}()
		log.Fatal(http.ListenAndServe(":8080", mux))
	} else {
		log.Println("Listening on " + port)
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}
}
