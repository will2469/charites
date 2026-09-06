package mcp_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/will2469/charites/internal/mcp"
	"github.com/will2469/charites/internal/rules"
)

func TestComputeIssueHash(t *testing.T) {
	// 1. Hash deterministik
	h1 := mcp.ComputeIssueHash("theme.color", "Title A", "Desc A", "snippet A", "false-positive")
	h2 := mcp.ComputeIssueHash("theme.color", "Title A", "Desc A", "snippet A", "false-positive")
	if h1 != h2 {
		t.Fatalf("expected identical hash for identical inputs, got %q vs %q", h1, h2)
	}

	// 2. Hash berbeda jika ada field berbeda
	hDiff := mcp.ComputeIssueHash("theme.color", "Title B", "Desc A", "snippet A", "false-positive")
	if h1 == hDiff {
		t.Fatalf("expected different hashes for different titles")
	}

	// 3. Null-byte separator delimiter collision defense
	// Misal: ruleID="a", title="bc" vs ruleID="ab", title="c"
	hCol1 := mcp.ComputeIssueHash("a", "bc", "d", "e", "f")
	hCol2 := mcp.ComputeIssueHash("ab", "c", "d", "e", "f")
	if hCol1 == hCol2 {
		t.Fatalf("delimiter collision detected: hash should differ with null byte separator")
	}
}

func TestApprovalManager_Lifecycle(t *testing.T) {
	am := mcp.NewApprovalManager()
	hash := mcp.ComputeIssueHash("theme.color", "Title", "Desc", "", "false-positive")

	// 1. Create Token
	tok, err := am.CreateToken(hash, 10*time.Minute)
	if err != nil {
		t.Fatalf("failed creating token: %v", err)
	}
	if !strings.HasPrefix(tok, "appr_") {
		t.Fatalf("expected token prefix 'appr_', got %q", tok)
	}

	// 2. Consume Token with wrong hash (Tampering defense)
	tamperedHash := mcp.ComputeIssueHash("theme.color", "Tampered", "Desc", "", "false-positive")
	errTamper := am.ConsumeToken(tok, tamperedHash)
	if !errors.Is(errTamper, mcp.ErrHashMismatch) {
		t.Fatalf("expected ErrHashMismatch, got %v", errTamper)
	}

	// 3. Consume Token with valid hash
	errValid := am.ConsumeToken(tok, hash)
	if errValid != nil {
		t.Fatalf("expected successful consume, got %v", errValid)
	}

	// 4. Replay attack: second consume must fail
	errReplay := am.ConsumeToken(tok, hash)
	if !errors.Is(errReplay, mcp.ErrTokenNotFound) {
		t.Fatalf("expected ErrTokenNotFound on replayed token, got %v", errReplay)
	}

	// 5. Expiration handling
	tokExp, err := am.CreateToken(hash, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("failed creating expiring token: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	errExpired := am.ConsumeToken(tokExp, hash)
	if !errors.Is(errExpired, mcp.ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", errExpired)
	}
}

func setupDirectTestServer(t *testing.T, workspace string) *mcp.Server {
	reg := rules.DefaultRegistry()
	srv := mcp.NewServer(workspace, strings.NewReader(""), nil, nil, reg, "1.0.0")
	srv.SetState(mcp.StateReady)
	return srv
}

func TestReportIssue_Phase1_DraftPreview(t *testing.T) {
	tmpDir := t.TempDir()
	srv := setupDirectTestServer(t, tmpDir)

	snippet := "<button class=\"bg-[#ff0000]```injection\">Click</button>"
	argsMap := map[string]any{
		"rule_id":     "theme.hardcode-color",
		"title":       "False positive on custom variable",
		"description": "Rule flagged valid CSS color in third party component",
		"snippet":     snippet,
		"category":    "false-positive",
	}
	argsBytes, _ := json.Marshal(argsMap)

	req := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"req-1"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytes))),
	}

	resp := srv.Dispatch(req)
	if resp.Error != nil {
		t.Fatalf("unexpected dispatch error: %+v", resp.Error)
	}

	result, ok := resp.Result.(mcp.CallToolResult)
	if !ok || len(result.Content) == 0 {
		t.Fatalf("expected CallToolResult with content")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(result.Content[0].Text), &payload); err != nil {
		t.Fatalf("failed parsing preview payload JSON: %v", err)
	}

	if payload["status"] != "pending_approval" {
		t.Errorf("expected status 'pending_approval', got %v", payload["status"])
	}

	token, _ := payload["token"].(string)
	if !strings.HasPrefix(token, "appr_") {
		t.Errorf("expected token prefix 'appr_', got %q", token)
	}

	preview, _ := payload["preview"].(map[string]any)
	title, _ := preview["title"].(string)
	body, _ := preview["body"].(string)

	if !strings.Contains(title, "[false-positive] theme.hardcode-color: False positive on custom variable") {
		t.Errorf("unexpected title: %q", title)
	}

	// Pastikan escaping code fence aktif (triple backtick diubah)
	if strings.Contains(body, "```injection") {
		t.Errorf("body contains unescaped triple backticks inside snippet")
	}
	if !strings.Contains(body, "`\u200B``injection") {
		t.Errorf("body missing zero-width escaped code fence")
	}
}

