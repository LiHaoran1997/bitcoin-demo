package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

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
func (pow *ProofOfWork)Run()([]byte,uint64){
	var nonce uint64
	block:=pow.block
	var hash [32]byte
	fmt.Println("开始挖矿....")
	for{
		//1.拼装数据（区块数据，不断变化随机数）
		tmp:=[][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkleRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Diffculty),
			Uint64ToByte(nonce),
			block.Hash,
			block.Data,
		}
		blockInfo:=bytes.Join(tmp,[]byte{})
		//2.哈希运算
		hash = sha256.Sum256(blockInfo)
		block.Hash = hash[:]
		//3.与pow中的target比较
		tmpInt:=big.Int{}
		tmpInt.SetBytes(hash[:])
		//比较当前哈希值和目标哈希，若当前小于目标哈希，找到；否则继续找
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if tmpInt.Cmp(&pow.target)==-1{
			fmt.Printf("挖矿成功！ hash : %x ,nonce : %d\n",hash,nonce)
			break
		}else{
			nonce++
		}
		//a.找到了退出返回
		//b.没找到，继续找，随机数+1
	}
	return hash[:],nonce
}
//Run()
//4. 提供⼀个校验函数
//IsValid()

