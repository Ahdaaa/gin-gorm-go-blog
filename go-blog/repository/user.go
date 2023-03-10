package repository

import (
	"context"
	"fmt"
	"oprec/go-blog/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	// transaksi database

	// functional
	CreateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
	CreateBlog(ctx context.Context, tx *gorm.DB, blog entity.Blog) (entity.Blog, error)
	GetAllBlog(ctx context.Context, tx *gorm.DB) ([]entity.Blog, error)
	GetBlogByID(tx *gorm.DB, id uint64) (entity.User, error)
	UpdateName(tx *gorm.DB, id uint64, name string) (entity.User, error)
	CreateComment(tx *gorm.DB, comment entity.Komentar) (entity.Komentar, error)
	GetDetailBlog(tx *gorm.DB, id uint64) (entity.Blog, error)
	GiveLikeByID(tx *gorm.DB, id uint64) error
	DeleteUser(tx *gorm.DB, id uint64) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	var err error
	if tx == nil {
		r.db.WithContext(ctx).Debug().Create(&user)
	} else {
		err = tx.WithContext(ctx).Debug().Create(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	var err error
	var user entity.User

	if tx == nil {
		r.db.WithContext(ctx).Debug().Where(("email = ?"), email).Take(&user)
	} else {
		err = tx.WithContext(ctx).Debug().Where(("email = ?"), email).Take(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CreateBlog(ctx context.Context, tx *gorm.DB, blog entity.Blog) (entity.Blog, error) {
	var err error
	if tx == nil {
		tx = r.db.Debug().Create(&blog)
		err = tx.Error
	} else {
		err = tx.Debug().Create(&blog).Error
	}

	if err != nil {
		return entity.Blog{}, err
	}

	return blog, nil
}

func (r *userRepository) GetAllBlog(ctx context.Context, tx *gorm.DB) ([]entity.Blog, error) {
	var blogs []entity.Blog
	var err error
	if tx == nil {
		r.db.WithContext(ctx).Debug().Find(&blogs)
	} else {
		err = tx.Debug().Find(&blogs).Error
	}
	fmt.Println(blogs)
	if err != nil {
		return []entity.Blog{}, err
	}

	return blogs, nil
}

func (r *userRepository) GetBlogByID(tx *gorm.DB, id uint64) (entity.User, error) {
	var user entity.User
	var err error
	if tx == nil {
		r.db.Where("id = ?", id).Preload("ListBlog").Take(&user)
	} else {
		err = tx.Where("id = ?", id).Preload("ListBlog").Take(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil

}

func (r *userRepository) UpdateName(tx *gorm.DB, id uint64, name string) (entity.User, error) {
	var user entity.User
	var err error

	if tx == nil {
		tx = r.db.Model(&user).Where(&entity.User{ID: id}).Update("name", name)
		err = tx.Error
	} else {
		err = tx.Model(&user).Where(&entity.User{ID: id}).Update("name", name).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CreateComment(tx *gorm.DB, comment entity.Komentar) (entity.Komentar, error) {
	var err error
	if tx == nil {
		tx = r.db.Debug().Create(&comment)
		err = tx.Error
		fmt.Println(err)
	} else {
		err = tx.Debug().Create(&comment).Error
	}

	if err != nil {
		return entity.Komentar{}, err
	}

	return comment, nil
}

func (r *userRepository) GetDetailBlog(tx *gorm.DB, id uint64) (entity.Blog, error) {
	var blog entity.Blog
	var err error

	if tx == nil {
		tx = r.db.Where("id = ?", id).Preload("ListKomentar").Take(&blog)
		err = tx.Error
	} else {
		err = tx.Where("id = ?", id).Preload("ListKomentar").Take(&blog).Error
	}

	if err != nil {
		return entity.Blog{}, err
	}

	return blog, nil
}

func (r *userRepository) GiveLikeByID(tx *gorm.DB, id uint64) error {
	var blog entity.Blog
	var err error
	fmt.Println("tes")
	if tx == nil {
		cek := r.db.Model(&blog).Where("id = ?", id).First(&blog)
		err = cek.Error
		fmt.Println(err)
		if err == nil {
			r.db.Debug().Model(&blog).Where(&entity.Blog{ID: id}).Update("jumlah_like", (blog.JumlahLike + 1))
			return nil
		}
		return err
	} else {
		err = tx.Model(&blog).Where(&entity.Blog{ID: id}).Update("jumlah_like", (blog.JumlahLike + 1)).Error
		return err
	}

}

func (r *userRepository) DeleteUser(tx *gorm.DB, id uint64) error {
	var blog []entity.Blog
	var err error

	if tx == nil {
		checkBlog := r.db.Where(&entity.Blog{UserID: id}).Find(&blog)
		var BlogID []uint64
		for _, cek := range blog {
			BlogID = append(BlogID, cek.ID)
		}
		fmt.Println(BlogID)
		if checkBlog.Error == nil {
			r.db.Debug().Model(&entity.Komentar{}).Where("blog_id IN ?", BlogID).Delete(&entity.Komentar{})
			r.db.Where(&entity.Blog{UserID: id}).Delete(&entity.Blog{})
		}

		tx = r.db.Unscoped().Delete(&entity.User{ID: id})
		err = tx.Error

	} else {
		err = tx.Unscoped().Delete(&entity.User{ID: id}).Error
	}

	if err != nil {
		return err
	}

	return nil
}
