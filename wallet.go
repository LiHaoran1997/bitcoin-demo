package main

import (
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
	publicKey []byte
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
	pubKey := w.publicKey
	rip160HashValue:=HashPubKey(pubKey)
	//版本
	version := byte(00)
	//拼接version
	payload := append([]byte{version}, rip160HashValue...)
	//checkSum
	checkCode:=checkSum(payload)
	payload = append(payload, checkCode...)
	//golang有base58
	address:= base58.Encode(payload)
	return address
}

func HashPubKey(data []byte)[]byte{
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

func checkSum(payload []byte)[]byte{
	//两次sha256
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	//前4字节校验码
	checkCode := hash2[:4]
	return checkCode
}
