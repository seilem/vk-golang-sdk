package vkapi

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const LastLPVersion = 3

type MsgReq struct {
	UserID          int
	RandomID        int64
	PeerID          int
	Domain          string
	ChatID          int64
	UsersID         []int
	Message         string
	Lat             float64
	Long            float64
	Attachment      []string
	ReplyTo         int
	ForwardMessages []int
	StickerID       int
	GroupID         int
	Keyboard        *Keyboard
	Payload         string
	DontParseLinks  bool
	DisableMentions bool
	Intent          string
}

func (MsgReq) Name() string {
	return "messages.send"
}

func (m *MsgReq) Values() url.Values {
	v := url.Values{}
	if m.UserID != 0 {
		v.Set("user_id", strconv.Itoa(m.UserID))
	}

	if m.RandomID != 0 {
		v.Set("random_id", strconv.FormatInt(m.RandomID, 10))
	}

	if m.PeerID != 0 {
		v.Set("peer_id", strconv.Itoa(m.PeerID))
	}

	if len(m.Domain) > 0 {
		v.Set("domain", m.Domain)
	}

	if m.ChatID != 0 {
		v.Set("chat_id", strconv.FormatInt(m.ChatID, 10))
	}

	if len(m.UsersID) > 0 {
		v.Set("user_ids", sliceToStr(m.UsersID))
	}

	if len(m.Message) > 0 {
		v.Set("message", m.Message)
	}

	if m.Lat != 0 {
		v.Set("lat", fmt.Sprintf("%f", m.Lat))
	}

	if m.Long != 0 {
		v.Set("long", fmt.Sprintf("%f", m.Long))
	}

	if len(m.Attachment) > 0 {
		v.Set("attachment", strings.Join(m.Attachment, ","))
	}

	if m.ReplyTo > 0 {
		v.Set("reply_to", strconv.Itoa(m.ReplyTo))
	}

	if len(m.ForwardMessages) > 0 {
		v.Set("forward_messages", sliceToStr(m.ForwardMessages))
	}

	if m.StickerID > 0 {
		v.Set("sticker_id", strconv.Itoa(m.StickerID))
	}

	if m.GroupID > 0 {
		v.Set("group_id", strconv.Itoa(m.GroupID))
	}

	if m.Keyboard != nil {
		k, err := json.Marshal(m.Keyboard)
		if err != nil {
			log.Print("marshalling error: %w", err)
		}
		v.Set("keyboard", string(k))
	}

	if len(m.Payload) > 0 {
		v.Set("payload", m.Payload)
	}

	v.Set("dont_parse_links", strconv.FormatBool(m.DontParseLinks))
	v.Set("disable_mentions", strconv.FormatBool(m.DisableMentions))

	if len(m.Intent) > 0 {
		v.Set("Intent", m.Intent)
	}

	return v
}

type MsgSetActivityReq struct {
	UserID  int
	Type    string
	PeerID  int
	GroupID int
}

func (MsgSetActivityReq) Name() string {
	return "messages.setActivity"
}

func (m *MsgSetActivityReq) Values() url.Values {
	v := url.Values{}
	if m.UserID != 0 {
		v.Set("user_id", strconv.Itoa(m.UserID))
	}

	if len(m.Type) > 0 {
		v.Set("type", m.Type)
	}

	if m.PeerID != 0 {
		v.Set("peer_id", strconv.Itoa(m.PeerID))
	}

	if m.GroupID > 0 {
		v.Set("group_id", strconv.Itoa(m.GroupID))
	}

	return v
}

type MsgEditReq struct {
	PeerID              int
	Message             string
	MessageID           int
	Lat                 float64
	Long                float64
	Attachment          []string
	GroupID             int
	KeepForwardMessages bool
	KeepSnippets        bool
	DontParseLinks      bool
}

func (MsgEditReq) Name() string {
	return "messages.edit"
}

