package stack

type Stack []interface{}

// IsEmpty check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(str interface{}) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Pop Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		return ""
	} else {
		// Get the index of the top most element.
		// Index into the slice and obtain the element.
		// Remove it from the stack by slicing it off.
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element
	}
}
