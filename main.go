package main

import (
	"fmt"

)

//1. 定义结构




//4. 引⼊区块链

//7. 重构代码
func main() {
	blockchain:=NewBlockchain()
	blockchain.AddBlock("我爱蒋林志")
	for i,block:=range blockchain.blocks{
		fmt.Printf("=========当前区块高度:  %d===============\n", i)
		fmt.Printf("前区块哈希值:%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值:%x\n", block.Hash)
		fmt.Printf("区块数据:  %s\n", block.Data)
	}
}
