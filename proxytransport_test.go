package proxytransport_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/deploys-app/proxytransport"
)

func TestTransport(t *testing.T) {
	client := http.Client{
		Transport: &proxytransport.Transport{},
	}

	resp, err := client.Get("https://icanhazip.com/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	t.Logf("ip: %s", buf.String())
	if !strings.Contains(buf.String(), "2a06:98c0:3600::103") {
		t.Errorf("wrong ip: %s", buf.String())
	}
}
