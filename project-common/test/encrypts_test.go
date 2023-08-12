package test

import (
	"fmt"
	"github.com/spxzx/project-common/encrypts"
	"testing"
)

func TestEncrypt(t *testing.T) {
	plain := "39"
	// AES 规定有3种长度的key: 16, 24, 32分别对应AES-128, AES-192, or AES-256
	key := "suhvcthjkpdcvgtrfujnhgzw"
	// 加密
	cipherByte, err := encrypts.Encrypt(plain, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %s\n", plain, cipherByte)
	// 解密
	plainText, err := encrypts.Decrypt(cipherByte, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %s\n", cipherByte, plainText)
}
