package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Person struct{
	//大写
	Name string
	Age int
}

func main() {
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	lily:=Person{"Lily",28}
	err:=encoder.Encode(&lily)
	if err!=nil{
		fmt.Println("encode failed!!!",err)
	}
	fmt.Println("after serialize :",buffer)
	var LILY Person
	decoder:=gob.NewDecoder(&buffer)
	err=decoder.Decode(&LILY)
	if err!=nil{
		fmt.Println("decode failed!!!",err)
	}
	fmt.Println(LILY)

}