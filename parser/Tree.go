package parser

//Refered from https://github.com/golang/go/blob/master/src/text/template/parse/parse.go and
//             https://github.com/golang/go/blob/master/src/text/template/parse/node.go

/****************************** Tree *********************************/

// Tree is the representation of a single parsed template.
type Tree struct {
	Name      string // name of the template represented by the tree.
	ParseName string // name of the top-level template during parsing, for error messages.
	Root      []Node // top-level root of the tree.
	text      string // text parsed to create the template (or its parent)
	// Parsing only; cleared after parse.
	funcs     []map[string]interface{}
	lex       *lexer
	token     [3]tokenItem // three-token lookahead for parser.
	peekCount int
	vars      []string // variables defined at the moment.
	treeSet   map[string]*Tree
}

// Copy returns a copy of the Tree. Any parsing state is discarded.
func (t *Tree) Copy() *Tree {
	if t == nil {
		return nil
	}

	root := []Node{}
	for i := 0; i < len(t.Root); i++ {
		root = append(root, t.Root[i].Copy())
	}

	return &Tree{
		Name:      t.Name,
		ParseName: t.ParseName,
		Root:      root,
		text:      t.text,
	}
}

/***************************** Node ***************************************/

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
	Type() NodeType
	String() string
	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node
	Position() Pos // byte position of start of node in full original input string
	// tree returns the containing *Tree.
	// It is unexported so all implementations of Node are in this package.
	tree() *Tree
}

// NodeType identifies the type of a parse tree node.
type NodeType int

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}
