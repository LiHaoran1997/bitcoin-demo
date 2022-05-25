package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

//1.演示如何使用ecdsa生成公私钥
//2.签名校验
func main() {
	//创建曲线
	curve:=elliptic.P256()
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		log.Panic()
	}

	//生成公钥
	publicKey:=privateKey.PublicKey
	data:="hello world!"
	hash:=sha256.Sum256([]byte(data))
	//签名
	r,s,err:=ecdsa.Sign(rand.Reader,privateKey,hash[:])
	if err!=nil{
		log.Panic()
	}
	fmt.Printf("pubkey: %v\n",publicKey)
	fmt.Printf("r: %v\n",r.Bytes())
	fmt.Printf("s: %v\n",s.Bytes())

	//把r,s进行序列化传输
	signature:=append(r.Bytes(),s.Bytes()...)

	//1.定义两个辅助的Bigint
	r1:=big.Int{}
	s1:=big.Int{}
	//2.拆分signature，均分
	r1.SetBytes(signature[0:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])
	//校验需要数据、签名、公钥
	res:=ecdsa.Verify(&publicKey,hash[:] ,&r1,&s1)
	fmt.Printf("校验结果:%v\n",res)


}
