package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.head != nil {
		l.head.Prev = item
		item.Next = l.head
	} else {
		l.tail = item
	}

	l.head = item
	l.len++
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.tail != nil {
		l.tail.Next = item
		item.Prev = l.tail
	} else {
		l.head = item
	}

	l.tail = item
	l.len++

	return l.head
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head != i && l.len > 1 {
		if i != l.tail {
			i.Prev.Next, i.Next.Prev = i.Next, i.Prev
		} else {
			i.Prev.Next = nil
		}

		l.head.Prev = i
		i.Next = l.head
		l.head = i
	}
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func NewList() List {
	return new(list)
}
