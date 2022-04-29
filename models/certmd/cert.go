/*
@Time : 2022/4/27
@Author : jzd
@Project: ocb-mep-mepmgr
*/
package certmd

import "time"

const (
	ClientCert = 0
	ServerCert = 1
	RootCert   = 2

	StatusValid       = 0
	StatusToBeInvalid = 1
	StatusInvalid     = 2
	StatusDeprecated  = 3
)

type Cert struct {
	Id        string    `gorm:"column:id;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Type      int       `gorm:"column:cert_type;INDEX:cert_type,UNIQUE" json:"type"`
	Cert      string    `gorm:"column:cert" json:"cert"`
	Key       string    `gorm:"column:key" json:"key"`
	Status    int       `gorm:"column:status" json:"status"`
	ValidTime time.Time `gorm:"column:valid_time" json:"validTime"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}
