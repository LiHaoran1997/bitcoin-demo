package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/decred/dcrd/crypto/ripemd160"
	"log"
)

//这里的钱包是一个结构，每个钱包保存了公钥和私钥对
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	//PubKey *ecdsa.PublicKey
	//这里PubKey不存储原始公钥，存储X和Y拼接字符串，在校验端重新拆分（参考rs）
	PublicKey []byte
}

//创建钱包
func NewWallet() *Wallet {
	//创建曲线
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}

	//生成公钥
	publicKeyOrigin := privateKey.PublicKey
	publicKey := append(publicKeyOrigin.X.Bytes(), publicKeyOrigin.Y.Bytes()...)
	return &Wallet{privateKey, publicKey}
}

//生成地址
func (w *Wallet) NewAddress() string {
	pubKey := w.PublicKey
	rip160HashValue := HashPubKey(pubKey)
	//版本
	version := byte(00)
	//拼接version
	payload := append([]byte{version}, rip160HashValue...)
	//checkSum
	checkCode := checkSum(payload)
	payload = append(payload, checkCode...)
	//golang有base58
	address := base58.Encode(payload)
	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)
	//编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	//返回rip160的hash结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

func checkSum(payload []byte) []byte {
	//两次sha256
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	//前4字节校验码
	checkCode := hash2[:4]
	return checkCode
}

//通过地址返回公钥哈希
func GetPubKeyFromAddress(address string) []byte {
	//1.解码
	addressByte := base58.Decode(address)
	//2.截取出公钥哈希，去除version，去除校验码（4字节）
	hashLen := len(addressByte)

	pubKeyHash := addressByte[1:hashLen-4]
	return pubKeyHash
}

func IsValidAddress(address string)bool{
	//1.解码
	addressByte := base58.Decode(address)
	//2.截取出公钥哈希，去除version，去除校验码（4字节）
	payload := addressByte[:len(addressByte)-4]
	checkSum1:=addressByte[len(addressByte)-4:]
	checkSum2:=checkSum(payload)
	//fmt.Println("checkSum1 : %s\n",checkSum1)
	//fmt.Println("checkSum2 : %s\n",checkSum2)
	return bytes.Equal(checkSum1,checkSum2)

}