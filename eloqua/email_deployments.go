package eloqua

import (
	"fmt"
)

// EmailDeploymentService provides access to all the endpoints related
// to email deployment assets within eloqua
type EmailDeploymentService struct {
	client *Client
}

type SendOptions struct {
	AllowResend              string   `json:"allowResend,omitempty,string"`
	AllowSendToUnsubscribe   string   `json:"allowSendToUnsubscribe,omitempty,string"`
}

// Eloqua API docs: https://goo.gl/BaqLvm

// EmailDeployment represents an Eloqua email deployment object.
type EmailDeployment struct {
	ContactId             int           `json:"contactId,omitempty,string"`
	Email                 Email         `json:"email,omitempty"`
	Name                  string        `json:"name,omitempty,string"`
	Type                  string        `json:"type,omitempty,string"`
}

// Create and send a new email deployment in eloqua
func (e *EmailDeploymentService) CreateAndSend(emailDeployment *EmailDeployment) (*EmailDeployment, *Response, error) {
	if emailDeployment == nil {
		emailDeployment = &EmailDeployment{}
	}
	endpoint := "/assets/email/deployment"
	resp, err := e.client.postRequestDecode(endpoint, emailDeployment)
	return emailDeployment, resp, err
}

// Get an email object via its ID
func (e *EmailDeploymentService) Retrieve(id int) (*EmailDeployment, *Response, error) {
	endpoint := fmt.Sprintf("/assets/email/deployment/%d?depth=complete", id)
	emailDeployment := &EmailDeployment{}
	resp, err := e.client.getRequestDecode(endpoint, emailDeployment)
	return emailDeployment, resp, err
}

