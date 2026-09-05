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
	default:
		return NewErrorResponse(req.ID, CodeMethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name))
	}
}

func (s *Server) handleScan(reqID json.RawMessage, argsRaw json.RawMessage) *JSONRPCResponse {
	var args ScanArgs
	if len(argsRaw) > 0 && string(argsRaw) != "null" {
		if err := json.Unmarshal(argsRaw, &args); err != nil {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid params: %v", err))
		}
	}

	// 1. Validasi Keberadaan Parameter Path
	pathParam := strings.TrimSpace(args.Path)
	if pathParam == "" {
		return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: missing required path")
	}

	// 2. Enkapsulasi Trust Boundary & Proteksi Path Traversal
	workspaceClean, absErr := filepath.Abs(s.workspace)
	if absErr != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to resolve workspace: %v", absErr))
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
		return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: path traversal denied")
	}

	// Evaluasi Symlink jika file/folder ada
	if evalPath, symErr := filepath.EvalSymlinks(targetPath); symErr == nil {
		evalRel, evalRelErr := filepath.Rel(workspaceClean, evalPath)
		if evalRelErr != nil || evalRel == ".." || strings.HasPrefix(evalRel, ".."+string(filepath.Separator)) {
			return NewErrorResponse(reqID, CodeInvalidParams, "Invalid Params: path traversal denied")
		}
	}

	// 3. Verifikasi Keberadaan Target Path di Disk
	if _, statErr := os.Stat(targetPath); statErr != nil {
		return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid Params: scan target %q does not exist: %v", pathParam, statErr))
	}

	// 4. Verifikasi Direct-Target Safety (Kekebalan Direktori Built-in Terlarang)
	matcher := config.NewIgnoreMatcher(nil)
	if matcher.HasBuiltinAncestor(targetPath) {
		return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid Params: scan target %q is within excluded directory (builtin hard exclusion)", pathParam))
	}

	// 5. Validasi Category dan Rule jika diberikan
	if args.Category != "" {
		if len(s.reg.ByCategory(args.Category)) == 0 {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid Params: unknown category %q", args.Category))
		}
	}

	if args.Rule != "" {
		r, exists := s.reg.Get(args.Rule)
		if !exists {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid Params: unknown rule %q", args.Rule))
		}
		if args.Category != "" && r.Category() != args.Category {
			return NewErrorResponse(reqID, CodeInvalidParams, fmt.Sprintf("Invalid Params: rule %q does not belong to category %q", args.Rule, args.Category))
		}
	}

	// 6. Siapkan Batas Waktu Eksekusi 30 Detik & Kanal Pembatalan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	key := canonicalIDKey(reqID)
	if key != "" {
		s.activeScans.Store(key, cancel)
		defer s.activeScans.Delete(key)
	}

	// 7. Resolusi Konfigurasi charites.yaml
	var cfg *config.Config
	workspaceConfig := filepath.Join(workspaceClean, "charites.yaml")
	if _, cfgErr := os.Stat(workspaceConfig); cfgErr == nil {
		cfg, _ = config.Load(workspaceConfig)
	} else {
		workspaceConfigYml := filepath.Join(workspaceClean, "charites.yml")
		if _, ymlErr := os.Stat(workspaceConfigYml); ymlErr == nil {
			cfg, _ = config.Load(workspaceConfigYml)
		} else {
			cfg, _ = config.Load("")
		}
	}

	activeRules := cfg.ResolveActiveRules(s.reg, args.Category, args.Rule)

	// 8. Penyiapan Ignore Matcher
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

	// 9. Normalisasi Ekstensi Berkas
	var exts []string
	if strings.TrimSpace(args.Ext) != "" {
		rawExts := strings.Split(args.Ext, ",")
		for _, e := range rawExts {
			trimmed := strings.ToLower(strings.TrimSpace(e))
			if trimmed != "" {
				if !strings.HasPrefix(trimmed, ".") {
					trimmed = "." + trimmed
				}
				exts = append(exts, trimmed)
			}
		}
	} else {
		exts = []string{".astro", ".tsx", ".jsx"}
	}

	// 10. Eksekusi Engine Pemindaian
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

	diagsJSON, jsonErr := json.MarshalIndent(diags, "", "  ")
	if jsonErr != nil {
		return NewErrorResponse(reqID, CodeInternalToolError, fmt.Sprintf("Failed to encode diagnostics: %v", jsonErr))
	}

	result := CallToolResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: string(diagsJSON),
			},
		},
	}
	return NewResultResponse(reqID, result)
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
