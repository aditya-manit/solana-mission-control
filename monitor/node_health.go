package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PrathyushaLakkireddy/solana-prometheus/alerter"
	"github.com/PrathyushaLakkireddy/solana-prometheus/config"
	"github.com/PrathyushaLakkireddy/solana-prometheus/types"
)

func GetNodeHealth(cfg *config.Config) (types.NodeHealth, error) {
	ops := types.HTTPOptions{
		Endpoint: cfg.Endpoints.RPCEndpoint,
		Method:   http.MethodPost,
		Body:     types.Payload{Jsonrpc: "2.0", Method: "getHealth", ID: 1},
	}

	var result types.NodeHealth
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	if strings.EqualFold(result.Result, "ok") {
		log.Printf("Node health : %s", result.Result)
	} else {
		_ = alerter.SendTelegramAlert(fmt.Sprintf("Your node is not running"), cfg)
	}

	return result, nil
}
