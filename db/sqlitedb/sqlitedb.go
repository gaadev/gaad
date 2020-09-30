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

func QueryPage(curPage int, pageSize int, models interface{}, query interface{}, args ...interface{}) int {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var count int = 0

	// 获取取指page，指定pagesize的记录
	db.Where(query, args...).Limit(pageSize).Offset((curPage - 1) * pageSize).Order("updated_at desc").Find(models)

	// 获取总条数
	db.Model(models).Where(query, args...).Count(&count)
	return count
}

func QueryList(models interface{}, query interface{}, args ...interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	// 获取取指page，指定pagesize的记录
	db.Find(models).Where(query, args)
}

func Update(model interface{}, attrs ...interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()
	// 更新 - 更新product的price为2000
	db.Model(model).Update(attrs...)
}

func Delete(model interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Delete(model)
}
