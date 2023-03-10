package entity

type Blog struct {
	ID           uint64     `json:"id" gorm:"primaryKey"`
	Judul        string     `json:"judul" binding:"required"`
	TanggalPost  string     `json:"tgl_post" binding:"required"`
	Isi          string     `json:"isi" binding:"required"`
	JumlahLike   uint64     `json:"jml_like"`
	UserID       uint64     `gorm:"foreignKey" json:"user_id"`
	User         *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	ListKomentar []Komentar `json:"list_komentar,omitempty"`
}
