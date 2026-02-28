package cmux

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tasuku43/kra/internal/cmuxmap"
	"github.com/tasuku43/kra/internal/infra/cmuxctl"
)

type Client interface {
	Capabilities(ctx context.Context) (cmuxctl.Capabilities, error)
	CreateWorkspaceWithCommand(ctx context.Context, command string) (string, error)
	RenameWorkspace(ctx context.Context, workspace string, title string) error
	SelectWorkspace(ctx context.Context, workspace string) error
	ListWorkspaces(ctx context.Context) ([]cmuxctl.Workspace, error)
	Identify(ctx context.Context, workspace string, surface string) (map[string]any, error)
}

type NewClientFunc func() Client
type NewStoreFunc func(root string) cmuxmap.Store

type Service struct {
	NewClient NewClientFunc
	NewStore  NewStoreFunc
	Now       func() time.Time
}

func NewService(newClient NewClientFunc, newStore NewStoreFunc) *Service {
	return &Service{
		NewClient: newClient,
		NewStore:  newStore,
		Now:       time.Now,
	}
}

type OpenTarget struct {
	WorkspaceID   string
	WorkspacePath string
	Title         string
}

type OpenResultItem struct {
	WorkspaceID     string
	WorkspacePath   string
	CMUXWorkspaceID string
	Ordinal         int
	Title           string
}

type OpenFailure struct {
	WorkspaceID string
	Code        string
	Message     string
}

type OpenResult struct {
	Results  []OpenResultItem
	Failures []OpenFailure
}

func (s *Service) Open(ctx context.Context, root string, targets []OpenTarget, concurrency int, multi bool) (OpenResult, string, string) {
	if s.NewClient == nil || s.NewStore == nil {
		return OpenResult{}, "internal_error", "cmux service is not initialized"
	}
	client := s.NewClient()
	caps, err := client.Capabilities(ctx)
	if err != nil {
		return OpenResult{}, "cmux_capability_missing", fmt.Sprintf("read cmux capabilities: %v", err)
	}
	for _, method := range []string{"workspace.create", "workspace.rename", "workspace.select"} {
		if _, ok := caps.Methods[method]; !ok {
			return OpenResult{}, "cmux_capability_missing", fmt.Sprintf("cmux capability missing: %s", method)
		}
	}

	store := s.NewStore(root)
	mapping, err := store.Load()
	if err != nil {
		return OpenResult{}, "state_write_failed", fmt.Sprintf("load cmux mapping: %v", err)
	}

	result := OpenResult{
		Results:  make([]OpenResultItem, 0, len(targets)),
		Failures: make([]OpenFailure, 0),
	}
	if multi && concurrency > 1 {
		result = s.openConcurrent(ctx, targets, concurrency, &mapping)
	} else {
		result = s.openSequential(ctx, client, targets, &mapping)
	}
	if len(result.Results) > 0 {
		if err := store.Save(mapping); err != nil {
			return OpenResult{}, "state_write_failed", fmt.Sprintf("save cmux mapping: %v", err)
		}
	}
	return result, "", ""
}

func (s *Service) openSequential(ctx context.Context, client Client, targets []OpenTarget, mapping *cmuxmap.File) OpenResult {
	res := OpenResult{
		Results:  make([]OpenResultItem, 0, len(targets)),
		Failures: make([]OpenFailure, 0),
	}
	var mapMu sync.Mutex
	for _, target := range targets {
		item, code, msg := s.openOne(ctx, client, target, mapping, &mapMu)
		if code != "" {
			res.Failures = append(res.Failures, OpenFailure{WorkspaceID: target.WorkspaceID, Code: code, Message: msg})
			return res
		}
		res.Results = append(res.Results, item)
	}
	return res
}

func (s *Service) openConcurrent(ctx context.Context, targets []OpenTarget, concurrency int, mapping *cmuxmap.File) OpenResult {
	type task struct {
		index  int
		target OpenTarget
	}
	type outItem struct {
		index int
		item  OpenResultItem
		fail  *OpenFailure
	}
	tasks := make([]task, 0, len(targets))
	for i, target := range targets {
		tasks = append(tasks, task{index: i, target: target})
	}
	jobs := make(chan task)
	out := make(chan outItem, len(tasks))
	var wg sync.WaitGroup
	var mapMu sync.Mutex
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := s.NewClient()
			for job := range jobs {
				item, code, msg := s.openOne(ctx, client, job.target, mapping, &mapMu)
				if code != "" {
					out <- outItem{
						index: job.index,
						fail: &OpenFailure{
							WorkspaceID: job.target.WorkspaceID,
							Code:        code,
							Message:     msg,
						},
					}
					continue
				}
				out <- outItem{index: job.index, item: item}
			}
		}()
	}
	go func() {
		for _, t := range tasks {
			jobs <- t
		}
		close(jobs)
		wg.Wait()
		close(out)
	}()
	collected := make([]outItem, 0, len(tasks))
	for o := range out {
		collected = append(collected, o)
	}
	sort.Slice(collected, func(i, j int) bool { return collected[i].index < collected[j].index })
	result := OpenResult{
		Results:  make([]OpenResultItem, 0, len(tasks)),
		Failures: make([]OpenFailure, 0),
	}
	for _, c := range collected {
		if c.fail != nil {
			result.Failures = append(result.Failures, *c.fail)
			continue
		}
		result.Results = append(result.Results, c.item)
	}
	return result
}

