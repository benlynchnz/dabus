package main

import "testing"

type MockHTTPClient struct {
	Message *SlackMessage
}

func (c *MockHTTPClient) PostJSON(url string, data interface{}) error {
	c.Message, _ = data.(*SlackMessage)
	return nil
}

var mockHTTPclient = &MockHTTPClient{}

var slack = &Slack{
	Team:    "SampleTeam",
	Channel: "SampleRoom",
	Token:   "foo",
	Active:  true,
	Failed:  true,
	Restart: true,
}

func Test_SendSlackCommon(t *testing.T) {
	slack.SendWithClient(mockHTTPclient, event)

	if mockHTTPclient.Message.Channel != "SampleRoom" {
		t.Errorf("Invalid slack room.")
	}

	if mockHTTPclient.Message.Username != "Systemd" {
		t.Errorf("Invalid slack username.")
	}
}

func Test_SendSlackActive(t *testing.T) {
	slack.SendWithClient(mockHTTPclient, event)

	if mockHTTPclient.Message.Attachments[0].Color != "good" {
		t.Errorf("Invalid slack color.")
	}

	if mockHTTPclient.Message.Attachments[0].Text != "Service *foo.service* is active" {
		t.Errorf("Invalid slack text.")
	}

	if mockHTTPclient.Message.Attachments[0].Fallback != "Service *foo.service* is active" {
		t.Errorf("Invalid slack fallback.")
	}
}

func Test_SendSlackFailed(t *testing.T) {
	event.ActiveStatus = "failed"
	slack.SendWithClient(mockHTTPclient, event)

	if mockHTTPclient.Message.Attachments[0].Color != "danger" {
		t.Errorf("Invalid slack color.")
	}

	if mockHTTPclient.Message.Attachments[0].Text != "Service *foo.service* failed" {
		t.Errorf("Invalid slack text.")
	}

	if mockHTTPclient.Message.Attachments[0].Fallback != "Service *foo.service* failed" {
		t.Errorf("Invalid slack fallback.")
	}
}

func Test_SendSlackRestart(t *testing.T) {
	event.ActiveStatus = "activating"
	event.SubStatus = "auto-restart"
	slack.SendWithClient(mockHTTPclient, event)

	if mockHTTPclient.Message.Attachments[0].Color != "warning" {
		t.Errorf("Invalid slack color.")
	}

	if mockHTTPclient.Message.Attachments[0].Text != "Service *foo.service* is auto-restarted" {
		t.Errorf("Invalid slack text.")
	}

	if mockHTTPclient.Message.Attachments[0].Fallback != "Service *foo.service* is auto-restarted" {
		t.Errorf("Invalid slack fallback.")
	}
}
