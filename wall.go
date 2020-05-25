package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type WallGetReq struct {
	OwnerID  int
	Domain   string
	Offset   int
	Count    int
	Filter   string
	Extended bool
	Fields   []string
}

type WallPostReq struct {
	OwnerID    int
	FriendOnly bool
	FromGroup  bool
	Message    string
	Copyright  string
	Attachments []string
}

type WallPostResp struct {
	PostID int64 `json:"post_id"`
}

func (w WallPostReq) Name() string {
	return "wall.post"
}

func (w WallPostReq) Values() url.Values {
	v := url.Values{}
	if w.OwnerID != 0 {
		v.Set("owner_id", strconv.Itoa(w.OwnerID))
	}
	v.Set("friend_only", strconv.Itoa(btoi(w.FriendOnly)))
	v.Set("from_group", strconv.Itoa(btoi(w.FromGroup)))
	if len(w.Message) > 0 {
		v.Set("message", w.Message)
	}

	if len(w.Copyright) > 0 {
		v.Set("copyright", w.Copyright)
	}

	if len(w.Attachments) > 0 {
		v.Set("attachments", strings.Join(w.Attachments, ","))
	}
	return v
}

func (WallGetReq) Name() string {
	return "wall.get"
}

func (w *WallGetReq) Values() url.Values {
	v := url.Values{}
	if w.OwnerID != 0 {
		v.Set("owner_id", strconv.Itoa(w.OwnerID))
	}

	if len(w.Domain) > 0 {
		v.Set("domain", w.Domain)
	}

	v.Set("offset", strconv.Itoa(w.Offset))

	if w.Count != 0 {
		v.Set("count", strconv.Itoa(w.Count))
	}

	if len(w.Filter) > 0 {
		v.Set("filter", w.Filter)
	}

	v.Set("extended", strconv.Itoa(btoi(w.Extended)))

	if len(w.Fields) > 0 {
		v.Set("fields", strings.Join(w.Fields, ","))
	}

	return v
}

func (vk *VkAPI) WallGet(r *WallGetReq) (*Wall, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}
	var w Wall
	if err := json.Unmarshal(resp.Response, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (vk *VkAPI) WallPost(r *WallPostReq) (int64, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return 0, err
	}
	var w WallPostResp
	if err := json.Unmarshal(resp.Response, &w); err != nil {
		return 0, err
	}
	return w.PostID, nil
}
