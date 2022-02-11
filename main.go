package main

import "fmt"
//1. 定义结构

type Block struct{
	//i. 前区块哈希
	PrevHash []byte
	Hash []byte
	Data []byte
	//ii. 当前区块哈希
	//iii. 数据
}

//2. 创建区块

func NewBlock(data string,prevBlockHash []byte) *Block{
	block :=Block{
		PrevHash: prevBlockHash,
		Hash: []byte{}, //TODO:Hash计算
		Data: []byte(data),
	}
	return &block
}

//3. ⽣成哈希
//4. 引⼊区块链
//5. 添加区块
//6. 重构代码
func main() {
	block:=NewBlock("老师转班长一枚比特币！",[]byte{})
	fmt.Printf("前区块哈希值:%x\n",block.PrevHash)
	fmt.Printf("当前区块哈希值:%x\n",block.Hash)
	fmt.Printf("区块数据:  %s\n",block.Data)

}
