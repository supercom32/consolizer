package scrambler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

/*
Encrypt is a method which encrypts data using AES-GCM.

Example:

	encrypted, err := Encrypt([]byte("secret"), key)
*/
func Encrypt(dataToEncrypt, key []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	encryptedData := gcm.Seal(nonce, nonce, dataToEncrypt, nil)
	return encryptedData, nil
}

/*
Decrypt is a method which decrypts data that was encrypted using AES-GCM.

Example:

	decrypted, err := Decrypt(encryptedData, key)
*/
func Decrypt(dataToDecrypt, key []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce, encryptedData := dataToDecrypt[:gcm.NonceSize()], dataToDecrypt[gcm.NonceSize():]
	unencryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}
	return unencryptedData, nil
}

/*
GenerateKey is a method which generates a random 32-byte key suitable for AES-256 encryption.

Example:

	key, err := GenerateKey()
*/
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

/*
func main() {
	data := []byte("our super secret text")
	key, err := GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	ciphertext, err := Encrypt(key, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ciphertext: %s\n", hex.EncodeToString(ciphertext))
	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("plaintext: %s\n", plaintext)
}
*/
