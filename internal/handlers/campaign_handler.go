package handlers

import (
	"fmt"

	models "bitbucket.org/latonaio/agrimedia-models"
	"bitbucket.org/latonaio/aion-core/pkg/log"
)

func HandleCampaign(metadata map[string]interface{}) error {
	campaigns, err := models.MetadataToCampaigns(metadata)
	if err != nil {
		return fmt.Errorf("failed to convert campaigns: %v", err)
	}
	for _, campaign := range campaigns {
		if campaign.SfCampaignID == nil {
			continue
		}
		c, err := models.CampaignByID(*campaign.SfCampaignID)
		if err != nil {
			log.Printf("failed to get campaign: %v", err)
			continue
		}
		if c != nil {
			log.Printf("update campaign: %s\n", *campaign.SfCampaignID)
			if err := campaign.Update(); err != nil {
				log.Printf("failed to update campaign: %v", err)
				continue
			}
		} else {
			log.Printf("register campaign: %s\n", *campaign.SfCampaignID)
			if err := campaign.Register(); err != nil {
				log.Printf("failed to register campaign: %v", err)
			}
		}
	}
	return nil
}
