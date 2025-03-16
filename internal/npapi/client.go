package npapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"go.albinodrought.com/neptunes-pride/internal/types"
)

type Request struct {
	GameNumber string
	APIKey     string
}

type NeptunesPrideClient interface {
	State(ctx context.Context, request *Request) (*types.APIResponse, error)
}

var ErrUnsuccessfulResponse = errors.New("response returned non-200 status code")

func NewClient(base *http.Client) NeptunesPrideClient {
	return &httpClient{base}
}

type httpClient struct {
	base *http.Client
}

func (c *httpClient) State(ctx context.Context, request *Request) (*types.APIResponse, error) {
	data := url.Values{}
	data.Set("api_version", "0.1")
	data.Set("game_number", request.GameNumber)
	data.Set("code", request.APIKey)

	content := data.Encode()

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://np.ironhelmet.com/api?"+content,
		nil,
	)

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("User-Agent", "NP Scanner by github.com/AlbinoDrought")

	httpResp, err := c.base.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, ErrUnsuccessfulResponse
	}
	defer httpResp.Body.Close()

	apiResponse := &types.APIResponse{}
	if err := json.NewDecoder(httpResp.Body).Decode(apiResponse); err != nil {
		return nil, err
	}

	if apiResponse.Error != "" {
		return nil, errors.New(apiResponse.Error)
	}

	return apiResponse, nil
}
