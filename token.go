package timaan

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
	"sync"

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
	Token map[string][]byte
	mu    sync.Mutex
}

// UT is a user's token methods
var UT = UserTokens{Token: make(map[string][]byte)}

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
func ExtractPayload(userName string) (TK, error) {
	if len(strings.TrimSpace(userName)) == 0 {
		return TK{}, errors.New("username is required")
	}
	var payLoad = TK{}
	tokBytes := UT.Get(userName)
	dec := gob.NewDecoder(bytes.NewReader(tokBytes))

	err := dec.Decode(&payLoad)
	if err != nil {
		return TK{}, errors.New("token not found for: " + userName)
	}
	return payLoad, nil
}

// Add insert the new token request to the 'UserTokens' map
func (t *UserTokens) Add(userName string, encBytes []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Token[userName] = encBytes
}

// Get gets the specific user's token filtered by the username
func (t *UserTokens) Get(userName string) []byte {
	tok, ok := t.Token[userName]
	if !ok {
		return []byte{}
	}
	return tok
}

// Remove any single stored token from the 'UserTokens' map
func (t *UserTokens) Remove(userName string) (bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.Token[userName]
	if !ok {
		return false, errors.New("username not found")
	}
	delete(t.Token, userName)
	return true, nil
}
