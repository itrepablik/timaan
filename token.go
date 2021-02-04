package timaan

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// TP type stands for 'Token Payload' that uses the map[string]interface{}
// as the custom token payload structure
type TP map[string]interface{}

// TK is a successful auth token requests after successful login
// the TokenKey might be anything, e.g username, random unique alpha-numeric strings, etc.
type TK struct {
	TokenKey string
	Payload  TP
	ExpireOn int64
}

// UserTokens is a users auth token requests stored in memory
type UserTokens struct {
	Token map[string][]byte
	mu    sync.Mutex
}

// UT is a user's token methods
var UT = UserTokens{Token: make(map[string][]byte)}

// GenerateToken generate a new timaan token which can be used mostly after successful authentication
// process, mostly use after successful login, this can also be used as you sessions for the entire
// duration of the token validity period.
func GenerateToken(tokenKey string, payLoad TK) ([]byte, error) {
	if len(strings.TrimSpace(tokenKey)) == 0 {
		return []byte{}, errors.New("token key is required")
	}
	encBytes, err := EncodePayload(payLoad)
	if err != nil {
		return []byte{}, err
	}
	_, isTokFound := UT.Get(tokenKey)
	if isTokFound {
		UT.Remove(tokenKey)
	}
	UT.Add(tokenKey, encBytes)
	return encBytes, nil
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

// DecodePayload extracts the token payload
func DecodePayload(tokenKey string) (TK, error) {
	if len(strings.TrimSpace(tokenKey)) == 0 {
		return TK{}, errors.New("token key is required")
	}
	var payLoad = TK{}
	tokBytes, _ := UT.Get(tokenKey)
	dec := gob.NewDecoder(bytes.NewReader(tokBytes))

	err := dec.Decode(&payLoad)
	if err != nil {
		return TK{}, errors.New("token key not found: " + tokenKey)
	}
	return payLoad, nil
}

// Add insert the new token request to the 'UserTokens' map
func (t *UserTokens) Add(tokenKey string, encBytes []byte) {
	if len(strings.TrimSpace(tokenKey)) > 0 {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.Token[tokenKey] = encBytes
	}
}

// Get gets the specific user's token filtered by the tokenKey
func (t *UserTokens) Get(tokenKey string) ([]byte, bool) {
	tok, ok := t.Token[tokenKey]
	if !ok {
		return []byte{}, ok
	}
	return tok, ok
}

// Remove any single stored token from the 'UserTokens' map
func (t *UserTokens) Remove(tokenKey string) (bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.Token[tokenKey]
	if !ok {
		return false, errors.New("token key not found: " + tokenKey)
	}
	delete(t.Token, tokenKey)
	return true, nil
}

// RandomToken uses the lower-case letters from 'a' to 'z' and positive numeric digits combination
func RandomToken() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
