package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
	"github.com/will2469/charites/internal/scanner"
	"github.com/will2469/charites/internal/wiki"
)

// Dispatch memproses satu request JSON-RPC dan mengembalikan response (atau nil untuk notifikasi).
func (s *Server) Dispatch(req *JSONRPCRequest) *JSONRPCResponse {
	if req == nil {
		return nil
	}

	// 1. Validasi Versi Protokol JSON-RPC
	if req.JSONRPC != "2.0" {
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeInvalidRequest, "Invalid Request: jsonrpc must be '2.0'")
	}

	// 2. Validasi Method Tidak Boleh Kosong
	if strings.TrimSpace(req.Method) == "" {
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeInvalidRequest, "Invalid Request: missing method")
	}

	currState := s.GetState()

	// 3. Evaluasi Berdasarkan State Machine Siklus Hidup
	switch currState {
	case StateNew:
		return s.dispatchNewState(req)
	case StateInitializing:
		return s.dispatchInitializingState(req)
	case StateReady:
		return s.dispatchReadyState(req)
	case StateTerminated:
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeInvalidRequest, "Invalid Request: server is terminated")
	default:
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeInvalidRequest, "Invalid Request: unknown server state")
	}
}

func (s *Server) dispatchNewState(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		s.SetState(StateInitializing)
		result := map[string]any{
			"protocolVersion": ProtocolVersion,
			"capabilities": map[string]any{
				"tools": map[string]any{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]any{
				"name":    ServerName,
				"version": s.version,
			},
		}
		return NewResultResponse(req.ID, result)

	case "ping":
		return NewResultResponse(req.ID, map[string]any{})

	default:
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeServerNotInitialized, "Server not initialized")
	}
}

func (s *Server) dispatchInitializingState(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "notifications/initialized":
		s.SetState(StateReady)
		return nil

	case "ping":
		return NewResultResponse(req.ID, map[string]any{})

	case "initialize":
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeInvalidRequest, "Invalid Request: already initializing")

	default:
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeServerNotInitialized, "Server not initialized")
	}
}

func (s *Server) dispatchReadyState(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		if req.IsNotification() {
			return nil
		}
		return NewErrorResponse(req.ID, CodeInvalidRequest, "Invalid Request: already initialized")

	case "notifications/initialized":
		// Redundant notification diabaikan secara senyap
		return nil

	case "notifications/cancelled":
		s.handleCancelled(req.Params)
		return nil

	case "ping":
		return NewResultResponse(req.ID, map[string]any{})

	case "tools/list":
		if req.IsNotification() {
			return nil
		}
		return NewResultResponse(req.ID, map[string]any{
			"tools": DefaultTools(),
		})

	case "tools/call":
		if req.IsNotification() {
			return nil
		}
		return s.handleToolCall(req)

	case "server/discover":
		if req.IsNotification() {
			return nil
		}
		return NewResultResponse(req.ID, map[string]any{
			"protocolVersion": ProtocolVersion,
			"tools":           DefaultTools(),
		})

	default:
		if req.IsNotification() {
			// Notifikasi tidak dikenal diabaikan secara senyap
			return nil
		}
		return NewErrorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) handleCancelled(paramsRaw json.RawMessage) {
	if len(paramsRaw) == 0 {
		return
	}
	var p struct {
		RequestID any `json:"requestId"`
	}
	if err := json.Unmarshal(paramsRaw, &p); err != nil {
		return
	}

	key := canonicalIDKey(p.RequestID)
	if key == "" {
		return
	}

	if val, ok := s.activeScans.Load(key); ok {
		if cancel, isCancel := val.(context.CancelFunc); isCancel {
			cancel()
		}
		s.activeScans.Delete(key)
	}
}

func (s *Server) handleToolCall(req *JSONRPCRequest) *JSONRPCResponse {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, CodeInvalidParams, fmt.Sprintf("Invalid params: %v", err))
	}

	switch params.Name {
	case "charites_scan":
		return s.handleScan(req.ID, params.Arguments)
	case "charites_explain_rule":
		return s.handleExplain(req.ID, params.Arguments)
	case "charites_list_rules":
		return s.handleListRules(req.ID, params.Arguments)
	case "charites_report_issue":
		return s.handleReportIssue(req.ID, params.Arguments)
	default:
		return NewErrorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name))
	}
}

