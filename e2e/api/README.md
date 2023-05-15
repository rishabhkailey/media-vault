# E2E tests
## What not to test in E2E
* request validations, e.g. bad request (tested in unit tests)

```bash
go test -timeout 30s -run ^<test-name>$ -tags e2e ./e2e/api/

# e.g.
go test -timeout 30s -run ^TestFlow$ -tags e2e ./e2e/api/
```


# Benchmark results

## Get Media

### #1
* Count - 1
* Without Cache
* file of 100MB
* upload once - chunk size 1MB
* Get Media - chunk size 1MB

```bash
go test -bench BenchmarkGetMedia -run '^$' ./e2e/api/ -count 1
goos: linux
goarch: amd64
pkg: github.com/rishabhkailey/media-service/e2e/api
cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
BenchmarkGetMedia-8            1        18299092271 ns/op
PASS
ok      github.com/rishabhkailey/media-service/e2e/api  18.310s
```

### #2
* Count - 5
* Without Cache
* file of 100MB
* upload once - chunk size 1MB
* Get Media - chunk size 1MB

```bash
go test -bench BenchmarkGetMedia -run '^$' ./e2e/api/ -count 5
--- FAIL: BenchmarkGetMedia
    media_flow_test.go:273: 
                Error Trace:    /workspaces/media-service/e2e/api/media_flow_test.go:273
                                                        /workspaces/media-service/e2e/api/media_flow_test.go:319
                Error:          Not equal: 
                                expected: 206
                                actual  : 500
                Test:           BenchmarkGetMedia
                Messages:       GetMediaRangeRequest status
FAIL
exit status 1
FAIL    github.com/rishabhkailey/media-service/e2e/api  102.777s
FAIL
```

### #3
> Huge Improvement :)
* Count - 5
* ARC Cache
* file of 100MB
* upload once - chunk size 1MB
* Get Media - chunk size 1MB

```bash
go test -bench BenchmarkGetMedia -run '^$' ./e2e/api/ -count 5
goos: linux
goarch: amd64
pkg: github.com/rishabhkailey/media-service/e2e/api
cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
BenchmarkGetMedia-8            1        2719891829 ns/op
BenchmarkGetMedia-8            1        2664753867 ns/op
BenchmarkGetMedia-8            1        2463087424 ns/op
BenchmarkGetMedia-8            1        2751681536 ns/op
BenchmarkGetMedia-8            1        2636278034 ns/op
PASS
ok      github.com/rishabhkailey/media-service/e2e/api  13.245s
```