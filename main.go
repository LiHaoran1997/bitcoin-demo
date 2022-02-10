package main

import "fmt"

func main() {
	fmt.Printf("hello")
	total:=0.0
	blockIntever:=21 //单位是w
	currentReward:=50.0
	for currentReward>0{
		amount1:=float64(blockIntever) * currentReward
		currentReward*=0.5
		total+=amount1
	}
	fmt.Println(total)
}
