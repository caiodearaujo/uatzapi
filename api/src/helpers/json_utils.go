package helpers

import (
	"encoding/json"
	"whatsgoingon/data"
)

// MarshalMessageToJSON converts a StoredMessage struct into a JSON byte array.
// Returns the JSON-encoded bytes or an error if marshalling fails.
func MarshalMessageToJSON(content data.StoredMessage) ([]byte, error) {
	return json.Marshal(content)
}
