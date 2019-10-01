package main

import(
	"log"
	"runtime"
	"fmt"

	"user-service/app/config"
	"user-service/app/model"
	"user-service/app/platform/nats"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envFlag    = "env"
	defaultEnv = "dev"
	serverPort = "server.port"
	natsServer = "nats.server"
)

func main() {

	// Get environment flag
	env := pflag.String(envFlag, defaultEnv, "environment config value to use")
	pflag.Parse()

	if err := config.LoadConfiguration(*env); err != nil {
		checkErr(err)
	}

	// Create new NATS server connection
	nc, err := natsclient.NewNATSServerConnection(viper.GetString(natsServer))
	checkErr(err)

	// Subscribe to user.create via channel
	sch := make(chan *model.User)
	nc.BindRecvChan("user.create", sch)
	u := <-sch
	fmt.Printf("Received a user: %+v\n", u)

	// Publish to user.create.completed via channel
	pch := make(chan *model.User)
	nc.BindSendChan("user.create.completed", pch)
	// Save User
	u.UserId = 1
	pch <- u

	runtime.Goexit()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}