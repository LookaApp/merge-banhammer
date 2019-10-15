package mergeban

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Notifier interface {
	Notify(responseURL string, responseBody string)
}

type slackNotifier struct {
	httpClient *http.Client
	logger     *log.Logger
}

func NewSlackNotifier(logger *log.Logger) *slackNotifier {
	return &slackNotifier{
		httpClient: &http.Client{},
		logger:     logger,
	}
}

func (n *slackNotifier) Notify(responseURL string, responseBody string) {
	responsePayload := map[string]string{"text": responseBody}
	responseJSON, err := json.Marshal(responsePayload)
	if err != nil {
		n.logger.Printf("Failed to marshal response payload %v: %v\n", responsePayload, err)
	}

	n.logger.Printf("Sending %s to %s\n", responseBody, responseURL)

	slackResponse, err := n.httpClient.Post(responseURL, "application/json", bytes.NewReader(responseJSON))
	if err != nil {
		n.logger.Printf("Failed to send response of %s to Slack URL %s: %v\n", responseBody, responseURL, err)
	}

	defer slackResponse.Body.Close()

	// net/http expects us to read to EOF & close the response body
	_, err = ioutil.ReadAll(slackResponse.Body)
	if err != nil {
		n.logger.Printf("Failed to read Slack response body: %v\n", err)
	}
}
