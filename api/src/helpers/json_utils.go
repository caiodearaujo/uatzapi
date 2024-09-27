package helpers

import (
	"encoding/json"
	"whatsgoingon/data"
)

func MarshalMessageToJSON(content data.StoredMessage) ([]byte, error) {
	return json.Marshal(content)
}