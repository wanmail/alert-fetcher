package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/searchapplication/search"
)

type Config struct {
	ClientConfig
	QueryConfig
}

type QueryConfig struct {
	Index     string `json:"index"`
	TimeQuery string `json:"timeQuery"`
	Query     string `json:"query"`
	Max       int    `json:"max"`
}

type ClientConfig struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Client struct {
	QueryConfig
	client *elasticsearch.Client
}

func NewElasticSource(ccfg ClientConfig, qcfg QueryConfig) (*Client, error) {
	ec, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{ccfg.Address},
			Username:  ccfg.Username,
			Password:  ccfg.Password,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Client{client: ec, QueryConfig: qcfg}, nil
}

var timeQuery = `
{
	"query": {
		"bool": {
			"must": [
				%s,
				{
                    "range":{
                        "@timestamp":{
                            "format":"strict_date_optional_time",
                            "gte":"%s",
                            "lte":"%s"
                        }
                    }
                }
			]
		}
	}
}
`

func (c *Client) FetchAll(ctx context.Context, from time.Time, now time.Time) ([]map[string]interface{}, error) {
	var query string

	// TODO: add timeQuery
	if c.TimeQuery == "" {
		query = fmt.Sprintf(timeQuery, c.Query, from.Format(time.RFC3339), now.Format(time.RFC3339))
	} else {
		query = c.Query
	}

	slog.DebugContext(ctx, "elasticsearch query", "query", query)
	resp, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.Index),
		c.client.Search.WithBody(strings.NewReader(query)),
		c.client.Search.WithSize(c.Max),
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to search")
	}

	if resp.StatusCode/100 != 2 {
		return nil, errors.Errorf("failed to search[%d][%s]", resp.StatusCode, resp.String())
	}

	response := search.NewResponse()

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response[%s]", resp.String())
	}

	outs := []map[string]interface{}{}
	for _, hit := range response.Hits.Hits {
		val := make(map[string]interface{})
		if err = json.Unmarshal(hit.Source_, &val); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal hit")
		}
		outs = append(outs, val)
	}

	return outs, nil
}

// TODO: implement this
func (c *Client) FetchOne(ctx context.Context, from time.Time, now time.Time) (map[string]interface{}, error) {
	return nil, nil
}
