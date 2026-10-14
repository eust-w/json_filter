# Release validation

Date: 2026-10-14. Checks are executed in an isolated working checkout. Screenshots use local synthetic/demo data. A source implementation, a passing automated test and external-service acceptance are separate claims.

## Commands

```bash
go test -race ./...
go build .
```

PASS: Go race tests and build; HTTP valid/invalid uploads, exact filtering, empty-result download and duplicate-row regression. Hardware, external credentials, signing and production deployment are separate setup/acceptance steps.
