package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `gorm:"primary_key; type:UUID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *Base) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type Deck struct {
	Base
	Shuffled bool   `gorm:"column:shuffled;not null;default:false;type:bool"`
	Cards    []Card `gorm:"foreignKey:DeckRefer;references:ID;constraint:OnDelete:CASCADE;"`
}

type Card struct {
	Base
	Value     string `gorm:"column:value;not null;type:varchar"`
	Suit      string `gorm:"column:remaining;not null;type:varchar"`
	Code      string `gorm:"column:shuffled;not null;type:varchar"`
	Revealed  bool   `gorm:"column:revealed;not null;default:false;type:bool"`
	DeckRefer uuid.UUID
}
