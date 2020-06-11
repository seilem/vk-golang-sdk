package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type OpenTopicReq struct {
	GroupID int64
	TopicID int64
}

func (OpenTopicReq) Name() string {
	return "board.openTopic"
}

func (r OpenTopicReq) Values() url.Values {
	v := url.Values{}
	v.Set("group_id", strconv.FormatInt(r.GroupID, 10))
	v.Set("topic_id", strconv.FormatInt(r.TopicID, 10))
	return v
}

// Re-opens a previously closed topic on a community's discussion board.
//
// See https://vk.com/dev/board.openTopic
func (vk *VkAPI) OpenTopic(r *OpenTopicReq) error {
	_, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return err
	}
	return nil
}

type CloseTopicReq struct {
	GroupID int64
	TopicID int64
}

func (CloseTopicReq) Name() string {
	return "board.closeTopic"
}

func (r CloseTopicReq) Values() url.Values {
	v := url.Values{}
	v.Set("group_id", strconv.FormatInt(r.GroupID, 10))
	v.Set("topic_id", strconv.FormatInt(r.TopicID, 10))
	return v
}

// Closes a topic on a community's discussion board so that comments cannot be posted.
//
// See https://vk.com/dev/board.closeTopic
func (vk *VkAPI) CloseTopic(r *CloseTopicReq) error {
	_, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return err
	}
	return nil
}

type CreateCommentReq struct {
	GroupID     int64
	TopicID     int64
	Message     string
	Attachments []string
	FromGroup   bool
	StickerID   int
	GUID        string
}

func (CreateCommentReq) Name() string {
	return "board.createComment"
}

func (r CreateCommentReq) Values() url.Values {
	v := url.Values{}
	v.Set("group_id", strconv.FormatInt(r.GroupID, 10))
	v.Set("topic_id", strconv.FormatInt(r.TopicID, 10))
	if len(r.Message) > 0 {
		v.Set("message", r.Message)
	}
	if len(r.Attachments) > 0 {
		v.Set("attachment", strings.Join(r.Attachments, ","))
	}
	v.Set("from_group", strconv.Itoa(btoi(r.FromGroup)))
	if r.StickerID > 0 {
		v.Set("sticker_id", strconv.Itoa(r.StickerID))
	}
	if len(r.GUID) > 0 {
		v.Set("guid", r.GUID)
	}
	return v
}

// Adds a comment on a topic on a community's discussion board.
//
// See https://vk.com/dev/board.createComment
func (vk *VkAPI) CreateComment(r *CreateCommentReq) (int64, error) {
	resp, err := vk.MakeRequest(r.Name(), r.Values())
	if err != nil {
		return 0, err
	}

	var id int64
	if err := json.Unmarshal(resp.Response, &id); err != nil {
		return 0, err
	}
	return id, nil
}
