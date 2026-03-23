package tests

import "testing"

type Item struct {
	id   int32
	name string
}

type ItemWrapper struct {
	*Item
}

func TestInline(t *testing.T) {
	item := &Item{
		id:   1,
		name: "item1",
	}
	itemWrapper := &ItemWrapper{
		Item: item,
	}

	t.Log("1. item:", item.name)
	t.Log("1. ItemWrapper:", itemWrapper.name)

	itemWrapper.name = "item2"

	t.Log("2. item:", item.name)
	t.Log("2. ItemWrapper:", itemWrapper.name)
}
