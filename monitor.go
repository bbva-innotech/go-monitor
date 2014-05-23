// Package monitor allows to print counters information into stdout, printing a
// row each second with the status of all Field objects and a title row each ten
// seconds with titles given to each Field objects.
//
// Fields column order is set by creation order.
//
// Fields can be added in realtime after a Start() signal is given. In this case
// a new row with titles will be printed
//
// Field titles must be longest than longest value that will be displayed as
// titles they are not adapted yet
package monitor

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"strconv"
	"time"
)

// Field struct gives control over information printed for this field on each
// row
type Field struct {
	title   string
	value   int
	values  []int
	accu    int
	average int
	reset   bool
	wait    chan int
}

// Add adds given value to Field value
func (f *Field) Add(number int) {
	f.wait <- 1
	f.value += number
	<-f.wait
}

// Set sets Field value to given number
func (f *Field) Set(number int) {
	f.wait <- 1
	f.value = number
	<-f.wait
}

var fields = []*Field{}
var previousFieldCount int
var stop = make(chan int)

// Start gives the order to start printing information on stdout
func Start() {
	go func() {
		tick := time.NewTicker(time.Second).C
		count := 0
		for {
			select {
			case <-tick:
				if count%10 == 0 || previousFieldCount != len(fields) {
					printTitles()
				}
				previousFieldCount = len(fields)

				printValues()
				count++
				if count%10 == 0 || previousFieldCount != len(fields) {
					printSubtotal()
					printAverage()
					fieldReset()
				}

			case <-stop:
				break
			}
		}
	}()
	go Listen()
}

// Stop gives the order to stop printing information on stdout
func Stop() {
	stop <- 1
}

func printTitles() {
	fmt.Println()
	ct.ChangeColor(ct.Black, false, ct.Cyan, true)
	first := true
	for _, field := range fields {
		if !first {
			fmt.Printf("  ")
		}
		fmt.Printf("%v", field.title)
		first = false
	}
	ct.ResetColor()
	fmt.Println()
}

func printValues() {
	first := true
	for _, field := range fields {
		field.wait <- 1
		field.values = append(field.values, field.value)
		if len(field.values) > 10 {
			field.values = field.values[1:]
		}

		s := strconv.Itoa(field.value)
		if !first {
			fmt.Printf("  ")
		}

		field.accu += field.value

		total := 0
		for _, v := range field.values {
			total += v
		}
		field.average = total / len(field.values)

		if field.average == 0 {
			ct.ChangeColor(ct.Yellow, false, ct.None, false)
		} else {
			if field.value >= field.average {
				ct.ChangeColor(ct.Green, false, ct.None, false)
			} else {
				ct.ChangeColor(ct.Red, false, ct.None, false)
			}
		}

		fmt.Printf("%[2]*[1]v", s, len(field.title))
		ct.ResetColor()
		if field.reset {
			field.value = 0
		}

		first = false
		<-field.wait
	}
	fmt.Println()
}

func printSubtotal() {
	first := true
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	for _, field := range fields {
		field.wait <- 1
		s := strconv.Itoa(field.accu)
		if !first {
			fmt.Printf(" ")
		}
		fmt.Printf("%[2]*[1]v+", s, len(field.title))
		first = false
		<-field.wait
	}
	ct.ResetColor()
	fmt.Println()
}

func printAverage() {
	first := true
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	for _, field := range fields {
		field.wait <- 1

		if !first {
			fmt.Printf(" ")
		}
		fmt.Printf("%[2]*[1]v~", field.average, len(field.title))

		first = false
		<-field.wait
	}
	ct.ResetColor()
	fmt.Println()
}

func fieldReset() {
	for _, field := range fields {
		if field.reset {
			field.accu = 0
			field.value = 0
		}
	}
}

// NewField creates a new monitor field that is shown as a column each second
//
// - title is the title printed on titles row each ten seconds
//
// - reset tells to reset field value after each time its printed
func NewField(title string, reset bool) *Field {
	field := new(Field)
	field.title = title
	field.reset = reset
	field.wait = make(chan int, 1)
	field.values = []int{}
	fields = append(fields, field)
	return field
}
