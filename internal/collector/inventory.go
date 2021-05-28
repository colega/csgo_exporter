package collector

import (
	"fmt"
	"strconv"

	"github.com/kinduff/csgo_exporter/internal/metrics"
	"github.com/kinduff/csgo_exporter/internal/model"

	log "github.com/sirupsen/logrus"
)

func (collector *collector) collectPlayerInventory() {
	inventory := model.Inventory{}
	inventoryEndpoint := fmt.Sprintf("https://steamcommunity.com/inventory/%s/730/2", collector.config.SteamID)
	if err := collector.client.DoCustomAPIRequest(inventoryEndpoint, collector.config, &inventory); err != nil {
		log.Fatal(err)
	}

	for _, s := range inventory.Assets {
		amount, _ := strconv.ParseInt(s.Amount, 10, 64)

		collector.playerInventory[s.ClassID] = model.PlayerInventory{
			ClassID: s.ClassID,
			Amount:  amount,
		}
	}

	for _, s := range inventory.Descriptions {
		t := collector.playerInventory[s.ClassID]

		collector.playerInventory[t.ClassID] = model.PlayerInventory{
			ClassID:    t.ClassID,
			Amount:     t.Amount,
			Tradable:   s.Tradable == 1,
			Marketable: s.Marketable == 1,
			MarketName: s.MarketName,
		}
	}

	for _, s := range collector.playerInventory {
		metrics.UserInventory.WithLabelValues(
			collector.config.SteamID,
			s.ClassID,
			s.MarketName,
			s.Currency,
			strconv.FormatBool(s.Tradable),
			strconv.FormatBool(s.Marketable),
		).Set(float64(s.Amount))
	}
}