package frontend

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetAvailablePorts(imageName string) ([]uint16, error) {
	var comunigoPeers []types.Container
	var availablePorts []uint16
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	for _, cont := range containers {
		if cont.Image == imageName {
			comunigoPeers = append(comunigoPeers, cont)
		}
	}

	for _, peer := range comunigoPeers {
		availablePorts = append(availablePorts, peer.Ports[0].PublicPort)
	}

	sort.Slice(availablePorts, func(i, j int) bool {
		return availablePorts[i] < availablePorts[j]
	})

	return availablePorts, nil
}
