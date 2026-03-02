package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	emailservice "github.com/marisasha/email-scheduler/internal/email"
	"github.com/marisasha/email-scheduler/internal/models"
	"github.com/marisasha/email-scheduler/internal/repository"
)

const (
	salt       = "vfzgz25f2sdf4gsf.fsg246ydhd.gh3ilof10"
	signingKey = "fnhj52..254nfslmnl8hfsvbnjs.2fjisg"
	tokenTTL   = 12 * time.Hour
	queueName  = "email_verification"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repos      repository.Authorization
	emailQueue emailservice.Publisher
}

func NewAuthService(repos repository.Authorization, emailQueue emailservice.Publisher) *AuthService {
	return &AuthService{
		repos:      repos,
		emailQueue: emailQueue,
	}
}

func (s *AuthService) CreateUser(user *models.User) error {
	user.Password = *generatePasswordHash(user.Password)
	return s.repos.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password *string) (string, error) {
	user, err := s.repos.GetUser(username, generatePasswordHash(*password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accesToken *string) (int, error) {
	token, err := jwt.ParseWithClaims(*accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func (s *AuthService) SendEmailVerification(userId *int) error {

	userEmail, err := s.repos.GetUserEmail(userId)
	if err != nil {
		return err
	}
	token := generateToken()

	err = s.repos.CreateEmailVerificationToken(userId, &token)
	if err != nil {
		return err
	}

	job := models.EmailJob{
		To:      userEmail,
		Subject: "Подтвердите Email",
		Body:    fmt.Sprintf("Перейдите по ссылке: \nlocalhost:8000/auth/verify-email/check?token=%s", token),
	}

	return s.emailQueue.PublishEmail(job, queueName)
}

func (s *AuthService) CheckEmailVerification(token *string) error {

	emailVerification, err := s.repos.CheckVerificationToken(token)
	if err != nil {
		return fmt.Errorf("invalid token : %s ", err)
	}

	if emailVerification.CreatedAt.Before(time.Now().Add(-10 * time.Minute)) {
		return fmt.Errorf("token is expired")
	}

	err = s.repos.ChangeEmailVerificationStatus(&emailVerification.UserId)
	if err != nil {
		return err
	}

	return nil

}

func generatePasswordHash(password string) *string {
	hash := sha1.New()
	hash.Write([]byte(password))
	passwordHash := fmt.Sprintf("%x", hash.Sum([]byte(salt)))
	return &passwordHash
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
