package main

import "fmt"

//1. 定义结构

//4. 引⼊区块链

//7. 重构代码
func main() {
	blockchain := NewBlockchain()
	blockchain.AddBlock("我爱蒋林志")
	blockchain.AddBlock("蒋林志是🐷")
	it := NewBlockchainiterator(blockchain)
	for {
		//调⽤迭代器访问函数， 返回当前block， 并且向左移动
		block := it.GetBlockAndMoveLeft()
		fmt.Println(" ============== =============")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Difficuty : %d\n", block.Diffculty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)
		if len(block.PrevHash)==0{
			fmt.Println("打印结束")
			break
		}
	}
}
