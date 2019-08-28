package cloudprober

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	pb "github.com/google/cloudprober/prober/proto"
	configpb "github.com/google/cloudprober/probes/proto"
	"google.golang.org/grpc"
)

type Client struct {
	Server string
	client pb.CloudproberClient
}

func New(server string) *Client {
	return &Client{
		Server: server,
	}
}

func (c *Client) Getclient() (client pb.CloudproberClient, err error) {
	if c.client != nil {
		client = c.client
		return
	}
	conn, err := grpc.Dial(c.Server, grpc.WithInsecure())
	if err != nil {
		err = fmt.Errorf("grpc dial err: %v", err)
		return
	}
	c.client = pb.NewCloudproberClient(conn)
	return c.client, nil
}

func (c *Client) RemoveProbe(name string) (err error) {
	client, err := c.Getclient()
	if err != nil {
		err = fmt.Errorf("get client err: %v", err)
		return
	}
	_, err = client.RemoveProbe(context.Background(), &pb.RemoveProbeRequest{ProbeName: &name})
	if err != nil {
		err = fmt.Errorf("remove probe %v err: %v", name, err)
		return
	}
	return
}

func (c *Client) AddProbe(cfg *configpb.ProbeDef) (err error) {
	if cfg == nil {
		err = fmt.Errorf("empty probedef")
		return
	}
	client, err := c.Getclient()
	if err != nil {
		err = fmt.Errorf("get client err: %v", err)
		return
	}
	_, err = client.AddProbe(context.Background(), &pb.AddProbeRequest{ProbeConfig: cfg})
	if err != nil {
		err = fmt.Errorf("add probe err: %v", err)
		return
	}
	return
}

func (c *Client) AddProbeWithFile(file string) (err error) {
	client, err := c.Getclient()
	if err != nil {
		err = fmt.Errorf("get client err: %v", err)
		return
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("Failed to read the config file: %v", err)
		return
	}

	// glog.Infof("Read probe config: %s", string(b))

	cfg := &configpb.ProbeDef{}
	if err = proto.UnmarshalText(string(b), cfg); err != nil {
		err = fmt.Errorf("unmarshal file: %v, err: %v", file, err)
		return
	}

	_, err = client.AddProbe(context.Background(), &pb.AddProbeRequest{ProbeConfig: cfg})
	if err != nil {
		err = fmt.Errorf("add probe err: %v", err)
		return
	}
	return
}
