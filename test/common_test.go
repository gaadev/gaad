package test

import (
	"fmt"
	"gaad/common"
	"gaad/db/boltdb"
	"gaad/db/sqlitedb"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	m.Run()
}

func TestUUID(t *testing.T) {
	// 创建
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)

	// 解析
	u2, err := uuid.FromString("f5394eef-e576-4709-9e4b-a7c231bd34a4")
	if err != nil {
		fmt.Printf("Something gone wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s", u2)
}

func TestBoltdb(t *testing.T) {
	boltdb.Update("age", "12")
	age := boltdb.View("age")
	fmt.Println(age)
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestSqlitedb(t *testing.T) {
	product := Product{Code: "L1212", Price: 1000}
	sqlitedb.Create(&product)
	// 读取

	var pro Product

	sqlitedb.First(&pro, "code = ?", "L1212")
	sqlitedb.Update(&product, "Price", 2000)
	fmt.Printf("%v", pro)
	sqlitedb.Delete(&product)
}

func TestSqlitedb2(t *testing.T) {

	// 读取

	var pro Product

	sqlitedb.First(&pro, "code = ?", "L1212")

	fmt.Printf("%v", pro)

}

func TestCreateFile(t *testing.T) {
	common.CreateFile("log/log.log")
}

func TestTruncation(t *testing.T) {
	var str = "aaaa/log/log.log/"
	pos := strings.LastIndex(str, "/")
	if pos != -1 {
		fmt.Println(str[:pos])
	}

}
