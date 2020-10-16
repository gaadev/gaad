package sqlitedb

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	dialect = "sqlite3"
	dbFile  = "gaad.db"
)

func Create(model interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	// 自动迁移模式
	db.AutoMigrate(model)

	// 创建
	db.Create(model)
}

func First(model interface{}, where ...interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.First(model, where...) // 查询id为1的product
}

func QueryPage(curPage int, PageRecord int, models interface{}, query interface{}, args ...interface{}) int {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var count int = 0

	// 获取取指page，指定PageRecord的记录
	db.Where(query, args...).Limit(PageRecord).Offset((curPage - 1) * PageRecord).Order("updated_at desc").Find(models)

	// 获取总条数
	db.Model(models).Where(query, args...).Count(&count)
	return count
}

func QueryList(models interface{}, where ...interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	// 获取取指page，指定PageRecord的记录
	db.Find(models, where...)
}

func Update(model interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	// 自动迁移模式
	db.AutoMigrate(model)
	//model为pointer
	db.Model(model).Update(model)
}

func Delete(model interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Delete(model)
}

func DeleteForce(model interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Unscoped().Delete(model)
}