func (s *Service) openOne(ctx context.Context, client Client, target OpenTarget, mapping *cmuxmap.File, mapMu *sync.Mutex) (OpenResultItem, string, string) {
	cmuxWorkspaceID, err := client.CreateWorkspaceWithCommand(ctx, fmt.Sprintf("cd %s", shellQuoteCDPath(target.WorkspacePath)))
	if err != nil {
		return OpenResultItem{}, "cmux_create_failed", fmt.Sprintf("create cmux workspace: %v", err)
	}
	mapMu.Lock()
	ordinal, err := cmuxmap.AllocateOrdinal(mapping, target.WorkspaceID)
	mapMu.Unlock()
	if err != nil {
		return OpenResultItem{}, "state_write_failed", fmt.Sprintf("allocate cmux ordinal: %v", err)
	}
	cmuxTitle, err := cmuxmap.FormatWorkspaceTitle(target.WorkspaceID, target.Title, ordinal)
	if err != nil {
		return OpenResultItem{}, "cmux_rename_failed", fmt.Sprintf("format cmux workspace title: %v", err)
	}
	if err := client.RenameWorkspace(ctx, cmuxWorkspaceID, cmuxTitle); err != nil {
		return OpenResultItem{}, "cmux_rename_failed", fmt.Sprintf("rename cmux workspace: %v", err)
	}
	if err := client.SelectWorkspace(ctx, cmuxWorkspaceID); err != nil {
		return OpenResultItem{}, "cmux_select_failed", fmt.Sprintf("select cmux workspace: %v", err)
	}
	now := s.Now().UTC().Format(time.RFC3339)
	mapMu.Lock()
	ws := mapping.Workspaces[target.WorkspaceID]
	ws.Entries = append(ws.Entries, cmuxmap.Entry{
		CMUXWorkspaceID: cmuxWorkspaceID,
		Ordinal:         ordinal,
		TitleSnapshot:   cmuxTitle,
		CreatedAt:       now,
		LastUsedAt:      now,
	})
	mapping.Workspaces[target.WorkspaceID] = ws
	mapMu.Unlock()
	return OpenResultItem{
		WorkspaceID:     target.WorkspaceID,
		WorkspacePath:   target.WorkspacePath,
		CMUXWorkspaceID: cmuxWorkspaceID,
		Ordinal:         ordinal,
		Title:           cmuxTitle,
	}, "", ""
}

func shellQuoteSingle(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func shellEscapeForDoubleQuotes(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "$", `\$`, "`", "\\`")
	return replacer.Replace(s)
}

func shellQuoteCDPath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil {
		if path == home {
			return `"$HOME"`
		}
		prefix := home + string(os.PathSeparator)
		if strings.HasPrefix(path, prefix) {
			suffix := strings.TrimPrefix(path, prefix)
			return `"$HOME/` + shellEscapeForDoubleQuotes(suffix) + `"`
		}
	}
	return shellQuoteSingle(path)
}

type ListRow struct {
	WorkspaceID string
	CMUXID      string
	Ordinal     int
	Title       string
	LastUsedAt  string
}

type ListResult struct {
	Rows            []ListRow
	RuntimeChecked  bool
	PrunedCount     int
	RuntimeWarnText string
}

func (s *Service) List(ctx context.Context, root string, workspaceID string) (ListResult, string, string) {
	store := s.NewStore(root)
	mapping, err := store.Load()
	if err != nil {
		return ListResult{}, "internal_error", fmt.Sprintf("load cmux mapping: %v", err)
	}
	result := ListResult{
		Rows: make([]ListRow, 0),
	}
	client := s.NewClient()
	cmuxList, lerr := client.ListWorkspaces(ctx)
	if lerr != nil {
		result.RuntimeWarnText = fmt.Sprintf("list cmux workspaces: %v", lerr)
	} else {
		result.RuntimeChecked = true
		reconciled, _, pruned, recErr := ReconcileMappingWithRuntime(store, mapping, cmuxList, true)
		if recErr != nil {
			return ListResult{}, "internal_error", fmt.Sprintf("save cmux mapping: %v", recErr)
		}
		mapping = reconciled
		result.PrunedCount = pruned
		if len(cmuxList) == 0 {
			probePruned, probeErr := probeAndPruneByID(ctx, store, &mapping, client)
			result.PrunedCount += probePruned
			if probeErr != "" {
				result.RuntimeWarnText = probeErr
			}
		}
	}

	workspaceIDs := make([]string, 0, len(mapping.Workspaces))
	for wsID := range mapping.Workspaces {
		if workspaceID != "" && workspaceID != wsID {
			continue
		}
		workspaceIDs = append(workspaceIDs, wsID)
	}
	sort.Strings(workspaceIDs)
	for _, wsID := range workspaceIDs {
		ws := mapping.Workspaces[wsID]
		for _, e := range ws.Entries {
			result.Rows = append(result.Rows, ListRow{
				WorkspaceID: wsID,
				CMUXID:      e.CMUXWorkspaceID,
				Ordinal:     e.Ordinal,
				Title:       e.TitleSnapshot,
				LastUsedAt:  e.LastUsedAt,
			})
		}
	}
	return result, "", ""
}

