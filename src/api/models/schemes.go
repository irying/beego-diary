package models

//Schemes 配色方案表
type Schemes struct {
	ID         int    `orm:"column(id)"`
	Name       string `orm:"column(name)"`
	ApplyPage  string
	Levels     int
	Style      string
	CreateTime int64
	UpdateTime int64
	TemplateID string `orm:"column(template_id)"`
	Background int
	Substrate  int
	SchemeID   string `orm:"column(scheme_id)"`
}

//SchemesItems 配色方案各层级数量表
type SchemesItems struct {
	ID            int    `orm:"column(id)"`
	SchemeID      string `orm:"column(scheme_id)"`
	TemplateID    string `orm:"column(template_id)"`
	Level         int
	Background    int
	Decoration    int
	DiagramSeries int
	Directitem    int
	ItemKey       string
}

//SchemesColor 存储颜色值
type SchemesColor struct {
	ID         int    `orm:"column(id)"`
	SchemeID   string `orm:"column(scheme_id)"`
	TemplateID string `orm:"column(template_id)"`
	Base       string
	Contrast   string
	Type       string
	ItemKey    string
}

//SchemesJSON 配色方案原始数据
type SchemesJSON struct {
	SchemeID   string `orm:"pk;column(scheme_id)"`
	TemplateID string `orm:"column(template_id)"`
	Category   string
	JSON       string `orm:"column(json)"`
}

//TableName 自定义表名（系统自动调用）
func (u *SchemesJSON) TableName() string {
	return "schemes_json"
}
