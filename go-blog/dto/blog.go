package dto

type BlogCreateRequest struct { // seseorang ingin melakukan login ke akun blog
	Judul       string `json:"judul" binding:"required"`
	TanggalPost string `json:"tgl_post" binding:"required"`
	Isi         string `json:"isi" binding:"required"`
	UserID      uint64 `json:"user_id"`
}
