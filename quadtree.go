// Package tree provides spatial search tree implementations
package tree

import (
	"fmt"
)

const (
	MAX_OBJECTS = 12
	MAX_LEVELS  = 5
)

const (
	TOP_LEFT = iota
	TOP_RIGHT
	BOTTOM_LEFT
	BOTTOM_RIGHT
)

type QuadTree struct {
	bounds  *Rectangle
	objects []*Rectangle
	nodes   map[int]*QuadTree
	level   int
}

// NewQuadTree creates a new quad tree at level pLevel and bounds
func NewQuadTree(pLevel int, bounds *Rectangle) *QuadTree {
	return &QuadTree{
		level:   pLevel,
		bounds:  bounds,
		objects: make([]*Rectangle, 0),
		nodes:   make(map[int]*QuadTree, 4),
	}
}

/*
 * Insert the object into the QuadTree. If the node exceeds the capacity, it
 * will split and push objects into the sub nodes
 */
func (q *QuadTree) Insert(rectangle *Rectangle) {
	// There are sub-nodes so push the object down to one of them if it
	// can fit in any of them
	if len(q.nodes) != 0 {
		index := q.index(rectangle)
		if index != -1 {
			q.nodes[index].Insert(rectangle)
			return
		}
	}
	// object didn't fit into the any of the sub-nodes, so push in into this
	// node
	q.objects = append(q.objects, rectangle)

	// We hit the objects limit and are still allowed to create more nodes,
	// split and push objects into sub-nodes
	if len(q.objects) > MAX_OBJECTS && q.level < MAX_LEVELS {
		// there are no sub nodes in this node
		if len(q.nodes) == 0 {
			q.split()
		}
		for i := len(q.objects) - 1; i >= 0; i-- {
			index := q.index(q.objects[i])
			if index != -1 {
				q.nodes[index].Insert(q.objects[i])
				q.objects = append(q.objects[:i], q.objects[i+1:]...)
			}
		}
	}
}

/**
 * Retrieve will populate the result with the triangles that intersects with
 * the rectangle
 */
func (q *QuadTree) Retrieve(rectangle *Rectangle, result *[]*Rectangle) {
	*result = append(*result, q.objects...)
	for _, node := range q.nodes {
		if node.bounds.Intersects(rectangle) {
			node.Retrieve(rectangle, result)
		}
	}
}

/**
 * Clear will remove all objects and sub-nodes from this node and sub-nodes
 * recursively
 */
func (q *QuadTree) Clear() {
	q.objects = make([]*Rectangle, 0)
	for _, node := range q.nodes {
		node.Clear()
	}
	q.nodes = make(map[int]*QuadTree, 4)
}

/**
 * Boundaries will return all "rectangles" in this tree
 */
func (q *QuadTree) Boundaries() []*Rectangle {
	boundaries := make([]*Rectangle, 0)
	boundaries = append(boundaries, q.bounds)
	for _, node := range q.nodes {
		boundaries = append(boundaries, node.Boundaries()...)
	}
	return boundaries
}

/**
 * Debug will print a debug out to the console
 */
func (q *QuadTree) Debug() {
	for _, obj := range q.objects {
		for i := 0; i < q.level; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("object x:%d y:%d w:%d h:%d\n", int(obj.position.X), int(obj.position.Y), int(obj.maxX-obj.minX), int(obj.maxY-obj.minY))
	}
	for index, node := range q.nodes {
		for i := 0; i < node.level; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("node %d.%d\n", node.level, index)
		node.Debug()
	}
}

/*
 * index determines which node the object belongs to. nil means
 * object cannot completely fit within a child node and is part
 * of the parent node
 */
func (q *QuadTree) index(r *Rectangle) int {
	// Object can completely fit within the left quadrants
	left := (r.position.X < q.bounds.position.X) && (r.maxX < q.bounds.position.X)
	// Object can completely fit within the top quadrants
	top := (r.position.Y < q.bounds.position.Y) && (r.maxY < q.bounds.position.Y)
	if top && left {
		return TOP_LEFT
	}
	// Object can completely fit within the right quadrants
	right := (r.position.X > q.bounds.position.X) && (r.minX > q.bounds.position.X)
	if top && right {
		return TOP_RIGHT
	}
	// Object can completely fit within the bottom quadrants
	bottom := (r.position.Y > q.bounds.position.Y) && (r.minY > q.bounds.position.Y)
	if bottom && left {
		return BOTTOM_LEFT
	}
	if bottom && right {
		return BOTTOM_RIGHT
	}
	// object can't fit in any of the sub nodes
	return -1
}

// split creates four sub-nodes for this node
func (q *QuadTree) split() {
	subWidth := q.bounds.maxX - q.bounds.minX
	subHeight := q.bounds.maxY - q.bounds.minY
	x := q.bounds.position.X
	y := q.bounds.position.Y
	q.nodes[TOP_LEFT] = NewQuadTree(q.level+1, NewRectangle(x-subWidth, y-subHeight, subWidth, subHeight))
	q.nodes[TOP_RIGHT] = NewQuadTree(q.level+1, NewRectangle(x+subWidth, y-subHeight, subWidth, subHeight))
	q.nodes[BOTTOM_LEFT] = NewQuadTree(q.level+1, NewRectangle(x-subWidth, y+subHeight, subWidth, subHeight))
	q.nodes[BOTTOM_RIGHT] = NewQuadTree(q.level+1, NewRectangle(x+subWidth, y+subHeight, subWidth, subHeight))
}
