package main

type Blockchain struct {
	//定义一个区块链数组
	blocks []*Block
}

//5.定义一个区块链
func NewBlockchain() *Blockchain {
	//创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock:=GenesisBlock()
	return &Blockchain{
		blocks:[] *Block{genesisBlock},
	}
}

//创世块
func GenesisBlock() *Block{
	return NewBlock("这个是创世块",[]byte{})
}
//6. 添加区块
func (bc *Blockchain)AddBlock(data string){
	prevHash:=bc.blocks[len(bc.blocks)-1].Hash
	//a.创建新的区块
	block:=NewBlock(data,prevHash)
	//b.添加到区块链数组中
	bc.blocks=append(bc.blocks,block)
}
