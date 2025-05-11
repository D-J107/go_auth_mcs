package main

import (
	"context"
	"fmt"
	auth "go_jwt_mcs/gen/go"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	auth.UnimplementedAuthServer
	DB *DB
}

func NewAuthHandler(db *DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	fmt.Println("get register grpc request with such credentials:", req.Username, req.Password, req.Email)
	if _, _, err := h.DB.GetUserByEmail(ctx, req.GetEmail()); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already registered")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 7)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash error: %v", err)
	}

	if err := h.DB.CreateNewUser(ctx, req.GetUsername(), string(hashedPwd), req.GetEmail()); err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	token := generateToken(req.Username, req.Password, req.Email)

	return &auth.RegisterResponse{Username: req.Username, AccessToken: token}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	fmt.Println("get login grpc request with such credentials:", req.Password, req.Email)
	username, hashedPwd, err := h.DB.GetUserByEmail(ctx, req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.NotFound, "user with such email does not exists")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(req.GetPassword())); err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong password!")
	}

	token := generateToken(username, hashedPwd, req.GetEmail())

	return &auth.LoginResponse{Username: username, AccessToken: token}, nil
}

func generateToken(username, password, email string) string {
	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"email":    email,
		"username": username,
		"password": password,
		"exp":      time.Now().Add(24 * 2 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := t.SignedString([]byte(secret))
	return signedToken
}
