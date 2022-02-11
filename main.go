package main

import (
	"crypto/sha256"
	"fmt"
)

//1. 定义结构

type Block struct {
	//i. 前区块哈希
	PrevHash []byte
	Hash     []byte
	Data     []byte
	//ii. 当前区块哈希
	//iii. 数据
}

//2. 创建区块

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

//3. ⽣成哈希
func (block *Block) SetHash() {
	//1.拼装数据
	blockInfo := append(block.PrevHash, block.Data...)
	//2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

//4. 引⼊区块链
type Blockchain struct {
	//定义一个区块链数组
	blocks []*Block
}

//5.定义一个区块链
func NewBlockchain() *Blockchain {
	//创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock:=GenesisBlock()
	return &Blockchain{
		blocks:[] *Block{genesisBlock},
	}
}

//创世块
func GenesisBlock() *Block{
	return NewBlock("创世块",[]byte{})
}
//6. 添加区块
//7. 重构代码
func main() {
	blockchain:=NewBlockchain()
	for i,block:=range blockchain.blocks{
		fmt.Printf("=========当前区块高度:  %d===============\n", i)
		fmt.Printf("前区块哈希值:%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值:%x\n", block.Hash)
		fmt.Printf("区块数据:  %s\n", block.Data)
	}
	block := NewBlock("老师转班长一枚比特币！", []byte{})
	fmt.Printf("前区块哈希值:%x\n", block.PrevHash)
	fmt.Printf("当前区块哈希值:%x\n", block.Hash)
	fmt.Printf("区块数据:  %s\n", block.Data)
}
