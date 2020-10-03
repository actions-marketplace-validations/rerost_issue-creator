// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/rerost/issue-creator/domain/issue"
	"github.com/rerost/issue-creator/domain/schedule"
	"github.com/rerost/issue-creator/repo"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"time"
)

// Injectors from wire.go:

func InitializeCmd(ctx context.Context, cfg Config) (*cobra.Command, error) {
	client := NewGithubClient(ctx, cfg)
	issueRepository := repo.NewIssueRepository(client)
	time := CurrentTime(cfg)
	issueService := NewIssueService(cfg, issueRepository, time)
	v := NewK8sCommand(cfg)
	scheduleRepository := repo.NewScheduleRepository(v)
	scheduleService := NewScheduleService(cfg, scheduleRepository)
	string2 := NewTemplateFile(cfg)
	command := NewCmdRoot(ctx, issueService, scheduleService, string2)
	return command, nil
}

// wire.go:

func CurrentTime(cfg Config) time.Time {
	return time.Now()
}

func NewGithubClient(ctx context.Context, cfg Config) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)
	return c
}

func NewK8sCommand(cfg Config) []string {
	return cfg.K8sCommands
}

func NewTemplateFile(cfg Config) string {
	return cfg.ManifestTemplateFile
}

func NewIssueService(cfg Config, issueRepo repo.IssueRepository, ct time.Time) issue.IssueService {
	return issue.NewIssueService(
		issueRepo,
		ct,
		cfg.CloseLastIssue,
	)
}

func NewScheduleService(cfg Config, scheduleRepository repo.ScheduleRepository) schedule.ScheduleService {
	return schedule.NewScheduleService(scheduleRepository, cfg.CloseLastIssue)
}
