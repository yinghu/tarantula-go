package util

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/auth/credentials"
	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

type GcpComputeEngine struct {
	ServiceAccount string
	ProjectId      string
	Zone           string
	client         *compute.InstancesClient
}

func (g *GcpComputeEngine) Auth() error {
	data, err := os.ReadFile(g.ServiceAccount)
	if err != nil {
		return err
	}
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:          []string{"https://www.googleapis.com/auth/compute"},
		CredentialsJSON: data,
	})
	if err != nil {
		return err
	}
	client, err := compute.NewInstancesRESTClient(context.Background(), option.WithAuthCredentials(creds))
	if err != nil {
		return err
	}
	g.client = client
	return nil
}

func (g *GcpComputeEngine) Close() error {
	return g.client.Close()
}

func (g *GcpComputeEngine) List() error {
	req := &computepb.ListInstancesRequest{
		Project: g.ProjectId,
		Zone:    g.Zone,
	}
	it := g.client.List(context.Background(), req)
	for {
		instance, err := it.Next()
		if err != nil {
			//fmt.Printf("iter err %s", err.Error())
			break
		}
		fmt.Printf("Name : %s\n", instance.GetName())
		fmt.Printf("Status : %s\n", instance.GetStatus())
		fmt.Printf("Inter IP : %s\n", instance.GetNetworkInterfaces()[0].GetNetworkIP())
		fmt.Printf("Public IP : %s\n", instance.GetNetworkInterfaces()[0].AccessConfigs[0].GetNatIP())
	}
	return nil
}

func (g *GcpComputeEngine) Insert(name string) error {
	req := &computepb.InsertInstanceRequest{
		Project: g.ProjectId,
		Zone:    g.Zone,
		InstanceResource: &computepb.Instance{
			Name:        proto.String(name),
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/e2-micro", g.Zone)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
					AccessConfigs: []*computepb.AccessConfig{
						{
							Type:        proto.String("ONE_TO_ONE_NAT"),
							Name:        proto.String("External NAT"),
							NetworkTier: proto.String("STANDARD"),
						},
					},
				},
			},
			Tags: &computepb.Tags{
				Items: []string{"http-server", "https-server"}, // Matching tag for firewall
			},
			Scheduling: &computepb.Scheduling{
				ProvisioningModel: proto.String(computepb.Scheduling_SPOT.String()),
			},
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						SourceImage: proto.String("projects/cos-cloud/global/images/family/cos-113-lts"),
						DiskSizeGb:  proto.Int64(10),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
		},
	}
	opt, err := g.client.Insert(context.Background(), req)
	if err != nil {
		return err
	}
	return opt.Wait(context.Background())
}

func (g *GcpComputeEngine) Delete(name string) error {
	req := &computepb.DeleteInstanceRequest{
		Project:  g.ProjectId,
		Zone:     g.Zone,
		Instance: name,
	}
	opt, err := g.client.Delete(context.Background(), req)
	if err != nil {
		return err
	}
	return opt.Wait(context.Background())
}
