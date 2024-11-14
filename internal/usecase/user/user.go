package user

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ramadhia/mnc-test/internal/config"
	"github.com/ramadhia/mnc-test/internal/lib"
	"github.com/ramadhia/mnc-test/internal/model"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/repository"
	"github.com/ramadhia/mnc-test/internal/usecase"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type UserImpl struct {
	config      config.Config
	userRepo    repository.UserRepository
	balanceRepo repository.BalanceRepository
}

func NewUser(p *provider.Provider) *UserImpl {
	return &UserImpl{
		config:      p.Config(),
		userRepo:    p.UserRepo(),
		balanceRepo: p.BalanceRepo(),
	}
}

func (u *UserImpl) RegisterUser(ctx context.Context, req usecase.RegisterUserRequest) (*usecase.RegisterUserResponse, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.RegisterUser")

	data := model.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
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

	logger.Info("about to create a balance data if register is done")
	_, err = u.balanceRepo.UpsertBalance(ctx, model.Balance{
		UserID: register.ID,
		Amount: lo.ToPtr(float64(0)),
	})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	logger.Info("Register successfully")
	return &usecase.RegisterUserResponse{
		UserID:      register.ID,
		FirstName:   register.FirstName,
		LastName:    register.LastName,
		PhoneNumber: register.PhoneNumber,
		Address:     register.Address,
		CreatedDate: register.CreatedDate,
	}, nil
}

func (u *UserImpl) LoginUser(ctx context.Context, req usecase.LoginUserRequest) (*usecase.LoginUserResponse, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.LoginUser")

	user, err := u.userRepo.GetUser(ctx, repository.GetUserRequest{
		PhoneNumber: req.PhoneNumber,
		Pin:         u.getmd5(*req.Pin),
	})
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}
	if user == nil {
		logger.Warning("user not found")
		return nil, errors.New("phone Number and PIN doesn’t match")
	}

	accessToken, refreshToken, err := lib.GenerateTokens(model.Claim{
		ID: *user.ID,
	})
	if err != nil {
		logger.Warning("error when generate the tokens")
		return nil, err
	}

	logger.Info("Login successfully")
	return &usecase.LoginUserResponse{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}

func (u *UserImpl) UpdateUser(ctx context.Context, claim model.Claim, req usecase.UpdateProfileRequest) (res *usecase.UpdateProfileResponse, err error) {
	logger := logrus.WithField("method", "usecase.UserImpl.UpdateUser")

	data := model.User{
		ID:          &claim.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Address:     req.Address,
		UpdatedDate: lo.ToPtr(time.Now()),
	}
	update, err := u.userRepo.UpdateUser(ctx, data)
	if err != nil {
		logger.Warning(err.Error())
		return nil, err
	}

	res = &usecase.UpdateProfileResponse{
		UserID:      update.ID,
		FirstName:   update.FirstName,
		LastName:    update.LastName,
		PhoneNumber: update.PhoneNumber,
		Address:     update.Address,
		UpdatedDate: update.UpdatedDate,
	}

	return
}

func (u *UserImpl) GetUser(ctx context.Context, user model.Claim) (*model.User, error) {
	logger := logrus.WithField("method", "usecase.UserImpl.GetUser")

	res, err := u.userRepo.GetUser(ctx, repository.GetUserRequest{
		UserID: &user.ID,
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

func (u *UserImpl) getmd5(s string) *string {
	h := md5.New()
	h.Write([]byte(s))
	r := hex.EncodeToString(h.Sum(nil))
	return &r

}
