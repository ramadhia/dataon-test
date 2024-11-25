package user

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/ramadhia/dataon-test/internal/config"
	"github.com/ramadhia/dataon-test/internal/entity"
	"github.com/ramadhia/dataon-test/internal/provider"
	"github.com/ramadhia/dataon-test/internal/repository"
	"github.com/ramadhia/dataon-test/internal/usecase"

	"github.com/sirupsen/logrus"
)

type UserImpl struct {
	config   config.Config
	userRepo repository.UserRepository
	group    repository.GroupRepository
}

func NewUser(p *provider.Provider) *UserImpl {
	return &UserImpl{
		config:   p.Config(),
		userRepo: p.UserRepo(),
		group:    p.GroupRepo(),
	}
}

func (u *UserImpl) RegisterUser(ctx context.Context, req usecase.RegisterUserRequest) (*usecase.RegisterUserResponse, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.RegisterUser")

	data := entity.User{
		GroupID:     req.GroupID,
		EmployeeID:  req.EmployeeID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Pin:         u.getmd5(*req.Pin),
	}

	user, err := u.userRepo.GetUser(ctx, repository.GetUserRequest{
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	if user != nil {
		logger.Warning("user already exists")
		return nil, errors.New("user already exists")
	}

	logger.Debugf("about to register user data")
	register, err := u.userRepo.Register(ctx, data)
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	logger.Info("Register successfully")
	return &usecase.RegisterUserResponse{
		UserID:      register.ID,
		GroupID:     register.GroupID,
		EmployeeID:  register.EmployeeID,
		Name:        register.Name,
		PhoneNumber: register.PhoneNumber,
		CreatedDate: register.CreatedDate,
	}, nil
}

//func (u *UserImpl) LoginUser(ctx context.Context, req usecase.LoginUserRequest) (*usecase.LoginUserResponse, error) {
//	logger := logrus.WithField("method", "usecase.UserImpl.LoginUser")
//
//	user, err := u.userRepo.GetUser(ctx, repository.GetUserRequest{
//		PhoneNumber: req.PhoneNumber,
//		Pin:         u.getmd5(*req.Pin),
//	})
//	if err != nil {
//		logger.Warning(err.Error())
//		return nil, err
//	}
//	if user == nil {
//		logger.Warning("user not found")
//		return nil, errors.New("phone Number and PIN doesn’t match")
//	}
//
//	accessToken, refreshToken, err := lib.GenerateTokens(model.Claim{
//		ID: *user.ID,
//	})
//	if err != nil {
//		logger.Warning("error when generate the tokens")
//		return nil, err
//	}
//
//	logger.Info("Login successfully")
//	return &usecase.LoginUserResponse{
//		AccessToken:  &accessToken,
//		RefreshToken: &refreshToken,
//	}, nil
//}

func (u *UserImpl) UpdateUser(ctx context.Context, data entity.User) (res *usecase.UpdateProfileResponse, err error) {
	logger := logrus.WithField("method", "usecase.UserImpl.UpdateUser")

	update, err := u.userRepo.UpdateUser(ctx, data)
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	res = &usecase.UpdateProfileResponse{
		UserID:      update.ID,
		GroupID:     update.GroupID,
		EmployeeID:  update.EmployeeID,
		Name:        update.Name,
		PhoneNumber: update.PhoneNumber,
	}

	return
}

func (u *UserImpl) GetUser(ctx context.Context, user entity.User) (*entity.User, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.GetUser")

	res, err := u.userRepo.GetUser(ctx, repository.GetUserRequest{
		UserID: user.ID,
	})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}
	if res == nil {
		logger.Warning("user not found")
		return nil, errors.New("user not found")
	}

	return res, nil
}

func (u *UserImpl) FetchUser(ctx context.Context, req usecase.FetchUserRequest) ([]*entity.User, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.GetUser")

	res, err := u.userRepo.FetchUser(ctx, repository.FetchUserRequest{
		OrganizationID: req.OrganizationID,
	})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}
	if res == nil {
		logger.Warning("user not found")
		return nil, errors.New("user not found")
	}

	return res, nil
}

func (u *UserImpl) DeleteUser(ctx context.Context, userId string) (bool, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.DeleteUser")

	res, err := u.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		logger.Warning(err.Error())
		return false, err
	}

	return res, nil
}

func (u *UserImpl) getmd5(s string) *string {
	h := md5.New()
	h.Write([]byte(s))
	r := hex.EncodeToString(h.Sum(nil))
	return &r

}
