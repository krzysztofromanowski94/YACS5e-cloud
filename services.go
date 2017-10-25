package main

import (
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"log"
)

var (
	sql_hostname = "{{ sql_hostname }}"
	sql_user     = "{{ sql_user }}"
	sql_password = "{{ sql_password }}"
)

type YACS5eServer struct {
}

func newServer() *YACS5eServer {
	return new(YACS5eServer)

}

// rpc Registration (User) returns (Empty)
// ERROR CODES:
// 100: UNKNOWN ERROR
// 101: INVALID LOGIN
// 102: INVALID PASSWORD
// 103: USER EXISTS
func (server *YACS5eServer) Registration(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	log.Println("Registration Context: ", ctx)
	return &pb.Empty{}, status.Errorf(0, "Got user: ", user.Login)
}

// rpc Login (User) returns (Empty)
// ERROR CODES:
// 110: UNKNOWN ERROR
// 111: INVALID LOGIN
// 112: INVALID PASSWORD
// 113: USER EXISTS
func (server *YACS5eServer) Login(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	log.Println("Login Context: ", ctx)
	return &pb.Empty{}, status.Errorf(0, "User ", user.Login, " may exists. I don't know yet ;x")
}
