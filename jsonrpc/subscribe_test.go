package jsonrpc

import (
	"strings"
	"testing"
)

func TestSubscribeNewHead(t *testing.T) {
	addr := "wss://arbitrum-one-rpc.publicnode.com"
	if strings.HasPrefix(addr, "http") {
		return
	}

	c, _ := NewClient(addr)
	defer c.Close()

	data := make(chan []byte)
	params := map[string]interface{}{
		"address": "0x641C00A822e8b671738d32a431a4Fb6074E5c79d",
		"topics":  []string{"0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"},
	}
	_, err := c.Subscribe("logs", params, func(b []byte) {
		data <- b
	})
	if err != nil {
		t.Fatal(err)
	}

}
