package panel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"github.com/wyx2685/v2node/common/crypt"
	"github.com/wyx2685/v2node/conf"
)

type Client struct {
	client           *resty.Client
	APIHost          string
	Token            string
	NodeId           int
	ApiPrefix        string
	workingKey       []byte
	nodeEtag         string
	userEtag         string
	responseBodyHash string
	UserList         *UserListBody
	AliveMap         *AliveMap
}

func New(c *conf.NodeConfig) (*Client, error) {
	client := resty.New()
	retryCount := conf.DefaultNodeRetryCount
	if c.RetryCount != nil {
		retryCount = *c.RetryCount
	}
	client.SetRetryCount(retryCount)
	client.SetHeader("User-Agent", fmt.Sprintf("v2node go-resty/%s (https://github.com/go-resty/resty)", resty.Version))
	if c.Timeout > 0 {
		client.SetTimeout(time.Duration(c.Timeout) * time.Second)
	} else {
		client.SetTimeout(time.Duration(conf.DefaultNodeTimeout) * time.Second)
	}
	client.OnError(func(req *resty.Request, err error) {
		var v *resty.ResponseError
		if errors.As(err, &v) {
			logrus.Error(v.Err)
		}
	})
	client.SetBaseURL(c.APIHost)
	prefix := normalizeApiPrefix(c.ApiPrefix)
	if prefix == "" {
		return nil, fmt.Errorf("ApiPrefix is required")
	}
	if strings.TrimSpace(c.Key) == "" {
		return nil, fmt.Errorf("ApiKey is required")
	}
	return &Client{
		client:     client,
		Token:      c.Key,
		APIHost:    c.APIHost,
		NodeId:     c.NodeID,
		ApiPrefix:  prefix,
		workingKey: crypt.DeriveNodeWorkingKey(c.Key),
		UserList:   &UserListBody{},
		AliveMap:   &AliveMap{},
	}, nil
}

func normalizeApiPrefix(prefix string) string {
	p := strings.TrimSpace(prefix)
	if p == "" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	for len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}
	return p
}

func (c *Client) actionPath(action string) string {
	return c.ApiPrefix + "/" + action
}

func (c *Client) buildE() (string, error) {
	identity := map[string]any{
		"k": c.Token,
		"i": c.NodeId,
		"t": "vn",
	}
	raw, err := json.Marshal(identity)
	if err != nil {
		return "", err
	}
	return crypt.EncryptCompact(raw, c.workingKey)
}

type sm4Envelope struct {
	IV      string `json:"iv"`
	Payload string `json:"payload"`
}

func (c *Client) decryptResponseBody(body []byte) ([]byte, error) {
	var env sm4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("decode sm4 envelope: %w", err)
	}
	if env.IV == "" || env.Payload == "" {
		return nil, fmt.Errorf("sm4 envelope missing iv/payload")
	}
	return crypt.DecryptEnvelope(env.IV, env.Payload, c.workingKey)
}

func (c *Client) encryptRequestBody(v any) (map[string]string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	iv, payload, err := crypt.EncryptEnvelope(raw, c.workingKey)
	if err != nil {
		return nil, err
	}
	return map[string]string{"iv": iv, "payload": payload}, nil
}

// NodeIDString helpers for tests / logging.
func (c *Client) NodeIDString() string {
	return strconv.Itoa(c.NodeId)
}
