package vkapi

import "encoding/json"

type Message struct {
	ID                    int          `json:"id"`
	Date                  int64        `json:"date"`
	PeerID                int          `json:"peer_id"`
	FromID                int          `json:"from_id"`
	Text                  string       `json:"text"`
	ConversationMessageID int          `json:"conversation_message_id"`
	RandomID              int64        `json:"random_id"`
	Ref                   string       `json:"ref"`
	RefSource             string       `json:"ref_source"`
	Attachments           []Attachment `json:"attachments"`
	Important             bool         `json:"important"`
	Place                 *Place       `json:"place"`
	Payload               string       `json:"payload"`
	Keyboard              *Keyboard    `json:"keyboard"`
	FwdMessages           []*Message   `json:"fwd_messages"`
	ReplyMessage          *Message     `json:"reply_message"`
	Action                *Action      `json:"action"`
}

type Geo struct {
	Type        string      `json:"type"`
	Coordinates Coordinates `json:"coordinates"`
	Place       *Place      `json:"place"`
}

type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Place struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Coordinates
	Created int64  `json:"created"`
	Icon    string `json:"icon"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type Action struct {
	Type     string  `json:"type"`
	MemberID int     `json:"member_id"`
	Text     string  `json:"text"`
	Email    string  `json:"email"`
	APhoto   *APhoto `json:"photo"`
}

type APhoto struct {
	Photo50  string `json:"photo_50"`
	Photo100 string `json:"photo_100"`
	Photo200 string `json:"photo_200"`
}

type ClientInfo struct {
	ButtonActions  []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

type NewMessage struct {
	Message    Message    `json:"message"`
	ClientInfo ClientInfo `json:"client_info"`
}

type Post struct {
	ID             int             `json:"id"`
	ToID           int             `json:"to_id"`
	FromID         int             `json:"from_id"`
	CreatedBy      int             `json:"created_by"`
	Date           int64           `json:"date"`
	Text           string          `json:"text"`
	ReplyOwnerID   int             `json:"reply_owner_id"`
	ReplyPostID    int             `json:"reply_post_id"`
	FriendsOnly    bool            `json:"friends_only"`
	Comments       Comments        `json:"comments"`
	Likes          Likes           `json:"likes"`
	Reposts        Reposts         `json:"reposts"`
	Views          Views           `json:"views"`
	PostSource     *PostSource     `json:"post_source"`
	AttachmentsRaw json.RawMessage `json:"attachments"`
	Attachments    map[string][]interface{}
	Geo            *Geo `json:"geo"`
	SignerID       int  `json:"signer_id"`
	CopyHistory    []*Post
	CanPin         int    `json:"can_pin"`
	CanEdit        int    `json:"can_edit"`
	IsPinned       int    `json:"is_pinned"`
	MarkedAsAds    int    `json:"marked_as_ads"`
	IsFavorite     bool   `json:"is_favorite"`
	AccessKey      string `json:"access_key"`
	PostponedID    int    `json:"postponed_id"`
}

func (p *Post) LoadAttachments() error {
	if len(string(p.AttachmentsRaw)) <= 0 {
		return nil
	}
	var attachments []attach
	if err := json.Unmarshal(p.AttachmentsRaw, &attachments); err != nil {
		return err
	}
	p.Attachments = make(map[string][]interface{}, len(attachments))
	for _, a := range attachments {
		switch a.Type {
		case "photo":
			p.Attachments[a.Type] = append(p.Attachments[a.Type], a.Photo)
		case "video":
			p.Attachments[a.Type] = append(p.Attachments[a.Type], a.Video)

		}
	}
	return nil
}

type attach struct {
	Type  string `json:"type"`
	Photo *Photo `json:"photo"`
	Video *Video `json:"video"`
}

type Comments struct {
	Count         int  `json:"count"`
	CanPost       int  `json:"can_post"`
	GroupsCanPost bool `json:"groups_can_post"`
	CanClose      bool `json:"can_close"`
	CanOpen       bool `json:"can_open"`
}

type Likes struct {
	Count      int `json:"count"`
	UserLikes  int `json:"user_likes"`
	CanLike    int `json:"can_like"`
	CanPublish int `json:"can_publish"`
}

type Reposts struct {
	Count        int `json:"count"`
	UserReposted int `json:"user_reposted"`
}

type Views struct {
	Count int `json:"count"`
}

type PostSource struct {
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Data     string `json:"data"`
	Url      string `json:"url"`
}

type Attachment struct {
	Type string `json:"type"`
	Wall *Post  `json:"wall"`
}

type MessagesWithCount struct {
	Count int        `json:"count"`
	Items []*Message `json:"items"`
}

type NewMessageResp struct {
	PeerID    int    `json:"peer_id"`
	MessageID int    `json:"message_id"`
	Error     *Error `json:"error"`
}

type LPServer struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     string `json:"ts"`
	PTS    string `json:"pts"`
}

type MsgLPServer struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     int    `json:"ts"`
	PTS    int    `json:"pts"`
}

type MsgLPEvent struct {
	Failed  int             `json:"failed"`
	TS      int             `json:"ts"`
	Updates [][]interface{} `json:"updates"`
}

type GroupLPEvent struct {
	Failed  int              `json:"failed"`
	TS      string           `json:"ts"`
	Updates []GroupLPUpdates `json:"updates"`
}

type GroupLPUpdates struct {
	Type    string          `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
}

