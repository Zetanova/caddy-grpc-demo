package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Log the env variables required for a reverse proxy
// func logSetup() {
// 	a_condtion_url := os.Getenv("A_CONDITION_URL")
// 	b_condtion_url := os.Getenv("B_CONDITION_URL")
// 	default_condtion_url := os.Getenv("DEFAULT_CONDITION_URL")

// 	log.Printf("Server will run on: %s\n", getListenAddress())
// 	log.Printf("Redirecting to A url: %s\n", a_condtion_url)
// 	log.Printf("Redirecting to B url: %s\n", b_condtion_url)
// 	log.Printf("Redirecting to Default url: %s\n", default_condtion_url)
// }

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http2.Transport{
		// So http2.Transport doesn't complain the URL scheme isn't 'https'
		AllowHTTP: true,
		// Pretend we are dialing a TLS endpoint. (Note, we ignore the passed tls.Config)
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	//requestPayload := parseRequestBody(req)
	//url := getProxyUrl(requestPayload.ProxyCondition)
	url := "http://localhost:50051"

	//logRequestPayload(requestPayload, url)

	serveReverseProxy(url, res, req)
}

func main() {
	// Log setup values
	//logSetup()

	h2s := &http2.Server{}

	handler := http.HandlerFunc(handleRequestAndRedirect)

	server := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: h2c.NewHandler(handler, h2s),
	}

	fmt.Printf("Listening [0.0.0.0:80]...\n")

	// start server
	//http.HandleFunc("/", handleRequestAndRedirect)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
