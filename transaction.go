package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//1.定义交易结构
type Transaction struct {
	TXID     []byte
	TXInput  []TXInput  //交易输入数组
	TXOutPut []TXOutput //交易输出脚本
}

type TXInput struct {
	//1.引用交易ID
	TXid []byte
	//2.引用Output索引值
	Index int64
	//3.解锁脚本，用地址模拟
	Sig string
}

type TXOutput struct {
	//接受的金额
	Value int64
	//锁定脚本Hash值
	PubKeyHash string
}

//设置交易ID
func (tx *Transaction)SetHash(){
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	data:=buffer.Bytes()
	hash:=sha256.Sum256(data)
	tx.TXID=hash[:]
}

//2.提供创建交易方法
//3.创建挖矿交易
