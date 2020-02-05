package kobodict

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// Crypter represents a symmetric dictionary encryption method.
type Crypter interface {
	Encrypter
	Decrypter
}

// CryptMethodAES represents AES-128-ECB encryption with PKCS#7 padding.
const CryptMethodAES string = "aes"

// NewCrypter creates the specified type of Crypter with the specified key.
func NewCrypter(method string, key []byte) (Crypter, error) {
	switch method {
	case CryptMethodAES:
		c, err := newCryptAES(key)
		return c, err
	default:
		return nil, fmt.Errorf("unknown encryption method %#v", method)
	}
}

type cryptAES struct {
	b cipher.Block
}

func newCryptAES(key []byte) (*cryptAES, error) {
	if b, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		return &cryptAES{b}, nil
	}
}

// Encrypt implements Encrypter.
func (c *cryptAES) Encrypt(buf []byte) ([]byte, error) {
	if dst, err := cryptPKCS7Pad(buf, aes.BlockSize); err != nil {
		return nil, err
	} else if dst, err = cryptAES128ECBEncrypt(c.b, dst); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

// Decrypt implements Decrypter.
func (c *cryptAES) Decrypt(buf []byte) ([]byte, error) {
	if dst, err := cryptAES128ECBDecrypt(c.b, buf); err != nil {
		return nil, err
	} else if dst, err := cryptPKCS7Unpad(dst, aes.BlockSize); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

func cryptPKCS7Unpad(src []byte, blockSize int) ([]byte, error) {
	if blockSize > 0xFF || blockSize < 0x00 {
		return nil, fmt.Errorf("block size %d out of bounds", blockSize)
	} else if len(src)%blockSize != 0 || len(src) == 0 {
		return nil, fmt.Errorf("data length %d is empty or not a multiple of block size %d", len(src), blockSize)
	}
	plen := int(src[len(src)-1])
	if len(src) <= plen {
		return nil, fmt.Errorf("invalid padding: padding length %d out of bounds", plen)
	}
	for _, v := range src[len(src)-plen:] {
		if int(v) != plen {
			return nil, fmt.Errorf("invalid padding: expected %d, got %d", plen, v)
		}
	}
	return src[:len(src)-plen], nil
}

func cryptPKCS7Pad(src []byte, blockSize int) ([]byte, error) {
	if blockSize > 0xFF || blockSize < 0x00 {
		return nil, fmt.Errorf("block size %d out of bounds", blockSize)
	}
	plen := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(plen)}, plen)...), nil
}

func cryptAES128ECBDecrypt(cb cipher.Block, src []byte) ([]byte, error) {
	if len(src)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("src not a multiple of block size %d", aes.BlockSize)
	}
	dst := make([]byte, len(src))
	for i := aes.BlockSize; i <= len(src); i += aes.BlockSize {
		cb.Decrypt(dst[i-aes.BlockSize:i], src[i-aes.BlockSize:i])
	}
	return dst, nil
}

func cryptAES128ECBEncrypt(cb cipher.Block, src []byte) ([]byte, error) {
	if len(src)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("src not a multiple of block size %d", aes.BlockSize)
	}
	dst := make([]byte, len(src))
	for i := aes.BlockSize; i <= len(src); i += aes.BlockSize {
		cb.Encrypt(dst[i-aes.BlockSize:i], src[i-aes.BlockSize:i])
	}
	return dst, nil
}
