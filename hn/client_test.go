package hn

import "testing"

func TestClient_defaultify(t *testing.T) {
	var c Client
	c.defaultify()

	got := c.apiBase
	want := apiBase

	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}
