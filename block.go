package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"os"
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
	//Data []byte
	//真实交易数据
	Transactions []*Transaction
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

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:      00,
		PrevHash:     prevBlockHash,
		MerkleRoot:   []byte{},
		TimeStamp:    uint64(time.Now().Unix()),
		Diffculty:    0, //无效值
		Nonce:        0, //无效值
		Hash:         []byte{},
		Transactions: txs,
	}
	block.MerkleRoot=block.MakeMerkleRoot()
	//block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	//查找目标随机数，不停的进行哈希运算
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

////3. ⽣成哈希
//func (block *Block) SetHash() {
//	var blockInfo []byte
//	//1.拼装数据
//	/*	blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
//		blockInfo = append(blockInfo, block.PrevHash...)
//		blockInfo = append(blockInfo, block.MerkleRoot...)
//		blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
//		blockInfo = append(blockInfo, Uint64ToByte(block.Diffculty)...)
//		blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
//		blockInfo = append(blockInfo, block.Hash...)
//		blockInfo = append(blockInfo,block.Data...)*/
//
//	tmp := [][]byte{
//		Uint64ToByte(block.Version),
//		block.PrevHash,
//		block.MerkleRoot,
//		Uint64ToByte(block.TimeStamp),
//		Uint64ToByte(block.Diffculty),
//		Uint64ToByte(block.Nonce),
//		block.Hash,
//		block.Transactions,
//	}
//	blockInfo = bytes.Join(tmp, []byte{})
//
//	//2.sha256
//	hash := sha256.Sum256(blockInfo)
//	block.Hash = hash[:]
//}

func (block *Block) Serialize() []byte {
	//将block数据转换为字节流
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		fmt.Println("encode failed!!!", err)
		os.Exit(1)
	}
	return buffer.Bytes()
}

func Deserialize(data []byte) Block {
	var block Block
	var buffer bytes.Buffer
	_, err := buffer.Write(data)
	if err != nil {
		fmt.Println("buffer.Read failed!", err)
		os.Exit(1)
	}
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&block)
	if err != nil {
		fmt.Println("decode failed!", err)
		os.Exit(1)
	}
	return block
}

//模拟默克尔根，对交易数据做简单拼接
func (block *Block)MakeMerkleRoot()[]byte{
	var info []byte
	//var finalInfo [][]byte
	for _, tx := range block.Transactions {
		//将交易的哈希值拼接起来，再整体做哈希处理
		info = append(info, tx.TXID...)
		//finalInfo = [][]byte{tx.TXID}
	}

	hash := sha256.Sum256(info)
	return hash[:]
}