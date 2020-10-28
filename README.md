# zero-tools
go 一些有用公共方法，比如批量处理的任务、安全的goroutines、http并发call、数据流等

# 测试用例
```
zero-tools [main] ⚡ go test ./...
ok      github.com/binwen/zero-tools/executors  1.192s
?       github.com/binwen/zero-tools/lang       [no test files]
?       github.com/binwen/zero-tools/proc       [no test files]
ok      github.com/binwen/zero-tools/rescue     0.053s
ok      github.com/binwen/zero-tools/stream     0.245s
ok      github.com/binwen/zero-tools/stringx    0.016s
ok      github.com/binwen/zero-tools/syncx      6.049s
ok      github.com/binwen/zero-tools/threading  0.019s
ok      github.com/binwen/zero-tools/timex      0.084s

zero-tools [main] ⚡ go test ./... -bench=. -benchmem -count=3
goos: darwin
goarch: amd64
pkg: github.com/binwen/zero-tools/executors
BenchmarkBulkExecutor-8              922           1264664 ns/op              36 B/op          0 allocs/op
BenchmarkBulkExecutor-8              921           1265714 ns/op              36 B/op          0 allocs/op
BenchmarkBulkExecutor-8              906           1281604 ns/op              37 B/op          0 allocs/op
BenchmarkChunkExecutor-8             920           1267705 ns/op              68 B/op          1 allocs/op
BenchmarkChunkExecutor-8             921           1269165 ns/op              68 B/op          1 allocs/op
BenchmarkChunkExecutor-8             940           1266706 ns/op              68 B/op          1 allocs/op
BenchmarkExecutor-8                  100          40171585 ns/op              34 B/op          0 allocs/op
BenchmarkExecutor-8                  100          40126026 ns/op              35 B/op          0 allocs/op
BenchmarkExecutor-8                  100          40132625 ns/op              43 B/op          0 allocs/op

pkg: github.com/binwen/zero-tools/stream
BenchmarkMapReduce-8      170424              6691 ns/op             776 B/op         13 allocs/op
BenchmarkMapReduce-8      175506              6536 ns/op             776 B/op         13 allocs/op
BenchmarkMapReduce-8      181098              6480 ns/op             776 B/op         13 allocs/op

pkg: github.com/binwen/zero-tools/stringx
BenchmarkRandString-8           15989094                68.9 ns/op            16 B/op          1 allocs/op
BenchmarkRandString-8           17623638                68.5 ns/op            16 B/op          1 allocs/op
BenchmarkRandString-8           16802840                68.0 ns/op            16 B/op          1 allocs/op

```
