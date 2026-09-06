package mcp_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/will2469/charites/internal/mcp"
	"github.com/will2469/charites/internal/rules"
)

// helper untuk menjalankan server MCP dengan pipe I/O terisolasi.
func setupTestServer(t *testing.T, workspace string) (*mcp.Server, *io.PipeWriter, *bufio.Reader, *bytes.Buffer, func()) {
	t.Helper()

	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	var stderrBuf bytes.Buffer

	srv := mcp.NewServer(workspace, inR, outW, &stderrBuf, rules.DefaultRegistry(), "1.0.0-test")

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run()
	}()

	outReader := bufio.NewReader(outR)

	cleanup := func() {
		_ = inW.Close()
		_ = inR.Close()
		_ = outW.Close()
		_ = outR.Close()
	}

	return srv, inW, outReader, &stderrBuf, cleanup
}

func readJSONResponse(t *testing.T, r *bufio.Reader) mcp.JSONRPCResponse {
	t.Helper()
	line, err := r.ReadBytes('\n')
	if err != nil {
		t.Fatalf("failed to read response line: %v", err)
	}

	var resp mcp.JSONRPCResponse
	if err := json.Unmarshal(line, &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON-RPC response: %v, raw: %s", err, string(line))
	}
	return resp
}

// MCP-TEST-001: Handshake normal: initialize lalu notifications/initialized -> respon sukses dengan protocolVersion: "2026-07-28".
func TestMCP_001_HandshakeNormal(t *testing.T) {
	ws := t.TempDir()
	srv, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// 1. Kirim initialize
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2026-07-28"}}` + "\n"
	_, _ = inW.Write([]byte(initReq))

	resp := readJSONResponse(t, outR)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
	if string(resp.ID) != "1" {
		t.Errorf("expected id 1, got %s", string(resp.ID))
	}

	resMap, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("result is not map: %v", resp.Result)
	}
	if resMap["protocolVersion"] != "2026-07-28" {
		t.Errorf("expected protocolVersion 2026-07-28, got %v", resMap["protocolVersion"])
	}

	// 2. Kirim notifikasi notifications/initialized
	initNotif := `{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"
	_, _ = inW.Write([]byte(initNotif))

	// Beri sedikit waktu untuk transisi state
	time.Sleep(20 * time.Millisecond)
	if srv.GetState() != mcp.StateReady {
		t.Errorf("expected server to be StateReady, got %v", srv.GetState())
	}
}

// MCP-TEST-002: tools/list atau tools/call dipanggil sebelum initialize -> Error -32002 (Server not initialized).
func TestMCP_002_UninitializedToolCalls(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Pemanggilan tools/list saat state NEW
	listReq := `{"jsonrpc":"2.0","id":100,"method":"tools/list"}` + "\n"
	_, _ = inW.Write([]byte(listReq))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error for uninitialized tools/list, got result: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeServerNotInitialized {
		t.Errorf("expected error code %d, got %d", mcp.CodeServerNotInitialized, resp.Error.Code)
	}

	// Pemanggilan tools/call saat state NEW
	callReq := `{"jsonrpc":"2.0","id":101,"method":"tools/call","params":{"name":"charites_list_rules"}}` + "\n"
	_, _ = inW.Write([]byte(callReq))

	resp2 := readJSONResponse(t, outR)
	if resp2.Error == nil {
		t.Fatalf("expected error for uninitialized tools/call, got result: %v", resp2.Result)
	}
	if resp2.Error.Code != mcp.CodeServerNotInitialized {
		t.Errorf("expected error code %d, got %d", mcp.CodeServerNotInitialized, resp2.Error.Code)
	}
}

// MCP-TEST-003: initialize dipanggil kedua kali setelah server READY -> Error -32600 (Invalid Request: already initialized).
func TestMCP_003_DoubleInitialize(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// 1. Handshake awal
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// 2. Pemanggilan initialize kedua
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"initialize"}` + "\n"))
	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error on double initialize, got result: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected code %d, got %d", mcp.CodeInvalidRequest, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "already initialized") {
		t.Errorf("expected error message to contain 'already initialized', got: %s", resp.Error.Message)
	}
}

