package entity

import "sync"

type Client struct {
	username string
	hostname string
}

type SafeClientList struct {
	mu      sync.Mutex
	members []Client
}

func (c *SafeClientList) Add(user string, host string) {
	c.mu.Lock()
	c.members = append(c.members, Client{
		username: user,
		hostname: host,
	})
	c.mu.Unlock()
}

func (c *SafeClientList) HowMany() int {
	return len(c.members)
}

// func (c *ClientList) Send(stream pb.Registry_SignServer) error {
// 	c.mu.Lock()

// 	for _, client := range c.members {
// 		payload := &pb.ClientInfo{
// 			Username: client.username,
// 			Hostname: client.hostname,
// 		}

// 		if err := stream.Send(payload); err != nil {
// 			return err
// 		}
// 	}

// 	c.mu.Unlock()

// 	return nil
// }
