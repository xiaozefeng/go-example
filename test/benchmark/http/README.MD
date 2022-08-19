# http server 压测

## 我使用 wrk 这个工具来做压测
```bash
wrk -t12 -c400 -d30s http://127.0.0.1:8080
# 在我机器是这个表现
Running 30s test @ http://127.0.0.1:8080
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.92ms    5.45ms  75.14ms   80.24%
    Req/Sec     2.91k     1.59k   21.44k    87.28%
  1045475 requests in 30.09s, 126.62MB read
  Socket errors: connect 155, read 0, write 0, timeout 0
Requests/sec:  34740.16
Transfer/sec:      4.21MB
```

## 查询分析结果
```bash
go tool pprof -http=:9000 http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -http=:9000 http://127.0.0.1:6060/debug/pprof/mutex
go tool pprof -http=:9000 http://127.0.0.1:6060/debug/pprof/block
```
