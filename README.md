rocketchat-jira-webhook
=======================
It's [Golang](https://golang.org/) port of [rocketchat-jira-trigger](https://github.com/gustavkarlsson/rocketchat-jira-trigger).  
Outgoing [Rocket.Chat](https://rocket.chat) webhook integration that summarizes mentioned 
[JIRA](https://www.atlassian.com/software/jira) issues.

Installation
------------
You need go 1.16 or newer.

```bash
go get -u github.com/pru-mike/rocketchat-jira-webhook
vi config.toml
rocketchat-jira-webhook -config config.toml
```

Running with docker
-------------------

```bash
git clone git@github.com:pru-mike/rocketchat-jira-webhook.git
cd rocketchat-jira-webhook
docker build . -t rocketchat-jira-webhook
cp example/minimal.toml config.toml
docker run -p 4567:4567 -v $PWD:/app rocketchat-jira-webhook
```

Configuration
-------------
Configuration is slightly differ from original [rocketchat-jira-trigger](https://github.com/gustavkarlsson/rocketchat-jira-trigger)
So it isn't substituted one by one.  
There is plenty of configuration options, but the only required parameters is jira credentials. 
For all other options there are reasonable defaults.  
For [minimal](https://github.com/pru-mike/rocketchat-jira-webhook/blob/master/example/minimal.toml)
and [all](https://github.com/pru-mike/rocketchat-jira-webhook/blob/master/example/everything.toml) 
options see configuration examples.

Usage
-----
The same as [rocketchat-jira-trigger](https://github.com/gustavkarlsson/rocketchat-jira-trigger).  

First you need to start **rocketchat-jira-webhook**

Second you need to going to Rocket.Chat administration panel and setting up outgoing webhook pointing 
at **rocketchat-jira-webhook** instance.  

And then you can write a message containing some JIRA issues. For example: `TEST-1234`  
Then **rocketchat-jira-webhook** will try to gather details about issues and reply it to Rocket.Chat if found some.

Health testing
--------------
JIRA connection can be checked with /health route.  

```bash
curl http://localhost:4567/health
{"ok":true,"jira":{"name":"JIRA Bot","error":""}}
```
