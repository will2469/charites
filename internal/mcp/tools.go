package mcp

import (
	"encoding/json"
)

// PropertySchema mendefinisikan skema atribut pada JSON schema parameter tool.
type PropertySchema struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ToolInputSchema mendefinisikan skema masukan JSON Schema dari sebuah MCP tool.
type ToolInputSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]PropertySchema `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

// ToolDefinition mendefinisikan katalog tool MCP yang diekspos ke AI agent.
type ToolDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema ToolInputSchema `json:"inputSchema"`
}

// ContentItem merepresentasikan item keluaran dari pemanggilan MCP tool.
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// CallToolParams merepresentasikan payload masukan untuk method tools/call.
type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// CallToolResult merepresentasikan respon kanonikal MCP untuk tools/call.
type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ScanArgs mendefinisikan parameter masukan untuk tool charites_scan.
type ScanArgs struct {
	Path     string `json:"path"`
	Category string `json:"category,omitempty"`
	Rule     string `json:"rule,omitempty"`
	Ext      string `json:"ext,omitempty"`
}

// ExplainArgs mendefinisikan parameter masukan untuk tool charites_explain_rule.
type ExplainArgs struct {
	RuleID string `json:"rule_id"`
}

// ListRulesArgs mendefinisikan parameter masukan untuk tool charites_list_rules.
type ListRulesArgs struct {
	Category string `json:"category,omitempty"`
}

// RuleInfo merepresentasikan entri katalog aturan dalam hasil charites_list_rules.
type RuleInfo struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

// CancelledParams mendefinisikan parameter untuk notifikasi notifications/cancelled.
type CancelledParams struct {
	RequestID any `json:"requestId"`
}

// ScanDiagnosticItem merepresentasikan temuan diagnostik yang diperkaya dengan tautan wiki kanonikal dan petunjuk explain untuk AI agent.
type ScanDiagnosticItem struct {
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Rule        string `json:"rule"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Hint        string `json:"hint,omitempty"`
	WikiURL     string `json:"wiki_url"`
	ExplainHint string `json:"explain_hint"`
}

// DefaultTools mengembalikan daftar 3 tool kanonikal Charites MCP.
func DefaultTools() []ToolDefinition {
	return []ToolDefinition{
		{
			Name: "charites_scan",
			Description: "MANDATORY FRONTEND AUDITOR: Run this tool to audit Astro (.astro) and React (.tsx/.jsx) components " +
				"for design token violations, accessibility regressions (WCAG 2.2), touch target ergonomics, and responsive layout blowouts. " +
				"Returns structured diagnostics with line numbers, error descriptions, online wiki URLs, and explain hints.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"path": {
						Type:        "string",
						Description: "Workspace-relative or absolute path to the file or directory to scan (required).",
					},
					"category": {
						Type:        "string",
						Description: "Optional category filter (e.g. theme, a11y, responsive, ergonomy, inp, cls, lcp, perf).",
					},
					"rule": {
						Type:        "string",
						Description: "Optional canonical rule identifier filter (<category>.<slug>, e.g. theme.hardcode-opacity-color).",
					},
					"ext": {
						Type:        "string",
						Description: "Optional comma-separated file extension filter (default: astro,tsx,jsx).",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name: "charites_explain_rule",
			Description: "Authoritative 8-Pillars documentation and remediation engine for any Charites static analysis rule. " +
				"Returns complete architectural overview, technical grounding, risk taxonomy, non-compliant code examples, " +
				"compliant implementations, and inline suppression directives.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"rule_id": {
						Type:        "string",
						Description: "Canonical Semgrep rule identifier (<category>.<slug>, e.g. theme.hardcode-color).",
					},
				},
				Required: []string{"rule_id"},
			},
		},
		{
			Name: "charites_list_rules",
			Description: "Lists all 90 available Charites static analysis rules across theme, a11y, responsive, ergonomy, " +
				"inp, cls, lcp, and perf categories, including severity and description.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"category": {
						Type:        "string",
						Description: "Optional filter by category name (e.g. theme, a11y, responsive).",
					},
				},
			},
		},
		{
			Name: "charites_report_issue",
			Description: "Two-Phase Human-in-the-Loop (HITL) issue reporter for filing false-positives, rule gaps, or bugs directly to GitHub. " +
				"Phase 1: Call without 'token' to generate a cryptographic draft preview and single-use approval token. " +
				"Phase 2: Present the draft preview to the user. Once explicitly approved, call this tool again with the exact same parameters and the 'token' to submit via GitHub CLI or prefilled browser URL.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"rule_id": {
						Type:        "string",
						Description: "Canonical rule identifier involved in the issue (required, e.g. theme.hardcode-color).",
					},
					"title": {
						Type:        "string",
						Description: "Short summary of the issue or false-positive (required).",
					},
					"description": {
						Type:        "string",
						Description: "Detailed description of why this rule triggered inappropriately or what the bug is (required).",
					},
					"snippet": {
						Type:        "string",
						Description: "Minimal reproducible code snippet demonstrating the false positive or issue (optional).",
					},
					"category": {
						Type:        "string",
						Description: "Issue category: false-positive, bug, rule-gap, or enhancement (default: false-positive).",
					},
					"token": {
						Type:        "string",
						Description: "Cryptographic approval token from Phase 1. Leave empty for Phase 1 draft generation; provide token for Phase 2 submission.",
					},
				},
				Required: []string{"rule_id", "title", "description"},
			},
		},
	}
}
