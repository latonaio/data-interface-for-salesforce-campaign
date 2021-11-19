package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	models "bitbucket.org/latonaio/agrimedia-models"
	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/pkg/log"
	"bitbucket.org/latonaio/data-interface-for-salesforce-campaign/internal/handlers"
	"bitbucket.org/latonaio/data-interface-for-salesforce-campaign/internal/resources"
)

func main() {
	// Create Kanban client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := models.NewDBConPool(ctx); err != nil {
		panic(err)
	}
	if err := newKanbanClient(ctx); err != nil {
		log.Fatalf("failed to get kanban client: %v", err)
	}
	log.Printf("successful get kanban client")
	defer kanbanClient.Close()

	// Get Kanban channel by Kanban client
	kanbanCh := kanbanClient.GetKanbanCh()
	log.Printf("successful get kanban channel")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)
	for {
		select {
		case s := <-signalCh:
			fmt.Printf("received signal: %s", s.String())
			goto END
		case k := <-kanbanCh:
			if k == nil {
				goto NODATA
			}

			// Get metadata from Kanban
			fromMetadata, err := msclient.GetMetadataByMap(k)
			if err != nil {
				log.Printf("failed to get metadata.err: %d", err)
				continue
			}

			toMetadata, err := handle(fromMetadata)
			if err != nil {
				log.Printf("failed to handle: %v", err)
				continue
			}

			if toMetadata == nil {
				continue
			}

			// Write metadata for kanban
			if err := writeKanban(toMetadata); err != nil {
				log.Printf("failed to write kanban: %v", err)
				continue
			}
			log.Printf("write metadata to kanban")
			log.Printf("metadata: %v", toMetadata)
		}
	NODATA:
	}
END:
}

func handle(fromMetadata map[string]interface{}) (map[string]interface{}, error) {
	log.Printf("got metadata from kanban")
	log.Printf("metadata: %v", fromMetadata)

	ck, ok := fromMetadata["connection_type"].(string)
	if !ok {
		return nil, errors.New("invalid connection key")
	}

	if ck == "response" {
		if err := handlers.HandleCampaign(fromMetadata); err != nil {
			return nil, fmt.Errorf("failed to handle campaign: %w", err)
		}
		return nil, nil
	}

	if ck == "request" {
		// Get resource from metadata
		campaign, err := resources.NewCampaign(fromMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to construct resource: %w", err)
		}

		// Build metadata for Kanban
		toMetadata, err := campaign.BuildMetadata()
		if err != nil {
			return nil, fmt.Errorf("failed to build metadata: %w", err)
		}
		return toMetadata, nil
	}

	return nil, fmt.Errorf("unknown connection_type: %v", ck)
}
