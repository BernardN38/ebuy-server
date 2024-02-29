package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/BernardN38/ebuy-server/authentication-service/messaging"
	users_sql "github.com/BernardN38/ebuy-server/authentication-service/sqlc/users"
)

type AuthService struct {
	userQueries    queries
	passwordHasher *passwordHasher
	messageEmitter messaging.MessageEmitter
}

// for testing purposes
type queries interface {
	CreateUser(context.Context, users_sql.CreateUserParams) error
	GetAll(context.Context) ([]users_sql.User, error)
	GetUser(context.Context, int32) (users_sql.GetUserRow, error)
	GetUserPassword(context.Context, string) (users_sql.GetUserPasswordRow, error)
}

func New(db *sql.DB, me messaging.MessageEmitter) *AuthService {
	userQueries := users_sql.New(db)
	ph := NewPasswordHasher()
	return &AuthService{
		userQueries:    userQueries,
		passwordHasher: ph,
		messageEmitter: me,
	}
}

func (a *AuthService) CreateUser(ctx context.Context, createUserParams CreateUserParams) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*2000)
	defer cancel()
	successCh := make(chan bool)
	errCh := make(chan error)
	go func() {
		encodedPassword, err := a.passwordHasher.CreateEncodedHash(createUserParams.Password)
		if err != nil {
			errCh <- err
			return
		}
		err = a.userQueries.CreateUser(timeoutCtx, users_sql.CreateUserParams{
			Username:        createUserParams.Username,
			Email:           createUserParams.Email,
			EncodedPassword: encodedPassword,
		})
		if err != nil {
			errCh <- err
			return
		}
		createUserMsg, err := json.Marshal(messaging.CreateUserMessage{
			FirstName: createUserParams.FirstName,
			LastName:  createUserParams.LastName,
			Username:  createUserParams.Username,
			Email:     createUserParams.Email,
			Dob:       createUserParams.DOB.Format(time.RFC3339),
		})
		if err == nil {
			err := a.messageEmitter.SendMessage(timeoutCtx, createUserMsg, "user_events", "user.created", "user.created")
			if err != nil {
				log.Println(err)
			}
		}
		// a.messageEmitter
		successCh <- true
	}()

	select {
	case <-successCh:
		return nil
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

func (a *AuthService) GetUser(ctx context.Context, userId int) (users_sql.GetUserRow, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	successCh := make(chan users_sql.GetUserRow)
	errCh := make(chan error)
	go func() {
		user, err := a.userQueries.GetUser(ctx, int32(userId))
		if err != nil {
			errCh <- err
			return
		}
		successCh <- user
	}()
	select {
	case user := <-successCh:
		return user, nil
	case err := <-errCh:
		return users_sql.GetUserRow{}, err
	case <-timeoutCtx.Done():
		return users_sql.GetUserRow{}, timeoutCtx.Err()
	}
}

func (a *AuthService) VerifyUser(ctx context.Context, email string, password string) (int, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*2000)
	defer cancel()
	successCh := make(chan bool)
	errCh := make(chan error)
	userId := 0
	go func() {
		encodedPassword, err := a.userQueries.GetUserPassword(timeoutCtx, email)
		if err != nil {
			log.Println("could not get user from db")
			errCh <- err
			return
		}
		verified, err := a.passwordHasher.VerifyPassword(encodedPassword.EncodedPassword, password)
		if err != nil {
			log.Println(err)
			errCh <- err
			return
		}
		if !verified {
			log.Println("password incorrect")
			errCh <- errors.New("password is not correct")
			return
		}
		userId = int(encodedPassword.ID)
		successCh <- true
	}()
	select {
	case <-successCh:
		return userId, nil
	case err := <-errCh:
		return 0, err
	case <-timeoutCtx.Done():
		return 0, timeoutCtx.Err()
	}
}
