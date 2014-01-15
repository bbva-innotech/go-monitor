package monitor

import (
	"fmt"
	"strconv"
	"time"
)

type Field struct {
	title string
	value int
	reset bool
	wait  chan int
}

func (f *Field) Add(n int) {
	f.wait <- 1
	f.value += n
	<-f.wait
}

func (f *Field) Set(n int) {
	f.wait <- 1
	f.value = n
	<-f.wait
}

var fields []*Field
var previousFieldCount int

func init() {
	fields = []*Field{}
}

func Start() {
	go func() {
		tick := time.NewTicker(time.Second).C
		count := 0
		for {
			<-tick
			if count%10 == 0 || previousFieldCount != len(fields) {
				printTitles()
			}
			previousFieldCount = len(fields)

			printValues()
			count++
		}
	}()
}

func printTitles() {
	for _, field := range fields {
		fmt.Printf("%v ", field.title)
	}
	fmt.Println()
}

func printValues() {
	for _, field := range fields {
		field.wait <- 1
		s := strconv.Itoa(field.value)
		fmt.Printf("%[2]*[1]v ", s, len(field.title))
		if field.reset {
			field.value = 0
		}
		<-field.wait
	}
	fmt.Println()
}

// creates a new monitor field that is shown as a column each second
func NewField(title string, reset bool) *Field {
	field := new(Field)
	field.title = title
	field.reset = reset
	field.wait = make(chan int, 1)
	fields = append(fields, field)
	return field
}
