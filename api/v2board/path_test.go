package panel

import (
	"testing"

	"github.com/wyx2685/v2node/conf"
)

func TestNormalizeApiPrefix(t *testing.T) {
	if got := normalizeApiPrefix("n/abc/"); got != "/n/abc" {
		t.Fatalf("got %s", got)
	}
	if got := normalizeApiPrefix(""); got != "" {
		t.Fatalf("empty should stay empty")
	}
}

func TestNewRequiresApiPrefix(t *testing.T) {
	_, err := New(&conf.NodeConfig{
		APIHost: "https://example.com",
		NodeID:  1,
		Key:     "1234567890123456",
	})
	if err == nil {
		t.Fatal("expected error without ApiPrefix")
	}
}

func TestActionPath(t *testing.T) {
	c, err := New(&conf.NodeConfig{
		APIHost:   "https://example.com",
		NodeID:    1,
		Key:       "1234567890123456",
		ApiPrefix: "/n/testhostpath",
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.actionPath("u") != "/n/testhostpath/u" {
		t.Fatalf("path=%s", c.actionPath("u"))
	}
	e, err := c.buildE()
	if err != nil || e == "" || e == c.Token {
		t.Fatalf("buildE failed: %v %s", err, e)
	}
}
