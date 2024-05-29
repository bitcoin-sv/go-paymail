package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// ExampleClient_GetOutputsTemplate example using GetOutputsTemplate()
//
// See more examples in /examples/
func ExampleClient_GetOutputsTemplate() {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Setup a mock HTTP client
	client := newTestClient(nil)
	mockPIKEOutputs(http.StatusOK)

	// Assume we have a PIKE Outputs URL
	pikeOutputsURL := "https://test.com/v1/bsvalias/pike/outputs/{alias}@{domain.tld}"

	// Prepare the payload
	payload := &PikePaymentOutputsPayload{
		SenderPaymail: "joedoe@example.com",
		Amount:        1000, // Example amount in satoshis
	}

	// Get the outputs template from PIKE
	outputs, err := client.GetOutputsTemplate(pikeOutputsURL, "alias", "domain.tld", payload)
	if err != nil {
		fmt.Printf("error getting outputs template: %s", err.Error())
		return
	}
	fmt.Printf("found outputs template: %+v", outputs)
	// Output: found outputs template: &{URL:https://example.com/outputs}
}

// TestClient_GetOutputsTemplate will test the method GetOutputsTemplate()
func TestClient_GetOutputsTemplate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := newTestClient(t)

	t.Run("successful PIKE outputs response", func(t *testing.T) {
		mockPIKEOutputs(http.StatusOK)

		outputsURL := "https://" + testDomain + "/v1/bsvalias/pike/outputs/{alias}@{domain.tld}"
		payload := &PikePaymentOutputsPayload{
			SenderPaymail: "joedoe@example.com",
			Amount:        1000,
		}
		response, err := client.GetOutputsTemplate(outputsURL, "alias", "domain.tld", payload)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "https://example.com/outputs", response.URL)
	})

	t.Run("PIKE outputs response error", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, "https://example.com/v1/bsvalias/pike/outputs/alias@domain.tld",
			httpmock.NewStringResponder(http.StatusBadRequest, `{"message": "bad request"}`),
		)

		outputsURL := "https://example.com/v1/bsvalias/pike/outputs/{alias}@{domain.tld}"
		payload := &PikePaymentOutputsPayload{
			SenderPaymail: "joedoe@example.com",
			Amount:        1000,
		}
		response, err := client.GetOutputsTemplate(outputsURL, "alias", "domain.tld", payload)
		require.Error(t, err)
		require.Nil(t, response)
	})
}

// mockPIKEOutputs is used for mocking the PIKE outputs response
func mockPIKEOutputs(statusCode int) {
	httpmock.RegisterResponder(http.MethodPost, "https://"+testDomain+"/v1/bsvalias/pike/outputs/alias@domain.tld",
		httpmock.NewStringResponder(
			statusCode,
			`{
                "url": "https://example.com/outputs"
            }`,
		),
	)
}

// BenchmarkClient_GetOutputsTemplate benchmarks the method GetOutputsTemplate()
func BenchmarkClient_GetOutputsTemplate(b *testing.B) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := newTestClient(nil)
	mockPIKEOutputs(http.StatusOK)

	outputsURL := "https://example.com/v1/bsvalias/pike/outputs/{alias}@{domain.tld}"
	payload := &PikePaymentOutputsPayload{
		SenderPaymail: "joedoe@example.com",
		Amount:        1000, // Example amount in satoshis
	}

	for i := 0; i < b.N; i++ {
		_, _ = client.GetOutputsTemplate(outputsURL, "alias", "domain.tld", payload)
	}
}
