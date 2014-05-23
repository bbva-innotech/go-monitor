package monitor_test

import (
	"github.com/bbva-innotech/go-monitor"
	"time"
)

func Example() {
	monValueA := monitor.NewField("Value_A", false)
	monValueA.Set(1111)

	monitor.NewField("Value_B", true)

	monValueC := monitor.NewField("Value_C", true)
	monValueC.Add(2222)

	monitor.Start()
	time.Sleep(time.Millisecond * 2100)
	monitor.Stop()

	// // Output:
	// // Value_A  Value_B
	// //    1111     2222
	// //    1111        0
}
