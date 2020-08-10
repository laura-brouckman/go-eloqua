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
	AccessedAt            string           `json:"accessedAt,omitempty"`
	ContactId             string           `json:"contacdId,omitempty,string"`
	Contacts            	[]Contact        `json:"contacdIds,omitempty"`
	CreatedAt             string           `json:"createdAt,omitempty"`
	CreatedBy             string           `json:"createdBy,omitempty,string"`
	CurrentStatus         string           `json:"currentStatus,omitempty,string"`
	Depth                 string           `json:"depth,omitempty,string"`
	Description           string           `json:"description,omitempty"`
	Email                 *Email           `json:"email,omitempty"`
	EndAt	                string           `json:"endAt,omitempty,string"`
	FailedSendCount       string           `json:"failedSendCount,omitempty"`
	FolderId              string           `json:"folderId,omitempty"`
	ID                    string           `json:"id,omitempty,string"`
	Name                  string           `json:"name,omitempty,string"`
	Permissions           []string         `json:"permissions,omitempty"`
	ScheduledFor          string           `json:"scheduledFor,omitempty"`
	SentContent           string           `json:"sentContent,omitempty"`
	SendOptions           *SendOptions     `json:"sendOptions,omitempty"` 
	SourceTemplateId      string           `json:"sourceTemplateId,omitempty,string"`
	SuccessfulSendCount   string           `json:"successfulSendCount,omitempty,string"`
	Type                  string           `json:"type,omitempty,string"`
	UpdatedAt             string           `json:"updatedAt,omitempty,string"`
	UpdatedBy             string           `json:"updatedBy,omitempty"`
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

