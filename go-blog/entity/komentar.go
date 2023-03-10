package entity

type Komentar struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Username string `json:"username"` // username nanti saya dapatkan melalui getNamebyToken
	IsiKomen string `json:"komen" binding:"required"`
	BlogID   uint64 `gorm:"foreignKey" json:"blog_id" binding:"required"`
	Blog     *Blog  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"blog,omitempty"`
}
