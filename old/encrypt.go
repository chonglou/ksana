package ksana

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"log"
)

func Md5(src []byte) [16]byte {
	return md5.Sum(src)
}

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

type Hmac struct {
	key []byte
}

func (h *Hmac) Sum(src []byte) []byte {
	mac := hmac.New(sha512.New, h.key)
	mac.Write(src)
	return mac.Sum(nil)
}

func (h *Hmac) Equal(src, dst []byte) bool {
	return hmac.Equal(src, dst)
}

type Aes struct {
	cip cipher.Block
}

//16、24或者32位的[]byte，分别对应AES-128, AES-192或AES-256算法
func (a *Aes) Init(key []byte) {
	c, e := aes.NewCipher(key)
	if e != nil {
		log.Fatalf("Error on new aes cipher: %v", e)
	}
	a.cip = c
}

func (a *Aes) Encrypt(src []byte) ([]byte, []byte) {
	iv := RandomBytes(aes.BlockSize)
	cfb := cipher.NewCFBEncrypter(a.cip, iv)
	ct := make([]byte, len(src))
	cfb.XORKeyStream(ct, src)
	return ct, iv

}

func (a *Aes) Decrypt(src, iv []byte) []byte {
	cfb := cipher.NewCFBDecrypter(a.cip, iv)
	pt := make([]byte, len(src))
	cfb.XORKeyStream(pt, src)
	return pt
}
