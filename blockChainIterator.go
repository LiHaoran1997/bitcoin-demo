package main

import (
	"../bolt"
	"log"
)

type BlockchainIterator struct {
	db            *bolt.DB
	current_point []byte
}

func (bc *Blockchain) NewIterator() *BlockchainIterator {
	return &BlockchainIterator{
		bc.db,
		//最初指向区块链的最后一个区块，随着Next的调用，不断变化
		bc.tail,
	}
}

//迭代器是属于区块链的
//Next方式是属于迭代器的
//1. 返回当前的区块
//2. 指针前移
func (it *BlockchainIterator) Next() *Block {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时bucket不应该为空，请检查!")
		}

		blockTmp := bucket.Get(it.current_point)
		//解码动作
		block = Deserialize(blockTmp)
		//游标哈希左移
		it.current_point = block.PrevHash

		return nil
	})

	return &block
}