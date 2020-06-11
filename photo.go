package vkapi

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type File struct {
	Name string
	Data io.Reader
}

type UploadServer struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type GetWallUploadServerReq struct {
	GroupID int64
}

func (GetWallUploadServerReq) Name() string {
	return "photos.getWallUploadServer"
}

type GetWallUploadServerResp struct {
	UploadURL string `json:"upload_url"`
	AlbumID   int    `json:"album_id"`
	UserID    int64  `json:"user_id"`
}

func (r GetWallUploadServerReq) Values() url.Values {
	v := url.Values{}
	if r.GroupID > 0 {
		v.Set("group_id", strconv.FormatInt(r.GroupID, 10))
	}
	return v
}

func (vk *VkAPI) GetWallUploadServer(r *GetWallUploadServerReq) (*GetWallUploadServerResp, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}
	var w GetWallUploadServerResp
	if err := json.Unmarshal(resp.Response, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

type SaveWallPhotoReq struct {
	UserID    int
	GroupID   int64
	Photo     string
	Server    int
	Hash      string
	Longitude float64
	Latitude  float64
	Caption   string
}

func (SaveWallPhotoReq) Name() string {
	return "photos.saveWallPhoto"
}

func (s SaveWallPhotoReq) Values() url.Values {
	v := url.Values{}
	if s.UserID > 0 {
		v.Set("user_id", strconv.Itoa(s.UserID))
	}
	if s.GroupID > 0 {
		v.Set("group_id", strconv.FormatInt(s.GroupID, 10))
	}
	v.Set("photo", s.Photo)
	v.Set("server", strconv.Itoa(s.Server))
	v.Set("hash", s.Hash)

	if s.Longitude > 0 {
		v.Set("longitude", strconv.FormatFloat(s.Longitude, 'f', -1, 64))
	}

	if s.Latitude > 0 {
		v.Set("latitude", strconv.FormatFloat(s.Latitude, 'f', -1, 64))
	}

	if len(s.Caption) > 0 {
		v.Set("caption", s.Caption)
	}
	return v
}

type GetMessagesUploadServerReq struct {
	PeerID int64
}

func (GetMessagesUploadServerReq) Name() string {
	return "photos.getMessagesUploadServer"
}

func (r GetMessagesUploadServerReq) Values() url.Values {
	v := url.Values{}
	if r.PeerID > 0 {
		v.Set("peer_id", strconv.FormatInt(r.PeerID, 10))
	}
	return v
}

type SaveMessagesPhotoReq struct {
	Photo  string
	Server int
	Hash   string
}

func (SaveMessagesPhotoReq) Name() string {
	return "photos.saveMessagesPhoto"
}

func (r SaveMessagesPhotoReq) Values() url.Values {
	v := url.Values{}
	v.Set("photo", r.Photo)
	v.Set("server", strconv.Itoa(r.Server))
	v.Set("hash", r.Hash)
	return v
}

type GetMessagesUploadServerResp struct {
	UploadURL string `json:"upload_url"`
	AlbumID   int    `json:"album_id"`
	GroupID   int64  `json:"user_id"`
}

func (vk *VkAPI) GetMessagesUploadServer(r *GetMessagesUploadServerReq) (*GetMessagesUploadServerResp, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}
	var w GetMessagesUploadServerResp
	if err := json.Unmarshal(resp.Response, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (vk *VkAPI) SaveWallPhoto(r *SaveWallPhotoReq) ([]Photo, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}

	var p []Photo
	if err := json.Unmarshal(resp.Response, &p); err != nil {
		return nil, err
	}
	return p, nil
}

func (vk *VkAPI) SaveMessagesPhoto(r *SaveMessagesPhotoReq) ([]Photo, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}

	var p []Photo
	if err := json.Unmarshal(resp.Response, &p); err != nil {
		return nil, err
	}
	return p, nil
}

func MakeUploadPhotoRequest(uploadURL string, files []File) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i := 0; i < len(files); i++ {
		part, err := writer.CreateFormFile("file"+strconv.Itoa(i+1), files[i].Name)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, files[i].Data); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
