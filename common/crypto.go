package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"

	"bytes"
	"encoding/base64"
	"fmt"
	"io"
)

func Md5String(s string) []byte {
	h := md5.New()
	io.WriteString(h, s)
	return h.Sum(nil)
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) (ret string, err error) {
	var buf []byte
	buf, err = base64.StdEncoding.DecodeString(s)
	ret = string(buf)
	return
}

/* https://dablelv.blog.csdn.net/article/details/89387460 */

// PKCS7Padding fills plaintext as an integral multiple of the block length
func PKCS7Padding(p []byte, blockSize int) []byte {
	pad := blockSize - len(p)%blockSize
	padtext := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(p, padtext...)
}

// PKCS7UnPadding removes padding data from the tail of plaintext
func PKCS7UnPadding(p []byte) []byte {
	length := len(p)
	paddLen := int(p[length-1])
	return p[:(length - paddLen)]
}

// AESCBCEncrypt encrypts data with AES algorithm in CBC mode
// Note that key length must be 16, 24 or 32 bytes to select AES-128, AES-192, or AES-256
// Note that AES block size is 16 bytes
func AesCBCEncrypt(p, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	p = PKCS7Padding(p, block.BlockSize())
	ciphertext := make([]byte, len(p))
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(ciphertext, p)
	return ciphertext, nil
}

// AESCBCDecrypt decrypts cipher text with AES algorithm in CBC mode
// Note that key length must be 16, 24 or 32 bytes to select AES-128, AES-192, or AES-256
// Note that AES block size is 16 bytes
func AesCBCDecrypt(c, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(c))
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(plaintext, c)
	return PKCS7UnPadding(plaintext), err
}

func Sha1Encode(s string) []byte {
	h := sha1.New()
	io.WriteString(h, s)
	return h.Sum(nil)
}

/* digital signature */
func SignEncode(salt, mac, time string) string {
	fmt.Println("======== sign encode =====")
	fmt.Printf("salt str: %s\n", salt)
	fmt.Printf("mac str: %s\n", mac)
	fmt.Printf("time str: %s\n", time)
	fmt.Println("-----------------------------")

	/* format */
	format := fmt.Sprintf("%s+%s+%s", mac, salt, time)
	fmt.Printf("format src str: %s\n", format)

	/* md5 */
	md5data := fmt.Sprintf("%X", Md5String(format))
	fmt.Printf("md5 hex: %s\n", md5data)

	/* hash */
	hdata := fmt.Sprintf("%X", Sha1Encode(md5data))
	fmt.Printf("sha-1 encode hex: %s\n", hdata)

	/* base64 encode */
	b64data := Base64Encode(hdata)
	fmt.Printf("base64 encode str: %s\n", b64data)
	fmt.Printf("*sign str: %s\n", b64data)
	fmt.Println("===========================")

	return b64data
}
