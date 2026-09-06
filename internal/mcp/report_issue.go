package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ReportIssueArgs mendefinisikan parameter input untuk tool charites_report_issue.
type ReportIssueArgs struct {
	RuleID      string `json:"rule_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Snippet     string `json:"snippet,omitempty"`
	Category    string `json:"category,omitempty"` // default: "false-positive"
	Token       string `json:"token,omitempty"`
}

// IssueSubmitter mendefinisikan fungsi pengiriman issue ke penyedia eksternal (misal gh CLI).
type IssueSubmitter func(ctx context.Context, repo, title, body, labels string) (string, error)

func defaultGHSubmitter(ctx context.Context, repo, title, body, labels string) (string, error) {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return "", err
	}

	args := []string{"issue", "create", "--repo", repo, "--title", title, "--body", body}
	if labels != "" {
		args = append(args, "--label", labels)
	}

	//nolint:gosec // ghPath is checked via exec.LookPath and arguments are structured GitHub CLI parameters
	cmd := exec.CommandContext(ctx, ghPath, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// escapeCodeFence mencegah markdown injection yang menutup code block secara prematur.
func escapeCodeFence(s string) string {
	return strings.ReplaceAll(s, "```", "`\u200B``")
}

// FormatIssueTitle memformat judul issue secara seragam dan informatif.
func FormatIssueTitle(category, ruleID, title string) string {
	cat := strings.TrimSpace(category)
	if cat == "" {
		cat = "false-positive"
	}
	rID := strings.TrimSpace(ruleID)
	t := strings.TrimSpace(title)
	if rID != "" {
		return fmt.Sprintf("[%s] %s: %s", cat, rID, t)
	}
	return fmt.Sprintf("[%s] %s", cat, t)
}

// FormatIssueBody merender badan issue dalam format Markdown standar Charites.
func FormatIssueBody(args ReportIssueArgs, version string) string {
	var b strings.Builder
	cat := strings.TrimSpace(args.Category)
	if cat == "" {
		cat = "false-positive"
	}

	b.WriteString("### Metadata\n")
	b.WriteString(fmt.Sprintf("- **Rule ID**: `%s`\n", strings.TrimSpace(args.RuleID)))
	b.WriteString(fmt.Sprintf("- **Category**: `%s`\n", cat))
	b.WriteString(fmt.Sprintf("- **Charites Version**: `%s`\n", version))
	b.WriteString(fmt.Sprintf("- **Platform**: `%s/%s`\n\n", runtime.GOOS, runtime.GOARCH))

	b.WriteString("### Description\n")
	b.WriteString(strings.TrimSpace(args.Description))
	b.WriteString("\n\n")

	if strings.TrimSpace(args.Snippet) != "" {
		b.WriteString("### Reproducible Code Snippet\n")
		b.WriteString("```html\n")
		b.WriteString(escapeCodeFence(strings.TrimSpace(args.Snippet)))
		b.WriteString("\n```\n\n")
	}

	b.WriteString("---\n")
	b.WriteString("_Reported automatically via Charites MCP Two-Phase Protocol_\n")
	return b.String()
}

// BuildGitHubNewIssueURL membuat prefilled URL untuk browser jika gh CLI tidak tersedia.
func BuildGitHubNewIssueURL(repo, title, body, labels string) string {
	v := url.Values{}
	v.Set("title", title)
	v.Set("body", body)
	if labels != "" {
		v.Set("labels", labels)
	}
	return fmt.Sprintf("https://github.com/%s/issues/new?%s", repo, v.Encode())
}

func resolveIssueLabels(category string) string {
	switch strings.ToLower(strings.TrimSpace(category)) {
	case "false-positive":
		return "false-positive,mcp"
	case "bug":
		return "bug,mcp"
	case "rule-gap":
		return "rule-gap,mcp"
	case "enhancement":
		return "enhancement,mcp"
	default:
		return "mcp"
	}
}

func validateReportIssueInput(args *ReportIssueArgs) string {
	if strings.TrimSpace(args.RuleID) == "" {
		return "missing required rule_id"
	}
	if strings.TrimSpace(args.Title) == "" {
		return "missing required title"
	}
	if strings.TrimSpace(args.Description) == "" {
		return "missing required description"
	}
	return ""
}