type GroupLPSubs struct {
	events map[string]func(m *GroupLPUpdates)
}

type MsgLPUserTyping struct {
	UserID int
}

type MsgLPNewMessage struct {
	MessageID  int
	Flags      int
	Type       string
	PeerID     int
	Timestamp  int64
	Text       string
	Attachment map[string]string
}

type MsgLPSubs struct {
	events map[int]func(m interface{})
}

type Error struct {
	Code        int
	Description string
}

type Button struct {
	Action map[string]interface{} `json:"action"`
	Color  interface{}            `json:"color"`
}

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
	Inline  bool       `json:"inline"`
}

type User struct {
	ID              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Deactivated     string `json:"deactivated"`
	IsClosed        bool   `json:"is_closed"`
	CanAccessClosed bool   `json:"can_access_closed"`
	// todo
}

type Wall struct {
	Count int    `json:"count"`
	Items []Post `json:"items"`
}

type Photo struct {
	ID      int    `json:"id"`
	AlbumID int    `json:"album_id"`
	OwnerID int    `json:"owner_id"`
	UserID  int    `json:"user_id"`
	Text    string `json:"text"`
	Date    int64  `json:"date"`
	Sizes   []Size `json:"sizes"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type Video struct {
	ID            int          `json:"id"`
	OwnerID       int          `json:"owner_id"`
	Title         string       `json:"title"`
	TrackCode     string       `json:"track_code"`
	Type          string       `json:"type"`
	Views         int          `json:"views"`
	LocalViews    int          `json:"local_views"`
	Platform      string       `json:"platform"`
	AccessKey     string       `json:"access_key"`
	CanComment    int          `json:"can_comment"`
	CanEdit       int          `json:"can_edit"`
	CanLike       int          `json:"can_like"`
	CanRepost     int          `json:"can_repost"`
	CanSubscribe  int          `json:"can_subscribe"`
	CanAddToFaves int          `json:"can_add_to_faves"`
	CanAdd        int          `json:"can_add"`
	CanAttachLink int          `json:"can_attach_link"`
	Comments      int          `json:"comments"`
	Date          int          `json:"date"`
	Description   string       `json:"description"`
	Duration      int          `json:"duration"`
	Image         []VideoImage `json:"image"`
	Files         []VideoFile  `json:"files"`
}

type VideoImage struct {
	Height      int    `json:"height"`
	URL         string `json:"url"`
	Width       int    `json:"width"`
	WithPadding int    `json:"with_padding"`
}

type VideoFile struct {
	Mp4240 string `json:"mp4_240"`
	Mp4360 string `json:"mp4_360"`
	Mp4480 string `json:"mp4_480"`
	Mp4720 string `json:"mp4_720"`
}

type Size struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Videos struct {
	Count int     `json:"count"`
	Items []Video `json:"items"`
}
