package vkapi

import (
	"encoding/json"
	"net/url"
	"strings"
)

type UsersGetReq struct {
	UserIDs  []string
	Fields   []string
	NameCase string
}

func (UsersGetReq) Name() string {
	return "users.get"
}

func (u *UsersGetReq) Values() url.Values {
	v := url.Values{}
	v.Set("user_ids", strings.Join(u.UserIDs, ","))
	v.Set("fields", strings.Join(u.Fields, ","))
	if len(u.NameCase) > 0 {
		v.Set("name_case", u.NameCase)
	}
	return v
}

func (vk *VkAPI) UsersGet(v *UsersGetReq) ([]User, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	var u []User
	if err := json.Unmarshal(resp.Response, &u); err != nil {
		return nil, err
	}

	return u, nil
}