func (s *Server) handleReportPhase1(reqID json.RawMessage, args ReportIssueArgs, hash string) *JSONRPCResponse {
	tok, err := s.getApprovalManager().CreateToken(hash, DefaultTokenTTL)
	if err != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to generate approval token: %v", err))
	}

	normCat := strings.TrimSpace(args.Category)
	if normCat == "" {
		normCat = "false-positive"
	}
	title := FormatIssueTitle(normCat, args.RuleID, args.Title)
	body := FormatIssueBody(args, s.version)

	previewPayload := map[string]any{
		"status":     "pending_approval",
		"token":      tok,
		"expires_in": "10m",
		"preview": map[string]string{
			"title":       title,
			"body":        body,
			"target_repo": "https://github.com/will2469/charites",
		},
		"instructions": fmt.Sprintf("Present this preview to the user for explicit approval. To submit this issue, call 'charites_report_issue' again with the exact same parameters and token: %q.", tok),
	}

	data, err := json.MarshalIndent(previewPayload, "", "  ")
	if err != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to encode preview payload: %v", err))
	}

	result := CallToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: string(data),
			},
		},
	}
	return NewResultResponse(reqID, result)
}

func (s *Server) handleReportPhase2(ctx context.Context, reqID json.RawMessage, args ReportIssueArgs, hash string) *JSONRPCResponse {
	if err := s.getApprovalManager().ConsumeToken(args.Token, hash); err != nil {
		return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Approval failed: %v", err))
	}

	normCat := strings.TrimSpace(args.Category)
	if normCat == "" {
		normCat = "false-positive"
	}
	title := FormatIssueTitle(normCat, args.RuleID, args.Title)
	body := FormatIssueBody(args, s.version)
	labels := resolveIssueLabels(normCat)
	repo := "will2469/charites"

	submitter := s.getIssueSubmitter()
	issueURL, submitErr := submitter(ctx, repo, title, body, labels)

	var responsePayload map[string]any
	if submitErr == nil && strings.HasPrefix(issueURL, "http") {
		responsePayload = map[string]any{
			"status": "submitted",
			"url":    issueURL,
			"method": "gh_cli",
		}
	} else {
		fallbackURL := BuildGitHubNewIssueURL(repo, title, body, labels)
		responsePayload = map[string]any{
			"status":  "ready_for_submission",
			"url":     fallbackURL,
			"method":  "browser_fallback",
			"message": "gh CLI not available or not authenticated. Please open the prefilled URL to create the issue.",
		}
	}

	data, err := json.MarshalIndent(responsePayload, "", "  ")
	if err != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to encode submission response: %v", err))
	}

	result := CallToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: string(data),
			},
		},
	}
	return NewResultResponse(reqID, result)
}

func (s *Server) handleReportIssue(reqID json.RawMessage, argsRaw json.RawMessage) *JSONRPCResponse {
	var args ReportIssueArgs
	if len(argsRaw) > 0 && string(argsRaw) != "null" {
		if err := json.Unmarshal(argsRaw, &args); err != nil {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid params: %v", err))
		}
	}

	if errStr := validateReportIssueInput(&args); errStr != "" {
		return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: "+errStr)
	}

	// Evaluasi status telemetri / izin pelaporan
	workspaceClean, _ := filepath.Abs(s.workspace)
	cfg, _ := resolveWorkspaceConfigAndMatcher(workspaceClean)
	if !cfg.IsTelemetryEnabled() {
		payload := map[string]any{
			"status":  "disabled",
			"message": "Issue reporting is disabled by policy (CHARITES_TELEMETRY=false or telemetry: false in charites.yaml).",
		}
		data, _ := json.MarshalIndent(payload, "", "  ")
		return NewResultResponse(reqID, CallToolResult{
			Content: []ContentItem{
				{
					Type: "text",
					Text: string(data),
				},
			},
		})
	}

	normCat := strings.TrimSpace(args.Category)
	if normCat == "" {
		normCat = "false-positive"
	}
	hash := ComputeIssueHash(
		strings.TrimSpace(args.RuleID),
		strings.TrimSpace(args.Title),
		strings.TrimSpace(args.Description),
		strings.TrimSpace(args.Snippet),
		normCat,
	)

	if strings.TrimSpace(args.Token) == "" {
		return s.handleReportPhase1(reqID, args, hash)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return s.handleReportPhase2(ctx, reqID, args, hash)
}
