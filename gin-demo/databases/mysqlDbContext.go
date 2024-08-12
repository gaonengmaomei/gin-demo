package databases

import (
	"gin-demo/models"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func AutoMigrate(data map[string]*models.ModelStruct) {
	connStr := "root:Mysql925299@tcp(192.168.192.188:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	for _, v := range data {
		modelStruct := *v
		rType := reflect.StructOf(modelStruct.StructField)
		boxPtr := reflect.New(rType).Interface()
		zeroVal := &boxPtr
		db.Table(modelStruct.TableName).AutoMigrate(zeroVal)
	}
}
