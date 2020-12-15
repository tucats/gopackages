package functions

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

type Token struct {
	Name    string
	Data    string
	TokenID uuid.UUID
	Expires time.Time
	AuthID  uuid.UUID
}

// Hash implements the _cipher.hash() function
func Hash(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return util.Hash(util.GetString(args[0])), nil
}

// Encrypt implements the _cipher.hash() function
func Encrypt(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	b, err := util.Encrypt(util.GetString(args[0]), util.GetString(args[1]))
	if err != nil {
		return b, err
	}
	return hex.EncodeToString([]byte(b)), nil

}

// Decrypt implements the _cipher.hash() function
func Decrypt(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	b, err := hex.DecodeString(util.GetString(args[0]))
	if err != nil {
		return nil, err
	}
	return util.Decrypt(string(b), util.GetString(args[1]))
}

// Validate creates a new token with a username and a data payload
func Validate(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	var err error

	// Take the token value, and de-hexify it.
	b, err := hex.DecodeString(util.GetString(args[0]))
	if err != nil {
		return false, nil
	}

	// Decrypt the token into a json string
	key := persistence.Get("token-key")
	j, err := util.Decrypt(string(b), key)
	if err != nil {
		return false, nil
	}
	var t = Token{}
	err = json.Unmarshal([]byte(j), &t)
	if err != nil {
		return false, nil
	}

	// Has the expiration passed?
	d := time.Since(t.Expires)
	if d.Seconds() > 0 {
		return false, nil
	}

	return true, nil
}

// Validate creates a new token with a username and a data payload
func Extract(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	var err error

	// Take the token value, and de-hexify it.
	b, err := hex.DecodeString(util.GetString(args[0]))
	if err != nil {
		return nil, err
	}

	// Decrypt the token into a json string
	key := persistence.Get("token-key")
	j, err := util.Decrypt(string(b), key)
	if err != nil {
		return nil, err
	}
	var t = Token{}
	err = json.Unmarshal([]byte(j), &t)
	if err != nil {
		return nil, err
	}

	// Has the expiration passed?
	d := time.Since(t.Expires)
	if d.Seconds() > 0 {
		return nil, errors.New("token expired")
	}

	r := map[string]interface{}{}
	r["name"] = t.Name
	r["data"] = t.Data
	r["session"] = t.AuthID.String()
	r["id"] = t.TokenID.String()

	return r, nil
}

// CreateToken creates a new token with a username and a data payload
func CreateToken(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	var err error

	// Create a new token object, with the username and an ID. If there was a
	// data payload as well, add that to the token.
	t := Token{
		Name:    util.GetString(args[0]),
		TokenID: uuid.New(),
	}
	if len(args) == 2 {
		t.Data = util.GetString(args[1])
	}

	// Get the session ID of the current Ego program and add it to
	// the token. A token can only be validated on the same system
	// that created it.
	if session, ok := s.Get("_session"); ok {
		t.AuthID, err = uuid.Parse(util.GetString(session))
		if err != nil {
			return nil, err
		}
	}

	// Fetch the default interval, or use 15 minutes as the default.
	// Calculate a time value for when this token expires
	interval := persistence.Get("token-expiration")
	if interval == "" {
		interval = "15m"
	}
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	t.Expires = time.Now().Add(duration)

	// Make the token into a json string
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	// Encrypt the string value
	key := persistence.Get("token-key")
	encryptedString, err := util.Encrypt(string(b), key)
	if err != nil {
		return b, err
	}
	return hex.EncodeToString([]byte(encryptedString)), nil
}
