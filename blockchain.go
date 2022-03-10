package main

import (
	"../bolt"
	"fmt"
	"os"
)

type Blockchain struct {
	db *bolt.DB
	//尾巴， 存储最后⼀个区块的哈希
	tail []byte
}

const blockChainDb="blockChain.db"
const blockBucket="blockBucket"
const last="LastHashKey"
//5.定义一个区块链
func NewBlockchain() *Blockchain {
	var lastHash []byte
	/*
	   1. 打开数据库(没有的话就创建)
	   2. 找到抽屉（bucket） ， 如果找到， 就返回bucket， 如果没有找到， 我们要创建bucket， 通过名字创建
	   a. 找到了
	   1. 通过"last"这个key找到我们最好⼀个区块的哈希。
	   b. 没找到创建
	   1. 创建bucket， 通过名字
	   2. 添加创世块数据
	   3. 更新"last"这个key的value（创世块的哈希值）
	*/
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		fmt.Println("bolt.Open failed!", err)
		os.Exit(1)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		var err error
		//如果是空的， 表明这个bucket没有创建， 我们就要去创建它， 然后再写数据。
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				fmt.Println("createBucket failed!", err)
				os.Exit(1)
			}
			//创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock:=GenesisBlock()
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//TODO
			bucket.Put([]byte(last), genesisBlock.Hash)
			//这个别忘了， 我们需要返回它
			lastHash = genesisBlock.Hash
			return nil
			//抽屉已经存在， 直接读取即可
		} else {
			//获取最后⼀个区块的哈希
			lastHash = bucket.Get([]byte(last))
		}
		return nil
	})
	return &Blockchain{db, lastHash}
}

//创世块
func GenesisBlock() *Block{
	return NewBlock("这个是创世块",[]byte{})
}
//6. 添加区块
func (bc *Blockchain)AddBlock(data string){
/*	prevHash:=bc.blocks[len(bc.blocks)-1].Hash
	//a.创建新的区块
	block:=NewBlock(data,prevHash)
	//b.添加到区块链数组中
	bc.blocks=append(bc.blocks,block)*/
	lastBlockHash:=bc.tail
	newBlock:=NewBlock(data,lastBlockHash)
	bc.db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			fmt.Println("未找到区块链Bucket！！！")
		}else{
			//添加区块
			bucket.Put(newBlock.Hash, newBlock.Serialize())
			bucket.Put([]byte(last), newBlock.Hash)
			bc.tail=newBlock.Hash
		}
		return nil
	})


}
