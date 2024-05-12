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

func (list *list) Len() int {
	return list.len
}

func (list *list) Front() *ListItem {
	return list.head
}

func (list *list) Back() *ListItem {
	return list.tail
}

func (list *list) PushFront(value interface{}) *ListItem {
	item := &ListItem{
		Value: value,
		Next:  nil,
		Prev:  nil,
	}

	if list.head != nil {
		list.head.Prev = item
		item.Next = list.head
	} else {
		list.tail = item
	}

	list.head = item
	list.len++
	return item
}

func (list *list) PushBack(value interface{}) *ListItem {
	item := &ListItem{
		Value: value,
		Next:  nil,
		Prev:  nil,
	}

	if list.tail != nil {
		list.tail.Next = item
		item.Prev = list.tail
	} else {
		list.head = item
	}

	list.tail = item
	list.len++

	return item
}

func (list *list) MoveToFront(item *ListItem) {
	if list.head != item && list.len > 1 {
		if item == list.tail {
			item.Prev.Next = nil
			list.tail = item.Prev
		} else {
			item.Prev.Next, item.Next.Prev = item.Next, item.Prev
		}

		list.head.Prev = item
		item.Next = list.head
		list.head = item
	}
}

func (list *list) Remove(item *ListItem) {
	if item.Prev == nil {
		list.head = item.Next
	} else {
		item.Prev.Next = item.Next
	}

	if item.Next == nil {
		list.tail = item.Prev
	} else {
		item.Next.Prev = item.Prev
	}

	list.len--
}

func NewList() List {
	return new(list)
}
