package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/adler32"
)

func Md5(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func HmacSHA1(keyStr, value string) string {
	mac := hmac.New(sha1.New, []byte(keyStr))
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

func Hash(data string) uint32 {
	Sha1Inst := adler32.New()
	Sha1Inst.Write([]byte(data))
	return Sha1Inst.Sum32()
}

// AesDecrypt AES解密
func AesDecrypt(str string, key []byte) ([]byte, error) {

	ciphertext, _ := base64.StdEncoding.DecodeString(str)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	if len(ciphertext) < aes.BlockSize {
		fmt.Println("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		fmt.Println("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks可以原地更新
	mode.CryptBlocks(ciphertext, ciphertext)

	fmt.Printf("%s\n", ciphertext)

	return ciphertext, nil
}
