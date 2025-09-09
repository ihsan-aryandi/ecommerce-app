package repository

import (
	"bytes"
	"ecommerce-app/internal/api/model"
	"ecommerce-app/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type RajaOngkirRepository struct {
	host      string
	endpoints *config.RajaOngkirEndpointsConfig
	apiKey    string
}

func NewRajaOngkirRepository(cfg *config.Config) *RajaOngkirRepository {
	return &RajaOngkirRepository{
		host:      cfg.RajaOngkir.Host,
		endpoints: cfg.RajaOngkir.Endpoints,
		apiKey:    cfg.RajaOngkir.APIKey,
	}
}

func (r RajaOngkirRepository) fetch(reqConfig *config.RequestConfig, params QueryParams, body []byte) (*http.Response, []byte, error) {
	url := r.host + reqConfig.Path + params.ToString(true)
	headers := HttpHeaders{
		"X-API-KEY": r.apiKey,
	}

	return fetch(reqConfig.Method, url, bytes.NewBuffer(body), headers)
}

func (r RajaOngkirRepository) CalculateShippingCost(costModel *model.RajaOngkirCalculateShippingCost) (*model.RajaOngkirResponse[*model.RajaOngkirCalculateResponse], error) {
	params := QueryParams{
		"shipper_destination_id":  fmt.Sprintf("%d", costModel.ShipperDestinationId),
		"receiver_destination_id": fmt.Sprintf("%d", costModel.ReceiverDestinationId),
		"weight":                  fmt.Sprintf("%.4f", costModel.Weight),
		"item_value":              costModel.ItemValue.String(),
	}

	resp, respBody, err := r.fetch(r.endpoints.CalculateShippingCost, params, nil)
	if err != nil {
		return nil, err
	}

	if err = r.checkError(resp, respBody); err != nil {
		return nil, err
	}

	result := new(model.RajaOngkirResponse[*model.RajaOngkirCalculateResponse])
	if err = json.Unmarshal(respBody, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (RajaOngkirRepository) checkError(resp *http.Response, respBody []byte) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var (
		content model.RajaOngkirResponse[interface{}]
		req     = resp.Request
	)

	_ = json.Unmarshal(respBody, &content)

	message := content.Meta.Message
	if message == "" {
		message = resp.Status
	}

	return fmt.Errorf(
		"request failed: repository=RajaOngkir | method=%s | uri=%s | status_code=%d | message = %s",
		req.Method,
		req.RequestURI,
		resp.StatusCode,
		message,
	)
}