func (m *MsgEditReq) Values() url.Values {
	v := url.Values{}
	if m.PeerID != 0 {
		v.Set("peer_id", strconv.Itoa(m.PeerID))
	}

	if len(m.Message) > 0 {
		v.Set("message", m.Message)
	}

	if m.MessageID != 0 {
		v.Set("message_id", strconv.Itoa(m.MessageID))
	}

	if m.Lat != 0 {
		v.Set("lat", fmt.Sprintf("%f", m.Lat))
	}

	if m.Long != 0 {
		v.Set("long", fmt.Sprintf("%f", m.Long))
	}

	if len(m.Attachment) > 0 {
		v.Set("attachment", strings.Join(m.Attachment, ","))
	}

	if m.GroupID > 0 {
		v.Set("group_id", strconv.Itoa(m.GroupID))
	}

	v.Set("keep_forward_messages", strconv.FormatBool(m.KeepForwardMessages))
	v.Set("keep_snippets", strconv.FormatBool(m.KeepSnippets))
	v.Set("dont_parse_links", strconv.FormatBool(m.DontParseLinks))

	return v
}

type MsgDeleteReq struct {
	MessageIDs   []int
	Spam         bool
	GroupID      int
	DeleteForAll bool
}

func (MsgDeleteReq) Name() string {
	return "messages.delete"
}

func (m *MsgDeleteReq) Values() url.Values {
	v := url.Values{}

	if len(m.MessageIDs) > 0 {
		v.Set("message_ids", sliceToStr(m.MessageIDs))
	}

	v.Set("spam", strconv.FormatBool(m.Spam))

	if m.GroupID != 0 {
		v.Set("group_id", strconv.Itoa(m.GroupID))
	}

	v.Set("delete_for_all", strconv.Itoa(btoi(m.DeleteForAll)))

	return v
}

type GetByConversationMessageIDReq struct {
	PeerID                 []int
	ConversationMessageIDs []int
	Extended               bool
	Fields                 []string
	GroupID                int
}

func (GetByConversationMessageIDReq) Name() string {
	return "messages.getByConversationMessageId"
}

func (g *GetByConversationMessageIDReq) Values() url.Values {
	v := url.Values{}
	if len(g.PeerID) > 0 {
		v.Set("peer_id", sliceToStr(g.PeerID))
	}

	if len(g.PeerID) > 0 {
		v.Set("conversation_message_ids", sliceToStr(g.ConversationMessageIDs))
	}

	v.Set("extended", strconv.Itoa(btoi(g.Extended)))

	if len(g.Fields) > 0 {
		v.Set("fields", strings.Join(g.Fields, ","))
	}

	return v
}

type MsgGetLPServerReq struct {
	NeedPTS   bool
	GroupID   int
	LPVersion int
}

func (MsgGetLPServerReq) Name() string {
	return "messages.getLongPollServer"
}

func (m *MsgGetLPServerReq) Values() url.Values {
	v := url.Values{}
	v.Set("need_pts", strconv.Itoa(btoi(m.NeedPTS)))

	if m.GroupID > 0 {
		v.Set("group_id", strconv.Itoa(m.GroupID))
	}

	v.Set("lp_version", strconv.Itoa(m.LPVersion))
	return v
}

// MsgSend sends a message.
//
// See https://vk.com/dev/messages.send
func (vk *VkAPI) MsgSend(v *MsgReq) ([]NewMessageResp, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	if len(v.UsersID) > 0 {
		var m []NewMessageResp
		if err := json.Unmarshal(resp.Response, &m); err != nil {
			return nil, err
		}
		return m, nil
	}

	m, err := resp.Response.MarshalJSON()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(m)
	id, err := binary.ReadVarint(buf)
	if err != nil {
		return nil, err
	}
	return []NewMessageResp{{MessageID: int(id)}}, nil

}

// MsgSetActivity changes the status of a user as typing in a conversation.
//
// See https://vk.com/dev/messages.setActivity
func (vk *VkAPI) MsgSetActivity(v *MsgSetActivityReq) error {
	_, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return err
	}
	return nil
}

