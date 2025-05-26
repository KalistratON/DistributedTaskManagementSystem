package repository

import (
	"context"
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-redis/redis/v8"
)

const TOKEN_EXPIRATION_TIME = 15 * time.Minute
const ID_TOKEN_HASH = "user_id:token"

func generateJwtToken(base string) (string, error) {
	srtKey := os.Getenv("JWT_TOKEN_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": base,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(TOKEN_EXPIRATION_TIME + time.Minute).Unix(),
	})
	return token.SignedString([]byte(srtKey))
}

type authRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewAuthRepository(rdb *redis.Client, ctx context.Context) (AuthRepository, error) {
	if rdb == nil {
		return nil, errors.NilRdb{}
	}
	return &authRepository{
		rdb: rdb,
		ctx: ctx,
	}, nil
}

var _ AuthRepository = &authRepository{}

func (r *authRepository) Create(auth *domain.Auth) (*domain.Auth, error) {
	token, err := generateJwtToken(auth.Id)
	if err != nil {
		log.Printf("error while generate jwt")
		return nil, err
	}

	if err = r.rdb.Set(r.ctx, token, auth.Id, TOKEN_EXPIRATION_TIME).Err(); err != nil {
		log.Printf("error while create auth info")
		return nil, err
	}

	res := r.rdb.Get(r.ctx, token)
	if err = res.Err(); err != nil {
		log.Printf("error while trying check-get result")
		return nil, err
	}

	if err = r.rdb.HSet(r.ctx, ID_TOKEN_HASH, auth.Id, token).Err(); err != nil {
		log.Printf("error while create auth info")
		return nil, err
	}

	return &domain.Auth{
		Id:    res.Val(),
		Token: token,
	}, nil
}

func (r *authRepository) Update(token string) (*domain.Auth, error) {
	result, err := r.Get(token)
	if err != nil {
		log.Printf("cash with such token does not exist")
	}

	if err = r.rdb.Expire(r.ctx, token, TOKEN_EXPIRATION_TIME).Err(); err != nil {
		log.Printf("can't update TTL for %s: %v", token, err)
		return nil, err
	}
	return result, nil
}

func (r *authRepository) Get(token string) (*domain.Auth, error) {
	token = strings.TrimSpace(token)
	if len(token) < 1 {
		return nil, errors.EmptyToken{}
	}

	res := r.rdb.Get(r.ctx, token)
	if err := res.Err(); err != nil {
		log.Printf("error while trying to get cash by token: %v", err)
		return nil, err
	}

	return &domain.Auth{
		Id:    res.Val(),
		Token: token,
	}, nil
}

func (r *authRepository) GetToken(id string) (*domain.Auth, error) {
	id = strings.TrimSpace(id)
	if len(id) < 1 {
		return nil, errors.EmptyId{}
	}

	res := r.rdb.HGet(r.ctx, ID_TOKEN_HASH, id)
	if err := res.Err(); err != nil {
		log.Printf("error while trying to get cash by id: %v", err)
		return nil, err
	}

	result, err := r.Get(res.Val())
	if err != nil {
		log.Printf("token was expired: %v", err)
		return nil, errors.ExpiredToken{}
	}

	return result, nil
}

func (r *authRepository) Delete(token string) (*domain.Auth, error) {
	token = strings.TrimSpace(token)
	if len(token) < 1 {
		return nil, errors.EmptyToken{}
	}

	target, err := r.Get(token)
	if err != nil {
		log.Printf("can't find auth by token: %s:%v", token, err)
		return nil, err
	}

	res := r.rdb.Del(r.ctx, token)
	if err = res.Err(); err != nil {
		log.Printf("error while trying delete token %s:%v", token, err)
		return nil, err
	}

	if res.Val() < 1 {
		log.Printf("no data was deleted")
		return nil, errors.EmptyDelete{}
	}

	if err = r.rdb.HDel(r.ctx, ID_TOKEN_HASH, target.Id).Err(); err != nil {
		log.Printf("error while create auth info")
		return nil, err
	}

	return target, nil
}
