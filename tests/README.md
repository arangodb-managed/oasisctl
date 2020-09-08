# Writing new tests

These tests should be ignored by regular `go test ./...`. This is achieved by a [build constraints](https://golang.org/cmd/go/#hdr-Build_constraints).

The test must have a `// +build e2e` tag at the top to make sure they aren't executed via CI or a regular
`go test ./...` run.

# Running the tests

**NOTE**: Make sure the oasisctl binary exists before running the tests.

To run all tests execute the following:

```bash
go test ./... -tags=e2c
```

To run a specific test, `cd` into the folder you wish to run tests for and execute for example:

```bash
go test -run=TestCreateCertificate -tags=e2c
```

# Constraints

These tests were designed to be run with either a user which has a single organization
and project, or defaults being set properly via environment properties:
`OASIS_ORGANIZATION` and `OASIS_PROJECT`.
