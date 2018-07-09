package txfun

import (
	"bytes"
	"fmt"
)

type node struct {
	key     []byte
	value   []byte
	next    *node
	prev    *node
	created uint64
}

type list struct {
	root *node
}

func newList() *list {
	return &list{}
}

func (l *list) insert(key, value []byte) *node {
	newNode := &node{
		key:   key,
		value: value,
	}

	if l.root == nil {
		l.root = newNode
		return newNode
	}

	var n *node

	for n = l.root; n != nil; n = n.next {
		cmp := bytes.Compare(key, n.key)
		switch {
		case cmp < 0:
			// new node goes before n
			if n == l.root {
				// insert before root
				l.root = newNode
				newNode.next = n
				n.prev = newNode
			} else {
				prev := n.prev
				newNode.next = n
				newNode.prev = prev
				prev.next = newNode
				n.prev = newNode
			}
			return newNode
		case cmp == 0:
			n.value = value
			return newNode
		}

		if n.next == nil {
			n.next = newNode
			newNode.prev = n
			return newNode
		}
	}

	return newNode
}

func (l *list) String() string {
	str := ""
	for n := l.root; n != nil; n = n.next {
		str += fmt.Sprintf("[%s, %s, %d]", string(n.key), string(n.value), n.created)
		if n.next != nil {
			str += "->"
		}
	}

	return str
}