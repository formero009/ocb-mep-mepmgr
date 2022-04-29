package mepmd

type Preference struct {
	Id       int64  `gorm:"column:id;primary_key;auto_increment;not null" json:"-"`
	NodeType string `gorm:"column:node_type" json:"nodeType,omitempty"`
	NodeId   string `gorm:"column:node_id" json:"nodeId,omitempty"`
	IconId   string `gorm:"column:icon_id" json:"iconId,omitempty"`
	UserName string `gorm:"column:user_name" json:"-"`
}
