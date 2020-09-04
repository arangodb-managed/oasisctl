# Testing oasisctl

To run the tests execute the following from oasisctl root:

```bash
make e2e-test
```

To run a specific test run from root:

```bash
# Make sure that the oasisctl binary exists
make
# Run the test
go test tests/crypto/crypto_create_test.go
```
