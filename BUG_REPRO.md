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
?   	wedding-sign/cmd/wedding-sign	[no test files]
?   	wedding-sign/internal/catalog	[no test files]
?   	wedding-sign/internal/config	[no test files]
ok  	wedding-sign/internal/audit	0.016s
ok  	wedding-sign/internal/display	0.003s
ok  	wedding-sign/internal/model	0.002s
--- FAIL: TestWorkflow24 (0.01s)
    workflow24_test.go:37: dispatch should be idempotent, got 2 audit entries
FAIL
FAIL	wedding-sign/internal/service	0.029s
ok  	wedding-sign/internal/store	0.020s
ok  	wedding-sign/internal/workflow	0.033s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/wedding-sign): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/wedding-sign): exit `0`
- Frontend build (web): exit `0`
