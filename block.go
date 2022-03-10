package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	//版本号
	Version uint64
	//前区块哈希
	PrevHash []byte
	//Merkle根
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度目标
	Diffculty uint64
	//随机数,也就是挖矿要找的数据
	Nonce uint64
	//a.当前区块哈希，正常比特币区块中没有当前区块哈希，为了实现方便就简化
	Hash []byte
	//b.数据
	Data []byte
}

//辅助函数，实现将uint转为[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

//2. 创建区块

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Diffculty:  0, //无效值
		Nonce:      0, //无效值
		Hash:       []byte{},
		Data:       []byte(data),
	}
	//block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	//查找目标随机数，不停的进行哈希运算
	block.Hash = hash
	block.Nonce = nonce
	return &block
}
func (block *Block) toByte() []byte {
	return []byte{}
}

//3. ⽣成哈希
func (block *Block) SetHash() {
	var blockInfo []byte
	//1.拼装数据
	/*	blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
		blockInfo = append(blockInfo, block.PrevHash...)
		blockInfo = append(blockInfo, block.MerkleRoot...)
		blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Diffculty)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
		blockInfo = append(blockInfo, block.Hash...)
		blockInfo = append(blockInfo,block.Data...)*/

	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkleRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Diffculty),
		Uint64ToByte(block.Nonce),
		block.Hash,
		block.Data,
	}
	blockInfo = bytes.Join(tmp, []byte{})

	//2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
