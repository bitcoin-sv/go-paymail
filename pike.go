package paymail

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type PikeContactRequestResponse struct {
	StandardResponse
}

type PikeContactRequestPayload struct {
	FullName      string
	PaymailAdress string
}

func (c *Client) AddContactRequest(url, alias, domain string, request *PikeContactRequestPayload) (response *PikeContactRequestResponse, err error) {

	if err = c.validateUrlWithPaymail(url, alias, domain); err != nil {
		return
	}

	if err = request.validate(); err != nil {
		return
	}

	// Set the base url and path, assuming the url is from the prior GetCapabilities() request
	// https://<host-discovery-target>/{alias}@{domain.tld}/id
	reqURL := replaceAliasDomain(url, alias, domain)

	response = &PikeContactRequestResponse{}
	if response.StandardResponse, err = c.postRequest(reqURL, request); err != nil {
		return
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		// Paymail address not found?
		if response.StatusCode == http.StatusNotFound {
			err = errors.New("paymail address not found")
		} else {
			err = c.prepareServerErrorResponse(&response.StandardResponse)
		}

		return nil, err
	}

	return response, nil
}

func (c *Client) validateUrlWithPaymail(url, alias, domain string) error {
	if len(url) == 0 || !strings.Contains(url, "https://") {
		return fmt.Errorf("invalid url: %s", url)
	} else if len(alias) == 0 {
		return errors.New("missing alias")
	} else if len(domain) == 0 {
		return errors.New("missing domain")
	}
	return nil
}

func (c *Client) prepareServerErrorResponse(response *StandardResponse) (err error) {
	serverError := &ServerError{}
	if err = json.Unmarshal(response.Body, serverError); err != nil {
		return
	}
	if len(serverError.Message) == 0 {
		err = fmt.Errorf("bad response from paymail provider: code %d, body: %s", response.StatusCode, string(response.Body))
	} else {
		err = fmt.Errorf("bad response from paymail provider: code %d, message: %s", response.StatusCode, serverError.Message)
	}

	return
}

func (r *PikeContactRequestPayload) validate() error {
	if r.FullName == "" {
		return errors.New("missing full name")
	}
	if r.PaymailAdress == "" {
		return errors.New("missing paymail address")
	}

	return ValidatePaymail(r.PaymailAdress)
}
