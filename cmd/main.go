package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/pru-mike/rocketchat-jira-webhook/app"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config")
	flag.Parse()
	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetLevelFromString(cfg.App.LogLevel)
	logger.Debugf("configuration '%+v'", *cfg)

	app := app.New(cfg)

	http.Handle("/health", http.HandlerFunc(app.Health))
	http.Handle("/jira", http.HandlerFunc(app.Jira))
	http.Handle("/confluence", http.HandlerFunc(app.Confluence))
	http.Handle("/jiraconfluence", http.HandlerFunc(app.JiraConfluence))
	logger.Infof("start listening on '%s'", cfg.ListenAddr())
	logger.Fatal(http.ListenAndServe(cfg.ListenAddr(), nil))
}
