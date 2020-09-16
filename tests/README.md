# Writing new tests

These tests should be ignored by regular `go test ./...`. This is achieved by a [build constraints](https://golang.org/cmd/go/#hdr-Build_constraints).

The test must have a `// +build e2e` tag at the top to make sure they aren't executed via CI or a regular
`go test ./...` run.

The command line options must contain `--organization-id=1234` because otherwise the parallel nature
of these tests might step on each other. You could get an error like, you are a member of 2 or more organizations...

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

The following assumptions are based on the system:

- There is a user with an organization and a project.
- The project has a Default Certificate. The certificate will not be deleted.

The tests are designed so that they can be run in parallel which is the default of `go test`.