type StatusRow struct {
	WorkspaceID string
	CMUXID      string
	Ordinal     int
	Title       string
	Exists      bool
}

type StatusResult struct {
	Rows []StatusRow
}

func (s *Service) Status(ctx context.Context, root string, workspaceID string) (StatusResult, string, string) {
	store := s.NewStore(root)
	mapping, err := store.Load()
	if err != nil {
		return StatusResult{}, "internal_error", fmt.Sprintf("load cmux mapping: %v", err)
	}
	runtime, err := s.NewClient().ListWorkspaces(ctx)
	if err != nil {
		return StatusResult{}, "cmux_list_failed", fmt.Sprintf("list cmux workspaces: %v", err)
	}
	_, exists, _, recErr := ReconcileMappingWithRuntime(store, mapping, runtime, false)
	if recErr != nil {
		return StatusResult{}, "internal_error", fmt.Sprintf("reconcile cmux mapping: %v", recErr)
	}
	workspaceIDs := make([]string, 0, len(mapping.Workspaces))
	for wsID := range mapping.Workspaces {
		if workspaceID != "" && workspaceID != wsID {
			continue
		}
		workspaceIDs = append(workspaceIDs, wsID)
	}
	sort.Strings(workspaceIDs)
	out := StatusResult{Rows: make([]StatusRow, 0)}
	for _, wsID := range workspaceIDs {
		ws := mapping.Workspaces[wsID]
		for _, e := range ws.Entries {
			out.Rows = append(out.Rows, StatusRow{
				WorkspaceID: wsID,
				CMUXID:      e.CMUXWorkspaceID,
				Ordinal:     e.Ordinal,
				Title:       e.TitleSnapshot,
				Exists:      exists[e.CMUXWorkspaceID],
			})
		}
	}
	return out, "", ""
}

func ReconcileMappingWithRuntime(store cmuxmap.Store, mapping cmuxmap.File, runtime []cmuxctl.Workspace, prune bool) (cmuxmap.File, map[string]bool, int, error) {
	exists := map[string]bool{}
	for _, row := range runtime {
		id := strings.TrimSpace(row.ID)
		if id != "" {
			exists[id] = true
		}
	}
	if !prune || len(exists) == 0 {
		return mapping, exists, 0, nil
	}
	prunedCount := 0
	for wsID, ws := range mapping.Workspaces {
		keep := make([]cmuxmap.Entry, 0, len(ws.Entries))
		for _, e := range ws.Entries {
			if exists[strings.TrimSpace(e.CMUXWorkspaceID)] {
				keep = append(keep, e)
				continue
			}
			prunedCount++
		}
		ws.Entries = keep
		mapping.Workspaces[wsID] = ws
	}
	if prunedCount > 0 {
		if err := store.Save(mapping); err != nil {
			return mapping, exists, prunedCount, err
		}
	}
	return mapping, exists, prunedCount, nil
}

func probeAndPruneByID(ctx context.Context, store cmuxmap.Store, mapping *cmuxmap.File, client Client) (int, string) {
	statusByID := map[string]int{}
	for _, ws := range mapping.Workspaces {
		for _, e := range ws.Entries {
			id := strings.TrimSpace(e.CMUXWorkspaceID)
			if id == "" {
				continue
			}
			if _, ok := statusByID[id]; ok {
				continue
			}
			_, err := client.Identify(ctx, id, "")
			if err == nil {
				statusByID[id] = 1
				continue
			}
			if IsNotFoundError(err) {
				statusByID[id] = -1
				continue
			}
			statusByID[id] = 0
		}
	}
	probeReachable := false
	for _, st := range statusByID {
		if st != 0 {
			probeReachable = true
			break
		}
	}
	if !probeReachable {
		return 0, "cmux probe could not verify any workspace; skipped stale pruning"
	}
	prunedCount := 0
	for wsID, ws := range mapping.Workspaces {
		keep := make([]cmuxmap.Entry, 0, len(ws.Entries))
		for _, e := range ws.Entries {
			id := strings.TrimSpace(e.CMUXWorkspaceID)
			st, ok := statusByID[id]
			if !ok || st >= 0 {
				keep = append(keep, e)
				continue
			}
			prunedCount++
		}
		ws.Entries = keep
		mapping.Workspaces[wsID] = ws
	}
	if prunedCount > 0 {
		if err := store.Save(*mapping); err != nil {
			return 0, fmt.Sprintf("save cmux mapping after probe prune: %v", err)
		}
	}
	return prunedCount, ""
}

func IsNotFoundError(err error) bool {
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}
	return strings.Contains(msg, "not found") || strings.Contains(msg, "unknown workspace")
}
