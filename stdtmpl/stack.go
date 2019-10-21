package stdtmpl

type Stack interface {
	// Size is the number of elements in the stack
	Size() int // max 2,147,483,647

	// Cap is the actual size of the stack
	Cap() int

	// Pop takes the next element off the stack
	Pop() interface{}

	// Push pushes the element onto the stack
	Push(x interface{})

	// Peek shows you the next element to be removed (by either Pop or Push)
	Peek() interface{}
}
