package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
)

//定义一个Wallets结构，保存所有的wallet和地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

//创建方法
func NewWallets() *Wallets {
	var ws Wallets
	//ws.WalletsMap = make(map[string]*Wallet)

	//ws := loadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	ws.WalletsMap = make(map[string]*Wallet)
	ws.WalletsMap[address] = wallet
	ws.saveToFile()
	return address
}

//保存方法，把新建的wallet添加进去
func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err:=encoder.Encode(&ws)
	if err!=nil{
		log.Panic(err)
	}
	ioutil.WriteFile("wallet.dat", buffer.Bytes(), 0600)
}

//读取文件方法，把所有的wallet读出来
