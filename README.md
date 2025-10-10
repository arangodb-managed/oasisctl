# oasisctl

Commandline utility for accessing the **Arango Managed Platform (AMP)**,
formerly called Oasis and ArangoGraph Insights Platform.

This utility is used to access AMP resources (such as projects, deployments, certificates) without the needs for a graphical dashboard.
It is often used in scripts.

## Maintainers

This utility is maintained by the team at [Arango.ai](https://arango.ai/).

## Installation

Downloading the [latest released binaries](https://github.com/arangodb-managed/oasisctl/releases),
extract the zip archive and install the binary for your platform in your preferred location.

Or to build from source, run:

```bash
git clone https://github.com/arangodb-managed/oasisctl.git
make
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

oasisctl uses an authentication token to authenticate with the Arango Managed Platform.

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

More information and a getting started guide about oasisctl is available at <https://docs.arango.ai/amp/>.
