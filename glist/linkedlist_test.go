package glist

import (
	"math/rand"
	"reflect"
	"testing"
)

func checkListElement[T comparable](list *LinkedList[T], items []T, t *testing.T) {
	if list.len != uint(len(items)) {
		t.Errorf("list len= %d, items len = %d.", list.len, len(items))
	}
	i := 0
	for it := list.Iterator(); it.Has(); {
		e := it.Next()
		if e != items[i] {
			t.Errorf("index=%d,except:%v, but got:%v", i, items[i], e)
		}
		i++
	}
}

func TestLinkedListIterator_Has(t *testing.T) {
	//空列表
	t.Run("empty list", func(t *testing.T) {
		emptyList := NewLinkedList[int]()
		it := emptyList.Iterator()
		if it.Has() {
			t.Errorf("iterator should empty, but not")
		}
	})
	//非空列表
	t.Run("non-empty list", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Append(1)
		list.Append(2)
		it := list.Iterator()
		if !(it.Has()) {
			t.Errorf("iterator should not empty, but not")
		}
		it.Next()
		it.Next()
		if it.Has() {
			t.Errorf("iterator should empty, but not")
		}
	})
}

func TestLinkedListIterator_Next(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		defer func() {
			if p := recover(); p == nil {
				t.Errorf("iterator is empty, should panic.")
			}
		}()
		emptyList := NewLinkedList[int]()
		it := emptyList.Iterator()
		it.Next()
	})

	t.Run("non-empty list", func(t *testing.T) {
		defer func() {
			if p := recover(); p != nil {
				t.Errorf("iterator is non-empty, should no panic.")
			}
		}()
		list := NewLinkedList[int]()
		times := rand.Int() % 20 //20次以内
		for i := times; i >= 0; i-- {
			list.Append(i)
		}
		it := list.Iterator()
		for i := times; i >= 0; i-- {
			it.Next()
		}
	})
}

func TestLinkedList_Append(t *testing.T) {
	var sliceWithTenNum [10]int
	for i := 0; i < 10; i++ {
		sliceWithTenNum[i] = i
	}
	tests := []struct {
		name  string
		input []int
	}{
		{"Append 0 value", []int{}},
		{"Append 10 value", sliceWithTenNum[:]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int]()
			for _, v := range tt.input {
				list.Append(v)
			}
			checkListElement(list, tt.input, t)
		})
	}
}

func TestLinkedList_Size(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		{"empty list", 0},
		{"10 values list", 10},
		{"20 values list", 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num := tt.want
			list := NewLinkedList[int]()
			for i := uint(0); i < num; i++ {
				list.Append(int(i))
			}
			if got := list.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_IsEmpty(t *testing.T) {
	nonEmpty := NewLinkedList[int]()
	nonEmpty.Append(1)

	tests := []struct {
		name string
		list *LinkedList[int]
		want bool
	}{
		{"empty list", NewLinkedList[int](), true},
		{"non-empty list", nonEmpty, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.list
			if got := l.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_Remove(t *testing.T) {
	tests := []struct {
		name    string
		listNum int
		args    uint
		want    bool
		values  []int
	}{
		{"empty list", 0, 0, false, nil},
		{"index gather than len", 0, 1, false, nil},
		{"1 value list:delete index 0", 1, 0, true, []int{}},
		{"10 value list:delete index 0", 10, 0, true, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"10 value list:delete index 5", 10, 5, true, []int{0, 1, 2, 3, 4, 6, 7, 8, 9}},
		{"10 value list:delete index 9", 10, 9, true, []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLinkedList[int]()
			for i := 0; i < tt.listNum; i++ {
				l.Append(i)
			}
			if got := l.Remove(tt.args); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
			if tt.want {
				checkListElement(l, tt.values, t)
			}
		})
	}
}

func TestLinkedList_Get(t *testing.T) {
	tests := []struct {
		name    string
		listNum int
		args    uint
		want    int
	}{
		{"empty list,index 0", 0, 0, -1},
		{"empty list,index 1", 0, 1, -1},
		{"1 value list, index 0", 1, 0, 0},
		{"10 value list, index 0", 10, 0, 0},
		{"10 value list, index 5", 10, 5, 5},
		{"10 value list, index 9", 10, 9, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLinkedList[int]()
			for i := 0; i < tt.listNum; i++ {
				l.Append(i)
			}
			want := new(int)
			if tt.want == -1 {
				want = nil
			} else {
				*want = tt.want
			}
			if got := l.Get(tt.args); !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}
}
