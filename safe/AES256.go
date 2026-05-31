package safe

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"misakadb/clilog"
	"misakadb/network/RegisterCenter"
	"os"
)

// EncryptAES 军工级加密：把明文(JSON字节)变成密文乱码
func EncryptAES(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("密钥损坏！")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成一个随机的 Nonce (随机数，用于保证即使相同明文每次加密结果也不同)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据，并把 nonce 拼在密文的最前面（解密时需要用到）
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES 军工级解密：把密文乱码还原成明文(JSON字节)
func DecryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("密钥损坏！")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	// 如果文件内容还没 nonce 长，说明文件坏了
	if len(ciphertext) < nonceSize {
		return nil, errors.New("密文太短，数据已损坏")
	}

	// 拆分出 nonce 和真正的密文
	nonce, cipherData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// 初始化密钥 到 master.key
func InitPassword() {
	clilog.Info("初始化密钥中，profiles/master.mikey极其重要切勿泄漏")
	newPassword := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newPassword); err != nil {
		clilog.Error("安全随机密钥生成失败:", err)
		panic("error!")
	}

	err := os.WriteFile("./profiles/master.mikey", newPassword, 0400)
	if err != nil {
		clilog.Error("无法写入文件，请检查根目录中是否存在profiles文件夹。")
		panic("error!")
	}
	clilog.Success("初始化密钥成功！")
}

func EncryptByte(plaintext []byte) ([]byte, error) {
	rc := RegisterCenter.NewRegisterCenter()
	key := []byte(rc.MasterKey)
	return EncryptAES(key, plaintext)
}

func DecryptByte(ciphertext []byte) ([]byte, error) {
	rc := RegisterCenter.NewRegisterCenter()
	key := []byte(rc.MasterKey)

	return DecryptAES(key, ciphertext)
}
