package tree

import "github.com/AsynkronIT/protoactor-go/actor"

//TODO set caller in CLI/
// Actor Node --------------------------------------------

type Node struct {
	maxValues int
	inner     Inner
	leaf      Leaf
}

type Inner struct {
	left    actor.PID
	right   actor.PID
	maxLeft int
}

type Leaf struct {
	values map[int]string
}

type searchMSG struct {
	caller actor.PID
	key    int
}

func (node *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *searchMSG:
		node.search(msg, context)

	}
}

func (node *Node) search(msg *searchMSG, context actor.Context) {
	if (node.inner != Inner{}) { //IF is inner node
		if msg.key > node.inner.maxLeft { // bigger -> keep searching on the right
			context.Send(&node.inner.right, msg)
		} else {
			context.Send(&node.inner.left, msg) // smaller -> keep searching on the left
		}
	} else { // IF leaf
		elem, ok := node.leaf.values[msg.key]
		if ok {
			context.Send(msg.caller)
		}
	}
}

// Messages -----------------------------------------------

type Insert struct {
	key    int
	value  string
	caller *actor.PID
}

type Error struct {
	originalMsg interface{}
}

type Success struct {
	key         int
	value       string
	originalMsg interface{}
}

// Actions ------------------------------------------------

func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Insert:
		state.insert(msg, context)
	}
}

func (state *Node) insert(msg *Insert, context actor.Context) {
	if state.inner != (Inner{}) {
		switch {
		case msg.key > state.inner.maxLeft:
			context.Send(state.inner.right, msg)
		case msg.key < state.inner.maxLeft:
			context.Send(state.inner.left, msg)
		case msg.key == state.inner.maxLeft:
			context.Send(msg.caller, &Error{originalMsg: msg})
		}
	}
}
