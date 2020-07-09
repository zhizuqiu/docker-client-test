package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"os"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://112.125.89.9:2376"))
	// cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		fmt.Printf("%s %s\n", c.ID[:10], c.Image)
	}

	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "blackicebird/2048",
			Tty:   true,
		},
		&container.HostConfig{
			Resources: container.Resources{},
		},
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case s := <-statusCh:
		fmt.Println(s)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}