func (s *Server) resolveScanTarget(pathParam string) (string, string) {
	pathParam = strings.TrimSpace(pathParam)
	if pathParam == "" {
		return "", "missing required path"
	}

	workspaceClean, absErr := filepath.Abs(s.workspace)
	if absErr != nil {
		return "", fmt.Sprintf("failed to resolve workspace: %v", absErr)
	}
	workspaceClean = filepath.Clean(workspaceClean)

	var targetPath string
	if filepath.IsAbs(pathParam) {
		targetPath = filepath.Clean(pathParam)
	} else {
		targetPath = filepath.Clean(filepath.Join(workspaceClean, pathParam))
	}

	rel, relErr := filepath.Rel(workspaceClean, targetPath)
	if relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "path traversal denied"
	}

	if evalPath, symErr := filepath.EvalSymlinks(targetPath); symErr == nil {
		evalRel, evalRelErr := filepath.Rel(workspaceClean, evalPath)
		if evalRelErr != nil || evalRel == ".." || strings.HasPrefix(evalRel, ".."+string(filepath.Separator)) {
			return "", "path traversal denied"
		}
	}

	if _, statErr := os.Stat(targetPath); statErr != nil {
		return "", fmt.Sprintf("scan target %q does not exist: %v", pathParam, statErr)
	}

	matcher := config.NewIgnoreMatcher(nil)
	if matcher.HasBuiltinAncestor(targetPath) {
		return "", fmt.Sprintf("scan target %q is within excluded directory (builtin hard exclusion)", pathParam)
	}

	return targetPath, ""
}

func (s *Server) validateScanCategoryAndRule(category, rule string) string {
	if category != "" && len(s.reg.ByCategory(category)) == 0 {
		return fmt.Sprintf("unknown category %q", category)
	}
	if rule != "" {
		r, exists := s.reg.Get(rule)
		if !exists {
			return fmt.Sprintf("unknown rule %q", rule)
		}
		if category != "" && r.Category() != category {
			return fmt.Sprintf("rule %q does not belong to category %q", rule, category)
		}
	}
	return ""
}

func resolveWorkspaceConfigAndMatcher(workspaceClean string) (*config.Config, *config.IgnoreMatcher) {
	var cfg *config.Config
	for _, cand := range []string{"charites.yaml", "charites.yml"} {
		p := filepath.Join(workspaceClean, cand)
		if _, err := os.Stat(p); err == nil {
			cfg, _ = config.Load(p)
			break
		}
	}
	if cfg == nil {
		cfg, _ = config.Load("")
	}

	var matcher *config.IgnoreMatcher
	targetIgnore := filepath.Join(workspaceClean, ".charitesignore")
	if _, ignErr := os.Stat(targetIgnore); ignErr == nil {
		matcher, _ = config.LoadIgnore(targetIgnore)
	}
	if matcher == nil {
		matcher = config.NewIgnoreMatcher(nil)
	}
	if cfg != nil && len(cfg.Ignore) > 0 {
		matcher.AddPatterns(cfg.Ignore)
	}
	return cfg, matcher
}

func parseScanExts(extParam string) []string {
	if strings.TrimSpace(extParam) == "" {
		return []string{".astro", ".tsx", ".jsx"}
	}
	rawExts := strings.Split(extParam, ",")
	var exts []string
	for _, e := range rawExts {
		trimmed := strings.ToLower(strings.TrimSpace(e))
		if trimmed != "" {
			if !strings.HasPrefix(trimmed, ".") {
				trimmed = "." + trimmed
			}
			exts = append(exts, trimmed)
		}
	}
	return exts
}

