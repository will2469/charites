package ir

import "iter"

// NodeType merepresentasikan klasifikasi jenis node di dalam pohon IR Charites.
type NodeType uint8

const (
	// NodeElement merepresentasikan elemen HTML atau JSX/Astro (contoh: <div>, <Card>, <button>).
	NodeElement NodeType = iota
	// NodeText merepresentasikan teks literal di dalam elemen.
	NodeText
	// NodeComment merepresentasikan komentar HTML/JSX (<!-- ... --> atau {/* ... */}).
	NodeComment
	// NodeFragment merepresentasikan container fragment (<> ... </>).
	NodeFragment
)

// Span menyimpan posisi rentang kode sumber asli berbasis 1-indexed.
type Span struct {
	Line      int // Posisi baris awal (1-indexed)
	Column    int // Posisi kolom awal (1-indexed)
	EndLine   int // Posisi baris akhir (1-indexed)
	EndColumn int // Posisi kolom akhir (1-indexed)
}

// Node adalah representasi pohon sintaks seragam (Unified AST/IR) untuk Astro, React TSX, dan HTML.
// Memiliki tata letak memori teroptimasi (136 bytes pada target 64-bit) dengan 7-byte padding eksplisit.
type Node struct {
	Span              Span              // 32 bytes (4 int @ 8 bytes)
	Parent            *Node             // 8 bytes (pointer ke node induk)
	Tag               string            // 16 bytes (pointer + panjang string)
	RawClasses        string            // 16 bytes (string asli atribut class sebelum pemisahan)
	Attributes        map[string]string // 8 bytes (pointer map atribut mentah)
	Classes           []string          // 24 bytes (slice token class yang sudah di-split & di-trim)
	Children          []*Node           // 24 bytes (slice pointer ke node anak)
	Type              NodeType          // 1 byte (uint8)
	HasDynamicClasses bool              // 1 byte (boolean penanda keberadaan kelas dinamis interpolasi ${...})
	_                 [6]byte           // 6 bytes explicit padding agar rata dengan word boundary 8-byte
}

// Walk melakukan pre-order depth-first traversal pada pohon IR.
// Menggunakan iterator natif Go 1.26 (iter.Seq) untuk iterasi zero-copy tanpa alokasi slice sementara.
func (n *Node) Walk() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		var walk func(curr *Node) bool
		walk = func(curr *Node) bool {
			if curr == nil {
				return true
			}
			if !yield(curr) {
				return false
			}
			for _, child := range curr.Children {
				if !walk(child) {
					return false
				}
			}
			return true
		}
		walk(n)
	}
}

// HasClass memeriksa apakah kelas CSS tertentu ada di dalam slice Classes yang sudah ditokenisasi.
func (n *Node) HasClass(name string) bool {
	for _, c := range n.Classes {
		if c == name {
			return true
		}
	}
	return false
}

// GetAttr mengambil nilai atribut berdasarkan nama dari peta Attributes secara aman.
func (n *Node) GetAttr(name string) (string, bool) {
	if n.Attributes == nil {
		return "", false
	}
	val, ok := n.Attributes[name]
	return val, ok
}

// IsElement memeriksa apakah node merupakan NodeElement. Jika parameter tags diberikan,
// node harus memiliki tag yang cocok dengan salah satu nilai dalam tags.
func (n *Node) IsElement(tags ...string) bool {
	if n == nil || n.Type != NodeElement {
		return false
	}
	if len(tags) == 0 {
		return true
	}
	for _, t := range tags {
		if n.Tag == t {
			return true
		}
	}
	return false
}
