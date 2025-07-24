package did

import (
	"github.com/google/uuid"
)

// Generate creates a simple DID using the example method.
func Generate() string {
	return "did:example:" + uuid.NewString()
}