func (s *Server) handleScan(reqID json.RawMessage, argsRaw json.RawMessage) *JSONRPCResponse {
	var args ScanArgs
	if len(argsRaw) > 0 && string(argsRaw) != "null" {
		if err := json.Unmarshal(argsRaw, &args); err != nil {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid params: %v", err))
		}
	}

	targetPath, errStr := s.resolveScanTarget(args.Path)
	if errStr != "" {
		return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: "+errStr)
	}

	if ruleErr := s.validateScanCategoryAndRule(args.Category, args.Rule); ruleErr != "" {
		return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: "+ruleErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	key := canonicalIDKey(reqID)
	if key != "" {
		s.activeScans.Store(key, cancel)
		defer s.activeScans.Delete(key)
	}

	workspaceClean, _ := filepath.Abs(s.workspace)
	cfg, matcher := resolveWorkspaceConfigAndMatcher(workspaceClean)
	activeRules := cfg.ResolveActiveRules(s.reg, args.Category, args.Rule)
	exts := parseScanExts(args.Ext)

	walker := scanner.NewWalker(matcher, exts)
	eng := analyzer.NewEngine(activeRules)
	pool := scanner.NewPool(0)

	diags, err := pool.Run(ctx, walker, targetPath, eng)
	if ctx.Err() != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, "Scan execution cancelled or timed out")
	}
	if err != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Scan execution failed: %v", err))
	}

	outputText, fmtErr := formatScanDiagnostics(diags)
	if fmtErr != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to encode diagnostics: %v", fmtErr))
	}

	result := CallToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: outputText,
			},
		},
	}
	return NewResultResponse(reqID, result)
}

func formatScanDiagnostics(diags []ir.Diagnostic) (string, error) {
	if len(diags) == 0 {
		return " CLEAN - No violations found across scanned components.", nil
	}

	items := make([]ScanDiagnosticItem, 0, len(diags))
	for _, d := range diags {
		items = append(items, ScanDiagnosticItem{
			File:        filepath.ToSlash(d.File),
			Line:        d.Line,
			Column:      d.Column,
			Rule:        d.Rule,
			Severity:    string(d.Severity),
			Message:     d.Message,
			Hint:        d.Hint,
			WikiURL:     "https://github.com/will2469/charites/wiki/" + d.Rule,
			ExplainHint: fmt.Sprintf("Call MCP tool 'charites_explain_rule' with rule_id: %q for full 8-Pillars documentation, good/bad examples, and remediation guidance.", d.Rule),
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleExplain(reqID json.RawMessage, argsRaw json.RawMessage) *JSONRPCResponse {
	var args ExplainArgs
	if len(argsRaw) > 0 && string(argsRaw) != "null" {
		if err := json.Unmarshal(argsRaw, &args); err != nil {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid params: %v", err))
		}
	}

	ruleID := strings.TrimSpace(args.RuleID)
	if ruleID == "" {
		return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: missing required rule_id")
	}

	rule, ok := s.reg.Get(ruleID)
	if !ok {
		return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid Params: rule %q not found", ruleID))
	}

	docText, err := wiki.RenderRuleDoc(rule)
	if err != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to render rule documentation: %v", err))
	}

	result := CallToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: docText,
			},
		},
	}
	return NewResultResponse(reqID, result)
}

func (s *Server) handleListRules(reqID json.RawMessage, argsRaw json.RawMessage) *JSONRPCResponse {
	var args ListRulesArgs
	if len(argsRaw) > 0 && string(argsRaw) != "null" {
		_ = json.Unmarshal(argsRaw, &args)
	}

	var rulesList []rules.Rule
	if strings.TrimSpace(args.Category) != "" {
		rulesList = s.reg.ByCategory(strings.TrimSpace(args.Category))
	} else {
		rulesList = s.reg.All()
	}

	out := make([]RuleInfo, 0, len(rulesList))
	for _, r := range rulesList {
		out = append(out, RuleInfo{
			ID:          r.ID(),
			Category:    r.Category(),
			Severity:    string(r.DefaultSeverity()),
			Description: r.Description(),
		})
	}

	jsonBytes, jsonErr := json.MarshalIndent(out, "", "  ")
	if jsonErr != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to encode rules: %v", jsonErr))
	}

	result := CallToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: string(jsonBytes),
			},
		},
	}
	return NewResultResponse(reqID, result)
}

func canonicalIDKey(id any) string {
	switch v := id.(type) {
	case json.RawMessage:
		return strings.TrimSpace(string(v))
	case string:
		return strings.TrimSpace(v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		if id == nil {
			return ""
		}
		return fmt.Sprintf("%v", id)
	}
}
