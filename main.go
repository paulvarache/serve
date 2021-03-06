package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"serve/middlewares"
)

var (
	opts        *CliOptions
	rootHandler http.Handler
	version     bool = false
)

// These variables are provided at build time using ldflags
var (
	BuildVersion string = "dev"
	BuildHash    string = "dev"
)

func init() {
	opts = NewCliOptions()

	FlagBool(&version, "version", "v", false, "Display the version.")
	FlagString(&opts.SpaIndex, "spa", "s", "", "Path to a Single Page App, e.g. app.html.")
	FlagString(&opts.Dir, "directory", "d", "", "Root directory, defaults to the current directory.")
	FlagBool(&opts.Open, "open", "o", false, "Automatically open the default system browser.")
	FlagInt(&opts.Port, "port", "p", 8000, "The port number to listen to.")
	FlagString(&opts.Hostname, "hostname", "h", "localhost", "The hostname or IP to bind to. Defaults to 0.0.0.0 (any host).")
	FlagBool(&opts.Log, "log", "l", false, "Display requests logs.")
	FlagBool(&opts.Compress, "compress", "z", false, "Serve gzip-compressed resources, where applicable.")

	flag.Parse()
}

func main() {
	// Display version if asked for it
	if version {
		fmt.Printf("Serve version %s-%s", BuildVersion, BuildHash)
		return
	}
	dir, err := filepath.Abs(opts.Dir)
	if err != nil {
		panic(err)
	}
	rootHandler = middlewares.NullHandler(dir)
	if opts.SpaIndex != "" {
		rootHandler = middlewares.SpaMiddleware(dir, opts.SpaIndex, &rootHandler)
	}
	if opts.Compress {
		rootHandler = middlewares.CompressMiddleware(rootHandler)
	}
	if opts.Log {
		rootHandler = middlewares.LogMiddleware(rootHandler)
	}
	rootHandler = middlewares.CorsMiddleware(rootHandler)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opts.Hostname, opts.Port))
	if err != nil {
		panic(err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	fmt.Printf("Listening on: http://%s:%d", addr.IP, addr.Port)
	err = http.Serve(listener, rootHandler)
	if err != nil {
		panic(err)
	}
}
