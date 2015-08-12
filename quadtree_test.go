package tree

import (
	"math"
	"math/rand"
	"testing"
)

func randomRectangles(n int, world *Rectangle, avgSize float64) []*Rectangle {
	ret := make([]*Rectangle, n)
	for i := 0; i < len(ret); i++ {
		w := rand.NormFloat64() * avgSize
		h := rand.NormFloat64() * avgSize
		x := rand.Float64() * world.maxX
		y := rand.Float64() * world.maxY
		ret[i] = NewRectangle(x, math.Min(world.maxX, x+w), y, math.Min(world.maxY, y+h))
	}
	return ret
}

func TestClear(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	object := NewRectangle(50, 50, 10, 10)
	qt.Insert(object)
	if len(qt.objects) != 1 {
		t.Errorf("Expected 1 object to be returned")
	}
	qt.Clear()
	if len(qt.objects) != 0 {
		t.Errorf("Expected 0 object to be returned")
	}
}

func TestClearRecursive(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	for i := 0; i <= MAX_OBJECTS; i++ {
		qt.Insert(NewRectangle(25, 25, 1, 1))
	}
	qt.Clear()
	boundaries := qt.Boundaries()
	if len(boundaries) != 1 {
		t.Error("Expected nr of boundaries to be 1 after Clear(), not %d", len(boundaries))
	}
}

func TestBoundariesOne(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	qt.Insert(NewRectangle(25, 25, 1, 1))
	boundaries := qt.Boundaries()
	if len(boundaries) != 1 {
		t.Error("Expected # of boundaries to be 1, not %d", len(boundaries))
	}
}

func TestBoundariesMultiple(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	for i := 0; i < MAX_OBJECTS+1; i++ {
		qt.Insert(NewRectangle(float64(i*4), float64(i*4), 3, 3))
	}
	boundaries := qt.Boundaries()
	if len(boundaries) != 5 {
		t.Error("Expected nr of boundaries to be 5, not %d", len(boundaries))
	}

	// this will get the rect sitting in the root node
	objects := make([]*Rectangle, 0)
	qt.Retrieve(NewRectangle(25, 25, 25, 25), &objects)
	if len(objects) != MAX_OBJECTS+1 {
		t.Errorf("Expected to get %d objects, got %d", MAX_OBJECTS+1, len(objects))
	}

	objects = make([]*Rectangle, 0)
	qt.Retrieve(NewRectangle(50, 50, 25, 25), &objects)
	if len(objects) != MAX_OBJECTS+1 {
		t.Errorf("Expected to get %d objects, got %d", MAX_OBJECTS+1, len(objects))
	}

	// this will get the rect sitting in the root node
	objects = make([]*Rectangle, 0)
	qt.Retrieve(NewRectangle(75, 75, 1, 1), &objects)
	if len(objects) != 1 {
		t.Errorf("Expected to get 1 objects, got %d", len(objects))
	}
}

func TestFillAllSlots(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	for i := 0; i <= MAX_OBJECTS/4; i++ {
		qt.Insert(NewRectangle(25, 25, 1, 1))
		qt.Insert(NewRectangle(75, 25, 1, 1))
		qt.Insert(NewRectangle(25, 75, 1, 1))
		qt.Insert(NewRectangle(75, 75, 1, 1))
	}
	boundaries := qt.Boundaries()
	if len(boundaries) != 5 {
		t.Error("Expected nr of boundaries to be 5, not %d", len(boundaries))
	}
	qt.Insert(NewRectangle(0, 0, 50, 50))
}

func TestCanRetrieve(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	qt.Insert(NewRectangle(25, 25, 1, 1))

	objs := make([]*Rectangle, 0)
	qt.Retrieve(NewRectangle(50, 50, 1, 1), &objs)
	if len(objs) != 1 {
		t.Errorf("q.Retrieve should have return 1 object, got %d", len(objs))
	}
}

func TestCantRetrieve(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	for i := 0; i < MAX_OBJECTS+1; i++ {
		qt.Insert(NewRectangle(25, 25, 1, 1))
	}
	objects := make([]*Rectangle, 0)
	qt.Retrieve(NewRectangle(75, 75, 1, 1), &objects)
	if len(objects) != 0 {
		t.Errorf("q.Retrieve should have return 0 objects, got %d", len(objects))
	}
}

func TestRetrieve(t *testing.T) {
	qt := NewQuadTree(0, NewRectangle(50, 50, 50, 50))
	searchRect := NewRectangle(100, 100, 10, 10)
	qt.Insert(searchRect)

	objects := make([]*Rectangle, 0)
	qt.Retrieve(searchRect, &objects)

	if len(objects) != 1 {
		t.Errorf("Expected 1 object to be returned")
	}
}

var indexResult int

func benchIndexRand(qt *QuadTree, rects []*Rectangle, b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, rect := range rects {
			indexResult = qt.index(rect)
		}
	}
}

func BenchmarkIndexRand1(b *testing.B) {
	qt := NewQuadTree(0, NewRectangle(500, 500, 500, 500))
	rand.Seed(1)
	rects := randomRectangles(1, qt.bounds, 2)
	benchIndexRand(qt, rects, b)
}

func BenchmarkIndexRand10(b *testing.B) {
	qt := NewQuadTree(0, NewRectangle(500, 500, 500, 500))
	rand.Seed(1)
	rects := randomRectangles(10, qt.bounds, 2)
	benchIndexRand(qt, rects, b)
}

func BenchmarkIndexRand100(b *testing.B) {
	qt := NewQuadTree(0, NewRectangle(500, 500, 500, 500))
	rand.Seed(1)
	rects := randomRectangles(100, qt.bounds, 2)
	benchIndexRand(qt, rects, b)
}

func benchInsertRand(i int, b *testing.B) {
	b.StopTimer()
	qt := NewQuadTree(0, NewRectangle(500, 500, 500, 500))
	rand.Seed(1)
	rects := randomRectangles(i, qt.bounds, 2)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		for _, rect := range rects {
			qt.Insert(rect)
			qt.Clear()
		}
	}
}

func BenchmarkInsertRand1(b *testing.B)  { benchInsertRand(1, b) }
func BenchmarkInsertRand2(b *testing.B)  { benchInsertRand(2, b) }
func BenchmarkInsertRand4(b *testing.B)  { benchInsertRand(4, b) }
func BenchmarkInsertRand8(b *testing.B)  { benchInsertRand(8, b) }
func BenchmarkInsertRand16(b *testing.B) { benchInsertRand(16, b) }
func BenchmarkInsertRand32(b *testing.B) { benchInsertRand(32, b) }

var result float64

func BenchmarkDirectAccess100(b *testing.B) {
	b.StopTimer()
	rects := randomRectangles(100, NewRectangle(50, 50, 50, 50), 10)
	b.StartTimer()
	var t float64
	for n := 0; n < b.N; n++ {
		for _, rect := range rects {
			t = rect.position.X
		}
	}
	result = t
}
