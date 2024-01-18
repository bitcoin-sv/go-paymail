package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateServer will test the method CreateServer()
func TestCreateServer(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := testConfig(t, "localhost")
		config.Port = 12345
		config.Timeout = 10 * time.Second
		s := CreateServer(config)
		require.NotNil(t, s)
		assert.IsType(t, &http.Server{}, s)
		assert.Equal(t, fmt.Sprintf(":%d", config.Port), s.Addr)
		assert.Equal(t, config.Timeout, s.WriteTimeout)
		assert.Equal(t, config.Timeout, s.ReadTimeout)
	})
}

// TestStart will test the method Start()
func TestStart(t *testing.T) {
	t.Run("run server", func(t *testing.T) {
		config, _ := NewConfig(new(mockServiceProvider), WithDomain("domain.com"))
		config.Prefix = "http://"

		server := httptest.NewServer(Handlers(config))
		defer server.Close()

		config.AddDomain(server.URL)

		resp, err := http.Get(fmt.Sprintf("%s/.well-known/bsvalias", server.URL))
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		resp.Body.Close()

		capabilities := result["capabilities"].(map[string]interface{})
		assert.NotNil(t, capabilities)
		assert.Greater(t, len(capabilities), 0)

		for _, cap := range capabilities {
			capUrl, ok := cap.(string)
			if !ok {
				continue
			}

			capUrl = strings.ReplaceAll(capUrl, PaymailAddressTemplate, "example@domain.com")
			capUrl = strings.ReplaceAll(capUrl, PubKeyTemplate, "xpub")

			_, err := url.Parse(capUrl)
			assert.NoError(t, err, "Endpoint %s is not a valid URL", capUrl)

			_, err = http.Get(capUrl) //only check whether this endpoint is reachable even if "get" method is not allowed
			assert.NoError(t, err)
		}

		assert.Equal(t, result["bsvalias"], config.BSVAliasVersion)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
