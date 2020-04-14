package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type JWT struct{}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) ClaimSet(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid-format %q", token)
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decoding-token-failed %w", err)
	}

	var ret map[string]interface{}
	err = json.Unmarshal(decoded, &ret)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling-token-failed %w", err)
	}

	return ret, nil
}
