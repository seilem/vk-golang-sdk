package vkapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type GroupGetLPServerReq struct {
	GroupID int
}

func (GroupGetLPServerReq) Name() string {
	return "groups.getLongPollServer"
}

func (g *GroupGetLPServerReq) Values() url.Values {
	v := url.Values{}

	if g.GroupID > 0 {
		v.Set("group_id", strconv.Itoa(g.GroupID))
	}
	return v
}

// GroupGetLPServer returns data for Bots Long Poll API connection.
//
// See https://vk.com/dev/groups.getLongPollServer
func (vk *VkAPI) GroupGetLPServer(v *GroupGetLPServerReq) (*LPServer, error) {
	resp, err := vk.MakeRequest(v.Name(), v.Values())
	if err != nil {
		return nil, err
	}

	var apiResp LPServer
	if err := json.Unmarshal(resp.Response, &apiResp); err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func (vk *VkAPI) GroupLPServ(groupID int) error {
	return vk.groupLongPoll(context.Background(), groupID)
}

func (vk *VkAPI) GroupLPCallback(name string, f func(event *GroupLPUpdates)) {
	if _, exists := vk.groupLPSubs.events[name]; !exists {
		vk.groupLPSubs.events[name] = f
	}
}

func (vk *VkAPI) groupLongPoll(ctx context.Context, groupID int) error {
	server, err := vk.GroupGetLPServer(&GroupGetLPServerReq{
		GroupID: groupID,
	})
	if err != nil {
		log.Print("Get group lp server is failed: %w", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		serverURL := fmt.Sprintf(
			"%s?act=a_check&key=%s&ts=%s&wait=25",
			server.Server, server.Key, server.TS)
		req, err := http.NewRequest(http.MethodGet, serverURL, nil)
		if err != nil {
			log.Print("Group lp server connection failed: %w", err)
			continue
		}

		resp, err := vk.Client.Do(req)
		if err != nil {
			log.Print("Group lp request connection failed: %w", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)

		var e GroupLPEvent
		if err := json.Unmarshal(body, &e); err != nil {

		}

		switch true {
		case e.Failed == 0:
			for _, update := range e.Updates {
				vk.handleGroupLPCallback(update.Type, &update)
			}
			server.TS = e.TS
		case e.Failed == 1:
			server.TS = e.TS
		case e.Failed == 2 || e.Failed == 3:
			newServer, err := vk.GroupGetLPServer(&GroupGetLPServerReq{
				GroupID: groupID,
			})
			if err != nil {
				log.Print("Get group lp server is failed: %w", err)
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

func (vk *VkAPI) handleGroupLPCallback(name string, event *GroupLPUpdates) {
	_, exists := vk.groupLPSubs.events[name]

	if exists {
		vk.groupLPSubs.events[name](event)
	}
}
