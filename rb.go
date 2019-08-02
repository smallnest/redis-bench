package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/time/rate"

	"github.com/smallnest/redbench"
)

var (
	h        = flag.String("host", "127.0.0.1", "server address")
	p        = flag.Int("p", 6379, "server port")
	s        = flag.String("s", "", "server address (overrides host and port)")
	c        = flag.Int("c", 100, "number of concurrent connections")
	l        = flag.Float64("l", 10000.0, "max throughputs (requests/s)")
	cpus     = flag.Int("cpu", runtime.GOMAXPROCS(-1), "max cpus count to run (default logical cpu cores)")
	d        = flag.Int("d", 16, "data size of SET/GET/... value in bytes")
	r        = flag.Int("r", 10000, "use random keys for SET/GET")
	n        = flag.Int("n", 1000000, "total number of requests")
	t        = flag.String("t", "set", "test type. only support set|get")
	pipeline = flag.Int("P", 1, "pipeline <numreq> requests. (default 1 no pipeline).")
)

func main() {
	flag.Parse()

	// address
	if *s == "" {
		if *h == "" {
			*h = "127.0.0.1"
		}
		*s = *h + ":" + strconv.Itoa(*p)
	}

	// set max CPU
	runtime.GOMAXPROCS(*cpus)

	// prepare data
	key := func() string {
		return fmt.Sprintf("mystring:%012d", rand.Intn(*r))
	}
	numkey := func() string {
		return fmt.Sprintf("mynum:%012d", rand.Intn(*r))
	}

	value := strings.Repeat("A", *d)

	// bench options
	opts := *redbench.DefaultOptions
	opts.Clients = *c
	opts.Requests = *n
	opts.Pipeline = *pipeline
	if *l > 0 {
		opts.Limter = rate.NewLimiter(rate.Limit(*l), 1)
	}

	*t = strings.ToLower(*t)
	switch *t {
	case "ping":
		redbench.Bench("PING", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "PING")
		})
	case "set":
		redbench.Bench("SET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "SET", key(), value)
		})
	case "get":
		redbench.Bench("GET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "GET", key())
		})
	case "getset":
		redbench.Bench("SET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "GETSET", key(), value)
		})
	case "mset":
		redbench.Bench("MSET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "MSET", key(), value, key(), value, key(), value)
		})
	case "mget":
		redbench.Bench("MGET", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "MGET", key(), key(), key())
		})
	case "incr":
		redbench.Bench("INCR", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "INCR", numkey())
		})
	case "decr":
		redbench.Bench("DECR", *s, &opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "DECR", numkey())
		})
	}
}
