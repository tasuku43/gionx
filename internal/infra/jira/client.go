package jira

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	envJiraBaseURL  = "GIONX_JIRA_BASE_URL"
	envJiraEmail    = "GIONX_JIRA_EMAIL"
	envJiraAPIToken = "GIONX_JIRA_API_TOKEN"
)

var issueKeyRegexp = regexp.MustCompile(`(?i)\b([a-z][a-z0-9]+-\d+)\b`)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) FetchIssueByTicketURL(ctx context.Context, ticketURL string) (key string, summary string, err error) {
	cfg, err := loadEnvConfig()
	if err != nil {
		return "", "", err
	}
	issueKey, err := parseTicketURL(ticketURL)
	if err != nil {
		return "", "", err
	}

	endpoint := strings.TrimRight(cfg.baseURL.String(), "/") + "/rest/api/3/issue/" + url.PathEscape(issueKey) + "?fields=summary"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", "", fmt.Errorf("build jira request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+basicAuth(cfg.email, cfg.apiToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("jira request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		// continue
	case http.StatusUnauthorized, http.StatusForbidden:
		return "", "", fmt.Errorf("jira authentication failed: status=%d", resp.StatusCode)
	case http.StatusNotFound:
		return "", "", fmt.Errorf("jira issue not found: %s", issueKey)
	default:
		return "", "", fmt.Errorf("jira request failed: status=%d", resp.StatusCode)
	}

	var payload struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
		} `json:"fields"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", fmt.Errorf("decode jira response: %w", err)
	}
	resolvedKey := strings.TrimSpace(payload.Key)
	if resolvedKey == "" {
		resolvedKey = issueKey
	}
	return strings.ToUpper(resolvedKey), strings.TrimSpace(payload.Fields.Summary), nil
}

type envConfig struct {
	baseURL  *url.URL
	email    string
	apiToken string
}

func loadEnvConfig() (envConfig, error) {
	baseURLRaw := strings.TrimSpace(os.Getenv(envJiraBaseURL))
	email := strings.TrimSpace(os.Getenv(envJiraEmail))
	apiToken := strings.TrimSpace(os.Getenv(envJiraAPIToken))

	missing := make([]string, 0, 3)
	if baseURLRaw == "" {
		missing = append(missing, envJiraBaseURL)
	}
	if email == "" {
		missing = append(missing, envJiraEmail)
	}
	if apiToken == "" {
		missing = append(missing, envJiraAPIToken)
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return envConfig{}, fmt.Errorf("missing jira env vars: %s", strings.Join(missing, ", "))
	}

	baseURL, err := url.Parse(baseURLRaw)
	if err != nil || baseURL.Scheme == "" || baseURL.Host == "" {
		return envConfig{}, fmt.Errorf("invalid %s: %q", envJiraBaseURL, baseURLRaw)
	}

	return envConfig{baseURL: baseURL, email: email, apiToken: apiToken}, nil
}

func parseTicketURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid jira ticket URL: %q", raw)
	}
	m := issueKeyRegexp.FindStringSubmatch(strings.ToUpper(raw))
	if len(m) < 2 {
		return "", fmt.Errorf("invalid jira ticket URL: %q", raw)
	}
	return strings.ToUpper(m[1]), nil
}

func basicAuth(email string, token string) string {
	return base64.StdEncoding.EncodeToString([]byte(email + ":" + token))
}
