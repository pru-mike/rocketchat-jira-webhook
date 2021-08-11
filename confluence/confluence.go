package confluence

import (
	"bytes"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"github.com/pru-mike/rocketchat-jira-webhook/utils"
	"regexp"
	"text/template"
)

const findPagesByViewID = `{{ .URL }}[^\s]+?Id=(\d+)`
const findPagesBySpaceTitle = `{{ .URL }}/display/([^/]+)/([^\s]+)`

func makePagesRegexp(config *config.Confluence) *regexp.Regexp {
	if config.FindPagesByViewID == "" {
		config.FindPagesByViewID = findPagesByViewID
	}

	reTmpl := template.Must(template.New("pagesByViewID").Parse(config.FindPagesByViewID))
	var re bytes.Buffer
	err := reTmpl.Execute(&re, config)
	if err != nil {
		logger.Fatal("can't execute pagesByViewID %v", err)
	}

	return regexp.MustCompile(re.String())
}

func makeSpaceTitlesRegexp(config *config.Confluence) *regexp.Regexp {
	if config.FindPagesBySpaceTitle == "" {
		config.FindPagesBySpaceTitle = findPagesBySpaceTitle
	}

	reTmpl := template.Must(template.New("pagesBySpaceTitle").Parse(config.FindPagesBySpaceTitle))
	var re bytes.Buffer
	err := reTmpl.Execute(&re, config)
	if err != nil {
		logger.Fatal("can't execute pagesBySpaceTitle %v", err)
	}

	return regexp.MustCompile(re.String())
}

type SpaceTitle struct {
	Space string
	Title string
}

func (c *Confluence) FindPagesSpaceTitles(text string) []SpaceTitle {
	match := c.reSpaceTitle.FindAllStringSubmatch(text, -1)
	if match == nil {
		return []SpaceTitle{}
	}
	set := make(map[string]struct{})
	var res []SpaceTitle
	for _, m := range match {
		if _, ok := set[m[0]]; !ok {
			res = append(res, SpaceTitle{m[1], m[2]})
			set[m[0]] = struct{}{}
		}
	}
	return res
}

func (c *Confluence) FindPagesViewIDs(text string) []string {
	match := c.reViewID.FindAllStringSubmatch(text, -1)
	if match == nil {
		return []string{}
	}
	res := make([]string, len(match))
	for i, m := range match {
		res[i] = m[1]
	}
	return utils.Uniq(res)
}

func (c *Confluence) FindPagesIDs(text string) []string {
	ids := c.FindPagesViewIDs(text)
	for _, spaceTitle := range c.FindPagesSpaceTitles(text) {
		id, err := c.GetPageIdBySpaceTitle(spaceTitle)
		if err != nil {
			logger.Error(err)
		} else {
			ids = append(ids, id)
		}
	}
	return ids
}
