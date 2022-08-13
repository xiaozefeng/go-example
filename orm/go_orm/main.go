package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failedto conect database")
	}
	db.AutoMigrate(&Product{})
	db.Create(&Product{
		Code:  "D42",
		Price: 100,
	})

	var product Product
	db.First(&product, 1)
	db.First(&product, "code = ?", "D42")

	db.Model(&product).Update("Price", 200)
	db.Model(&product).Updates(Product{
		Code:  "F42",
		Price: 200,
	})
	db.Model(&product).Updates(map[string]interface{}{
		"Price": 200,
		"Code":  "F42",
	})

	db.Delete(&product, 1)
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}
