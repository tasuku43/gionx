package appports

import (
	"context"
	"fmt"

	"github.com/tasuku43/kra/internal/app/wscreate"
	"github.com/tasuku43/kra/internal/infra/jira"
)

type WSCreateJiraPort struct {
	client *jira.Client
}

func NewWSCreateJiraPort() *WSCreateJiraPort {
	return NewWSCreateJiraPortWithBaseURL("")
}

func NewWSCreateJiraPortWithBaseURL(baseURL string) *WSCreateJiraPort {
	return &WSCreateJiraPort{client: jira.NewClientWithBaseURL(baseURL)}
}

func (p *WSCreateJiraPort) FetchIssueByTicketURL(ctx context.Context, ticketURL string) (wscreate.JiraIssue, error) {
	key, summary, err := p.client.FetchIssueByTicketURL(ctx, ticketURL)
	if err != nil {
		return wscreate.JiraIssue{}, fmt.Errorf("fetch jira issue: %w", err)
	}
	return wscreate.JiraIssue{Key: key, Summary: summary}, nil
}
