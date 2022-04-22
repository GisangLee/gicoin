package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/gisanglee/gicoin/explorer"
	"github.com/gisanglee/gicoin/rest"
)

func usage() {
	fmt.Printf("Welcome to Gi Coin\n\n")
	fmt.Printf("Please use the following flags: \n\n")
	fmt.Printf("--port:    Set the port of the server\n")
	fmt.Printf("--mode:    Choose between 'html' and 'rest (recommended)'\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	if len(os.Args) < 2 {
		usage()
	}

	port := flag.Int("port", 4000, "Set the port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html', 'rest', 'both' ")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	case "both":
		go explorer.Start(*port)
		rest.Start(*port + 1)
	default:
		usage()
	}
}
