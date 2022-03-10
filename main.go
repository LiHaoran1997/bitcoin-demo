package main

//1. 定义结构

//4. 引⼊区块链

//7. 重构代码
func main() {
	blockchain := NewBlockchain()
	cli:=Cli{blockchain}
	cli.Run()
}