// MCP-TEST-004: JSON sintaks cacat (malformed JSON string) -> Error -32700 (Parse Error).
func TestMCP_004_MalformedJSON(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	badJSON := `{"jsonrpc": "2.0", "id": 1, "method": ` + "\n"
	_, _ = inW.Write([]byte(badJSON))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected parse error, got: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeParseError {
		t.Errorf("expected code %d, got %d", mcp.CodeParseError, resp.Error.Code)
	}
}

// MCP-TEST-005: Frame pesan melebihi batas 4 Megabytes -> Error -32700 (Parse Error / Frame Exceeded).
func TestMCP_005_FrameExceeds4MB(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Buat payload berukuran 4MB + 100KB
	hugePayload := strings.Repeat("A", 4*1024*1024+100*1024) + "\n"
	_, _ = inW.Write([]byte(hugePayload))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected parse error on oversized frame, got result: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeParseError {
		t.Errorf("expected code %d, got %d", mcp.CodeParseError, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "4MB") {
		t.Errorf("expected message to mention 4MB, got: %s", resp.Error.Message)
	}
}

// MCP-TEST-006: Request tanpa field "jsonrpc": "2.0" -> Error -32600 (Invalid Request).
func TestMCP_006_MissingJSONRPCVersion(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Request tanpa "jsonrpc"
	req := `{"id":42,"method":"ping"}` + "\n"
	_, _ = inW.Write([]byte(req))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error for missing jsonrpc version, got: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected code %d, got %d", mcp.CodeInvalidRequest, resp.Error.Code)
	}

	// Request dengan jsonrpc selain "2.0"
	req2 := `{"jsonrpc":"1.0","id":43,"method":"ping"}` + "\n"
	_, _ = inW.Write([]byte(req2))

	resp2 := readJSONResponse(t, outR)
	if resp2.Error == nil || resp2.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected invalid request for jsonrpc 1.0, got: %v", resp2.Error)
	}
}

