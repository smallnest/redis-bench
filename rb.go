package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/juju/ratelimit"
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
	f        = flag.Int("f", 100, "Use random fields for SADD/HSET/... (default 100)")
	n        = flag.Int("n", 1000000, "total number of requests")
	t        = flag.String("t", "set", "Only run the comma separated list of tests.")
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

	// bench options
	opts := *redbench.DefaultOptions
	opts.Clients = *c
	opts.Requests = *n
	opts.Pipeline = *pipeline
	if *l > 0 {
		rate := int64(*l) / 1000
		opts.Limter = ratelimit.NewBucketWithQuantum(time.Millisecond, rate, rate)
	}

	*t = strings.ToLower(*t)
	commands := strings.Split(*t, ",")
	for _, cmd := range commands {
		bench := benches[cmd]
		if bench != nil {
			bench(cmd, *s, &opts)
		}
	}
}

func keyN(n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", rand.Intn(*r))
}
func key() string {
	return fmt.Sprintf("%d", rand.Intn(*r))
}
func keyField() string {
	return fmt.Sprintf("%d", rand.Intn(*f))
}
func keyFieldN(n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", rand.Intn(*f))
}

var value = strings.Repeat("A", *d)

type benchFunc func(string, string, *redbench.Options)

var benches = map[string]benchFunc{
	"ping": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "PING")
		})
	},
	"set": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "SET", "mystring:"+keyN(12), value)
		})
	},
	"get": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "GET", "mystring:"+keyN(12))
		})
	},
	"getset": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "GETSET", "mystring:"+keyN(12), value)
		})
	},
	"mset": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "MSET", "mystring:"+keyN(12), value, "mystring:"+keyN(12), value, "mystring:"+keyN(12), value)
		})
	},
	"mget": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "MGET", "mystring:"+keyN(12), "mystring:"+keyN(12), "mystring:"+keyN(12))
		})
	},
	"incr": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "INCR", "mynum:"+keyN(12))
		})
	},
	"decr": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "DECR", "mynum:"+keyN(12))
		})
	},
	"hset": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HSET", "myhash:"+keyN(12), "field:"+keyFieldN(14), value)
		})
	},
	"hget": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HGET", "myhash:"+keyN(12), "field:"+keyFieldN(14))
		})
	},
	"hdel": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HDEL", "myhash:"+keyN(12), "field:"+keyFieldN(14))
		})
	},
	"hmset": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HMSET", "myhash:"+keyN(12), "field:"+keyFieldN(14), value, "field:"+keyFieldN(14), value, "field:"+keyFieldN(14), value)
		})
	},
	"hmget": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HMGET", "myhash:"+keyN(12), "field:"+keyFieldN(14), "field:"+keyFieldN(14), "field:"+keyFieldN(14))
		})
	},
	"hkeys": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HKEYS", "myhash:"+keyN(12))
		})
	},
	"hvals": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HVALS", "myhash:"+keyN(12))
		})
	},
	"hgetall": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "HGETALL", "myhash:"+keyN(12))
		})
	},
	"lpush": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "LPUSH", "mylist:"+keyN(12), value)
		})
	},
	"lpop": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "LPOP", "mylist:"+keyN(12))
		})
	},
	"rpush": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "RPUSH", "mylist:"+keyN(12), value)
		})
	},
	"rpop": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "RPOP", "mylist:"+keyN(12))
		})
	},
	"lrange": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "LRANGE", "mylist:"+keyN(12), "0", "1000")
		})
	},
	"sadd": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "SADD", "myset::"+keyN(12), keyFieldN(14))
		})
	},
	"spop": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "SPOP", "myset::"+keyN(12))
		})
	},
	"smember": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "smember", "myset::"+keyN(12))
		})
	},
	"sismember": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "SISMEMBERS", "myset::"+keyN(12), keyFieldN(14))
		})
	},
	"zadd": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "ZADD", "mysortedset::"+keyN(12), key(), keyFieldN(14))
		})
	},
	"zrem": func(name string, addr string, opts *redbench.Options) {
		redbench.Bench(strings.ToUpper(name), addr, opts, nil, func(buf []byte) []byte {
			return redbench.AppendCommand(buf, "ZREM", "mysortedset::"+keyN(12), keyFieldN(14))
		})
	},
}
