package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 50

//1.定义交易结构
type Transaction struct {
	TXID     []byte
	TXInput  []TXInput  //交易输入数组
	TXOutPut []TXOutput //交易输出脚本
}

type TXInput struct {
	//1.引用交易ID
	TXid []byte
	//2.引用Output索引值
	Index int64
	//3.解锁脚本，用地址模拟
	Sig string
}

type TXOutput struct {
	//接受的金额
	Value float64
	//锁定脚本Hash值
	PubKeyHash string
}

//设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//2.创建挖矿交易
func NewCoinbaseTx(address string, data string) *Transaction {
	//挖矿交易的特点：只有一个input，无需引用交易id，无需饮用index
	//矿工由于挖矿时无需指定签名，所以这个字段可以由矿工自由填写数据，一般填写矿池名字
	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	//对于挖矿交易来说，只有一个input和一个output
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()
	return &tx
}
func (tx *Transaction) IsCoinbase() bool {
	//1.交易input只有一个
	//2.交易ID为空
	//3.交易Index为-1
	if len(tx.TXInput) == 1 {
		input := tx.TXInput[0]
		if bytes.Equal(tx.TXInput[0].TXid, []byte{}) && input.Index == -1 {
			return true
		}
	}
	return false
}

//创建普通转账交易
//1.找到最合理的UTXO集合 map[string][]uint64
//2.将这些UTXO逐一转为inputs
//3.创建inputs、outputs
//4.如果有零钱，要找零
func NewTransaction(from string, to string, amount float64, bc *Blockchain) *Transaction {
	//1.找到最合理的UTXO集合 map[string][]uint64
	utxos, resValue := bc.FindNeedUTXOs(from, amount)
	if resValue < amount {
		fmt.Printf("余额不足，交易失败！")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	//2.将这些UTXO逐一转为inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), from}
			inputs = append(inputs, input)
		}
	}
	//创建交易输出
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	if resValue > amount {
		//找零
		outputs = append(outputs, TXOutput{resValue - amount, from})
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}
