package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/tjfoc/gmsm/sm4"
)

// DeriveNodeWorkingKey returns SHA-256(UTF-8(apiKey))[:16] — matches Java NodeSm4Codec.
func DeriveNodeWorkingKey(apiKey string) []byte {
	sum := sha256.Sum256([]byte(apiKey))
	key := make([]byte, 16)
	copy(key, sum[:16])
	return key
}

// EncryptEnvelope returns standard base64 iv/payload for HTTP JSON bodies.
func EncryptEnvelope(plaintext []byte, key []byte) (ivB64, payloadB64 string, err error) {
	iv := make([]byte, sm4.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", "", err
	}
	out, err := sm4CBCEncrypt(key, iv, plaintext)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(iv), base64.StdEncoding.EncodeToString(out), nil
}

// DecryptEnvelope decrypts standard base64 iv/payload body envelope.
func DecryptEnvelope(ivB64, payloadB64 string, key []byte) ([]byte, error) {
	iv, err := base64.StdEncoding.DecodeString(ivB64)
	if err != nil {
		return nil, err
	}
	payload, err := base64.StdEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, err
	}
	return sm4CBCDecrypt(key, iv, payload)
}

// EncryptCompact returns base64url(iv).base64url(ciphertext) for query param e.
func EncryptCompact(plaintext []byte, key []byte) (string, error) {
	iv := make([]byte, sm4.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	out, err := sm4CBCEncrypt(key, iv, plaintext)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(iv) + "." + base64.RawURLEncoding.EncodeToString(out), nil
}

// DecryptCompact decrypts base64url(iv).base64url(ciphertext).
func DecryptCompact(compact string, key []byte) ([]byte, error) {
	dot := strings.IndexByte(compact, '.')
	if dot <= 0 || dot >= len(compact)-1 {
		return nil, fmt.Errorf("invalid compact ciphertext")
	}
	iv, err := base64.RawURLEncoding.DecodeString(compact[:dot])
	if err != nil {
		return nil, err
	}
	payload, err := base64.RawURLEncoding.DecodeString(compact[dot+1:])
	if err != nil {
		return nil, err
	}
	return sm4CBCDecrypt(key, iv, payload)
}

func sm4CBCEncrypt(key, iv, plaintext []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != sm4.BlockSize {
		return nil, fmt.Errorf("iv must be %d bytes", sm4.BlockSize)
	}
	padded := pkcs7Pad(plaintext, sm4.BlockSize)
	out := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, padded)
	return out, nil
}

func sm4CBCDecrypt(key, iv, ciphertext []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != sm4.BlockSize {
		return nil, fmt.Errorf("iv must be %d bytes", sm4.BlockSize)
	}
	if len(ciphertext) == 0 || len(ciphertext)%sm4.BlockSize != 0 {
		return nil, fmt.Errorf("invalid ciphertext length")
	}
	out := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, ciphertext)
	return pkcs7Unpad(out, sm4.BlockSize)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid pkcs7 data")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > blockSize || pad > len(data) {
		return nil, fmt.Errorf("invalid pkcs7 padding")
	}
	for i := 0; i < pad; i++ {
		if data[len(data)-1-i] != byte(pad) {
			return nil, fmt.Errorf("invalid pkcs7 padding")
		}
	}
	return data[:len(data)-pad], nil
}