func TestReportIssue_Phase2_Submission_GHCLI(t *testing.T) {
	tmpDir := t.TempDir()
	srv := setupDirectTestServer(t, tmpDir)

	mockCalled := false
	expectedIssueURL := "https://github.com/will2469/charites/issues/101"
	srv.SetIssueSubmitter(func(ctx context.Context, repo, title, body, labels string) (string, error) {
		mockCalled = true
		if repo != "will2469/charites" {
			t.Errorf("expected repo will2469/charites, got %s", repo)
		}
		if !strings.Contains(labels, "false-positive") {
			t.Errorf("expected label false-positive, got %s", labels)
		}
		return expectedIssueURL, nil
	})

	// Phase 1: Buat draft preview
	argsMap := map[string]any{
		"rule_id":     "theme.hardcode-color",
		"title":       "Bug in AST walker",
		"description": "Explanation",
		"category":    "false-positive",
	}
	argsBytes, _ := json.Marshal(argsMap)
	reqP1 := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"req-p1"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytes))),
	}
	respP1 := srv.Dispatch(reqP1)
	resP1 := respP1.Result.(mcp.CallToolResult)
	var p1Payload map[string]any
	_ = json.Unmarshal([]byte(resP1.Content[0].Text), &p1Payload)
	tok := p1Payload["token"].(string)

	// Phase 2: Kirim dengan token
	argsMap["token"] = tok
	argsBytesP2, _ := json.Marshal(argsMap)
	reqP2 := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"req-p2"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytesP2))),
	}

	respP2 := srv.Dispatch(reqP2)
	if respP2.Error != nil {
		t.Fatalf("unexpected error in Phase 2: %+v", respP2.Error)
	}

	if !mockCalled {
		t.Errorf("expected mock issue submitter to be called")
	}

	resP2 := respP2.Result.(mcp.CallToolResult)
	var p2Payload map[string]any
	_ = json.Unmarshal([]byte(resP2.Content[0].Text), &p2Payload)

	if p2Payload["status"] != "submitted" {
		t.Errorf("expected status 'submitted', got %v", p2Payload["status"])
	}
	if p2Payload["url"] != expectedIssueURL {
		t.Errorf("expected issue url %q, got %v", expectedIssueURL, p2Payload["url"])
	}
	if p2Payload["method"] != "gh_cli" {
		t.Errorf("expected method 'gh_cli', got %v", p2Payload["method"])
	}
}

func TestReportIssue_Phase2_Submission_BrowserFallback(t *testing.T) {
	tmpDir := t.TempDir()
	srv := setupDirectTestServer(t, tmpDir)

	// Submitter mengembalikan error (simulasi gh CLI tidak terinstall / gagal)
	srv.SetIssueSubmitter(func(ctx context.Context, repo, title, body, labels string) (string, error) {
		return "", errors.New("exec: \"gh\": executable file not found in $PATH")
	})

	// Phase 1: Buat token
	argsMap := map[string]any{
		"rule_id":     "a11y.button-name",
		"title":       "Missing accessible name in SVG",
		"description": "SVG icon button was reported incorrectly",
		"category":    "bug",
	}
	argsBytes, _ := json.Marshal(argsMap)
	reqP1 := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"p1"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytes))),
	}
	respP1 := srv.Dispatch(reqP1)
	resP1 := respP1.Result.(mcp.CallToolResult)
	var p1Payload map[string]any
	_ = json.Unmarshal([]byte(resP1.Content[0].Text), &p1Payload)
	tok := p1Payload["token"].(string)

	// Phase 2: Kirim
	argsMap["token"] = tok
	argsBytesP2, _ := json.Marshal(argsMap)
	reqP2 := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"p2"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytesP2))),
	}

	respP2 := srv.Dispatch(reqP2)
	if respP2.Error != nil {
		t.Fatalf("unexpected error in Phase 2: %+v", respP2.Error)
	}

	resP2 := respP2.Result.(mcp.CallToolResult)
	var p2Payload map[string]any
	_ = json.Unmarshal([]byte(resP2.Content[0].Text), &p2Payload)

	if p2Payload["status"] != "ready_for_submission" {
		t.Errorf("expected status 'ready_for_submission', got %v", p2Payload["status"])
	}
	if p2Payload["method"] != "browser_fallback" {
		t.Errorf("expected method 'browser_fallback', got %v", p2Payload["method"])
	}
	rawURL, _ := p2Payload["url"].(string)
	if !strings.HasPrefix(rawURL, "https://github.com/will2469/charites/issues/new?") {
		t.Errorf("unexpected prefilled URL format: %s", rawURL)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed parsing generated url: %v", err)
	}
	q := parsed.Query()
	if !strings.Contains(q.Get("title"), "Missing accessible name") {
		t.Errorf("query missing title: %v", q.Get("title"))
	}
	if !strings.Contains(q.Get("labels"), "bug,mcp") {
		t.Errorf("query missing labels: %v", q.Get("labels"))
	}
}

