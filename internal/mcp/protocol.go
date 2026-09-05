package mcp

import (
	"encoding/json"
)

// ServerState mendefinisikan siklus hidup state machine server MCP.
type ServerState int

const (
	// StateNew adalah status awal server sebelum menerima request initialize.
	StateNew ServerState = iota
	// StateInitializing adalah status setelah menerima initialize namun menunggu notifications/initialized.
	StateInitializing
	// StateReady adalah status server siap melayani tools/list, tools/call, dsb.
	StateReady
	// StateTerminated adalah status setelah shutdown (EOF pada stdin).
	StateTerminated
)

const (
	// ProtocolVersion merepresentasikan versi protokol MCP yang didukung.
	ProtocolVersion = "2026-07-28"
	// ServerName merepresentasikan nama server MCP kanonikal.
	ServerName = "charites"
	// MaxFrameSize merepresentasikan batas aman ukuran frame satu pesan JSON-RPC (4 Megabytes).
	MaxFrameSize = 4 * 1024 * 1024
)

// Matriks Kode Error Standar JSON-RPC 2.0 & MCP.
const (
	CodeParseError           = -32700
	CodeInvalidRequest       = -32600
	CodeMethodNotFound       = -32601
	CodeInvalidParams        = -32602
	CodeInternalError        = -32603
	CodeInternalToolError    = -32000
	CodeServerNotInitialized = -32002
)

// JSONRPCRequest merepresentasikan struktur pesan request atau notifikasi masuk.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// IsNotification mengembalikan true jika pesan tidak memiliki field ID.
func (r *JSONRPCRequest) IsNotification() bool {
	return len(r.ID) == 0 || string(r.ID) == "null"
}

// RPCError merepresentasikan objek error terstandar JSON-RPC 2.0.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSONRPCResponse merepresentasikan struktur respon JSON-RPC 2.0 keluar.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// NewResultResponse membuat response sukses baru dengan preservasi Request ID presisi.
func NewResultResponse(id json.RawMessage, result any) *JSONRPCResponse {
	respID := id
	if len(respID) == 0 {
		respID = json.RawMessage("null")
	}
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      respID,
		Result:  result,
	}
}

// NewErrorResponse membuat response error baru dengan preservasi Request ID presisi.
func NewErrorResponse(id json.RawMessage, code int, message string) *JSONRPCResponse {
	respID := id
	if len(respID) == 0 {
		respID = json.RawMessage("null")
	}
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      respID,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
}
