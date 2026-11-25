package model

import (
	"context"
	"errors"

	"SolarNotesAPI/db"
)

// Planet 映射 planet 表中简介接口需要的列。
type Planet struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	Name  string `gorm:"column:name"`
	Title string `gorm:"column:title"`
	Intro string `gorm:"column:intro"`
}

// FindAll 按主键升序读取所有已创建的天体。
func (Planet) FindAll(ctx context.Context) ([]Planet, error) {
	if db.DB == nil {
		return nil, errors.New("数据库连接尚未初始化")
	}

	var rows []Planet
	err := db.DB.WithContext(ctx).
		Select("id", "name", "title", "intro").
		Order("id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
