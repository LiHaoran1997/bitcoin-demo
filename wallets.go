package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)
const walletFile="wallet.dat"

//定义一个Wallets结构，保存所有的wallet和地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

//创建方法
func NewWallets() *Wallets {
	//ws.WalletsMap = make(map[string]*Wallet)
	var ws Wallets
	ws.loadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	ws.WalletsMap[address] = wallet
	for i,j:=range ws.WalletsMap{
		fmt.Println(i)
		fmt.Println(j.PublicKey)
		fmt.Println(j.PrivateKey)
	}
	ws.saveToFile()
	return address
}

//保存方法，把新建的wallet添加进去
func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&ws)
	if err != nil {
		log.Panic(err)
	}
	ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
}

//读取文件方法，把所有的wallet读出来
func (ws *Wallets) loadFile() {
	//读取前，先确认文件是否存在，如果不存在，直接退出，否则继续
	_,err:=os.Stat(walletFile)
	if os.IsNotExist(err){
		ws.WalletsMap=make(map[string]*Wallet)
		return
	}
	//读取内容
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wsLocal Wallets
	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}
	//对于结构来说，里面有map的，要指定赋值，不要在外面赋值
	ws.WalletsMap=wsLocal.WalletsMap
}

func (ws *Wallets)ListAllAddresses()[]string{
	var addresses []string
	//遍历所有钱包
	for address:=range ws.WalletsMap{
		addresses=append(addresses, address)
	}
	return addresses
}

