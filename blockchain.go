package main

import (
	"../bolt"
	"bytes"
	"fmt"
	"os"
	"time"
)

type Blockchain struct {
	db *bolt.DB
	//尾巴， 存储最后⼀个区块的哈希
	tail []byte
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"
const last = "LastHashKey"
const genesisInfo = "这个是创世块"

//5.定义一个区块链
func NewBlockchain(address string) *Blockchain {
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
			genesisBlock := GenesisBlock(address)
			//fmt.Printf("genesisBlock :%s\n", genesisBlock)
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
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
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTx(address, genesisInfo)
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//6. 添加区块
func (bc *Blockchain) AddBlock(txs []*Transaction) {
	/*	prevHash:=bc.blocks[len(bc.blocks)-1].Hash
		//a.创建新的区块
		block:=NewBlock(data,prevHash)
		//b.添加到区块链数组中
		bc.blocks=append(bc.blocks,block)*/
	lastBlockHash := bc.tail
	newBlock := NewBlock(txs, lastBlockHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("未找到区块链Bucket！！！")
		} else {
			//添加区块
			bucket.Put(newBlock.Hash, newBlock.Serialize())
			bucket.Put([]byte(last), newBlock.Hash)
			bc.tail = newBlock.Hash
		}
		return nil
	})
}

func GetBlockChainObj() *Blockchain {
	var lastHash []byte
	if !dbExists() {
		fmt.Println("区块链未创建，请先创建区块链")
		os.Exit(1)
	}
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		fmt.Println("bolt.Open failed!", err)
		os.Exit(1)
	}
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("未找到Bucket")
			os.Exit(1)
		}
		lastHash = bucket.Get([]byte(last))
		return nil
	})
	return &Blockchain{db, lastHash}
}

func dbExists() bool {
	if _, err := os.Stat(blockChainDb); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *Blockchain) Printchain() {

	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockBucket"))

		//从第一个key-> value 进行遍历，到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}

			block := Deserialize(v)
			//fmt.Printf("key=%x, value=%s\n", k, v)
			fmt.Printf("=============== 区块高度: %d ==============\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
			fmt.Printf("梅克尔根: %x\n", block.MerkleRoot)
			timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("时间戳: %s\n", timeFormat)
			fmt.Printf("难度值(随便写的）: %d\n", block.Diffculty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInput[0].PubKey)
			return nil
		})
		return nil
	})
}

func (bc *Blockchain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutPut {
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

func (bc *Blockchain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxos集合
	utxos := make(map[string][]uint64)
	var calc float64

	txs := bc.FindUTXOTransactions(senderPubKeyHash)
	for _, tx := range txs {
		for i, output := range tx.TXOutPut {
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				if calc < amount {
					//1.把UTXO加进来
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					//2.统计当前总额
					calc += output.Value
					//3.比较是否满足转账需求
					if calc >= amount {
						fmt.Printf("找到了满足条件的金额:%f\n", calc)
						return utxos, calc
					}
					//a.满足直接返回 utxos,calc
					//b.不满足，继续统计
					//统计当前的
				} else {
					fmt.Printf("不满足转账金额，当前总额:%f\n，目标金额：%f\n", calc, amount)
				}
			}
		}
	}
	return utxos, calc

}

func (bc *Blockchain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含UTXO的交易
	//定义map保存消费过的output，key为output所在id，value这个交易索引数组
	//map[交易id][]uint64
	spentOutPuts := make(map[string][]int64)
	//1.遍历区块
	//创建迭代器
	it := bc.NewIterator()
	for {
		block := it.Next()
		//2.遍历交易
		for _, tx := range block.Transactions {
			//3.遍历output,找到与自己相关相关的utxo（在添加output之前检查下是否消耗过）
		OUTPUT:
			for i, output := range tx.TXOutPut {
				//这个output和我们目标地址相同，满足条件，加入到返回utxo数组中
				//做一次过滤，将所有消耗过的outputs和当前所添加的output对比一下
				//如果当前的被消耗，则跳过，否则添加
				//如果当前交易id已经存在于map，则说明交易中有消耗过的
				if spentOutPuts[string(tx.TXID)] != nil {
					for _, j := range spentOutPuts[string(tx.TXID)] {
						if int64(i) == j {
							//当前准备添加的output已经消耗过了，不用再加了
							continue OUTPUT
						}
					}
				}
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					//返回所有相关的交易
					txs = append(txs, tx)
				}
			}
			//如果当前交易是挖矿交易，不做遍历，直接跳过
			if !tx.IsCoinbase() {
				//4.遍历input,找到花费过的utxo集合（把自己消耗过的标识出来）ey
				for _, input := range tx.TXInput {
					//判断当前input和目标是否一致，如果相同，说明是消耗过的
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						spentOutPuts[string(input.TXid)] = append(spentOutPuts[string(input.TXid)], input.Index)
					}
				}
			} else {
				//fmt.Printf("这是coinbase交易，不需要遍历input!\n")
			}
		}
		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块遍历完成退出！")
		}
	}
	return txs
}

