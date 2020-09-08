# Testing oasisctl

**NOTE**: Make sure that the oasisctl binary exists before running the tests.

To run all tests execute the following:

```bash
go test ./...
```

To run a specific test:

```bash
go test crypto/crypto_create_test.go -run=TestCreateCertificate
```

# Constraints

These tests were designed to be run with either a user which has a single organization
and project, or defaults being set properly via environment properties:
`OASIS_ORGANIZATION` and `OASIS_PROJECT`.
