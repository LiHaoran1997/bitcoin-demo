package main

//1. 定义结构

//4. 引⼊区块链

//7. 重构代码
func main() {
	bc:=NewBlockchain("测试")
	cli:=CLI{bc}
	cli.Run()
}
