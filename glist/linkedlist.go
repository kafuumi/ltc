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
	Element E        //保存的内容
	Prev    *Node[E] //前一个结点
	Next    *Node[E] //后一个结点
}

// Clone 克隆Node,返回的Node的Prev和Next均为nil,Element保持不变
func (n *Node[E]) Clone() *Node[E] {
	node := &Node[E]{
		Element: n.Element,
		Prev:    nil,
		Next:    nil,
	}
	return node
}

// LinkedList 链表,实现了List
type LinkedList[E any] struct {
	Len   uint     //链表中元素个数
	First *Node[E] //头指针
	Last  *Node[E] //尾指针
}

// NewLinkedList 创建一个链表,
//列表的最大容量为uint类型的最大值
func NewLinkedList[E any]() *LinkedList[E] {
	return &LinkedList[E]{
		Len:   0,
		First: nil,
		Last:  nil,
	}
}

func (l *LinkedList[E]) Size() uint {
	return l.Len
}

func (l *LinkedList[E]) IsEmpty() bool {
	return l.Len == 0
}

func (l *LinkedList[E]) IsNotEmpty() bool {
	return l.Len != 0
}

func (l *LinkedList[E]) Append(element E) bool {
	//超出最大值无法添加
	if l.Len == CAPACITY {
		return false
	}
	node := &Node[E]{
		Element: element,
		Prev:    nil,
		Next:    nil,
	}
	//链表为空,头指针指向该结点
	if l.First == nil {
		l.First = node
		l.Last = node
	} else {
		//链表不为空,添加到尾部
		node.Prev = l.Last
		l.Last.Next = node
		l.Last = node
	}
	l.Len++
	return true
}

func (l *LinkedList[E]) Insert(index uint, element E) bool {
	//当前size已经达到最大值或者索引越界
	if l.Len == CAPACITY || index > l.Len {
		return false
	}
	node := &Node[E]{
		Element: element,
		Prev:    nil,
		Next:    nil,
	}
	//插入头部
	if index == 0 {
		if l.First == nil {
			//链表为空
			l.First = node
			l.Last = node
		} else {
			//链表不为空
			node.Next = l.First
			l.First.Prev = node
			l.First = node
		}
	} else if index == l.Len {
		//插入尾部
		l.Last.Next = node
		node.Prev = l.Last
		l.Last = node
	} else {
		var prev *Node[E]
		head := l.First
		for i := ZERO; i < index; i++ {
			prev = head
			head = head.Next
		}
		node.Next = head
		node.Prev = prev
		prev.Next = node
		head.Prev = node
	}
	l.Len++
	return true
}

func (l *LinkedList[E]) Remove(index uint) bool {
	if index >= l.Len {
		return false
	}
	head := l.First
	var prev *Node[E]
	for i := ZERO; i < index; i++ {
		prev = head
		head = head.Next
	}
	//删除第一个结点
	if head == l.First {
		l.First.Next = nil
		l.First = head.Next
	} else if head == l.Last {
		//删除最后一个结点
		l.Last = prev
		l.Last.Next = nil
	} else {
		prev.Next = head.Next
		head.Next.Prev = prev
	}
	l.Len--
	return true
}

func (l *LinkedList[E]) Get(index uint) *E {
	if index >= l.Len {
		return nil
	}
	node := l.First
	for i := ZERO; i < index; i++ {
		node = node.Next
	}
	return &(node.Element)
}

func (l *LinkedList[E]) Set(index uint, element E) bool {
	if index >= l.Len {
		return false
	}
	node := l.First
	for i := ZERO; i < index; i++ {
		node = node.Next
	}
	node.Element = element
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
	if l.Len == 0 {
		return nil
	}
	node := l.Last
	//只有一个元素
	if l.Len == 1 {
		l.Last = nil
		l.First = nil
	} else {
		l.Last = node.Prev
		l.Last.Next = nil
	}
	l.Len--
	return &(node.Element)
}

func (l *LinkedList[E]) PopFront() *E {
	if l.Len == 0 {
		return nil
	}
	node := l.First
	if l.Len == 1 {
		l.First = nil
		l.Last = nil
	} else {
		l.First = node.Next
		l.First.Prev = nil
	}
	l.Len--
	return &(node.Element)
}

func (l *LinkedList[E]) PullBack() *E {
	if l.Len == 0 {
		return nil
	}
	return &(l.Last.Element)
}

func (l *LinkedList[E]) PullFront() *E {
	if l.Len == 0 {
		return nil
	}
	return &(l.First.Element)
}

// Iterator 获取该链表的迭代器
func (l *LinkedList[E]) Iterator() Iterator[E] {
	return &LinkedListIterator[E]{
		reverse: false,
		next:    l.First,
	}
}

// ReverseIterator 获取反向迭代器
func (l *LinkedList[E]) ReverseIterator() Iterator[E] {
	return &LinkedListIterator[E]{
		reverse: true,
		next:    l.Last,
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
		l.next = e.Prev
	} else {
		l.next = e.Next
	}
	return e.Element
}
