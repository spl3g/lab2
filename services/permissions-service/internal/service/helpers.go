package service

import (
	"github.com/spl3g/lab2/internal/proxyproto"
	"github.com/spl3g/lab2/internal/userdb"
)

func ConnectRespondError(code uint32, msg string) (*proxyproto.ConnectResponse, error) {
	return &proxyproto.ConnectResponse{
		Error: &proxyproto.Error{
			Code:    code,
			Message: msg,
		},
	}, nil
}

func SubscribeRespondError(code uint32, msg string) (*proxyproto.SubscribeResponse, error) {
	return &proxyproto.SubscribeResponse{
		Error: &proxyproto.Error{
			Code:    code,
			Message: msg,
		},
	}, nil
}

func PublishRespondError(code uint32, msg string) (*proxyproto.PublishResponse, error) {
	return &proxyproto.PublishResponse{
		Error: &proxyproto.Error{
			Code:    code,
			Message: msg,
		},
	}, nil
}

func UserToCreateUserParams(user userdb.User) userdb.CreateUserParams {
	return userdb.CreateUserParams{
		ID:         user.ID,
		Username:   user.Username,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Enabled:    user.Enabled,
	}
}
