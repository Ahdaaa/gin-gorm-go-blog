package service

import (
	"context"
	"errors"
	"oprec/go-blog/dto"
	"oprec/go-blog/entity"
	"oprec/go-blog/repository"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

type userService struct {
	userRepo repository.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserRegisterRequest) (entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateBlog(ctx context.Context, blogDTO dto.BlogCreateRequest) (entity.Blog, error)
	GetAllBlog(ctx context.Context) ([]entity.Blog, error)
	GetBlogByID(id uint64) (entity.User, error)
	UpdateName(id uint64, name string) (entity.User, error)
	CreateComment(commentDTO dto.CreateCommentRequest) (entity.Komentar, error)
	GetBlogDetails(id uint64) (entity.Blog, error)
	GiveLikeByID(id uint64) error
	DeleteUser(id uint64) error
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepo: ur,
	}
}

func (s *userService) CreateUser(ctx context.Context, userDTO dto.UserRegisterRequest) (entity.User, error) {
	var user entity.User
	copier.Copy(&user, &userDTO)

	checkEmail, err := s.userRepo.FindUserByEmail(ctx, nil, user.Email)
	if err != nil {
		return entity.User{}, err
	}

	if !(cmp.Equal(checkEmail, entity.User{})) { // saya menggunakan library cmp dengan tujuan untuk membandingkan 2 struct, tidak bisa dengan = karena pada struct user terdapat []blog
		return entity.User{}, errors.New("email yang diinput sudah pernah digunakan")
	}

	berhasilRegis, err := s.userRepo.CreateUser(ctx, nil, user)
	if err != nil {
		return entity.User{}, err
	}

	return berhasilRegis, nil

}

func (s *userService) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, nil, email)
	if err != nil {
		return entity.User{}, err
	}

	if cmp.Equal(user, entity.User{}) {
		return entity.User{}, errors.New("email tidak valid")
	}

	return user, nil
}

func (s *userService) CreateBlog(ctx context.Context, blogDTO dto.BlogCreateRequest) (entity.Blog, error) {
	var blog entity.Blog
	copier.Copy(&blog, &blogDTO)

	berhasilUpload, err := s.userRepo.CreateBlog(ctx, nil, blog)
	if err != nil {
		return entity.Blog{}, err
	}

	return berhasilUpload, nil
}

func (s *userService) GetAllBlog(ctx context.Context) ([]entity.Blog, error) {
	berhasilGet, err := s.userRepo.GetAllBlog(ctx, nil)
	if err != nil {
		return []entity.Blog{}, err
	}

	return berhasilGet, nil
}

func (s *userService) GetBlogByID(id uint64) (entity.User, error) {
	berhasilGet, err := s.userRepo.GetBlogByID(nil, id)
	if err != nil {
		return entity.User{}, err
	}

	return berhasilGet, nil
}

func (s *userService) UpdateName(id uint64, name string) (entity.User, error) {
	if name == "" {
		return entity.User{}, errors.New("nama tidak boleh kosong")
	}

	berhasilUpdate, err := s.userRepo.UpdateName(nil, id, name)
	if err != nil {
		return entity.User{}, err
	}

	return berhasilUpdate, nil
}

func (s *userService) CreateComment(commentDTO dto.CreateCommentRequest) (entity.Komentar, error) {
	var comment entity.Komentar
	copier.Copy(&comment, &commentDTO)

	if comment.IsiKomen == "" {
		return entity.Komentar{}, errors.New("komen tidak boleh kosong")
	}

	berhasilComment, err := s.userRepo.CreateComment(nil, comment)
	if err != nil {
		return entity.Komentar{}, err
	}

	return berhasilComment, nil
}

func (s *userService) GetBlogDetails(id uint64) (entity.Blog, error) {
	berhasilGet, err := s.userRepo.GetDetailBlog(nil, id)
	if err != nil {
		return entity.Blog{}, err
	}

	return berhasilGet, nil
}

func (s *userService) GiveLikeByID(id uint64) error {
	checkErr := s.userRepo.GiveLikeByID(nil, id)
	if checkErr != nil {
		return checkErr
	}
	return nil
}

func (s *userService) DeleteUser(id uint64) error {
	checkDel := s.userRepo.DeleteUser(nil, id)
	if checkDel != nil {
		return checkDel
	}
	return nil
}
