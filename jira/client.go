package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"net/http"
	"strings"
	"sync"
)

const apiPrefix = "/rest/api/2"
const browsePrefix = "/browse"

type Jira struct {
	username, password string
	apiURL             string
	browseURL          string
	httpClient         *http.Client
	requestFields      string
}

func NewClient(config *config.Jira) *Jira {
	requestFields := []string{"summary", "description"}
	requestFields = append(requestFields, config.RequestFields()...)
	return &Jira{
		username:      config.Username,
		password:      config.Password,
		httpClient:    &http.Client{Timeout: config.Timeout},
		browseURL:     fmt.Sprintf("%s%s", config.URL, browsePrefix),
		apiURL:        fmt.Sprintf("%s%s", config.URL, apiPrefix),
		requestFields: strings.Join(requestFields, ","),
	}
}

func (j *Jira) makeRequest(ctx context.Context, url string, addFields bool, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("can't create request %q: %w", url, err)
	}

	if addFields {
		req.URL.Query().Add("fields", j.requestFields)
	}

	req.SetBasicAuth(j.username, j.password)
	resp, err := j.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed %q: %s", url, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("can't read json %q: %w", url, err)
	}
	return nil
}

func (j *Jira) GetMyself() (*Myself, error) {
	return j.GetMyselfCtx(context.Background())
}

func (j *Jira) GetMyselfCtx(ctx context.Context) (*Myself, error) {
	var myself Myself
	return &myself, j.makeRequest(ctx, fmt.Sprintf("%s/myself", j.apiURL), false, &myself)
}

func (j *Jira) issueAPIURL(issueKey string) string {
	return fmt.Sprintf("%s/issue/%s", j.apiURL, issueKey)
}

func (j *Jira) GetIssueCtx(ctx context.Context, issueKey string) (*Issue, error) {
	issue := Issue{
		browseURL: j.browseURL,
	}
	return &issue, j.makeRequest(ctx, j.issueAPIURL(issueKey), true, &issue)
}

func (j *Jira) GetIssues(issueKeys []string) ([]*Issue, error) {
	return j.GetIssuesCtx(context.Background(), issueKeys)
}

func (j *Jira) GetIssuesCtx(ctx context.Context, issueKeys []string) ([]*Issue, error) {
	var wg sync.WaitGroup
	var issues []*Issue
	var errs []error
	wg.Add(len(issueKeys))
	var mu sync.Mutex
	setResult := func(issue *Issue, err error) {
		mu.Lock()
		if err != nil {
			errs = append(errs, err)
		} else {
			issues = append(issues, issue)
		}
		mu.Unlock()
	}
	for _, issue := range issueKeys {
		go func(issue string) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					setResult(nil, fmt.Errorf("panic while getting %s: %v", issue, err))
				}
			}()
			setResult(j.GetIssueCtx(ctx, issue))
		}(issue)
	}
	wg.Wait()
	var err error
	if len(errs) != 0 {
		err = fmt.Errorf("can't get issues: %v", errs)
	}
	return issues, err
}
