package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/vx3r/wg-gen-web/model"
)

// apiError implements a top-level JSON-RPC error.
type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	Data interface{} `json:"data,omitempty"`
}

type apiRequest struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type apiResponse struct {
	Version string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *apiError       `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

func fetchWireGuardAPI(reqData apiRequest) (*apiResponse, error) {
	apiUrl := os.Getenv("WG_STATS_API")
	if apiUrl == "" {
		return nil, errors.New("Status API integration not configured")
	}

	apiClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	jsonData, _ := json.Marshal(reqData)
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "wg-gen-web")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	if os.Getenv("WG_STATS_API_TOKEN") != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("WG_STATS_API_TOKEN")))
	} else if os.Getenv("WG_STATS_API_USER") != "" {
		req.SetBasicAuth(os.Getenv("WG_STATS_API_USER"), os.Getenv("WG_STATS_API_PASS"))
	}

	res, getErr := apiClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	response := apiResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &response, nil
}

// ReadInterfaceStatus object, create default one
func ReadInterfaceStatus() (*model.InterfaceStatus, error) {
	interfaceStatus := &model.InterfaceStatus{
		Name:          "unknown",
		DeviceType:    "unknown",
		ListenPort:    0,
		NumberOfPeers: 0,
		PublicKey:     "",
	}

	data, err := fetchWireGuardAPI(apiRequest{
		Version: "2.0",
		Method:  "GetDeviceInfo",
		Params:  nil,
	})
	if err != nil {
		return interfaceStatus, err
	}

	resultData := data.Result.(map[string]interface{})
	device := resultData["device"].(map[string]interface{})
	interfaceStatus.Name = device["name"].(string)
	interfaceStatus.DeviceType = device["type"].(string)
	interfaceStatus.PublicKey = device["public_key"].(string)
	interfaceStatus.ListenPort = int(device["listen_port"].(float64))
	interfaceStatus.NumberOfPeers = int(device["num_peers"].(float64))

	return interfaceStatus, nil
}

// ReadClientStatus object, create default one, last recent active client is listed first
func ReadClientStatus() ([]*model.ClientStatus, error) {
	var clientStatus []*model.ClientStatus

	data, err := fetchWireGuardAPI(apiRequest{
		Version: "2.0",
		Method:  "ListPeers",
		Params:  []byte("{}"),
	})
	if err != nil {
		return clientStatus, err
	}

	resultData := data.Result.(map[string]interface{})
	peers := resultData["peers"].([]interface{})

	clients, err := ReadClients()
	withClientDetails := true
	if err != nil {
		withClientDetails = false
	}

	for _, tmpPeer := range peers {
		peer := tmpPeer.(map[string]interface{})
		peerHandshake, _ := time.Parse(time.RFC3339Nano, peer["last_handshake"].(string))
		peerIPs := peer["allowed_ips"].([]interface{})
		peerAddresses := make([]string, len(peerIPs))
		for i, peerIP := range peerIPs {
			peerAddresses[i] = peerIP.(string)
		}
		peerHandshakeRelative := time.Since(peerHandshake)
		peerActive := peerHandshakeRelative.Minutes() < 3 // TODO: we need a better detection... ping for example?

		newClientStatus := &model.ClientStatus{
			PublicKey:             peer["public_key"].(string),
			HasPresharedKey:       peer["has_preshared_key"].(bool),
			ProtocolVersion:       int(peer["protocol_version"].(float64)),
			Name:                  "UNKNOWN",
			Email:                 "UNKNOWN",
			Connected:             peerActive,
			AllowedIPs:            peerAddresses,
			Endpoint:              peer["endpoint"].(string),
			LastHandshake:         peerHandshake,
			LastHandshakeRelative: peerHandshakeRelative,
			ReceivedBytes:         int(peer["receive_bytes"].(float64)),
			TransmittedBytes:      int(peer["transmit_bytes"].(float64)),
		}

		if withClientDetails {
			for _, client := range clients {
				if client.PublicKey != newClientStatus.PublicKey {
					continue
				}

				newClientStatus.Name = client.Name
				newClientStatus.Email = client.Email
				break
			}
		}

		clientStatus = append(clientStatus, newClientStatus)
	}

	sort.Slice(clientStatus, func(i, j int) bool {
		return clientStatus[i].LastHandshakeRelative < clientStatus[j].LastHandshakeRelative
	})

	return clientStatus, nil
}