func TestReportIssue_TamperingRejection(t *testing.T) {
	tmpDir := t.TempDir()
	srv := setupDirectTestServer(t, tmpDir)

	// Phase 1
	argsMap := map[string]any{
		"rule_id":     "theme.hardcode-color",
		"title":       "Original Title",
		"description": "Original Description",
	}
	argsBytes, _ := json.Marshal(argsMap)
	reqP1 := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"t-p1"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytes))),
	}
	respP1 := srv.Dispatch(reqP1)
	resP1 := respP1.Result.(mcp.CallToolResult)
	var p1Payload map[string]any
	_ = json.Unmarshal([]byte(resP1.Content[0].Text), &p1Payload)
	tok := p1Payload["token"].(string)

	// Phase 2: Ubah title (Tampering)
	argsMap["title"] = "Malicious Altered Title"
	argsMap["token"] = tok
	argsBytesTampered, _ := json.Marshal(argsMap)
	reqP2 := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"t-p2"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytesTampered))),
	}

	respP2 := srv.Dispatch(reqP2)
	if respP2.Error == nil {
		t.Fatalf("expected error on tampered issue submission, got nil")
	}
	if !strings.Contains(respP2.Error.Message, "Approval failed") {
		t.Errorf("expected 'Approval failed' error message, got %q", respP2.Error.Message)
	}
}

func TestReportIssue_TelemetryDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	// Tulis charites.yaml dengan telemetry: false
	cfgPath := filepath.Join(tmpDir, "charites.yaml")
	_ = os.WriteFile(cfgPath, []byte("telemetry: false\n"), 0600)

	srv := setupDirectTestServer(t, tmpDir)

	argsMap := map[string]any{
		"rule_id":     "theme.hardcode-color",
		"title":       "Should be refused",
		"description": "Telemetry disabled",
	}
	argsBytes, _ := json.Marshal(argsMap)
	req := &mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"tel-1"`),
		Method:  "tools/call",
		Params: json.RawMessage(fmt.Sprintf(`{
			"name": "charites_report_issue",
			"arguments": %s
		}`, string(argsBytes))),
	}

	resp := srv.Dispatch(req)
	if resp.Error != nil {
		t.Fatalf("unexpected rpc error: %+v", resp.Error)
	}

	res := resp.Result.(mcp.CallToolResult)
	var payload map[string]any
	_ = json.Unmarshal([]byte(res.Content[0].Text), &payload)

	if payload["status"] != "disabled" {
		t.Errorf("expected status 'disabled', got %v", payload["status"])
	}
}

func TestReportIssue_ValidationErrors(t *testing.T) {
	srv := setupDirectTestServer(t, t.TempDir())

	tests := []struct {
		name    string
		args    map[string]any
		wantErr string
	}{
		{
			name:    "missing rule_id",
			args:    map[string]any{"title": "T", "description": "D"},
			wantErr: "missing required rule_id",
		},
		{
			name:    "missing title",
			args:    map[string]any{"rule_id": "theme.color", "description": "D"},
			wantErr: "missing required title",
		},
		{
			name:    "missing description",
			args:    map[string]any{"rule_id": "theme.color", "title": "T"},
			wantErr: "missing required description",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.args)
			req := &mcp.JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      json.RawMessage(`"val-err"`),
				Method:  "tools/call",
				Params: json.RawMessage(fmt.Sprintf(`{
					"name": "charites_report_issue",
					"arguments": %s
				}`, string(b))),
			}
			resp := srv.Dispatch(req)
			if resp.Error == nil {
				t.Fatalf("expected error, got result: %+v", resp.Result)
			}
			if !strings.Contains(resp.Error.Message, tc.wantErr) {
				t.Errorf("expected error containing %q, got %q", tc.wantErr, resp.Error.Message)
			}
		})
	}
}
