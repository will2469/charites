package mcp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"sync"

	"encoding/json"

	"github.com/will2469/charites/internal/rules"
)

// Server mengelola loop komunikasi I/O berbasis Stdio, siklus state machine, dan dispatching tool.
type Server struct {
	state       ServerState
	workspace   string
	version     string
	reg         *rules.Registry
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	stdoutMu    sync.Mutex
	activeScans sync.Map // string(key) -> context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

// NewServer membuat instance Server MCP baru.
func NewServer(workspace string, stdin io.Reader, stdout, stderr io.Writer, reg *rules.Registry, version string) *Server {
	if reg == nil {
		reg = rules.DefaultRegistry()
	}
	if version == "" {
		version = "1.0.0"
	}
	if workspace == "" {
		workspace = "."
	}

	return &Server{
		state:     StateNew,
		workspace: workspace,
		version:   version,
		reg:       reg,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
	}
}

// GetState mengembalikan status terkini dari state machine secara thread-safe.
func (s *Server) GetState() ServerState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetState memperbarui status state machine secara thread-safe.
func (s *Server) SetState(st ServerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = st
}

// Run memulai loop pemrosesan pesan dari stdin hingga mencapai EOF atau error I/O fatal.
func (s *Server) Run() error {
	reader := bufio.NewReader(s.stdin)
	var buf bytes.Buffer

	for {
		chunk, isPrefix, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if buf.Len() > 0 {
					s.handleRawLine(buf.Bytes())
				}
				s.wg.Wait()
				s.SetState(StateTerminated)
				return nil
			}
			return err
		}

		if buf.Len()+len(chunk) > MaxFrameSize {
			// Frame melebihi batas protokol 4MB
			for isPrefix {
				_, isPrefix, _ = reader.ReadLine()
			}
			buf.Reset()
			_ = s.WriteResponse(NewErrorResponse(nil, CodeParseError, "Parse error: frame size exceeds 4MB limit"))
			continue
		}

		buf.Write(chunk)
		if !isPrefix {
			line := make([]byte, buf.Len())
			copy(line, buf.Bytes())
			buf.Reset()
			s.handleRawLine(line)
		}
	}
}

func (s *Server) handleRawLine(line []byte) {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(trimmed, &req); err != nil {
		_ = s.WriteResponse(NewErrorResponse(nil, CodeParseError, "Parse error"))
		return
	}

	s.processRequest(&req)
}

func (s *Server) processRequest(req *JSONRPCRequest) {
	// Notifikasi pembatalan diproses segera secara sinkron
	if req.Method == "notifications/cancelled" {
		s.handleCancelled(req.Params)
		return
	}

	// Pemanggilan tool dijalankan pada goroutine agar pembatalan interaktif tidak terblokir
	if req.Method == "tools/call" {
		s.wg.Add(1)
		go func(r *JSONRPCRequest) {
			defer s.wg.Done()
			resp := s.Dispatch(r)
			if resp != nil {
				_ = s.WriteResponse(resp)
			}
		}(req)
		return
	}

	// Permintaan handshake dan metode lainnya dieksekusi secara berurutan
	resp := s.Dispatch(req)
	if resp != nil {
		_ = s.WriteResponse(resp)
	}
}

// WriteResponse menulis frame JSON-RPC keluar ke stdout secara thread-safe dan atomik diakhiri LF (\n).
// Invarian Mutlak: Hanya fungsi ini yang berhak menulis ke stdout untuk menjamin zero stream contamination.
func (s *Server) WriteResponse(resp *JSONRPCResponse) error {
	if resp == nil {
		return nil
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	s.stdoutMu.Lock()
	defer s.stdoutMu.Unlock()

	data = append(data, '\n')
	_, writeErr := s.stdout.Write(data)
	return writeErr
}
