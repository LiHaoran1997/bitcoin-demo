package main

import (
	"../bolt"
	"fmt"
	"os"
)

type BlockchainIterator struct {
	db            *bolt.DB
	current_point []byte
}

func NewBlockchainiterator(bc *Blockchain) BlockchainIterator {
	var it BlockchainIterator
	it.db = bc.db
	it.current_point = bc.tail
	return it
}

func (it *BlockchainIterator) GetBlockAndMoveLeft() Block {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("bucket为空！！！")
			os.Exit(1)
		} else {
			current_block_tmp := bucket.Get(it.current_point)
			if current_block_tmp==nil|| len(current_block_tmp)==0{
				return nil
			}
			current_block := Deserialize(current_block_tmp)
			block = current_block
			it.current_point = current_block.PrevHash
		}
		return nil
	})
	return block
}
