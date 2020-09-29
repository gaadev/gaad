package sqlitedb

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	dialect = "sqlite3"
	dbFile  = "gaad.db"
)

func Create(modle interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	// 自动迁移模式
	db.AutoMigrate(modle)

	// 创建
	db.Create(modle)
}

func First(modle interface{}, where ...interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.First(modle, where...) // 查询id为1的product
}

func QueryPage(curPage int, pageSize int, modles interface{}, query interface{}, where ...interface{}) int {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	var count int = 0

	// 获取取指page，指定pagesize的记录
	db.Where(query, where...).Limit(pageSize).Offset((curPage - 1) * pageSize).Order("updated_at,created_at desc").Find(modles)

	// 获取总条数
	db.Model(modles).Where(query, where...).Count(&count)
	return count
}

func Update(modle interface{}, where ...interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()
	// 更新 - 更新product的price为2000
	db.Model(modle).Update(where...)
}

func Delete(modle interface{}) {
	db, err := gorm.Open(dialect, dbFile)
	if err != nil {
		panic("连接数据库失败")
	}
	defer db.Close()

	db.Delete(modle)
}
