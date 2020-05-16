package timaan

import (
	"bytes"
	"encoding/gob"

	"github.com/itrepablik/tago"
)

// TP type stands for 'Token Payload' that uses the map[string]interface{}
// as the custom token payload structure
type TP map[string]interface{}

// TK is a token collections storage
type TK struct {
	UserName string
	Payload  TP
	ExpireOn int64
}

// UserTokens is a users token requests stored in memory storage
type UserTokens struct {
	Token map[interface{}][]byte
}

// UT is a user's token methods
var UT = UserTokens{}

// GenerateToken generate a new timaan token
func GenerateToken(userName, secretKey string, payLoad TK) (string, error) {
	newToken, err := tago.Encrypt(userName, secretKey)
	if err != nil {
		return "", err
	}
	encBytes, err := EncodePayload(payLoad)
	if err != nil {
		return "", err
	}
	UT.Add(userName, encBytes) // Add new requested token to the map
	return newToken, nil
}

// EncodePayload encodes the token payload using gob
func EncodePayload(payLoad TK) ([]byte, error) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(payLoad)
	if err != nil {
		return data.Bytes(), err
	}
	return data.Bytes(), nil
}

// ExtractPayload decodes the token payload
func ExtractPayload(enc []byte, payLoad *TK) error {
	dec := gob.NewDecoder(bytes.NewReader(enc))
	err := dec.Decode(&payLoad)
	if err != nil {
		return err
	}
	return nil
}

// Add insert the new token request to the 'UserTokens' map
func (t *UserTokens) Add(userName string, encBytes []byte) {
	t.Token[userName] = encBytes
}

// Remove any single stored token from the 'UserTokens' map
func (t *UserTokens) Remove(userName string, encBytes []byte) {
	delete(t.Token, userName)
}
