// Package intcomp provides an easy to use implementation of the Intcode computer as described by Advent of Code 2019
package intcomp

import (
	"fmt"
)

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

func (c *Computer) getValue(i int) (v int) {
	if len(c.memory) > c.pointer+i {
		v = c.memory[c.pointer+i]
		if c.memory[c.pointer] > 99 {
			mode := (c.memory[c.pointer] / pow(10, 1+i)) % 10
			switch mode {
			case 1:
				v = c.pointer + i
			case 2:
				v = c.relativeBase + v
			}
		}
	}
	return
}

func pow(i, n int) (res int) {
	res = i
	for j := 1; j < n; j++ {
		res *= i
	}
	return
}

func (c *Computer) getValues() (v1, v2, v3, v4 int) {
	v1 = c.memory[c.pointer]
	if v1 > 99 {
		v1 = v1 % 100
	}
	v2 = c.getValue(1)
	v3 = c.getValue(2)
	v4 = c.getValue(3)
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
