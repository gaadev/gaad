package boltdb

import (
	"github.com/boltdb/bolt"
	"log"
)

const (
	dbFile = "bolt.db"
)

func Update(key string, value string) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	// 创建或者打开数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表
	err = db.Update(func(tx *bolt.Tx) error {
		// 获取BlockBucket表单
		b, _ := tx.CreateBucketIfNotExists([]byte("GaadBucket"))

		// 往表里面存储数据
		if b != nil {
			err := b.Put([]byte(key), []byte(value))
			if err != nil {
				log.Panic("数据存储失败......")
			}
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}

}

func View(key string) (value string) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	// 创建或者打开数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查看数据
	err = db.View(func(tx *bolt.Tx) error {

		// 获取BlockBucket表对象
		b := tx.Bucket([]byte("GaadBucket"))

		// 往表里面存储数据
		if b != nil {
			data := b.Get([]byte(key))
			value = string(data)
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}
	return
}
