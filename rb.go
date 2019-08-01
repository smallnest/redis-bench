package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strings"

	"github.com/smallnest/redbench"
)

var (
	s        = flag.String("s", "10.41.15.226:6000", "server address")
	c        = flag.Int("c", 100, "number of concurrent connections (default 100)")
	threads  = flag.Int("T", runtime.GOMAXPROCS(-1), "threads count to run (default logical cpu cores)")
	d        = flag.Int("d", 16, "data size of SET/GET/... value in bytes (default 16)")
	r        = flag.Int("r", 10000, "use random keys for SET/GET (default 10000)")
	n        = flag.Int("n", 1000000, "total number of requests (default 1000000)")
	t        = flag.String("t", "set", "test type. only support set|get")
	pipeline = flag.Int("P", 1, "pipeline <numreq> requests. default 1 (no pipeline).")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*threads)

	key := func() string {
		return fmt.Sprintf("mystring:%012d", rand.Intn(*r))
	}
	value := strings.Repeat("A", *d)

	opts := *redbench.DefaultOptions
	opts.Clients = *c
	opts.Requests = *n
	opts.Pipeline = *pipeline

	*t = strings.ToLower(*t)

	switch *t {
	case "set":
		redbench.Bench("SET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "SET", key(), value)
		})
	case "get":
		redbench.Bench("GET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "GET", key())
		})
	}
}
