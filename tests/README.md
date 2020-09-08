# Testing oasisctl

**NOTE**: Make sure the oasisctl binary exists before running the tests.

To run all tests execute the following:

```bash
go test ./...
```

To run a specific test, `cd` into the folder you wish to run tests for and execute for example:

```bash
go test -run=TestCreateCertificate
```

# Constraints

These tests were designed to be run with either a user which has a single organization
and project, or defaults being set properly via environment properties:
`OASIS_ORGANIZATION` and `OASIS_PROJECT`.