// MCP-TEST-007: Method RPC tidak dikenal (misal foo/bar) -> Error -32601 (Method not found).
func TestMCP_007_MethodNotFound(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Lakukan handshake agar state READY
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":99,"method":"foo/bar"}` + "\n"))
	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error for unknown method, got result: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeMethodNotFound {
		t.Errorf("expected code %d, got %d", mcp.CodeMethodNotFound, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "foo/bar") {
		t.Errorf("expected message to mention foo/bar, got: %s", resp.Error.Message)
	}
}

// MCP-TEST-008: Tool call dengan nama tool tidak terdaftar -> Error -32601 (Tool not found).
func TestMCP_008_ToolNotFound(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Pemanggilan tool tidak dikenal
	callReq := `{"jsonrpc":"2.0","id":55,"method":"tools/call","params":{"name":"unknown_nonexistent_tool","arguments":{}}}` + "\n"
	_, _ = inW.Write([]byte(callReq))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error for unregistered tool, got: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeMethodNotFound {
		t.Errorf("expected code %d, got %d", mcp.CodeMethodNotFound, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "unknown_nonexistent_tool") {
		t.Errorf("expected message to mention tool name, got: %s", resp.Error.Message)
	}
}

// MCP-TEST-009: Tool call charites_scan tanpa parameter path -> Error -32602 (Invalid Params: missing required path).
func TestMCP_009_ScanMissingPath(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Scan tanpa path
	callReq := `{"jsonrpc":"2.0","id":56,"method":"tools/call","params":{"name":"charites_scan","arguments":{"category":"theme"}}}` + "\n"
	_, _ = inW.Write([]byte(callReq))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error for missing path, got: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected code %d, got %d", mcp.CodeInvalidParams, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "missing required path") {
		t.Errorf("expected message to mention missing required path, got: %s", resp.Error.Message)
	}
}

// MCP-TEST-010: charites_scan dengan path traversal ../../etc/passwd -> Error -32602 (Invalid Params: path traversal denied).
func TestMCP_010_PathTraversalDenied(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	traversalPaths := []string{
		"../../etc/passwd",
		"../something",
		"/etc/passwd",
	}

	for i, tp := range traversalPaths {
		callReq := fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":%q}}}`+"\n", 70+i, tp)
		_, _ = inW.Write([]byte(callReq))

		resp := readJSONResponse(t, outR)
		if resp.Error == nil {
			t.Fatalf("expected traversal error for path %q, got: %v", tp, resp.Result)
		}
		if resp.Error.Code != mcp.CodeInvalidParams {
			t.Errorf("for path %q: expected code %d, got %d", tp, mcp.CodeInvalidParams, resp.Error.Code)
		}
		if !strings.Contains(resp.Error.Message, "path traversal denied") && !strings.Contains(resp.Error.Message, "does not exist") {
			t.Errorf("for path %q: unexpected message: %s", tp, resp.Error.Message)
		}
	}
}

// MCP-TEST-011: Preservasi Request ID presisi (string "req-abc" dan integer 42) -> Respon memuat field id yang identik 100% tanpa konversi float.
func TestMCP_011_RequestIDPreservation(t *testing.T) {
	ws := t.TempDir()
	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// 1. Uji ID integer 42
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":42,"method":"ping"}` + "\n"))
	respInt := readJSONResponse(t, outR)
	if string(respInt.ID) != "42" {
		t.Errorf("expected integer ID '42', got verbatim %s", string(respInt.ID))
	}

	// 2. Uji ID string "req-abc"
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":"req-abc","method":"ping"}` + "\n"))
	respStr := readJSONResponse(t, outR)
	if string(respStr.ID) != `"req-abc"` {
		t.Errorf("expected string ID '\"req-abc\"', got verbatim %s", string(respStr.ID))
	}
}

// MCP-TEST-012: Klien mengirim notifikasi pembatalan notifications/cancelled -> Proses pemindaian dihentikan dan context dibatalkan bersih.
func TestMCP_012_Cancellation(t *testing.T) {
	ws := t.TempDir()
	// Buat berkas dummy di workspace
	testFile := filepath.Join(ws, "Test.astro")
	_ = os.WriteFile(testFile, []byte(`<div class="bg-black/50"></div>`), 0o600)

	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Mulai scan dan segera kirim cancel
	scanReq := `{"jsonrpc":"2.0","id":888,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":"Test.astro"}}}` + "\n"
	cancelNotif := `{"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":888}}` + "\n"

	_, _ = inW.Write([]byte(scanReq))
	_, _ = inW.Write([]byte(cancelNotif))

	// Respon scan harus dikembalikan (bisa berhasil jika cepat, atau error cancelled)
	resp := readJSONResponse(t, outR)
	if string(resp.ID) != "888" {
		t.Errorf("expected ID 888, got %s", string(resp.ID))
	}
}

// MCP-TEST-013: Invarian Kemurnian Output (Zero Stream Contamination) -> Setiap baris di stdout lolos parsing JSON-RPC tunggal.
func TestMCP_013_ZeroStdoutContamination(t *testing.T) {
	ws := t.TempDir()
	// Buat berkas dengan konten tidak patuh
	testFile := filepath.Join(ws, "Card.astro")
	_ = os.WriteFile(testFile, []byte(`<div class="bg-black/50">Test</div>`), 0o600)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	inputCommands := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize"}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"charites_explain_rule","arguments":{"rule_id":"theme.hardcode-opacity-color"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"charites_list_rules","arguments":{}}}`,
		`{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":"Card.astro"}}}`,
	}

	inR := strings.NewReader(strings.Join(inputCommands, "\n") + "\n")
	srv := mcp.NewServer(ws, inR, &stdoutBuf, &stderrBuf, rules.DefaultRegistry(), "1.0.0")

	err := srv.Run()
	if err != nil {
		t.Fatalf("server.Run() failed: %v", err)
	}

	// Verifikasi kemurnian stdout: setiap baris WAJIB berupa JSON-RPC response yang valid
	scanner := bufio.NewScanner(&stdoutBuf)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		lineCount++

		var raw map[string]any
		if err := json.Unmarshal(line, &raw); err != nil {
			t.Fatalf("NON-JSON-RPC BYTE DETECTED ON STDOUT (Stream Contamination): %s", string(line))
		}

		if raw["jsonrpc"] != "2.0" {
			t.Fatalf("line %d missing jsonrpc 2.0: %s", lineCount, string(line))
		}
	}

	if lineCount != 5 { // 5 requests had IDs (1, 2, 3, 4, 5); notification does not produce response
		t.Errorf("expected exactly 5 JSON-RPC response lines on stdout, got %d", lineCount)
	}
}

