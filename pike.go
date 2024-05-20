package paymail

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PikeContactRequestResponse struct {
	StandardResponse
}

type PikeContactRequestPayload struct {
	FullName string `json:"fullName"`
	Paymail  string `json:"paymail"`
}

// TODO: check if everything is needed after whole PIKE implementation
type PikePaymentDestinationsRequest struct {
	SenderName    string    `json:"senderName"`
	SenderPaymail string    `json:"senderPaymail"`
	Amount        uint64    `json:"amount"`
	Dt            time.Time `json:"dt"`
	Reference     string    `json:"reference"`
	Signature     string    `json:"signature"`
}

type PikePaymentDestinationsResponse struct {
	Outputs   []PikePaymentDestination `json:"outputs"`
	Reference string                   `json:"reference"`
}

type PikePaymentDestination struct {
	Script   string `json:"script"`
	Satoshis int    `json:"satoshis"`
}

func (c *Client) AddContactRequest(url, alias, domain string, request *PikeContactRequestPayload) (*PikeContactRequestResponse, error) {

	if err := c.validateUrlWithPaymail(url, alias, domain); err != nil {
		return nil, err
	}

	if err := request.validate(); err != nil {
		return nil, err
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/{alias}@{domain.tld}/id
	reqURL := replaceAliasDomain(url, alias, domain)

	response, err := c.postRequest(reqURL, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		if response.StatusCode == http.StatusNotFound {
			return nil, errors.New("paymail address not found")
		} else {
			return nil, c.prepareServerErrorResponse(&response)
		}
	}

	return &PikeContactRequestResponse{response}, nil
}

func (c *Client) validateUrlWithPaymail(url, alias, domain string) error {
	if len(url) == 0 || !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid url: %s", url)
	} else if alias == "" {
		return errors.New("missing alias")
	} else if domain == "" {
		return errors.New("missing domain")
	}
	return nil
}

func (c *Client) prepareServerErrorResponse(response *StandardResponse) error {
	var details string

	serverError := &ServerError{}
	if err := json.Unmarshal(response.Body, serverError); err != nil || serverError.Message == "" {
		details = fmt.Sprintf("body: %s", string(response.Body))
	} else {
		details = fmt.Sprintf("message: %s", serverError.Message)
	}

	return fmt.Errorf("bad response from paymail provider: code %d, %s", response.StatusCode, details)
}

func (r *PikeContactRequestPayload) validate() error {
	if r.FullName == "" {
		return errors.New("missing full name")
	}
	if r.Paymail == "" {
		return errors.New("missing paymail address")
	}

	return ValidatePaymail(r.Paymail)
}
