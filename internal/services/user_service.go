package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/qerdcv/qerdcv/internal/config"
	"github.com/qerdcv/qerdcv/internal/repositories"
	"github.com/qerdcv/qerdcv/pkg/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credential")

	ErrInvalidToken    = errors.New("invalid token")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errors.New("session not found")
)

const (
	sessionTokenTTL = 7 * 24 * time.Hour
)

type UserService struct {
	repo        *repositories.UserRepo
	tokenSecret []byte
}

func NewUserService(repo *repositories.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt generate from password: %w", err)
	}

	if err = s.repo.CreateUser(ctx, username, string(hashedPwd)); err != nil {
		if errors.Is(err, repositories.ErrUniqueConstraint) {
			return ErrUserAlreadyExists
		}

		return fmt.Errorf("repo create user: %w", err)
	}

	return nil
}

func (s *UserService) AuthorizeUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf("repo get user by username: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredential
		}

		return "", fmt.Errorf("bcrypt compare hash and password")
	}

	var strToken string
	if err = s.repo.WithTX(
		ctx,
		func(r *repositories.UserRepo) error {
			sessionExpiresAt := time.Now().Add(sessionTokenTTL)
			session, txErr := r.CreateUserSession(ctx, user, sessionExpiresAt)
			if txErr != nil {
				return fmt.Errorf("repo create user session: %w", err)
			}

			strToken, txErr = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
				Issuer:    config.AppName,
				Subject:   strconv.Itoa(session.UserID),
				ExpiresAt: jwt.NewNumericDate(sessionExpiresAt),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        strconv.Itoa(session.ID),
			}).SignedString(s.tokenSecret)
			if txErr != nil {
				return fmt.Errorf("token signed string: %w", txErr)
			}

			return nil
		},
	); err != nil {
		return "", fmt.Errorf("with tx: %w", err)
	}

	return strToken, nil
}

func (s *UserService) VerifySession(ctx context.Context, sessionToken string) (domain.UserSession, error) {
	t, err := jwt.Parse(sessionToken, func(_ *jwt.Token) (interface{}, error) {
		return s.tokenSecret, nil
	})
	if err != nil {
		return domain.UserSession{}, errors.Join(err, ErrInvalidToken)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return domain.UserSession{}, ErrInvalidToken
	}

	sID, err := strconv.Atoi(claims["jti"].(string))
	if err != nil {
		return domain.UserSession{}, errors.Join(err, ErrInvalidToken)
	}

	uID, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return domain.UserSession{}, errors.Join(err, ErrInvalidToken)
	}

	userSession, err := s.repo.GetUserSession(ctx, sID, uID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return domain.UserSession{}, ErrSessionNotFound
		}

		return domain.UserSession{}, fmt.Errorf("get user session: %w", err)
	}

	if time.Now().After(userSession.ExpiresAt) {
		return domain.UserSession{}, errors.Join(ErrSessionExpired, s.repo.DeleteUserSession(ctx, sID, uID))
	}

	return userSession, nil
}
