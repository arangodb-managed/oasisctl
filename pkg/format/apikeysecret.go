//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	iam "github.com/arangodb-managed/apis/iam/v1"
)

// APIKeySecret returns a single api key secret formatted for humans.
func APIKeySecret(x *iam.APIKeySecret, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"secret", x.GetSecret()},
	)
}
