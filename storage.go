package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type StorageSetReq struct {
	Key    string
	Value  string
	UserID int
}

func (StorageSetReq) Name() string {
	return "storage.set"
}

func (s *StorageSetReq) Values() url.Values {
	v := url.Values{}
	v.Set("key", s.Key)
	v.Set("value", s.Value)
	if s.UserID != 0 {
		v.Set("user_id", strconv.Itoa(s.UserID))
	}
	return v
}

type StorageGetKeysReq struct {
	UserID int
	Offset int
	Count  int
}

func (StorageGetKeysReq) Name() string {
	return "storage.getKeys"
}

func (s *StorageGetKeysReq) Values() url.Values {
	v := url.Values{}
	v.Set("user_id", strconv.Itoa(s.UserID))
	v.Set("offset", strconv.Itoa(s.Offset))
	if s.Count != 0 {
		v.Set("count", strconv.Itoa(s.Count))
	}
	return v
}

type StorageGetReq struct {
	Key    string
	Keys   []string
	UserID int
}

func (StorageGetReq) Name() string {
	return "storage.get"
}

func (s StorageGetReq) Values() url.Values {
	v := url.Values{}
	if len(s.Key) > 0 {
		v.Set("key", s.Key)
	}

	if len(s.Keys) > 0 {
		v.Set("keys", strings.Join(s.Keys, ","))
	}
	v.Set("user_id", strconv.Itoa(s.UserID))
	return v
}

type StorageGetResp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// StorageSet saves a value of variable with the name set by key parameter.
//
// See https://vk.com/dev/storage.set
func (vk *VkAPI) StorageSet(v *StorageSetReq) (bool, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return false, err
	}

	var r float64
	if err := json.Unmarshal(resp.Response, &r); err != nil {
		return false, err
	}

	if r == 0 {
		return false, nil
	}

	return true, nil
}

// StorageGetKeysReq returns the names of all variables.
//
// See https://vk.com/dev/storage.get
func (vk *VkAPI) StorageGetKeys(v *StorageGetKeysReq) ([]string, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	var r []string
	if err := json.Unmarshal(resp.Response, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// StorageGet returns a value of variable with the name set by key parameter.
//
// See https://vk.com/dev/storage.get
func (vk *VkAPI) StorageGet(v *StorageGetReq) ([]StorageGetResp, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	var r interface{}
	if err := json.Unmarshal(resp.Response, &r); err != nil {
		return nil, err
	}

	if _, ok := r.(string); ok {
		return []StorageGetResp{{Key: v.Key, Value: r.(string)}}, nil
	}

	var m []StorageGetResp
	if err := json.Unmarshal(resp.Response, &m); err != nil {
		return nil, err
	}
	return m, nil
}
