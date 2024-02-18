package service

import (
	"errors"
	"net/http"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/repository"
)

type ConnectUserService struct {
	wakatimeUserRepo repository.WakatimeUser
	slackUserRepo    repository.SlackUser
	userRepo         repository.User
}

func NewConnectUserService(wakatimeUserRepo repository.WakatimeUser, slackUserRepo repository.SlackUser, userRepo repository.User) ConnectUserService {
	return ConnectUserService{
		wakatimeUserRepo: wakatimeUserRepo,
		slackUserRepo:    slackUserRepo,
		userRepo:         userRepo,
	}
}

func (s ConnectUserService) CreateUser(slackUserID, wakatimeUserID, slackToken string, wakatimeOauth2Client *http.Client) (model.User, error) {
	slackUser, err := s.slackUserRepo.GetUser(slackUserID, slackToken)
	if err != nil {
		return model.User{}, err
	}

	wakatimeUser, err := s.wakatimeUserRepo.GetUser(wakatimeUserID, wakatimeOauth2Client)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Username:       slackUser.Name,
		SlackUserID:    slackUser.ID,
		WakatimeUserID: wakatimeUser.ID,
	}

	if err = s.userRepo.AddUser(user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s ConnectUserService) DisconnectWakatimeUser(userID, wakatimeUserID string) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	if user.WakatimeUserID != wakatimeUserID {
		return ErrorUserNotFound
	}

	user.WakatimeUserID = ""
	if err = s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s ConnectUserService) DisconnectSlackUser(userID, slackUserID string) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	if user.SlackUserID != slackUserID {
		return ErrorUserNotFound
	}

	user.SlackUserID = ""
	if err = s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

var ErrorUserNotFound = errors.New("user not found")
