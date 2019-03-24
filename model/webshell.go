package model

import (
	"fmt"
	"os"

	webshell "git.trapti.tech/CPCTF2019/webshell/rpc"
	"google.golang.org/grpc"
)

var webShellConn *grpc.ClientConn
var webShellCli webshell.WebShellClient

//InitWebShellCli Initialize Web Shell Client
func InitWebShellCli() error {
	conn, err := grpc.Dial(os.Getenv("WEBSHELL_GRPC_TARGET"), grpc.WithInsecure())
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
