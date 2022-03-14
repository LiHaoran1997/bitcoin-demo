package main

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

//2.提供创建交易方法
//3.创建挖矿交易
