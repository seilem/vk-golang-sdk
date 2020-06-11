package vkapi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type VideoGetReq struct {
	OwnerID  int
	Videos   []string
	AlbumID  int
	Offset   int
	Count    int
	Extended bool
}

type UploadedVideoResp struct {
	Size    int   `json:"size"`
	VideoID int64 `json:"video_id"`
}

func (VideoGetReq) Name() string {
	return "video.get"
}

func (r *VideoGetReq) Values() url.Values {
	v := url.Values{}
	if r.OwnerID != 0 {
		v.Set("owner_id", strconv.Itoa(r.OwnerID))
	}

	v.Set("videos", strings.Join(r.Videos, ","))

	if r.AlbumID != 0 {
		v.Set("album_id", strconv.Itoa(r.AlbumID))
	}

	v.Set("offset", strconv.Itoa(r.Offset))

	if r.Count != 0 {
		v.Set("count", strconv.Itoa(r.Count))
	}
	v.Set("extended", strconv.Itoa(btoi(r.Extended)))

	return v
}

type VideoSaveReq struct {
	VideoName   string
	Description string
	Link        string
	GroupID     int64 `json:"group_id"`
}

func (VideoSaveReq) Name() string {
	return "video.save"
}

func (r *VideoSaveReq) Values() url.Values {
	v := url.Values{}
	if len(r.Link) != 0 {
		v.Set("link", r.Link)
	}

	if len(r.VideoName) > 0 {
		v.Set("name", r.VideoName)
	}

	if len(r.Description) > 0 {
		v.Set("description", r.Description)
	}

	if r.GroupID > 0 {
		v.Set("group_id", strconv.FormatInt(r.GroupID, 10))
	}

	return v
}

type VideoSaveResp struct {
	AccessKey   string `json:"access_key"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`
	Title       string `json:"title"`
	UploadURL   string `json:"upload_url"`
	VideoID     int    `json:"video_id"`
}

func (vk *VkAPI) VideoGet(r *UsersGetReq) (*Videos, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}

	var v Videos
	if err := json.Unmarshal(resp.Response, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (vk *VkAPI) VideoSave(r *VideoSaveReq) (*VideoSaveResp, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return nil, err
	}

	var v VideoSaveResp
	if err := json.Unmarshal(resp.Response, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func MakeUploadVideoRequest(uploadURL string, file File) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("video_file", file.Name)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file.Data); err != nil {
		return nil, err
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

func (vk *VkAPI) UploadVideo(uploadURL string, file File) (*UploadedVideoResp, error) {
	req, err := MakeUploadVideoRequest(uploadURL, file)
	if err != nil {
		return nil, err
	}

	resp, err := vk.Client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v UploadedVideoResp
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
