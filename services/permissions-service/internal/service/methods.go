package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spl3g/lab2/internal/proxyproto"
	"github.com/spl3g/lab2/internal/userdb"
	"github.com/spl3g/lab2/services/permissions-service/internal/keycloak"
)

const (
	CentrifugoInternalServerError = 100
	CentrifugoUnauthorized        = 101
	CentrifugoPermissionDenied    = 103
	CentrifugoBadRequest          = 107
)

func (s *Service) fetchKeycloakUser(ctx context.Context, userId uuid.UUID) (userdb.User, error) {
	kcUser, err := s.kcClient.GetUserByID(ctx, userId.String())
	if err != nil {
		return userdb.User{}, err
	}

	user := userdb.User{
		ID:         pgtype.UUID{Valid: true, Bytes: userId},
		Username:   kcUser.Username,
		GivenName:  kcUser.FirstName,
		FamilyName: kcUser.LastName,
		Enabled:    kcUser.Enabled,
	}

	err = s.storage.CreateUser(ctx, UserToCreateUserParams(user))
	if err != nil {
		return userdb.User{}, err
	}

	return user, nil
}

func (s *Service) Subscribe(ctx context.Context, request *proxyproto.SubscribeRequest) (*proxyproto.SubscribeResponse, error) {
	userId, err := uuid.Parse(request.User)
	if err != nil {
		return SubscribeRespondError(CentrifugoBadRequest, "invalid user id")
	}

	user, err := s.storage.GetUserByID(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if errors.Is(err, sql.ErrNoRows) {
		user, err = s.fetchKeycloakUser(ctx, userId)
		if errors.Is(err, keycloak.UserNotFoundErr) {
			return SubscribeRespondError(CentrifugoUnauthorized, "unknown user")
		} else if err != nil {
			log.Println(err)
			return SubscribeRespondError(CentrifugoInternalServerError, "internal server error")
		}
	} else if err != nil {
		log.Println(err)
		return SubscribeRespondError(CentrifugoInternalServerError, "internal server error")
	}

	count, err := s.storage.UserCanSubscribe(ctx, userdb.UserCanSubscribeParams{
		ID:      user.ID,
		Channel: request.Channel,
	})

	if count == 0 {
		return SubscribeRespondError(CentrifugoPermissionDenied, "permission denied")
	}

	return &proxyproto.SubscribeResponse{}, nil
}

func (s *Service) Publish(ctx context.Context, request *proxyproto.PublishRequest) (*proxyproto.PublishResponse, error) {
	userId, err := uuid.Parse(request.User)
	if err != nil {
		return PublishRespondError(CentrifugoBadRequest, "invalid user id")
	}

	user, err := s.storage.GetUserByID(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if errors.Is(err, sql.ErrNoRows) {
		user, err = s.fetchKeycloakUser(ctx, userId)
		if errors.Is(err, keycloak.UserNotFoundErr) {
			return PublishRespondError(CentrifugoUnauthorized, "unknown user")
		} else if err != nil {
			return PublishRespondError(CentrifugoInternalServerError, "internal server error")
		}
	} else if err != nil {
		return PublishRespondError(CentrifugoInternalServerError, "internal server error")
	}

	count, err := s.storage.UserCanPublish(ctx, userdb.UserCanPublishParams{
		ID:      user.ID,
		Channel: request.Channel,
	})

	if count == 0 {
		return PublishRespondError(CentrifugoPermissionDenied, "permission denied")
	}

	return &proxyproto.PublishResponse{}, nil
}
