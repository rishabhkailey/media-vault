```bash
go test -timeout 30s -run ^<test-name>$ -tags e2e ./e2e/api/

# e.g.
go test -timeout 30s -run ^TestFlow$ -tags e2e ./e2e/api/
```