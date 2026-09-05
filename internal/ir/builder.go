package ir

import "strings"

// MaxStackDepth adalah batas kedalaman maksimum hierarki IR (256 tingkat).
// Elemen pada tingkat 257+ disematkan sebagai flat siblings di bawah node tingkat-256.
const MaxStackDepth = 256

// Builder merakit token mentah dari parser frontend menjadi pohon IR *ir.Node yang terpadu.
// Mengelola relasi Parent/Children, pemulihan input cacat (recovery semantics),
// dan pembatasan kedalaman stack (stack guard 256).
type Builder struct {
	root  *Node
	stack []*Node
}

// NewBuilder menginisialisasi Builder baru dengan kapasitas stack awal teralokasi.
func NewBuilder() *Builder {
	return &Builder{
		stack: make([]*Node, 0, 32),
	}
}

// OpenElement membuat node elemen baru, menyematkannya ke hierarki pohon,
// dan mendorongnya (push) ke stack jika batas kedalaman MaxStackDepth belum tercapai.
func (b *Builder) OpenElement(tag string, span Span, rawClasses string, classes []string, attrs map[string]string, hasDynamic bool) *Node {
	node := &Node{
		Type:              NodeElement,
		Tag:               tag,
		Span:              span,
		RawClasses:        rawClasses,
		Classes:           classes,
		Attributes:        attrs,
		HasDynamicClasses: hasDynamic,
	}
	b.attachAndPush(node)
	return node
}

// OpenFragment membuat node fragment (<> ... </>) dan mendorongnya ke stack jika batas kedalaman belum tercapai.
func (b *Builder) OpenFragment(span Span) *Node {
	node := &Node{
		Type: NodeFragment,
		Span: span,
	}
	b.attachAndPush(node)
	return node
}

// AddSelfClosingElement membuat elemen mandiri (tanpa closing tag, misal <img>, <input>, atau <Component />)
// dan menyematkannya ke parent aktif tanpa memasukkannya ke dalam stack.
func (b *Builder) AddSelfClosingElement(tag string, span Span, rawClasses string, classes []string, attrs map[string]string, hasDynamic bool) *Node {
	node := &Node{
		Type:              NodeElement,
		Tag:               tag,
		Span:              span,
		RawClasses:        rawClasses,
		Classes:           classes,
		Attributes:        attrs,
		HasDynamicClasses: hasDynamic,
	}
	b.attach(node)
	return node
}

// AddText menyematkan node teks literal di bawah elemen aktif tanpa memasukkannya ke stack.
func (b *Builder) AddText(text string, span Span) *Node {
	node := &Node{
		Type:       NodeText,
		RawClasses: text,
		Span:       span,
	}
	b.attach(node)
	return node
}

// AddComment menyematkan node komentar di bawah elemen aktif tanpa memasukkannya ke stack.
func (b *Builder) AddComment(comment string, span Span) *Node {
	node := &Node{
		Type:       NodeComment,
		RawClasses: comment,
		Span:       span,
	}
	b.attach(node)
	return node
}

// CloseElement mencari tag pencocok dari puncak stack ke bawah.
// Jika ditemukan: pop seluruh node dari puncak stack hingga elemen yang cocok (stack unwinding).
// Elemen perantara yang tidak ditutup secara eksplisit tetap menjadi anak sah dari parent masing-masing.
// Jika tidak ditemukan: buang token penutup secara hening, stack tidak berubah (unmatched discard).
func (b *Builder) CloseElement(tag string) {
	if len(b.stack) == 0 {
		return
	}

	if tag == "" || tag == "<>" {
		b.CloseFragment()
		return
	}

	// Cari kecocokan tag dari puncak stack ke bawah
	matchIdx := -1
	for i := len(b.stack) - 1; i >= 0; i-- {
		n := b.stack[i]
		if n.Type == NodeElement && (n.Tag == tag || strings.EqualFold(n.Tag, tag)) {
			matchIdx = i
			break
		}
	}

	if matchIdx != -1 {
		// Pop seluruh node dari puncak hingga matchIdx
		b.stack = b.stack[:matchIdx]
	}
	// Jika matchIdx == -1: unmatched closing tag dibuang secara hening, stack tidak berubah.
}

// CloseFragment menutup fragment container (<></>) terdekat di stack.
func (b *Builder) CloseFragment() {
	if len(b.stack) == 0 {
		return
	}

	matchIdx := -1
	for i := len(b.stack) - 1; i >= 0; i-- {
		if b.stack[i].Type == NodeFragment {
			matchIdx = i
			break
		}
	}

	if matchIdx != -1 {
		b.stack = b.stack[:matchIdx]
	}
}

// Root mengembalikan simpul akar dari pohon IR yang telah dirakit.
// Jika tidak ada elemen yang diproses, mengembalikan node fragment kosong.
func (b *Builder) Root() *Node {
	if b.root == nil {
		return &Node{Type: NodeFragment}
	}
	return b.root
}

// StackDepth mengembalikan jumlah elemen aktif di dalam stack saat ini.
func (b *Builder) StackDepth() int {
	return len(b.stack)
}

// Reset mengosongkan state builder untuk dapat digunakan kembali.
func (b *Builder) Reset() {
	b.root = nil
	b.stack = b.stack[:0]
}

// attachAndPush menyematkan node ke parent aktif dan memasukkannya ke stack jika kedalaman < 256.
func (b *Builder) attachAndPush(node *Node) {
	if len(b.stack) == 0 {
		if b.root == nil {
			b.root = node
		} else {
			// Dokumen memiliki lebih dari satu elemen akar: bungkus ke dalam NodeFragment
			if b.root.Type != NodeFragment || b.root.Tag != "" {
				oldRoot := b.root
				b.root = &Node{
					Type:     NodeFragment,
					Children: []*Node{oldRoot},
				}
				oldRoot.Parent = b.root
			}
			b.root.Children = append(b.root.Children, node)
			node.Parent = b.root
		}
		b.stack = append(b.stack, node)
		return
	}

	parent := b.stack[len(b.stack)-1]
	node.Parent = parent
	parent.Children = append(parent.Children, node)

	// Nesting guard 256: jika sudah mencapai 256 tingkat, jangan di-push ke stack
	// Elemen ini menjadi flat sibling di bawah node tingkat-256.
	if len(b.stack) < MaxStackDepth {
		b.stack = append(b.stack, node)
	}
}

// attach menyematkan node ke parent aktif tanpa memasukkannya ke stack.
func (b *Builder) attach(node *Node) {
	if len(b.stack) == 0 {
		if b.root == nil {
			b.root = node
		} else {
			if b.root.Type != NodeFragment || b.root.Tag != "" {
				oldRoot := b.root
				b.root = &Node{
					Type:     NodeFragment,
					Children: []*Node{oldRoot},
				}
				oldRoot.Parent = b.root
			}
			b.root.Children = append(b.root.Children, node)
			node.Parent = b.root
		}
		return
	}

	parent := b.stack[len(b.stack)-1]
	node.Parent = parent
	parent.Children = append(parent.Children, node)
}
