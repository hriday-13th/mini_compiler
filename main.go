package main

import (
	"log"
)

// The tokenizer

type token struct {
	kind  string
	value string
}

func tokenizer(input string) []token {
	input += "\n"
	current := 0

	tokens := []token{}

	for current < len([]rune(input)) {
		char := string([]rune(input)[current])

		if char == "(" {
			tokens = append(tokens, token{
				kind:  "paren",
				value: "(",
			})
			current++
			continue
		}

		if char == ")" {
			tokens = append(tokens, token{
				kind:  "paren",
				value: ")",
			})
			current++
			continue
		}

		if char == " " {
			current++
			continue
		}

		if isNumber(char) {
			value := ""
			for isNumber(char) {
				value += char
				current++
				char = string([]rune(input)[current])
			}
			tokens = append(tokens, token{
				kind:  "number",
				value: value,
			})
			continue
		}

		if isLetter(char) {
			value := ""
			for isLetter(char) {
				value += char
				current++
				char = string([]rune(input)[current])
			}
			tokens = append(tokens, token{
				kind:  "name",
				value: value,
			})
			continue
		}
		break
	}
	return tokens
}

func isNumber(char string) bool {
	if char == "" {
		return false
	}
	n := []rune(char)[0]
	if n >= '0' && n <= '9' {
		return true
	}
	return false
}

func isLetter(char string) bool {
	if char == "" {
		return false
	}
	n := []rune(char)[0]
	if n >= 'a' && n <= 'z' {
		return true
	}
	return false
}

// The parser

type node struct {
	kind       string
	value      string
	name       string
	callee     *node
	expression *node
	body       []node
	params     []node
	arguments  *[]node
	context    *[]node
}

type ast node

var pc int

var pt []token

func parser(tokens []token) ast {
	pc = 0
	pt = tokens

	ast := ast{
		kind: "Program",
		body: []node{},
	}

	for pc < len(pt) {
		ast.body = append(ast.body, walk())
	}

	return ast
}

func walk() node {
	token := pt[pc]

	if token.kind == "number" {
		pc++
		return node{
			kind:  "NumberLiteral",
			value: token.value,
		}
	}

	if token.kind == "paren" && token.value == "(" {
		pc++
		token = pt[pc]

		n := node{
			kind:   "CallExpression",
			name:   token.value,
			params: []node{},
		}

		pc++
		token = pt[pc]

		for token.kind != "paren" || (token.kind == "paren" && token.value != ")") {
			n.params = append(n.params, walk())
			token = pt[pc]
		}
		pc++
		return n
	}
	log.Fatal(token.kind)
	return node{}
}

// The traverser

type visitor map[string]func(n *node, p node)

func traverser(a ast, v visitor) {
	traverseNode(node(a), node{}, v)
}

func traverseArray(a []node, p node, v visitor){
	for _, child := range a{
		traverseNode(child, p, v)
	}
}

func traverseNode(n, p node, v visitor) {
	for k, va := range v {
		if k == n.kind {
			va(&n, p)
		}
	}
	switch n.kind {
	case "Program":
		traverseArray(n.body, n, v)
		break

	case "CallExpression":
		traverseArray(n.params, n, v)
		break

	case "NumberLiteral":
		break

	default:
		log.Fatal(n.kind)
	}
}

// The transformer

func transformer(a ast) ast {
	nast := ast{
		kind: "Program",
		body: []node{},
	}

	a.context = &nast.body

	traverser(a, map[string]func(n *node, p node)){
		"NumberLiteral": func(n *node, p node) {
			*p.context = append(*p.context, node(
				kind: "NumberLiteral",
				value = n.value,
			))
		}
		"CallExpression": func(n *node, p node) {
			e := node{
				kind: "CallExpression",
				callee: &node{
					kind: "Identifier",
					name: n.name,
				},
				arguments: new([]node),
			}
			n.context = e.argumentsconst
			if p.kind != "CallExpression" {
				es := node{
					kind: "ExpressionStatement",
					expression: &e,
				}
				*p.context = append(*p.context, es)
			} else {
				*p.context = append(*p.context, e)
			}
		},
	}
	return nast
}

// The Code generator
