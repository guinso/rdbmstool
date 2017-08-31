package parser

func (t *Tree) parseExpression() {
	//accept:
	//1. literal
	//2. number
	//3. string
	//4. logical operator
	//5. function

	//push parentesis into stack
}

func (t *Tree) parseParentesis() {
	//consider complete if well received
	//1. '(',
	//2. expression,
	//3. ')'

	//otherwise throw error
}

func (t *Tree) parseLogicalExpresion() {
	//TODO: handle various simplest logical expression
	//1. a operator b
	//2. operator a
}

func (t *Tree) parseField() {

}

func (t *Tree) parseFunction() {
	//TODO: expand this function to verify various functions
	//e.g. MAX(), MIN(), COUNT(), AVG(), etc.
}

func (t *Tree) parseJoin() {
	//TODO: handle all types of JOINs
}

func (t *Tree) parseSelect() {

}
