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

// DefaultTools mengembalikan daftar 3 tool kanonikal Charites MCP.
func DefaultTools() []ToolDefinition {
	return []ToolDefinition{
		{
			Name:        "charites_scan",
			Description: "Pindai berkas frontend (Astro, React TSX/JSX) di dalam workspace untuk audit kualitas token semantik, aksesibilitas, dan performa.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"path": {
						Type:        "string",
						Description: "Path relatif berkas atau direktori di dalam workspace yang akan dipindai (wajib).",
					},
					"category": {
						Type:        "string",
						Description: "Filter opsional kategori rule (contoh: theme, a11y, responsive, perf, tailwind).",
					},
					"rule": {
						Type:        "string",
						Description: "Filter opsional Charites Rule ID spesifik (<category>.<slug>, contoh: theme.hardcode-opacity-color).",
					},
					"ext": {
						Type:        "string",
						Description: "Filter opsional ekstensi berkas yang dipisahkan koma (default: astro,tsx,jsx).",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "charites_explain_rule",
			Description: "Ambil dokumentasi komprehensif 8-Pillars (alasan larangan, contoh buruk, contoh benar, rekomendasi perbaikan) untuk satu Charites Rule ID spesifik.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"rule_id": {
						Type:        "string",
						Description: "Charites Rule ID kanonikal (<category>.<slug>, contoh: theme.hardcode-opacity-color).",
					},
				},
				Required: []string{"rule_id"},
			},
		},
		{
			Name:        "charites_list_rules",
			Description: "Daftar seluruh rule analisis statis Charites yang terdaftar beserta kategori, tingkat keparahan, dan deskripsinya.",
			InputSchema: ToolInputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"category": {
						Type:        "string",
						Description: "Filter opsional berdasarkan nama kategori (contoh: theme).",
					},
				},
			},
		},
	}
}
