package main

import (
	"flag"
)

var (
	host = flag.String("host", "127.0.0.1", "The ip of the http server")
	port = flag.Uint("port", 8080, "The port the the http server")
	tls = flag.Bool("tls", false, "Whether to enable tls support")
)

func main() {
	flag.Parse()
	http_server(*host, *port, *tls)
}