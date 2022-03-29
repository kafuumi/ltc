package glist

import (
	"math"
)

const (
	// CAPACITY 链表的最大容量
	CAPACITY uint = math.MaxUint
	ZERO          = uint(0) //uint类型的0
)

// Node 链表中的一个结点
type Node[E any] struct {
	element E        //保存的内容
	prev    *Node[E] //前一个结点
	next    *Node[E] //后一个结点
}

// Clone 克隆Node,返回的Node的Prev和Next均为nil,Element保持不变
func (n *Node[E]) Clone() *Node[E] {
	node := &Node[E]{
		element: n.element,
		prev:    nil,
		next:    nil,
	}
	return node
}

// LinkedList 链表,实现了List
type LinkedList[E any] struct {
	len   uint     //链表中元素个数
	first *Node[E] //头指针
	last  *Node[E] //尾指针
}

// NewLinkedList 创建一个链表,
//列表的最大容量为uint类型的最大值
func NewLinkedList[E any]() *LinkedList[E] {
	return &LinkedList[E]{
		len:   0,
		first: nil,
		last:  nil,
	}
}

func (l *LinkedList[E]) Size() uint {
	return l.len
}

func (l *LinkedList[E]) IsEmpty() bool {
	return l.len == 0
}

func (l *LinkedList[E]) IsNotEmpty() bool {
	return l.len != 0
}

func (l *LinkedList[E]) Append(element E) bool {
	//超出最大值无法添加
	if l.len == CAPACITY {
		return false
	}
	node := &Node[E]{
		element: element,
		prev:    nil,
		next:    nil,
	}
	//链表为空,头指针指向该结点
	if l.first == nil {
		l.first = node
		l.last = node
	} else {
		//链表不为空,添加到尾部
		node.prev = l.last
		l.last.next = node
		l.last = node
	}
	l.len++
	return true
}

func (l *LinkedList[E]) Insert(index uint, element E) bool {
	//当前size已经达到最大值或者索引越界
	if l.len == CAPACITY || index > l.len {
		return false
	}
	node := &Node[E]{
		element: element,
		prev:    nil,
		next:    nil,
	}
	//插入头部
	if index == 0 {
		if l.first == nil {
			//链表为空
			l.first = node
			l.last = node
		} else {
			//链表不为空
			node.next = l.first
			l.first.prev = node
			l.first = node
		}
	} else if index == l.len {
		//插入尾部
		l.last.next = node
		node.prev = l.last
		l.last = node
	} else {
		var prev *Node[E]
		head := l.first
		for i := ZERO; i < index; i++ {
			prev = head
			head = head.next
		}
		node.next = head
		node.prev = prev
		prev.next = node
		head.prev = node
	}
	l.len++
	return true
}

func (l *LinkedList[E]) Remove(index uint) bool {
	if index >= l.len {
		return false
	}
	head := l.first
	var prev *Node[E]
	for i := ZERO; i < index; i++ {
		prev = head
		head = head.next
	}
	//删除第一个结点
	if head == l.first {
		l.first.next = nil
		l.first = head.next
	} else if head == l.last {
		//删除最后一个结点
		l.last = prev
		l.last.next = nil
	} else {
		prev.next = head.next
		head.next.prev = prev
	}
	l.len--
	return true
}

func (l *LinkedList[E]) Get(index uint) *E {
	if index >= l.len {
		return nil
	}
	node := l.first
	for i := ZERO; i < index; i++ {
		node = node.next
	}
	return &(node.element)
}

func (l *LinkedList[E]) Set(index uint, element E) bool {
	if index >= l.len {
		return false
	}
	node := l.first
	for i := ZERO; i < index; i++ {
		node = node.next
	}
	node.element = element
	return true
}

func (l *LinkedList[E]) PushBack(element E) bool {
	return l.Append(element)
}

func (l *LinkedList[E]) PushFront(element E) bool {
	return l.Insert(0, element)
}

func (l *LinkedList[E]) PopBack() *E {
	//链表为空
	if l.len == 0 {
		return nil
	}
	node := l.last
	//只有一个元素
	if l.len == 1 {
		l.last = nil
		l.first = nil
	} else {
		l.last = node.prev
		l.last.next = nil
	}
	l.len--
	return &(node.element)
}

func (l *LinkedList[E]) PopFront() *E {
	if l.len == 0 {
		return nil
	}
	node := l.first
	if l.len == 1 {
		l.first = nil
		l.last = nil
	} else {
		l.first = node.next
		l.first.prev = nil
	}
	l.len--
	return &(node.element)
}

func (l *LinkedList[E]) PullBack() *E {
	if l.len == 0 {
		return nil
	}
	return &(l.last.element)
}

func (l *LinkedList[E]) PullFront() *E {
	if l.len == 0 {
		return nil
	}
	return &(l.first.element)
}

// Iterator 获取该链表的迭代器
func (l *LinkedList[E]) Iterator() Iterator[E] {
	return &LinkedListIterator[E]{
		reverse: false,
		next:    l.first,
	}
}

// ReverseIterator 获取反向迭代器
func (l *LinkedList[E]) ReverseIterator() Iterator[E] {
	return &LinkedListIterator[E]{
		reverse: true,
		next:    l.last,
	}
}

type LinkedListIterator[E any] struct {
	//是否反向,如果为true,则是从尾部向头部迭代
	reverse bool
	next    *Node[E]
}

func (l *LinkedListIterator[E]) Has() bool {
	return l.next != nil
}

func (l *LinkedListIterator[E]) Next() E {
	e := l.next
	if e == nil {
		panic("iterator is empty.")
	}
	if l.reverse {
		l.next = e.prev
	} else {
		l.next = e.next
	}
	return e.element
}