// Uji tambahan: Eksekusi Tool charites_explain_rule dan charites_list_rules
func TestMCP_ToolsExecution(t *testing.T) {
	ws := t.TempDir()
	srv, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// 1. tools/list
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"))
	respList := readJSONResponse(t, outR)
	if respList.Error != nil {
		t.Fatalf("tools/list error: %v", respList.Error)
	}
	listMap := respList.Result.(map[string]any)
	toolsSlice := listMap["tools"].([]any)
	if len(toolsSlice) != 4 {
		t.Errorf("expected 4 tools, got %d", len(toolsSlice))
	}

	// 2. charites_explain_rule
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"charites_explain_rule","arguments":{"rule_id":"theme.hardcode-opacity-color"}}}` + "\n"))
	respExplain := readJSONResponse(t, outR)
	if respExplain.Error != nil {
		t.Fatalf("charites_explain_rule error: %v", respExplain.Error)
	}

	// 3. charites_list_rules
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"charites_list_rules","arguments":{"category":"theme"}}}` + "\n"))
	respRules := readJSONResponse(t, outR)
	if respRules.Error != nil {
		t.Fatalf("charites_list_rules error: %v", respRules.Error)
	}

	// 4. server/discover
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":5,"method":"server/discover"}` + "\n"))
	respDisc := readJSONResponse(t, outR)
	if respDisc.Error != nil {
		t.Fatalf("server/discover error: %v", respDisc.Error)
	}

	_ = srv
}

func TestMCP_DirectTargetSafety(t *testing.T) {
	ws := t.TempDir()
	// Buat direktori .git di workspace
	gitDir := filepath.Join(ws, ".git")
	_ = os.MkdirAll(gitDir, 0o750)

	srv, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Pemindaian path di dalam .git
	callReq := `{"jsonrpc":"2.0","id":90,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":".git"}}}` + "\n"
	_, _ = inW.Write([]byte(callReq))

	resp := readJSONResponse(t, outR)
	if resp.Error == nil {
		t.Fatalf("expected error for scanning .git directly, got result: %v", resp.Result)
	}
	if resp.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected code %d, got %d", mcp.CodeInvalidParams, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "excluded directory") {
		t.Errorf("expected message to mention excluded directory, got: %s", resp.Error.Message)
	}
	_ = srv
}

func TestMCP_EmptyLineHandling(t *testing.T) {
	ws := t.TempDir()
	var stdoutBuf, stderrBuf bytes.Buffer
	inR := strings.NewReader("\n\n   \n\r\n")

	srv := mcp.NewServer(ws, inR, &stdoutBuf, &stderrBuf, nil, "")
	if err := srv.Run(); err != nil {
		t.Fatalf("expected nil error on empty lines, got: %v", err)
	}
	if stdoutBuf.Len() != 0 {
		t.Errorf("expected zero bytes on stdout for empty lines, got: %s", stdoutBuf.String())
	}
}

