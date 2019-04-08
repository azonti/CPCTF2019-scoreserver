package model

import (
	"fmt"
	"os"

	webshell "git.trapti.tech/CPCTF2019/webshell/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var webShellConn *grpc.ClientConn
var webShellCli webshell.WebShellClient

//InitWebShellCli Initialize Web Shell Client
func InitWebShellCli() error {
	hostname := os.Getenv("WEBSHELL_GRPC_HOSTNAME")
	port := os.Getenv("WEBSHELL_GRPC_PORT")
	creds, err := credentials.NewClientTLSFromFile("lets-encrypt-x3-cross-signed.pem", hostname)
	if err != nil {
		return fmt.Errorf("failed to NewClientTLS: %v", err)
	}
	conn, err := grpc.Dial(hostname+":"+port, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("failed to connect webshell grpc: %v", err)
	}
	webShellConn = conn
	webShellCli = webshell.NewWebShellClient(conn)
	return nil
}

//TermWebShellCli Terminate Web Shell Client
func TermWebShellCli() {
	webShellConn.Close()
}
