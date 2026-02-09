package appports

import (
	"context"
	"fmt"

	"github.com/tasuku43/gionx/internal/app/wscreate"
	"github.com/tasuku43/gionx/internal/infra/jira"
)

type WSCreateJiraPort struct {
	client *jira.Client
}

func NewWSCreateJiraPort() *WSCreateJiraPort {
	return &WSCreateJiraPort{client: jira.NewClient()}
}

func (p *WSCreateJiraPort) FetchIssueByTicketURL(ctx context.Context, ticketURL string) (wscreate.JiraIssue, error) {
	key, summary, err := p.client.FetchIssueByTicketURL(ctx, ticketURL)
	if err != nil {
		return wscreate.JiraIssue{}, fmt.Errorf("fetch jira issue: %w", err)
	}
	return wscreate.JiraIssue{Key: key, Summary: summary}, nil
}
