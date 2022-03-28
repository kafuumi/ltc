package glist

// List 一个基本的列表
type List[E any] interface {
	// Size 获取列表中数据个数
	Size() uint
	//IsEmpty 判断列表是否为空，如果为空，返回true，否则返回false
	IsEmpty() bool
	// IsNotEmpty 判断列表是否非空,如果列表不为空,返回true,否则返回false
	IsNotEmpty() bool
	// Append 向列表尾部添加一个元素
	Append(element E) bool
	// Insert 向列表指定索引处插入一个元素,如果插入成功返回true,否则返回false
	Insert(index uint, element E) bool
	// Remove 从列表中移除元素element,如果元素不存在,则返回false
	Remove(index uint) bool
	// Get 从列表中获取索引为index元素的指针,索引从0开始,如果索引超出范围则返回nil
	Get(index uint) *E
	// Set 改变列表中索引为index的元素的值,如果索引超出范围则返回false
	Set(index uint, element E) bool
	// Iterator 获取列表的迭代器
	Iterator() Iterator[E]
}

// Queue 队列
type Queue[E any] interface {
	List[E]
	// PushBack 队列尾部添加元素,添加成功返回true
	PushBack(element E) bool
	// PushFront 队列头部添加元素,添加成功返回true
	PushFront(element E) bool
	// PopBack 删除队列尾部的元素,返回被删除的元素的指针,如果队列为空,则返回nil
	PopBack() *E
	// PopFront 删除队列头部的元素,返回被删除元素的指针,如果队列为空,返回nil
	PopFront() *E
	// PullBack 获取队列尾部的元素的指针,不会删除,如果队列为空,返回nil
	PullBack() *E
	// PullFront 获取队列头部的元素的指针,不会删除,如果队列为空,返回nil
	PullFront() *E
}

// Iterator 列表迭代器
type Iterator[E any] interface {
	// Has 是否还有元素
	Has() bool
	// Next 获取元素
	Next() E
}
