package wallet

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/encrypt"
)

// secrets key name
const (
	secretSeed     = "seed"
	secretLastSeed = "lastSeed"
)

type cryptor interface {
	Encrypt(data, password []byte) ([]byte, error)
	Decrypt(data, password []byte) ([]byte, error)
}

// CryptoType represents the type of crypto name
type CryptoType string

// StrToCryptoType converts string to CryptoType
func StrToCryptoType(s string) (CryptoType, error) {
	switch CryptoType(s) {
	case CryptoTypeSha256Xor:
		return CryptoTypeSha256Xor, nil
	case CryptoTypeScryptChacha20poly1305:
		return CryptoTypeScryptChacha20poly1305, nil
	default:
		return "", errors.New("unknow crypto type")
	}
}

// Crypto types
const (
	CryptoTypeSha256Xor              = CryptoType("sha256-xor")
	CryptoTypeScryptChacha20poly1305 = CryptoType("scrypt-chacha20poly1305")
)

// cryptoTable records all supported wallet crypto methods
// If want to support new crypto methods, register here.
var cryptoTable = map[CryptoType]cryptor{
	CryptoTypeSha256Xor:              encrypt.DefaultSha256Xor,
	CryptoTypeScryptChacha20poly1305: encrypt.DefaultScryptChacha20poly1305,
}

// ErrAuthenticationFailed wraps the error of decryption.
type ErrAuthenticationFailed struct {
	err error
}

func (e ErrAuthenticationFailed) Error() string {
	return e.err.Error()
}

// getCrypto gets crypto of given type
func getCrypto(cryptoType CryptoType) (cryptor, error) {
	c, ok := cryptoTable[cryptoType]
	if !ok {
		return nil, fmt.Errorf("can not find crypto %v in crypto table", cryptoType)
	}

	return c, nil
}

type secrets map[string]string

func (s secrets) get(key string) (string, bool) {
	v, ok := s[key]
	return v, ok
}

func (s secrets) set(key, v string) {
	s[key] = v
}

func (s secrets) serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s secrets) deserialize(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s secrets) erase() {
	for k := range s {
		s[k] = ""
		delete(s, k)
	}
}
