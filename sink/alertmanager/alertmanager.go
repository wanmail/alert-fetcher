package alertmanager

import (
	"log/slog"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/wanmail/alert-fetcher/label"
)

type ClientConfig struct {
	Host     string `json:"host"`
	BasePath string `json:"basePath"`
}

type Client struct {
	client *client.AlertmanagerAPI
}

func NewClient(cfg ClientConfig) (*Client, error) {
	c := client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
		Host:     cfg.Host,
		BasePath: cfg.BasePath,
	})

	return &Client{client: c}, nil
}

func (c *Client) Send(msg label.Message) error {
	if name, ok := msg.Labels["alertname"]; (!ok) || name == "" {
		msg.Labels["alertname"] = msg.ID
	}

	ok, err := c.client.Alert.PostAlerts(
		&alert.PostAlertsParams{
			Alerts: []*models.PostableAlert{
				{
					Annotations: msg.Annotations,
					Alert: models.Alert{
						Labels: msg.Labels,
					},
				},
			},
		},
	)

	if err != nil {
		return err
	}

	if !ok.IsSuccess() {
		return errors.Errorf("alertmanager returned [%d][%s]", ok.Code(), ok.Error())
	}

	return nil
}

func (c *Client) AsyncSend(ch <-chan label.Message) {
	for msg := range ch {
		err := c.Send(msg)
		if err != nil {
			slog.Warn("Failed to send alert", "ID", msg.ID, "error", err)
		}
	}
}
