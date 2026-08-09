package crypt

import (
	"encoding/json"
	"testing"
)

func TestDeriveNodeWorkingKeyLength(t *testing.T) {
	key := DeriveNodeWorkingKey("1234567890123456")
	if len(key) != 16 {
		t.Fatalf("key len=%d", len(key))
	}
}

func TestEnvelopeRoundTrip(t *testing.T) {
	key := DeriveNodeWorkingKey("abcdefghijklmnopqrstuvwxyz")
	plain := []byte(`{"users":[{"id":1}]}`)
	iv, payload, err := EncryptEnvelope(plain, key)
	if err != nil {
		t.Fatal(err)
	}
	back, err := DecryptEnvelope(iv, payload, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(back) != string(plain) {
		t.Fatalf("got %s", back)
	}
}

func TestCompactRoundTripIdentity(t *testing.T) {
	token := "1234567890123456"
	key := DeriveNodeWorkingKey(token)
	identity := map[string]any{"k": token, "i": 7, "t": "vn"}
	raw, _ := json.Marshal(identity)
	e, err := EncryptCompact(raw, key)
	if err != nil {
		t.Fatal(err)
	}
	if e == "" || e == string(raw) {
		t.Fatal("compact should be ciphertext")
	}
	back, err := DecryptCompact(e, key)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(back, &got); err != nil {
		t.Fatal(err)
	}
	if got["k"] != token || got["t"] != "vn" {
		t.Fatalf("identity mismatch: %#v", got)
	}
}
