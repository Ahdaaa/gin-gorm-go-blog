package controller

import (
	"net/http"
	"oprec/go-blog/entity"
	"oprec/go-blog/dto"
	"oprec/go-blog/service"
	"oprec/go-blog/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Upload(ctx *gin.Context)
	GetAllBlog(ctx *gin.Context)
	GetBlogByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Comment(ctx *gin.Context)
	BlogDetails(ctx *gin.Context)
	Like(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewUserController(us service.UserService) UserController {
	return &userController{
		userService: us,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var userDTO dto.UserRegisterRequest

	errDTO := ctx.ShouldBind(&userDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.userService.CreateUser(ctx, userDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("Akun go-blog anda berhasil dibuat", http.StatusCreated, user)
	ctx.JSON(http.StatusCreated, response)

}

func (c *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	
	errDTO := ctx.ShouldBind(&userDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.userService.FindUserByEmail(ctx, userDTO.Email)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	checkPass, _ := utils.ComparePassword(user.Password, []byte(userDTO.Password))
	if !checkPass {
		response := utils.BuildErrorResponse("Password Salah", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	tokenService := service.NewJWTService()
	tokenString := tokenService.GenerateToken(user.ID, user.Name)

	response := utils.BuildResponse("Anda berhasil login, berikut token anda:", http.StatusCreated, tokenString)
	ctx.JSON(http.StatusCreated, response)

}

func (c *userController) Upload(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1) // menghilangkan kata bearer dari token, karena mau di proses
	tokenService := service.NewJWTService()

	id, _ := tokenService.GetUserIDByToken(token) // id sudah pasti ada jika berhasil melewati validate, jadi tidak mungkin error
	
	var blogDTO dto.BlogCreateRequest
	errDTO := ctx.ShouldBind(&blogDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	//fmt.Println("tes")

	newDTO := dto.BlogCreateRequest{
		Judul:       blogDTO.Judul,
		TanggalPost: blogDTO.TanggalPost,
		Isi:         blogDTO.Isi,
		UserID:      id,
	}

	activeBlog, err := c.userService.CreateBlog(ctx, newDTO)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal membuat blog", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("blog anda berhasil dibuat", http.StatusCreated, activeBlog)
	ctx.JSON(http.StatusCreated, response)

}

func (c *userController) GetAllBlog(ctx *gin.Context) {
	getBlogs, err := c.userService.GetAllBlog(ctx)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal melihat blog", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil", http.StatusOK, getBlogs)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) GetBlogByID(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	id, _ := tokenService.GetUserIDByToken(token)

	getBlog, err := c.userService.GetBlogByID(id)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal melihat blog", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil", http.StatusOK, getBlog)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) Update(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	id, _ := tokenService.GetUserIDByToken(token)
	
	var name dto.UserChangeNameRequest
	err := ctx.ShouldBind(&name)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	updateName, err := c.userService.UpdateName(id, name.Name)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := utils.BuildResponse("berhasil ganti nama", http.StatusOK, updateName)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) Comment(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	username, _ := tokenService.GetNameByToken(token)

	idBlog, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("ID invalid", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var commentDTO dto.CreateCommentRequest
	errDTO := ctx.ShouldBind(&commentDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newDTO := dto.CreateCommentRequest{
		Username: username,
		IsiKomen: commentDTO.IsiKomen,
		BlogID:   idBlog,
	}

	comment, err := c.userService.CreateComment(newDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest) // semisal terjadi post komentar pada blog yang belum dibuat
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("komentar berhasil dibuat", http.StatusCreated, comment)
	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) BlogDetails(ctx *gin.Context) {
	idBlog, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("ID invalid", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	details, err := c.userService.GetBlogDetails(idBlog)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest) // semisal terjadi post komentar pada blog yang belum dibuat
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil", http.StatusOK, details)
	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) Like(ctx *gin.Context) {
	idBlog, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("ID invalid", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	checkErr := c.userService.GiveLikeByID(idBlog)
	if checkErr != nil {
		response := utils.BuildErrorResponse("gagal memberi like", http.StatusBadRequest) // semisal terjadi post komentar pada blog yang belum dibuat
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil memberi like", http.StatusCreated, nil)
	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) Delete(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	id, _ := tokenService.GetUserIDByToken(token)

	checkErr := c.userService.DeleteUser(id)
	if checkErr != nil {
		response := utils.BuildErrorResponse("gagal mendelete akun", http.StatusBadRequest) // semisal terjadi post komentar pada blog yang belum dibuat
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil mendelete", http.StatusCreated, nil)
	ctx.JSON(http.StatusCreated, response)
}
