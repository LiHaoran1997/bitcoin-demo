package main

import (
	"fmt"
	"../bolt"
	"log"
)

func main() {
	fmt.Println("helloworld")
	//1.打开数据库
	db,err:=bolt.Open("test.db",0600,nil)
	if err!=nil{
		log.Panic("打开数据库失败")
	}
	//操作数据库
	db.Update(func(tx *bolt.Tx) error {
		//2.找到Bucket,如果没有则创建
		bucket:=tx.Bucket([]byte("b1"))
		if bucket==nil{
			//没有抽屉,要创建
			bucket,err=tx.CreateBucket([]byte("b1"))
			if err!=nil{
				log.Panic("创建bucket(b1)失败")
			}
		}
		bucket.Put([]byte("11111"),[]byte("hello"))
		bucket.Put([]byte("22222"),[]byte("world"))
		return nil
	})
	//3.写数据
	//4.读数据
}
