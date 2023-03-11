package dto

type UserRegisterRequest struct { // seseorang ingin melakukan registrasi akun blog
	ID       uint64 `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct { // seseorang ingin melakukan login ke akun blog
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserChangeNameRequest struct {
	Name string `json:"name" binding:"required"`
}
