package main

import (
	"context"
	"fmt"
	"sync"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/proto/kanbanpb"
)

const msName = "data-interface-for-salesforce-campaign"

var (
	kanbanClient msclient.MicroserviceClient
	once         sync.Once
)

func newKanbanClient(ctx context.Context) error {
	var err error
	once.Do(func() {
		kanbanClient, err = msclient.NewKanbanClient(ctx, msName, kanbanpb.InitializeType_START_SERVICE)
	})
	if err != nil {
		return fmt.Errorf("failed to construct kanban client: %v", err)
	}
	return nil
}

func writeKanban(data map[string]interface{}) error {
	var options []msclient.Option
	options = append(options, msclient.SetMetadata(data))
	options = append(options, msclient.SetProcessNumber(kanbanClient.GetProcessNumber()))
	req, err := msclient.NewOutputData(options...)
	if err != nil {
		return fmt.Errorf("failed to construct output request: %v", err)
	}
	if err := kanbanClient.OutputKanban(req); err != nil {
		return fmt.Errorf("failed to output to kanban: %v", err)
	}
	return nil
}
