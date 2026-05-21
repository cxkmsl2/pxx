package service

import (
	"context"
	"errors"

	"pxx/internal/middleware"
	"pxx/internal/model"
	"pxx/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type LoginResult struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

func (s *UserService) Login(ctx context.Context, openID string) (*LoginResult, error) {
	var user model.User
	err := s.db.WithContext(ctx).Where("open_id = ?", openID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			OpenID:   openID,
			Nickname: "校园用户" + utils.MD5(openID)[:6],
			Role:     "user",
		}
		if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: &user, Token: token}, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (s *UserService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}