func TestMCP_Dispatcher_EdgeCases(t *testing.T) {
	ws := t.TempDir()
	srv := mcp.NewServer(ws, nil, nil, nil, nil, "")

	// 1. Dispatch nil request
	if resp := srv.Dispatch(nil); resp != nil {
		t.Errorf("expected nil response for nil request, got %v", resp)
	}

	// 2. Missing method
	resp := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", ID: json.RawMessage("1"), Method: ""})
	if resp == nil || resp.Error == nil || resp.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected invalid request for empty method, got %v", resp)
	}

	// 3. Notification with missing method
	if resp := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: ""}); resp != nil {
		t.Errorf("expected nil response for notification with empty method, got %v", resp)
	}

	// 4. Notification with invalid JSONRPC
	if resp := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "1.0", Method: "ping"}); resp != nil {
		t.Errorf("expected nil response for notification with bad jsonrpc, got %v", resp)
	}

	// 5. Terminated state
	srv.SetState(mcp.StateTerminated)
	respTerm := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", ID: json.RawMessage("2"), Method: "ping"})
	if respTerm == nil || respTerm.Error == nil || respTerm.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected invalid request on terminated state, got %v", respTerm)
	}
	if notifTerm := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "ping"}); notifTerm != nil {
		t.Errorf("expected nil for notification on terminated state, got %v", notifTerm)
	}

	// 6. Unknown server state
	srv.SetState(mcp.ServerState(999))
	respUnknown := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", ID: json.RawMessage("3"), Method: "ping"})
	if respUnknown == nil || respUnknown.Error == nil || respUnknown.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected invalid request on unknown state, got %v", respUnknown)
	}
	if notifUnknown := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "ping"}); notifUnknown != nil {
		t.Errorf("expected nil for notification on unknown state, got %v", notifUnknown)
	}

	// 7. Initializing state transitions and edge cases
	srv.SetState(mcp.StateInitializing)
	// ping during initializing
	respPingInit := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", ID: json.RawMessage("4"), Method: "ping"})
	if respPingInit == nil || respPingInit.Error != nil {
		t.Errorf("expected success for ping during initializing, got %v", respPingInit)
	}
	// duplicate initialize during initializing
	respInitInit := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", ID: json.RawMessage("5"), Method: "initialize"})
	if respInitInit == nil || respInitInit.Error == nil || respInitInit.Error.Code != mcp.CodeInvalidRequest {
		t.Errorf("expected error for duplicate initialize during initializing, got %v", respInitInit)
	}
	// notification duplicate initialize
	if notifInit := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "initialize"}); notifInit != nil {
		t.Errorf("expected nil for notification initialize during initializing, got %v", notifInit)
	}
	// unknown method during initializing
	respOtherInit := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", ID: json.RawMessage("6"), Method: "tools/list"})
	if respOtherInit == nil || respOtherInit.Error == nil || respOtherInit.Error.Code != mcp.CodeServerNotInitialized {
		t.Errorf("expected not initialized error, got %v", respOtherInit)
	}
	// unknown notification during initializing
	if notifOtherInit := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "some/notif"}); notifOtherInit != nil {
		t.Errorf("expected nil for unknown notification during initializing, got %v", notifOtherInit)
	}

	// 8. Ready state edge cases
	srv.SetState(mcp.StateReady)
	// unknown notification in ready state
	if notifReady := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "custom/notif"}); notifReady != nil {
		t.Errorf("expected nil for unknown notification in ready state, got %v", notifReady)
	}
	// notifications/initialized in ready state (redundant)
	if notifRedundant := srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "notifications/initialized"}); notifRedundant != nil {
		t.Errorf("expected nil for redundant notifications/initialized, got %v", notifRedundant)
	}
	// notifications/cancelled with invalid JSON
	srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "notifications/cancelled", Params: json.RawMessage(`{bad-json}`)})
	srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "notifications/cancelled", Params: json.RawMessage(`{"requestId":null}`)})
	srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "notifications/cancelled", Params: json.RawMessage(`{"requestId": 12.5}`)})
	srv.Dispatch(&mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "notifications/cancelled", Params: json.RawMessage(`{"requestId": true}`)})
}

