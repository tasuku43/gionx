package wscreate

import (
	"context"
	"fmt"
)

type JiraIssue struct {
	Key     string
	Summary string
}

type JiraIssuePort interface {
	FetchIssueByTicketURL(ctx context.Context, ticketURL string) (JiraIssue, error)
}

type JiraWorkspaceInput struct {
	ID        string
	Title     string
	SourceURL string
}

type Service struct {
	jiraPort JiraIssuePort
}

func NewService(jiraPort JiraIssuePort) *Service {
	return &Service{jiraPort: jiraPort}
}

func (s *Service) ResolveJiraWorkspaceInput(ctx context.Context, ticketURL string) (JiraWorkspaceInput, error) {
	if s.jiraPort == nil {
		return JiraWorkspaceInput{}, fmt.Errorf("jira issue port is not configured")
	}
	issue, err := s.jiraPort.FetchIssueByTicketURL(ctx, ticketURL)
	if err != nil {
		return JiraWorkspaceInput{}, err
	}
	if issue.Key == "" {
		return JiraWorkspaceInput{}, fmt.Errorf("jira issue key is empty")
	}
	return JiraWorkspaceInput{
		ID:        issue.Key,
		Title:     issue.Summary,
		SourceURL: ticketURL,
	}, nil
}
