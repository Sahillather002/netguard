package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/securecloud/auth-service/proto"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Printf("🔐 Auth Service starting on %s", port)

	// Start server in goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Auth Service...")
	s.GracefulStop()
	log.Println("✅ Auth Service exited")
}

// Login authenticates a user and returns JWT tokens
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// TODO: Validate credentials against database
	// TODO: Hash password comparison
	// TODO: Generate JWT access and refresh tokens

	log.Printf("Login attempt for email: %s", req.Email)

	// Mock response
	return &pb.LoginResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEyMyIsImV4cCI6MTYzOTk5OTk5OX0.signature",
		RefreshToken: "refresh_token_here",
		ExpiresIn:    3600,
		UserId:       "user_123",
	}, nil
}

// Register creates a new user account
func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// TODO: Validate email format
	// TODO: Check if email already exists
	// TODO: Hash password with bcrypt
	// TODO: Create user in database
	// TODO: Send verification email

	log.Printf("Registration attempt for email: %s", req.Email)

	return &pb.RegisterResponse{
		UserId:  "user_new",
		Message: "User registered successfully",
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// TODO: Validate refresh token
	// TODO: Check if token is blacklisted
	// TODO: Generate new access token

	log.Printf("Token refresh attempt")

	return &pb.RefreshTokenResponse{
		AccessToken: "new_access_token",
		ExpiresIn:   3600,
	}, nil
}

// ValidateToken validates a JWT token and returns user info
func (s *server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// TODO: Parse JWT token
	// TODO: Verify signature
	// TODO: Check expiration
	// TODO: Check if token is blacklisted

	log.Printf("Token validation attempt")

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: "user_123",
		Role:   "admin",
	}, nil
}

// Logout invalidates a user's tokens
func (s *server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// TODO: Add token to blacklist in Redis
	// TODO: Set expiration time for blacklist entry

	log.Printf("Logout attempt for user: %s", req.UserId)

	return &pb.LogoutResponse{
		Message: "Logged out successfully",
	}, nil
}

// ChangePassword changes a user's password
func (s *server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// TODO: Verify old password
	// TODO: Hash new password
	// TODO: Update password in database
	// TODO: Invalidate all existing tokens

	log.Printf("Password change attempt for user: %s", req.UserId)

	return &pb.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}

// ResetPassword initiates password reset process
func (s *server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// TODO: Generate reset token
	// TODO: Store reset token in Redis with expiration
	// TODO: Send reset email

	log.Printf("Password reset requested for email: %s", req.Email)

	return &pb.ResetPasswordResponse{
		Message: "Password reset email sent",
	}, nil
}
