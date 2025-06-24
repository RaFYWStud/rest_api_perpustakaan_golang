package entity

type Intro struct {
	ID            int    `gorm:"column:id;primaryKey;autoIncrement;not null;<-create"`
	Nama          string `gorm:"column:nama;type:varchar(255);not null"`
	NamaPanggilan string `gorm:"column:nama_panggilan;type:varchar(255)"`
	FunFact       string `gorm:"column:fun_fact;type:varchar(255)"`
	KeinginanBE   string `gorm:"column:keinginan_be;type:varchar(255)"`
	UpdatedAt     string `gorm:"column:updated_at;type:timestamp;not null;default:now()"`
	CreatedAt     string `gorm:"column:created_at;type:timestamp;not null;default:now()"`
}

func (e *Intro) TableName() string {
	return "intro"
}
