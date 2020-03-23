# ArangoDB Oasis

![ArangoDB Oasis](https://cloud.arangodb.com/assets/logos/arangodb-oasis-logo-whitebg-right.png)

Commandline utility for accessing ArangoDB Oasis.

This utility is used to access ArangoDB Oasis resources (such as projects, deployments, certificates) without the needs for a graphical dashboard.
It is often used in scripts.

## Maintainers

This utility is maintained by the team at [ArangoDB](https://www.arangodb.com/).

## Installation

Run:

```bash
go get github.com/arangodb-managed/oasisctl
```

## Usage

```bash
oasisctl [command...]
```

A list of commands and options can be shown using:

```bash
oasisctl -h
```

## Authentication

Oasisctl uses an authentication token to authenticate with the ArangoDB Oasis platform.

To get such an authentication token, create an [API key](https://cloud.arangodb.com/dashboard/user/api-keys) and run:

```bash
oasisctl login \
  --key-id=<your-key-id> \
  --key-secret=<your-key-secret>
```

The output is the authentication token.
The most convenient use is to put it in an `OASIS_TOKEN` environment variable.
In a script that can be done like this:

```bash
export OASIS_TOKEN=$(oasisctl login --key-id=<your-key-id> --key-secret=<your-key-secret>)
```

## More information

More information and a getting started guide about Oasisctl is available at [arangodb.com/docs/stable/oasis](https://www.arangodb.com/docs/stable/oasis/).