// MsgEdit edits the message.
//
// See https://vk.com/dev/messages.edit
func (vk *VkAPI) MsgEdit(v *MsgEditReq) error {
	_, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return err
	}
	return nil
}

// MsgDelete deletes one or more messages.
//
// See https://vk.com/dev/messages.delete
func (vk *VkAPI) MsgDelete(v *MsgDeleteReq) error {
	_, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return err
	}
	return nil
}

// MsgGetByConversationMessageID returns messages by their IDs.
//
// See https://vk.com/dev/messages.getByConversationMessageId
func (vk *VkAPI) MsgGetByConversationMessageID(v *GetByConversationMessageIDReq) (*MessagesWithCount, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	var m MessagesWithCount
	if err := json.Unmarshal(resp.Response, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// Edits the message.
//
// See https://vk.com/dev/messages.edit
func (vk *VkAPI) MsgGetLPServer(v *MsgGetLPServerReq) (*MsgLPServer, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	var apiResp MsgLPServer
	if err := json.Unmarshal(resp.Response, &apiResp); err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func (vk *VkAPI) MsgLPServ(groupID, mode int) error {
	return vk.msgLongPoll(context.Background(), groupID, LastLPVersion, mode)
}

func (vk *VkAPI) msgLongPoll(ctx context.Context, groupID, LPVersion, mode int) error {
	server, err := vk.MsgGetLPServer(&MsgGetLPServerReq{
		GroupID:   groupID,
		LPVersion: LPVersion,
	})
	if err != nil {
		log.Print("Get message lp server is failed: %w", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		serverURL := fmt.Sprintf(
			"https://%s?act=a_check&key=%s&ts=%d&wait=25&mode=%d&version=%d",
			server.Server, server.Key, server.TS, mode, LPVersion)

		req, err := http.NewRequest(http.MethodGet, serverURL, nil)
		if err != nil {
			log.Print("Message lp server connection failed: %w", err)
			continue
		}

		resp, err := vk.Client.Do(req)
		if err != nil {
			log.Print("Message lp request connection failed: %w", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		var e MsgLPEvent
		if err := json.Unmarshal(body, &e); err != nil {
			fmt.Println(err)
		}

		switch true {
		case e.Failed == 0:
			for _, u := range e.Updates {
				code := u[0].(float64)
				switch code {
				case 4:
					m := MsgLPNewMessage{
						MessageID:  int(u[1].(float64)),
						Flags:      int(u[2].(float64)),
						PeerID:     int(u[3].(float64)),
						Timestamp:  int64(u[4].(float64)),
						Text:       u[5].(string),
						Attachment: make(map[string]string),
					}

					for t, v := range u[7].(map[string]interface{}) {
						m.Attachment[t] = v.(string)
					}

					if m.Flags == 19 || m.Flags == 51 ||
						m.Flags == 531 || m.Flags == 563 ||
						m.Flags == 3 || m.Flags == 35 {
						m.Type = "message_new"
					} else {
						m.Type = "message_reply"
					}

					vk.handleMsgLPCallback(4, m)
				case 61:
					t := MsgLPUserTyping{UserID: int(u[1].(float64))}
					vk.handleMsgLPCallback(61, t)
				}
			}
			server.TS = e.TS
		case e.Failed == 1:
			server.TS = e.TS
		case e.Failed == 2 || e.Failed == 3:
			newServer, err := vk.MsgGetLPServer(&MsgGetLPServerReq{
				GroupID:   groupID,
				LPVersion: LPVersion,
			})
			if err != nil {
				log.Print("Get message lp server is failed: %w", err)
				continue
			}

			if e.Failed == 2 {
				server.Key = newServer.Key
			}

			if e.Failed == 3 {
				server.Key = newServer.Key
				server.TS = newServer.TS
			}
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

func (vk *VkAPI) handleMsgLPCallback(code int, event interface{}) {
	_, exists := vk.msgLPSubs.events[code]

	if exists {
		vk.msgLPSubs.events[code](event)
	}
}