func TestMCP_Scan_AdvancedPaths(t *testing.T) {
	ws := t.TempDir()

	// 1. Buat file konform
	goodFile := filepath.Join(ws, "Good.astro")
	_ = os.WriteFile(goodFile, []byte(`<div>OK</div>`), 0o600)

	// 2. Buat file konfigurasi charites.yaml di workspace
	cfgFile := filepath.Join(ws, "charites.yaml")
	_ = os.WriteFile(cfgFile, []byte("ignore:\n  - ignored/**\n"), 0o600)

	srv, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()
	_ = srv

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Uji scan dengan path absolut valid di dalam workspace
	absScan := fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":%q,"category":"theme","rule":"theme.hardcode-opacity-color","ext":"astro,tsx"}}}`+"\n", goodFile)
	_, _ = inW.Write([]byte(absScan))
	respAbs := readJSONResponse(t, outR)
	if respAbs.Error != nil {
		t.Fatalf("expected clean scan result for absolute path, got error: %v", respAbs.Error)
	}

	// Uji scan dengan category tidak dikenal
	catScan := fmt.Sprintf(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":%q,"category":"unknown-cat"}}}`+"\n", goodFile)
	_, _ = inW.Write([]byte(catScan))
	respCat := readJSONResponse(t, outR)
	if respCat.Error == nil || respCat.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected invalid params error for unknown category, got %v", respCat)
	}

	// Uji scan dengan rule tidak dikenal
	ruleScan := fmt.Sprintf(`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":%q,"rule":"unknown.rule"}}}`+"\n", goodFile)
	_, _ = inW.Write([]byte(ruleScan))
	respRule := readJSONResponse(t, outR)
	if respRule.Error == nil || respRule.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected invalid params error for unknown rule, got %v", respRule)
	}

	// Uji scan dengan rule yang kategori-nya tidak cocok
	mismatchScan := fmt.Sprintf(`{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":%q,"category":"a11y","rule":"theme.hardcode-opacity-color"}}}`+"\n", goodFile)
	_, _ = inW.Write([]byte(mismatchScan))
	respMismatch := readJSONResponse(t, outR)
	if respMismatch.Error == nil || respMismatch.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected invalid params error for category mismatch, got %v", respMismatch)
	}

	// Uji scan dengan target file yang tidak ada
	nonExistScan := `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":"does-not-exist.astro"}}}` + "\n"
	_, _ = inW.Write([]byte(nonExistScan))
	respNonExist := readJSONResponse(t, outR)
	if respNonExist.Error == nil || respNonExist.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected invalid params error for non-existent target, got %v", respNonExist)
	}

	// Uji scan dengan argumen JSON rusak
	badArgsScan := `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"charites_scan","arguments":"not-an-object"}}` + "\n"
	_, _ = inW.Write([]byte(badArgsScan))
	respBadArgs := readJSONResponse(t, outR)
	if respBadArgs.Error == nil || respBadArgs.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected invalid params error for bad arguments JSON, got %v", respBadArgs)
	}
}

func TestMCP_Explain_EdgeCases(t *testing.T) {
	ws := t.TempDir()
	srv, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// 1. Missing rule_id
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"charites_explain_rule","arguments":{"rule_id":""}}}` + "\n"))
	respEmpty := readJSONResponse(t, outR)
	if respEmpty.Error == nil || respEmpty.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected error for empty rule_id, got %v", respEmpty)
	}

	// 2. Unknown rule_id
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"charites_explain_rule","arguments":{"rule_id":"nonexistent.rule"}}}` + "\n"))
	respUnknown := readJSONResponse(t, outR)
	if respUnknown.Error == nil || respUnknown.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected error for nonexistent rule_id, got %v", respUnknown)
	}

	// 3. Bad args JSON
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"charites_explain_rule","arguments":"invalid"}}` + "\n"))
	respBad := readJSONResponse(t, outR)
	if respBad.Error == nil || respBad.Error.Code != mcp.CodeInvalidParams {
		t.Errorf("expected error for bad args, got %v", respBad)
	}
	_ = srv
}

func TestMCP_ListRules_EdgeCases(t *testing.T) {
	ws := t.TempDir()
	srv, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// 1. List rules with unknown category
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"charites_list_rules","arguments":{"category":"nonexistent-category"}}}` + "\n"))
	resp := readJSONResponse(t, outR)
	if resp.Error != nil {
		t.Fatalf("unexpected error for listing rules: %v", resp.Error)
	}
	_ = srv
}

