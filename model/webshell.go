package model

import (
	"fmt"
	webshell "git.trapti.tech/CPCTF2019/webshell/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

var webShellConn *grpc.ClientConn
var webShellCli webshell.WebShellClient

//InitWebShellCli Initialize Web Shell Client
func InitWebShellCli() error {
	creds, err := credentials.NewClientTLSFromFile("lets-encrypt-x3-cross-signed.pem", os.Getenv("WEBSHELL_GRPC_HOSTNAME"))
	if err != nil {
		return fmt.Errorf("failed to load credentials: %v", err)
	}
	conn, err := grpc.Dial(os.Getenv("WEBSHELL_GRPC_HOSTNAME")+":443", grpc.WithTransportCredentials(creds))
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
