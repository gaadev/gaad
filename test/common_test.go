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

func TestSqliteUpdate(t *testing.T) {
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

type A struct {
	Name string
}

//// 测试unit
//func TestReflect(t *testing.T)  {
//	reflectNew((*A)(nil))
//}
//
////反射创建新对象。
//func reflectNew(target interface{}) {
//	if target == nil {
//		fmt.Println("参数不能未空")
//		return
//	}
//
//	t := reflect.TypeOf(target)
//	if t.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
//		t = t.Elem()
//	}
//
//	newStruc := reflect.New(t)                            // 调用反射创建对象
//
//	rValues := make([]reflect.Value, 0)
//	t := reflect.TypeOf(rValues)
//	fmt.Println(t)
//	newStruc.Elem().FieldByName("Name").SetString("Lily") //设置值
//
//	newVal := newStruc.Elem().FieldByName("Name") //获取值
//	fmt.Println(newVal.String())
//}
//
//
///*
//   需import "fmt" "reflect"
//   通过reflect反射获取不定长的任意object对象的type数据类型
//   返回数据类型切片
//*/
//func TypesOf(args ...interface{}) []reflect.Type {
//	mTypes := make([]reflect.Type, 0, cap(args))
//	for _, arg := range args {
//		mType := reflect.TypeOf(arg)
//		fmt.Println(" object=", arg, " type=", mType)
//		mTypes = append(mTypes, mType)
//	}
//	return mTypes
//}
//
