package main

import (
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strconv"
	//"golang.org/x/crypto/acme/autocert"
	//"golang.org/x/crypto/acme"
	//"google.golang.org/genproto/googleapis/monitoring/v3"
	//"crypto/tls"
	"github.com/krzysztofromanowski94/YACS5e-cloud/utils"
)

var (
	useTls, parseErr = strconv.ParseBool("{{ use_tls }}")
	tlsCertFile      = "{{ tls_cert_file }}"
	tlsKeyFile       = "{{ tls_key_file }}"
	serverAddress    = ":{{ server_port }}"
	serverOpts       []grpc.ServerOption

	exitProgramChannel chan bool = make(chan bool, 1)
)

func main() {
	log.Println("Server starting...")

	//lis := autocert.NewListener("ptp-thingers.pl")
	//
	//

	if parseErr != nil {
		log.Fatal("Error parsing use_tls ansible variable. ERROR:", parseErr)
	}
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal("Failed to create listen service. ERROR: ", err)
	}

	if useTls {
		//cred, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
		//if err != nil {
		//	log.Fatal("Failed to generate credentials. ERROR: ", err)
		//}
		//serverOpts = append(serverOpts, grpc.Creds(cred))
		/*serverOpts = append(serverOpts,*/

		//ttt := &tls.Config{GetCertificate: asd.GetCertificate}

		tls := utils.GetTLS("ptp-thingers.pl", "/root/secret/certs")
		creds := credentials.NewTLS(tls)

		serverOpts = append(serverOpts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterYACS5EServer(grpcServer, newServer())

	go func() {
		log.Println("Server Started")
		grpcServer.Serve(lis)
	}()

	<-exitProgramChannel
	return
}
