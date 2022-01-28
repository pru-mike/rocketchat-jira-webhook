package confluence

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pru-mike/rocketchat-jira-webhook/client"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
)

const apiPrefix = "/rest/api"

type Confluence struct {
	baseURL      string
	reViewID     *regexp.Regexp
	reSpaceTitle *regexp.Regexp
	client       *client.Client
	expand       string
}

func NewClient(config *config.Confluence) *Confluence {
	expand := []string{"history.lastUpdated", "space", "body." + config.BodyExpand}
	return &Confluence{
		config.URL,
		makePagesRegexp(config),
		makeSpaceTitlesRegexp(config),
		client.New(config.Username, config.Password, config.Timeout),
		strings.Join(expand, ","),
	}
}

func (c *Confluence) makePageURL(pageID string) string {
	return fmt.Sprintf("%s%s/content/%s", c.baseURL, apiPrefix, pageID)
}

func (c *Confluence) makeRequest(ctx context.Context, url string, expand bool, data interface{}) error {
	var queryString map[string]string
	if expand {
		queryString = make(map[string]string)
		queryString["expand"] = c.expand
	}
	return c.client.MakeGETRequest(ctx, url, queryString, data)
}

func (c *Confluence) GetPageCtx(ctx context.Context, pageID string) (*Page, error) {
	var page Page
	return &page, c.makeRequest(ctx, c.makePageURL(pageID), true, &page)
}

var _ client.KeyGetter = (*Confluence)(nil)

func (c *Confluence) GetByKey(ctx context.Context, pageID string) (interface{}, error) {
	return c.GetPageCtx(ctx, pageID)
}

func (c *Confluence) GetPages(pageIDs []string) ([]*Page, error) {
	var pages []*Page
	data, err := c.client.BulkKeyRequests(context.Background(), c, pageIDs)
	for _, d := range data {
		pages = append(pages, d.(*Page))
	}
	return pages, err
}

func (c *Confluence) makeContentPageSearchURL() string {
	return fmt.Sprintf("%s%s/content", c.baseURL, apiPrefix)
}

func (c *Confluence) GetPageIdBySpaceTitle(st SpaceTitle) (string, error) {
	queryString := make(map[string]string)
	queryString["space"] = st.Space
	queryString["title"] = st.Title
	var search PageSearchResult
	err := c.client.MakeGETRequest(context.Background(), c.makeContentPageSearchURL(), queryString, &search)
	if err != nil {
		return "", err
	}
	if len(search.Results) != 1 {
		return "", fmt.Errorf("found %d records by space:%s title:%s", len(search.Results), st.Space, st.Title)
	}
	return search.Results[0].ID, nil
}

func (c *Confluence) makeCurrentUserURL() string {
	return fmt.Sprintf("%s%s/user/current", c.baseURL, apiPrefix)
}

func (c *Confluence) GetCurrentUser() (string, error) {
	var user User
	err := c.client.MakeGETRequest(context.Background(), c.makeCurrentUserURL(), nil, &user)
	if err != nil {
		return "", err
	}
	return user.DisplayName, nil
}