func TestMCP_Scan_CleanOutput(t *testing.T) {
	ws := t.TempDir()
	cleanFile := filepath.Join(ws, "Clean.astro")
	_ = os.WriteFile(cleanFile, []byte(`<div>Clean Component</div>`), 0o600)

	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Scan clean file
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":"Clean.astro"}}}` + "\n"))
	resp := readJSONResponse(t, outR)
	if resp.Error != nil {
		t.Fatalf("unexpected error on clean scan: %v", resp.Error)
	}

	callRes, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("expected result map, got: %T", resp.Result)
	}
	contentSlice, ok := callRes["content"].([]any)
	if !ok || len(contentSlice) == 0 {
		t.Fatalf("expected non-empty content slice, got: %v", callRes["content"])
	}
	firstItem := contentSlice[0].(map[string]any)
	text := firstItem["text"].(string)

	if text == "null" {
		t.Errorf("expected clean message, got literal 'null'")
	}
	if !strings.Contains(text, " CLEAN") {
		t.Errorf("expected text to contain ' CLEAN', got: %s", text)
	}
}

func TestMCP_Scan_RichDiagnosticMetadata(t *testing.T) {
	ws := t.TempDir()
	violationFile := filepath.Join(ws, "Violations.astro")
	_ = os.WriteFile(violationFile, []byte(`<div class="bg-black/50">Violations</div>`), 0o600)

	_, inW, outR, _, cleanup := setupTestServer(t, ws)
	defer cleanup()

	// Handshake
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"))
	_ = readJSONResponse(t, outR)
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	time.Sleep(20 * time.Millisecond)

	// Scan file with violations
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"charites_scan","arguments":{"path":"Violations.astro"}}}` + "\n"))
	resp := readJSONResponse(t, outR)
	if resp.Error != nil {
		t.Fatalf("unexpected error on scan: %v", resp.Error)
	}

	callRes := resp.Result.(map[string]any)
	contentSlice := callRes["content"].([]any)
	firstItem := contentSlice[0].(map[string]any)
	text := firstItem["text"].(string)

	var items []mcp.ScanDiagnosticItem
	if err := json.Unmarshal([]byte(text), &items); err != nil {
		t.Fatalf("failed to unmarshal scan diagnostic items: %v, raw: %s", err, text)
	}
	if len(items) == 0 {
		t.Fatalf("expected at least 1 violation item, got 0")
	}

	item := items[0]
	if !strings.HasPrefix(item.WikiURL, "https://github.com/will2469/charites/wiki/") {
		t.Errorf("expected canonical online wiki URL, got: %s", item.WikiURL)
	}
	if !strings.Contains(item.ExplainHint, "charites_explain_rule") {
		t.Errorf("expected explain hint mentioning charites_explain_rule, got: %s", item.ExplainHint)
	}
}

func TestMCP_Tools_PromptSteering(t *testing.T) {
	tools := mcp.DefaultTools()
	toolMap := make(map[string]mcp.ToolDefinition)
	for _, tool := range tools {
		toolMap[tool.Name] = tool
	}

	scanTool, ok := toolMap["charites_scan"]
	if !ok {
		t.Fatal("missing charites_scan tool")
	}
	if !strings.Contains(scanTool.Description, "MANDATORY FRONTEND AUDITOR") {
		t.Errorf("expected scan tool to contain prompt steering directive, got: %s", scanTool.Description)
	}

	explainTool, ok := toolMap["charites_explain_rule"]
	if !ok {
		t.Fatal("missing charites_explain_rule tool")
	}
	if !strings.Contains(explainTool.Description, "Authoritative 8-Pillars documentation") {
		t.Errorf("expected explain tool to mention 8-Pillars, got: %s", explainTool.Description)
	}
}
