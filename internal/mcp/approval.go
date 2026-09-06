package mcp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	// ErrTokenNotFound terjadi saat token persetujuan tidak ditemukan atau sudah dikonsumsi.
	ErrTokenNotFound = errors.New("approval token not found or already consumed")
	// ErrTokenExpired terjadi saat masa berlaku token (TTL) telah habis.
	ErrTokenExpired = errors.New("approval token has expired (validity is 10 minutes)")
	// ErrHashMismatch terjadi saat konten issue diubah setelah token dihasilkan (anti-tampering).
	ErrHashMismatch = errors.New("issue content has been tampered with or does not match approved token")
)

const (
	// DefaultTokenTTL mendefinisikan masa berlaku token persetujuan (10 menit).
	DefaultTokenTTL = 10 * time.Minute
)

// ApprovalToken menyimpan metadata token persetujuan pengguna.
type ApprovalToken struct {
	Token     string
	IssueHash string
	CreatedAt time.Time
	ExpiresAt time.Time
	Used      bool
}

// ApprovalManager mengelola token persetujuan kriptografis dua tahap (HITL).
type ApprovalManager struct {
	mu     sync.Mutex
	tokens map[string]*ApprovalToken
}

var (
	defaultApprovalOnce    sync.Once
	defaultApprovalManager *ApprovalManager
)

// DefaultApprovalManager mengembalikan instance singleton ApprovalManager default.
func DefaultApprovalManager() *ApprovalManager {
	defaultApprovalOnce.Do(func() {
		defaultApprovalManager = NewApprovalManager()
	})
	return defaultApprovalManager
}

// NewApprovalManager membuat instance ApprovalManager baru.
func NewApprovalManager() *ApprovalManager {
	return &ApprovalManager{
		tokens: make(map[string]*ApprovalToken),
	}
}

// ComputeIssueHash menghitung digest SHA-256 kanonikal untuk parameter issue.
// Menggunakan pemisah null byte (\x00) antar kolom guna mencegah delimiter collision attacks.
func ComputeIssueHash(ruleID, title, description, snippet, category string) string {
	h := sha256.New()
	h.Write([]byte(ruleID))
	h.Write([]byte{0})
	h.Write([]byte(title))
	h.Write([]byte{0})
	h.Write([]byte(description))
	h.Write([]byte{0})
	h.Write([]byte(snippet))
	h.Write([]byte{0})
	h.Write([]byte(category))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateToken menghasilkan token persetujuan kriptografis acak (appr_<hex>) dengan durasi TTL tertentu.
// Membersihkan token kadaluwarsa secara in-flight guna mencegah memory leak.
func (m *ApprovalManager) CreateToken(issueHash string, ttl time.Duration) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for k, tok := range m.tokens {
		if now.After(tok.ExpiresAt) || tok.Used {
			delete(m.tokens, k)
		}
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tokenStr := "appr_" + hex.EncodeToString(b)

	if ttl <= 0 {
		ttl = DefaultTokenTTL
	}

	m.tokens[tokenStr] = &ApprovalToken{
		Token:     tokenStr,
		IssueHash: issueHash,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		Used:      false,
	}

	return tokenStr, nil
}

// ConsumeToken memvalidasi dan mengonsumsi token secara atomik.
// Mengembalikan error jika token tidak ditemukan, sudah digunakan, kadaluwarsa, atau hash tidak cocok.
func (m *ApprovalManager) ConsumeToken(tokenStr string, expectedIssueHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tok, ok := m.tokens[tokenStr]
	if !ok || tok.Used {
		return ErrTokenNotFound
	}

	if time.Now().After(tok.ExpiresAt) {
		delete(m.tokens, tokenStr)
		return ErrTokenExpired
	}

	if tok.IssueHash != expectedIssueHash {
		return ErrHashMismatch
	}

	tok.Used = true
	delete(m.tokens, tokenStr)
	return nil
}
