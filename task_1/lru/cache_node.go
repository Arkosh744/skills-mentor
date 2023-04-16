package main

type Node struct {
	prev, next *Node
	value      *CacheItem
}

func (c *Cache) moveToFront(node *Node) {
	if node == c.head {
		return
	}

	c.removeNode(node)
	c.addNode(node)
}

func (c *Cache) removeNode(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

func (c *Cache) addNode(node *Node) {
	node.prev = nil
	node.next = c.head

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = c.head
	}
}
