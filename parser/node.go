package parser

const (
	NodeBool       NodeType = iota //a boolean constant
	NodeField                      //a db field such as table name, column name
	NodeIdentifier                 //a function name
	NodeNumber                     //a numerical constant
	NodeString                     //a string constant
	NodeParam                      //an SQL parameter
	NodeList                       //a list of SQL expression
	NodeCondition                  //a logical comparison expression; e.g. x > y
	NodeSelect                     //SQL select statement
	NodeFrom                       //SQL from statement
	NodeJoin                       //SQL join statement
	NodeWhere                      //SQL where statement
	NodeHaving                     //SQL having statement
	NodeGroupBy                    //SQL group by statement
	NodeLimit                      //SQl limit statement
	NodeUnion                      //SQL union statement
	NodeQuery                      //SQL query statement
)

/********************* Select Node *******************/

//SelectNode node for Select statement
type SelectNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *SelectNode) Type() NodeType {
	return NodeSelect
}

//String return node name
func (node *SelectNode) String() string {
	return node.name
}

//Position get current node position
func (node *SelectNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *SelectNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *SelectNode) Copy() Node {
	//TODO: clone SelectNode
	return &SelectNode{
		name:     "select",
		position: 0,
		root:     node.root.Copy()}
}

/********************* From Node *****************/

//FromNode node for From statement
type FromNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *FromNode) Type() NodeType {
	return NodeFrom
}

//String return node name
func (node *FromNode) String() string {
	return node.name
}

//Position get current node position
func (node *FromNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *FromNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *FromNode) Copy() Node {
	//TODO: clone FromNode
	return &FromNode{
		name:     "from",
		position: 0,
		root:     node.root.Copy()}
}

/************************* Where Node ******************/

//WhereNode node for Where statement
type WhereNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *WhereNode) Type() NodeType {
	return NodeWhere
}

//String return node name
func (node *WhereNode) String() string {
	return node.name
}

//Position get current node position
func (node *WhereNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *WhereNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *WhereNode) Copy() Node {
	//TODO: clone FromNode
	return &WhereNode{
		name:     "where",
		position: 0,
		root:     node.root.Copy()}
}

/********************** Having Node *********************/

//HavingNode node for Having statement
type HavingNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *HavingNode) Type() NodeType {
	return NodeHaving
}

//String return node name
func (node *HavingNode) String() string {
	return node.name
}

//Position get current node position
func (node *HavingNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *HavingNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *HavingNode) Copy() Node {
	//TODO: clone FromNode
	return &HavingNode{
		name:     "having",
		position: 0,
		root:     node.root.Copy()}
}

/****************************** GroupBy Node ******************/

//GroupByNode node for Group By statement
type GroupByNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *GroupByNode) Type() NodeType {
	return NodeGroupBy
}

//String return node name
func (node *GroupByNode) String() string {
	return node.name
}

//Position get current node position
func (node *GroupByNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *GroupByNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *GroupByNode) Copy() Node {
	//TODO: clone FromNode
	return &GroupByNode{
		name:     "group by",
		position: 0,
		root:     node.root.Copy()}
}

/******************* Limit Node ********************/

//LimitNode node for Limit statement
type LimitNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *LimitNode) Type() NodeType {
	return NodeLimit
}

//String return node name
func (node *LimitNode) String() string {
	return node.name
}

//Position get current node position
func (node *LimitNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *LimitNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *LimitNode) Copy() Node {
	//TODO: clone FromNode
	return &LimitNode{
		name:     "limit",
		position: 0,
		root:     node.root.Copy()}
}

/*************************** Union Node ******************/

//UnionNode node for Union statement
type UnionNode struct {
	name     string
	position Pos
	root     *Tree
}

//Type return node type
func (node *UnionNode) Type() NodeType {
	return NodeUnion
}

//String return node name
func (node *UnionNode) String() string {
	return node.name
}

//Position get current node position
func (node *UnionNode) Position() Pos {
	return node.position.Position()
}

//Tree get node's tree
func (node *UnionNode) tree() *Tree {
	return node.root
}

//Copy deep copy current node
func (node *UnionNode) Copy() Node {
	//TODO: clone FromNode
	return &UnionNode{
		name:     "union",
		position: 0,
		root:     node.root.Copy()}
}
