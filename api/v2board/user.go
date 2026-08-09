package panel

import (
	"context"
	"errors"
	"fmt"

	"encoding/json"
)

type OnlineUser struct {
	UID int
	IP  string
}

type UserInfo struct {
	Id          int    `json:"id"`
	Uuid        string `json:"uuid"`
	SpeedLimit  int    `json:"speed_limit"`
	DeviceLimit int    `json:"device_limit"`
}

type UserListBody struct {
	Users []UserInfo `json:"users"`
}

type AliveMap struct {
	Alive map[int]int `json:"alive"`
}

// GetUserList will pull user from panel ({prefix}/u).
func (c *Client) GetUserList(ctx context.Context) ([]UserInfo, error) {
	e, err := c.buildE()
	if err != nil {
		return nil, err
	}
	r, err := c.client.R().
		SetContext(ctx).
		SetHeader("If-None-Match", c.userEtag).
		SetQueryParam("e", e).
		ForceContentType("application/json").
		Get(c.actionPath("u"))
	if err != nil {
		return nil, err
	}
	if r == nil || r.RawResponse == nil {
		return nil, fmt.Errorf("received nil response or raw response")
	}
	defer r.RawResponse.Body.Close()

	if r.StatusCode() == 304 {
		return nil, nil
	}
	plain, err := c.decryptResponseBody(r.Body())
	if err != nil {
		return nil, fmt.Errorf("decrypt user list error: %w", err)
	}
	userlist := &UserListBody{}
	if err := json.Unmarshal(plain, userlist); err != nil {
		return nil, fmt.Errorf("decode user list error: %w", err)
	}
	c.userEtag = r.Header().Get("ETag")
	return userlist.Users, nil
}

// GetUserAlive will fetch the alive_ip count for users ({prefix}/l).
func (c *Client) GetUserAlive(ctx context.Context) (map[int]int, error) {
	c.AliveMap = &AliveMap{}
	e, err := c.buildE()
	if err != nil {
		c.AliveMap.Alive = make(map[int]int)
		return c.AliveMap.Alive, nil
	}
	r, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("e", e).
		ForceContentType("application/json").
		Get(c.actionPath("l"))
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		c.AliveMap.Alive = make(map[int]int)
		return c.AliveMap.Alive, nil
	}
	if r == nil || r.RawResponse == nil || r.StatusCode() >= 399 {
		c.AliveMap.Alive = make(map[int]int)
		return c.AliveMap.Alive, nil
	}
	defer r.RawResponse.Body.Close()
	plain, err := c.decryptResponseBody(r.Body())
	if err != nil {
		fmt.Printf("decrypt user alive list error: %s", err)
		c.AliveMap.Alive = make(map[int]int)
		return c.AliveMap.Alive, nil
	}
	if err := json.Unmarshal(plain, c.AliveMap); err != nil {
		fmt.Printf("unmarshal user alive list error: %s", err)
		c.AliveMap.Alive = make(map[int]int)
	}

	return c.AliveMap.Alive, nil
}

type UserTraffic struct {
	UID      int
	Upload   int64
	Download int64
}

// ReportUserTraffic reports the user traffic ({prefix}/p).
func (c *Client) ReportUserTraffic(ctx context.Context, userTraffic []UserTraffic) error {
	data := make(map[int][]int64, len(userTraffic))
	for i := range userTraffic {
		data[userTraffic[i].UID] = []int64{userTraffic[i].Upload, userTraffic[i].Download}
	}
	body, err := c.encryptRequestBody(data)
	if err != nil {
		return err
	}
	e, err := c.buildE()
	if err != nil {
		return err
	}
	_, err = c.client.R().
		SetContext(ctx).
		SetQueryParam("e", e).
		SetBody(body).
		ForceContentType("application/json").
		Post(c.actionPath("p"))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ReportNodeOnlineUsers(ctx context.Context, data *map[int][]string) error {
	body, err := c.encryptRequestBody(data)
	if err != nil {
		return err
	}
	e, err := c.buildE()
	if err != nil {
		return err
	}
	_, err = c.client.R().
		SetContext(ctx).
		SetQueryParam("e", e).
		SetBody(body).
		ForceContentType("application/json").
		Post(c.actionPath("a"))

	if err != nil {
		return err
	}

	return nil
}
