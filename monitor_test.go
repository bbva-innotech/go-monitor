package monitor_test

import (
	"github.com/bbva-innotech/go-monitor"
	"time"
)

func Example() {
	// A value that will not reset its value after each row
	var monValueA = monitor.NewField("Value_A", false)
	monValueA.Set(1111)

	// A value that will reset its value after each row
	var monValueB = monitor.NewField("Value_B", true)
	monValueB.Add(2222)

	monitor.Start()
	time.Sleep(time.Millisecond * 2100)
	monitor.Stop()

	// The output must be:
	// Value_A Value_B
	//    1111    2222
	//    1111       0
}
