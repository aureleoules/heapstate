package common

import (
	"fmt"

	"github.com/docker/docker/client"
)

var DockerClient *client.Client

func init() {
	var err error
	DockerClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("Initialized Docker client")
}
