// Package intcomp provides an easy to use implementation of the Intcode computer as described by Advent of Code 2019
package intcomp

import "fmt"

// Computer represents an Intcode Computer
type Computer struct {
	memory       []int
	pointer      int
	relativeBase int
	input        chan int
	output       chan int
}

// NewComputer creates a new Intcode Computer
// in and out are the input / output used for intcodes 3 and 4. They are both blocking!
// Memsize represents your desired memory size for the computer.
// The computer does not gracefully handle index out of range errors!
func NewComputer(memory []int, memSize int, in, out chan int) *Computer {
	if memSize < len(memory) {
		memSize = len(memory)
	}
	c := Computer{memory: make([]int, memSize), input: in, output: out}
	copy(c.memory, memory)
	return &c
}

func (c *Computer) getValues() (v1, v2, v3, v4 int) {
	v1 = c.memory[c.pointer]
	if len(c.memory) > c.pointer+1 {
		v2 = c.memory[c.pointer+1]
	}
	if len(c.memory) > c.pointer+2 {
		v3 = c.memory[c.pointer+2]
	}
	if len(c.memory) > c.pointer+3 {
		v4 = c.memory[c.pointer+3]
	}

	if c.memory[c.pointer] > 99 {
		val := v1
		v1 = val % 100
		val /= 100
		if len(c.memory) > c.pointer+1 {
			if val%10 == 1 {
				v2 = c.pointer + 1
			} else if val%10 == 2 {
				v2 = c.relativeBase + v2
			}
		}
		val /= 10
		if len(c.memory) > c.pointer+2 {
			if val%10 == 1 {
				v3 = c.pointer + 2
			} else if val%10 == 2 {
				v3 = c.relativeBase + v3
			}
		}
		val /= 10
		if len(c.memory) > c.pointer+3 {
			if val%10 == 1 {
				v4 = c.pointer + 3
			} else if val%10 == 2 {
				v4 = c.relativeBase + v4
			}
		}
	}
	return
}

// Run starts the intcode computer, it will return if the program exits or if it encouters an error
func (c *Computer) Run() (err error) {
	var done bool
	for !done {
		done, err = c.executeOpcode()
	}
	return err
}

func (c *Computer) executeOpcode() (bool, error) {
	opcode, v1, v2, v3 := c.getValues()
	switch opcode {
	case 99:
		return true, nil
	case 1:
		c.code1(v3, v1, v2)
	case 2:
		c.code2(v3, v1, v2)
	case 3:
		c.code3(v1)
	case 4:
		c.code4(v1)
	case 5:
		c.code5(v1, v2)
	case 6:
		c.code6(v1, v2)
	case 7:
		c.code7(v3, v1, v2)
	case 8:
		c.code8(v3, v1, v2)
	case 9:
		c.code9(v1)
	default:
		return true, fmt.Errorf("unexpected opcode %d at memory position %d", opcode, c.pointer)
	}
	return false, nil
}

func (c *Computer) code1(target, a, b int) {
	c.memory[target] = c.memory[a] + c.memory[b]
	c.pointer += 4
}

func (c *Computer) code2(target, a, b int) {
	c.memory[target] = c.memory[a] * c.memory[b]
	c.pointer += 4
}

func (c *Computer) code3(target int) {
	c.memory[target] = <-c.input
	c.pointer += 2
}

func (c *Computer) code4(target int) {
	c.output <- c.memory[target]
	c.pointer += 2
}

func (c *Computer) code5(condition, target int) {
	if c.memory[condition] != 0 {
		c.pointer = c.memory[target]
	} else {
		c.pointer += 3
	}
}

func (c *Computer) code6(condition, target int) {
	if c.memory[condition] == 0 {
		c.pointer = c.memory[target]
	} else {
		c.pointer += 3
	}
}

func (c *Computer) code7(target, a, b int) {
	if c.memory[a] < c.memory[b] {
		c.memory[target] = 1
	} else {
		c.memory[target] = 0
	}
	c.pointer += 4
}

func (c *Computer) code8(target, a, b int) {
	if c.memory[a] == c.memory[b] {
		c.memory[target] = 1
	} else {
		c.memory[target] = 0
	}
	c.pointer += 4
}

func (c *Computer) code9(a int) {
	c.relativeBase += c.memory[a]
	c.pointer += 2
}
