package entity

import "time"

type StatusKetersediaan string

const (
    Available    StatusKetersediaan = "available"
    NotAvailable StatusKetersediaan = "not_available"
)

type Perpus struct {
    ID                  int                `gorm:"column:id;primaryKey;autoIncrement;not null;<-create"`
    Judul              string             `gorm:"column:judul;type:varchar(255);not null"`
    Penulis            string             `gorm:"column:penulis;type:varchar(100);not null"`
    StatusKetersediaan StatusKetersediaan `gorm:"column:status_ketersediaan;type:ketersediaan;not null;default:available"`
    LinkBuku           string             `gorm:"column:link_buku;type:varchar(500)"`
    CreatedAt          time.Time          `gorm:"column:created_at;type:timestamp;not null;default:now()"`
    UpdatedAt          time.Time          `gorm:"column:updated_at;type:timestamp;not null;default:now()"`
}

func (e *Perpus) TableName() string {
    return "perpus"
}