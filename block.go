package main

import "crypto/sha256"

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
