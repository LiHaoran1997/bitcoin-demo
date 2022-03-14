package main

import (
	"fmt"
	"os"
)

type Cli struct {
}

const Usage = `
             createBlockChain --address ADDRESS "create a blockchain"
             addBlock --data DATA "add a block"
             printChain "print blockchain"`

func (cli *Cli) Run() {
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "createBlockChain":
		if len(os.Args) > 3 && os.Args[2] == "--address" {
			address := os.Args[3]
			if address == "" {
				fmt.Println("地址为空！")
				os.Exit(1)
			}
			cli.createBlockchain(address)
		} else {
			fmt.Println(Usage)
		}
	case "addBlock":
		if len(os.Args) > 3 && os.Args[2] == "--data" {
			data := os.Args[3]
			if data == "" {
				fmt.Println("区块数据不能为空")
				os.Exit(1)
			}
			cli.addBlock(data)
		} else {
			fmt.Println(Usage)
		}
	case "printChain":
		cli.printChain()
	default:
		fmt.Println(Usage)
	}
}

func (cli *Cli) addBlock(data string) {
	bc:=GetBlockChainObj()
	bc.AddBlock(data)
	bc.db.Close()
	fmt.Println("创建区块成功")
}

func (cli *Cli) printChain() {
	bc:=GetBlockChainObj()
	it := NewBlockchainiterator(bc)
	for {
		block := it.GetBlockAndMoveLeft()
		fmt.Println(" ============== =============")
		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
		fmt.Printf("梅克尔根: %x\n", block.MerkleRoot)
		fmt.Printf("时间戳: %d\n", block.TimeStamp)
		fmt.Printf("难度值(随便写的）: %d\n", block.Diffculty)
		fmt.Printf("随机数 : %d\n", block.Nonce)
		fmt.Printf("当前区块哈希值: %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Data)
		if len(block.PrevHash) == 0 {
			fmt.Println("区块数据遍历结束")
			break
		}
	}
}

func (cli *Cli) createBlockchain(address string) {
	bc := NewBlockchain(address)
	err := bc.db.Close()
	if err != nil {
		if dbExists() {
			os.Remove(blockChainDb)
		}
		fmt.Println("创建区块链成功")
	}
	fmt.Println("创建区块链成功")
}
