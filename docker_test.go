package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"testing"
)

var cli *client.Client
var ctx context.Context

func init() {
	cliTemp, err := client.NewClientWithOpts(client.WithHost("tcp://112.125.89.9:2376"))
	// cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli = cliTemp

	ctx = context.Background()
}

func TestRun(t *testing.T) {
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

	/*
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

	*/

}

var CPU_SHARES_PER_CPU = 1024
var MIN_CPU_SHARES = 2
var MIN_MEMORY = 4194304

func TestUpdate(t *testing.T) {

	cpus := 0.1
	memory := 4

	resp, err := cli.ContainerUpdate(
		ctx,
		"f88e3b4ae9bb6602d942931c38beec41010f9b983f8f7b6faa6889c062b55450",
		container.UpdateConfig{
			Resources: container.Resources{
				// 1,048,576‬ = 1M 最小为4M = 4,194,304
				Memory: max64(int64(1048576*memory), int64(MIN_MEMORY)),
				// swap设定为Memory的2倍
				MemorySwap: max64(int64(1048576*memory), int64(MIN_MEMORY)) * 2,
				// 1024的倍数，最小值为2
				CPUShares: max64(int64(cpus*(float64(CPU_SHARES_PER_CPU))), int64(MIN_CPU_SHARES)),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func max64(first int64, args ...int64) int64 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func TestKill(t *testing.T) {

	err := cli.ContainerKill(ctx, "mesos-9dd63492-0419-4c74-a2ce-33d9bce4b0a0", "SIGKILL")
	if err != nil {
		panic(err)
	}

}

func TestInspect(t *testing.T) {

	c, err := cli.ContainerInspect(ctx, "956b12490e78371ecdbe88794139a8460d6244cfbca47f0baf74b63b6d6a027b")
	if err != nil {
		panic(err)
	}

	b, _ := json.Marshal(c)
	fmt.Println(string(b))
}
