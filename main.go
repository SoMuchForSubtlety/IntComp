package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/SoMuchForSubtlety/intcomp/intcomp"
)

var directInput = flag.String("i", "", "directly input an intcode program")
var fileInput = flag.String("f", "", "execute intcode from a file")
var memorySize = flag.Int("m", -1, "the memory size for the computer")

func main() {
	flag.Parse()
	var inputText []string
	if *directInput != "" {
		inputText = strings.Split(*directInput, ",")
	} else if *fileInput != "" {
		data, err := ioutil.ReadFile(*fileInput)
		if err != nil {
			fmt.Print(err)
			return
		}
		inputText = strings.Split(string(data), ",")
	} else {
		fmt.Println("please provide intcode to run")
		return
	}

	initialMemory := make([]int, len(inputText))
	for i, value := range inputText {
		intValue, _ := strconv.Atoi(strings.TrimSpace(value))
		initialMemory[i] = intValue
	}

	out := make(chan int)
	in := make(chan int)

	go func() {
		for {
			fmt.Println(<-out)
		}
	}()

	err := intcomp.NewComputer(initialMemory, *memorySize, in, out).Run()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(10 * time.Millisecond)
}
