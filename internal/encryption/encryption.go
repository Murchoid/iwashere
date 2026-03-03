package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// GenerateKey creates a new random key for a note
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256
	_, err := rand.Read(key)
	return key, err
}

// Encrypt encrypts data with the given key
func Encrypt(data []byte, key []byte) (string, string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", err
	}

	// Encrypt
	ciphertext := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, data)

	// Return as base64 strings
	return base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(iv),
		nil
}

// Decrypt decrypts data with the given key
func Decrypt(ciphertextB64 string, ivB64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(ivB64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// EncryptKeyForRecipient simulates asymmetric encryption
// For v0.2.0, we'll use a simple hash of recipient's email as "public key"
// In production, you'd use real public key cryptography
func EncryptKeyForRecipient(key []byte, recipient string) (string, error) {
	// This is a SIMPLIFIED approach for v0.2.0
	// In reality, you'd use the recipient's public key
	recipientKey := sha256.Sum256([]byte(recipient))

	encrypted := make([]byte, len(key))
	for i := range key {
		encrypted[i] = key[i] ^ recipientKey[i%len(recipientKey)]
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptKeyForRecipient reverses the simple XOR
func DecryptKeyForRecipient(encryptedKeyB64 string, recipient string) ([]byte, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(encryptedKeyB64)
	if err != nil {
		return nil, err
	}

	recipientKey := sha256.Sum256([]byte(recipient))

	key := make([]byte, len(encryptedKey))
	for i := range encryptedKey {
		key[i] = encryptedKey[i] ^ recipientKey[i%len(recipientKey)]
	}

	return key, nil
}
