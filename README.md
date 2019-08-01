# redis-bench

redis-bench is a different tool with redis-benchmark tool. Its target is testing latency of redis under given throughputs.

There are some redis benchmark tool:

- [redis-benchmark](https://redis.io/topics/benchmarks)
- [vire-benchmark](https://github.com/vipshop/vire)
- [tidwall/redbench](https://github.com/tidwall/redbench)

But their target is testing throughputs of redis. In high concurrent clients case, the latency is very big (> 1 second), which is not acceptable in production. So we want to benchmark to get the acceptable throuphputs, latency and concurrency.

I searched but not found any tool for this case, so I create this project base on code of [tidwall/redbench](https://github.com/tidwall/redbench).