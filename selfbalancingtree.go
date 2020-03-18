///
/// Author Garr Godfrey
///
/// A self-balancing binary tree in GO. Guarantee a O(log N) lookup with an O(log^2 N) insert.
///
package main

import (
	"fmt"
	"strings"
)

type Compare interface {
	Compare(b Compare) int
}

type Node struct {
	value  interface{}
	left	* Node
	right * Node
	count   int
}

type tree struct {
	head	* Node
}

///
/// Allow strings, int or any type that supports the Compare interface.
/// Can easily add floats or other primitive types to this.
///
func CompareAny(a,b interface{}) int {
	v,ok := a.(Compare)
	if (ok) {
		return v.Compare(b.(Compare))
	}
	v2,ok2 := a.(string)
	if (ok2) {
		return strings.Compare(v2,b.(string))
	}
	v3,ok3 := a.(int)
	if (ok3) {
		return v3 - b.(int)
	}
	
	return 0
}


///
/// Find the smallest item in this subtree so we can move it
///
func (t * Node) leftmost() *Node {
	if t == nil {
		return nil
	}
	
	if (t.left == nil) {
		// current node is the smallest value, remove and replace with
		// the right hand side of the tree
		return t
	}
	
	leftmost := t.left.leftmost()
	t.left = leftmost.right
	t.count -= 1
	leftmost.right = t
	
	
	// now right hand side may be heavier as we've lowered the count 
	// on the left. Rebalance as we are rebalancing. In order to save
	// checking nil for left, we use the count we already have on t
	if t.right != nil && t.right.count > t.count-t.right.count {
		newcenter := t.right.leftmost()
		newcenter.left = t.left.add(t.value)
		newcenter.count += newcenter.left.count
		leftmost.right = newcenter
	}
	
	leftmost.count = t.count + 1
	return leftmost
}


///
/// find the right most (greatest) node in this subtree.
///
func (t * Node) rightmost() *Node {
	if t == nil {
		return nil
	}
	
	if (t.right== nil) {
		// current node is the biggest value, remove and replace with
		// the left hand side of the tree
		return t
	}
	
	rightmost := t.right.rightmost()
	t.count -= 1;
	t.right = rightmost.left
	rightmost.left = t
	
	// now left hand side may be heavier as we've lowered the count 
	// on the right. Rebalance as we are rebalancing. In order to save
	// checking nil for right, we use the count we already have on t
	if t.left != nil && t.left.count > t.count-t.left.count {
		newcenter := t.left.rightmost()
		newcenter.right = t.right.add(t.value)
		newcenter.count += newcenter.right.count
		rightmost.left = newcenter
	}
	
	// since rightmost never has a right left, and the left leg is t, its
	// count should always be t.count+1
	rightmost.count = t.count+1
	return rightmost
}

/// insert and return the new center node
func (t * Node) add( d interface{}) *Node {
	if (t == nil) {
		return &Node{value: d, count: 1}
	}
	
	var rcount, lcount int
	
	if CompareAny(t.value, d) < 0 {
		if t.right == nil {
			t.right = &Node{value: d, count: 1}
		} else {
			t.right = t.right.add(d)
		}
	} else if t.left == nil {
		t.left = &Node{value: d, count: 1}
	} else {
		t.left = t.left.add(d)
	}
	t.count += 1
	//
	// now rebalance
	if rcount=0; t.right != nil {
		rcount = t.right.count
	}
	if lcount=0; t.left != nil {
		lcount = t.left.count
	}
	
	if rcount - lcount >= 2 {
		// adjust tree, making our right node the new
		// center node and put our center node down the
		// left hand side
		newcenter := t.right.leftmost()
		newcenter.left = t.left.add(t.value)
		newcenter.count += newcenter.left.count
		
		return newcenter
	}
	if (lcount - rcount >= 2) {
		newcenter := t.left.rightmost()
		newcenter.right = t.right.add(t.value)
		newcenter.count += newcenter.right.count
		
		return newcenter
	}

	
	return t
}

func (t * Node) walk(c chan interface{}) {
	if (t == nil) {
		return
	}
	
	t.left.walk(c)
	c <- t.value
	t.right.walk(c)
}

///
/// send all elements through c in sorted order.
///
func (t * tree) walk(c chan interface{}) {
	t.head.walk(c)
	close(c)
}

///
/// add the new element to our tree, potentially changing our head node.
///
func (t * tree) add (d interface{}) {
	t.head = t.head.add(d)
}



///
/// Our main funciton tests the functionality of our self balancing tree
///
func main() {
	t := &tree{}
	
	t.add("a")
	t.add("ab")
	t.add("ac")
	t.add("ae")
	t.add("af")
	t.add("f")
	t.add("e")
	t.add("e9")
	t.add("e8")
	t.add("e7")
	t.add("e3")
	t.add("e2")
	t.add("d")
	t.add("b")
	t.add("c")
	t.add("c1")
	t.add("c2")
	t.add("c3")
	t.add("c4")
	t.add("c5")
	t.add("c6")
	t.add("b1")
	t.add("v6")
	t.add("v5")
	t.add("bx")
	t.add("qr")
	t.add("v4")
	t.add("v3")
	t.add("ba")
	t.add("v2")
	t.add("cx")
	
	c := make(chan interface{}, 3)
	
	go t.walk(c)
	
	for msg := range c {
        	fmt.Println(msg)
	}
}
