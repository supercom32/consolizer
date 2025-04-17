package scrambler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

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