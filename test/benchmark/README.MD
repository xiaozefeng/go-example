# 基础测试

## bench的参数
```bash
-count # 表示执行的次数
-benchtime # 表示一次执行多久，默认一秒 ，可以指定  2s 或者  多次  比如 500x
-benchmem # 输出内存分配情况

# 性能分析
-cpuprofile=cpu.out
-memprofile=mem.out

# 内置工具分析
go tool pprof -http=":8080" cpu.out 
go tool pprof -http=":8080" mem.out 

# 或者使用 google 的 pprof 功能
go install github.com/google/pprof@latest

# 然后执行
pprof -http=:8080 cpu.out 
pprof -http=:8080 mem.out 
```
