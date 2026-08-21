# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
--- FAIL: TestBusiness24Regression (0.01s)
    TestBusiness24Regression_test.go:31: batch ZX89024 reread returned incorrect state: []*domain.Record{(*domain.Record)(0x40000a8780)}
FAIL
FAIL	memorialstation	0.024s
?   	memorialstation/api	[no test files]
?   	memorialstation/archive	[no test files]
?   	memorialstation/cmd/memorial	[no test files]
?   	memorialstation/domain	[no test files]
?   	memorialstation/importx	[no test files]
?   	memorialstation/review	[no test files]
?   	memorialstation/search	[no test files]
?   	memorialstation/storage	[no test files]
?   	memorialstation/workflow	[no test files]
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/memorial): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/memorial): exit `0`
