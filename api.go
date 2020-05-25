package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	APIEndpoint = "https://api.vk.com/method/%s"
	APIVersion  = "5.103"
)

type Params map[string]interface{}

type VkAPI struct {
	Token       string
	Client      *http.Client
	groupLPSubs GroupLPSubs
	msgLPSubs   MsgLPSubs
}

type APIResponse struct {
	Response      json.RawMessage `json:"response"`
	ResponseError json.RawMessage `json:"error"`
}

type APIError struct {
	Code int    `json:"error_code"`
	Msg  string `json:"error_msg"`
	Raw  json.RawMessage
}

type Method interface {
	Name() string
	Values() url.Values
}

func (e APIError) Error() string {
	return e.Msg
}

func NewVkAPI(token string) *VkAPI {
	return NewVkAPIWithClient(token, http.DefaultClient)
}

func NewVkAPIWithClient(token string, client *http.Client) *VkAPI {
	return &VkAPI{
		Token:  token,
		Client: client,
		groupLPSubs: GroupLPSubs{
			events: make(map[string]func(m *GroupLPUpdates)),
		},
		msgLPSubs: MsgLPSubs{
			events: make(map[int]func(m interface{})),
		},
	}
}

func (vk *VkAPI) MakeRequest(method string, params url.Values) (*APIResponse, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("v", APIVersion)
	params.Set("access_token", vk.Token)

	resp, err := vk.Client.PostForm(fmt.Sprintf(APIEndpoint, method), params)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error: %w", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(data, &apiResponse); err != nil {
		return nil, err
	}

	if len(apiResponse.ResponseError) > 0 {
		var e APIError
		if err := json.Unmarshal(apiResponse.ResponseError, &e); err != nil {
			return nil, err
		}
		e.Raw = apiResponse.ResponseError
		return nil, e
	}

	return &apiResponse, nil
}
