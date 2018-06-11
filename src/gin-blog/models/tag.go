package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model

	Name string `json:"name"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

func ExistTagById(id int) bool  {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)

	if tag.ID > 0 {
		return true
	}

	return false
}

func AddTag(data map[string]interface{}) bool {
	db.Create(&Tag{
		Name:data["name"].(string),
		CreatedBy:data["created_by"].(string),
		State:data["state"].(int),
	})

	return true
}

func ExistTagByName(name string) bool  {
	var tag Tag
	db.Select("name").Where("name = ?", name).First(&tag)

	if tag.ID > 0 {
		return true
	}

	return false
}

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error  {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error  {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}

func GetTags(offset int, limit int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(offset).Limit(limit).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int)  {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

