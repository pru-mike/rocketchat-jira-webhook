package jira

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pru-mike/rocketchat-jira-webhook/client"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
)

const apiPrefix = "/rest/api/2"
const browsePrefix = "/browse"

type Jira struct {
	username, password string
	apiURL             string
	browseURL          string
	client             *client.Client
	requestFields      string
	re                 *regexp.Regexp
}

func NewClient(config *config.Jira) *Jira {
	requestFields := []string{"summary", "description"}
	requestFields = append(requestFields, config.RequestFields()...)

	re := findKeysRegexp
	if config.FindKeyRegexp != "" {
		re = regexp.MustCompile(config.FindKeyRegexp)
	}

	return &Jira{
		client:        client.New(config.Username, config.Password, config.Timeout),
		browseURL:     fmt.Sprintf("%s%s", config.URL, browsePrefix),
		apiURL:        fmt.Sprintf("%s%s", config.URL, apiPrefix),
		requestFields: strings.Join(requestFields, ","),
		re:            re,
	}
}

func (j *Jira) ParseKeys(text string) []string {
	return parseKeys(j.re, text)
}

func (j *Jira) makeRequest(ctx context.Context, url string, addFields bool, data interface{}) error {
	var queryString map[string]string
	if addFields {
		queryString = make(map[string]string)
		queryString["fields"] = j.requestFields
	}
	return j.client.MakeGETRequest(ctx, url, queryString, data)
}

func (j *Jira) issueAPIURL(issueKey string) string {
	return fmt.Sprintf("%s/issue/%s", j.apiURL, issueKey)
}

func (j *Jira) GetCurrentUser() (string, error) {
	var myself Myself
	err := j.makeRequest(context.Background(), fmt.Sprintf("%s/myself", j.apiURL), false, &myself)
	if err != nil {
		return "", err
	}
	return myself.DisplayName, nil
}

func (j *Jira) GetIssueCtx(ctx context.Context, issueKey string) (*Issue, error) {
	issue := Issue{
		browseURL: j.browseURL,
	}
	return &issue, j.makeRequest(ctx, j.issueAPIURL(issueKey), true, &issue)
}

var _ client.KeyGetter = (*Jira)(nil)

func (j *Jira) GetByKey(ctx context.Context, issueKey string) (interface{}, error) {
	return j.GetIssueCtx(ctx, issueKey)
}

func (j *Jira) GetIssues(issueKeys []string) ([]*Issue, error) {
	var issues []*Issue
	data, err := j.client.BulkKeyRequests(context.Background(), j, issueKeys)
	for _, d := range data {
		issues = append(issues, d.(*Issue))
	}
	return issues, err
}
