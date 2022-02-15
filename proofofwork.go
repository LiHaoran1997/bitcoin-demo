package main

import "math/big"

//定义⼀个⼯作量证明的结构ProofOfWork
type ProofOfWork struct {
	//a. block
	block *Block
	//b. ⽬标值 一个大数，有很多方法：比较、赋值等
	target big.Int
}

//2. 提供创建POW的函数
//NewProofOfWork(参数)
func NewProofOfWork(block *Block)*ProofOfWork{
	pow:=ProofOfWork{
		block:block,
	}
	//指定的难度值，需要进行转换
	targetStr:="0000100000000000000000000000000000000000000000000000000000000000"
	//引入辅助变量，目的是将上面的难度值转为big.Int
	tmpInt:=big.Int{}
	tmpInt.SetString(targetStr,16)
	pow.target=tmpInt
	return &pow
}
//3. 提供计算不断计算hash的哈数
//Run()
//4. 提供⼀个校验函数
//IsValid()

