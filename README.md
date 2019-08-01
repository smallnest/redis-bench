# redis-bench

redis-bench is a different tool with redis-benchmark tool. Its target is testing latency of redis under given throughputs.

There are some redis benchmark tool:

- [redis-benchmark](https://redis.io/topics/benchmarks)
- [vire-benchmark](https://github.com/vipshop/vire)
- [tidwall/redbench](https://github.com/tidwall/redbench)

But their target is testing throughputs of redis. In high concurrent clients case, the latency is very big (> 1 second), which is not acceptable in production. So we want to benchmark to get the acceptable throuphputs, latency and concurrency.

I searched but not found any tool for this case, so I create this project base on code of [tidwall/redbench](https://github.com/tidwall/redbench).

## Usage

```sh
rb -h
Usage of rb:
  -P int
    	pipeline <numreq> requests. (default 1 no pipeline). (default 1)
  -c int
    	number of concurrent connections (default 100)
  -cpu int
    	max cpus count to run (default logical cpu cores) (default 4)
  -d int
    	data size of SET/GET/... value in bytes (default 16)
  -host string
    	server address (default "127.0.0.1")
  -l float
    	max throughputs (requests/s) (default 10000)
  -n int
    	total number of requests (default 1000000)
  -p int
    	server port (default 6379)
  -r int
    	use random keys for SET/GET (default 10000)
  -s string
    	server address (overrides host and port) (default "10.41.15.226:6000")
  -t string
    	test type. only support set|get (default "set")
```

run the command:

```sh
rb -t get -r 10000000 -n 20000000 -s 127.0.0.1:6379  -cpu 16 -c 200 -l 100000
```