package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handleStats)
	go http.ListenAndServe(":1007", nil)
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func handleStats(rw http.ResponseWriter, req *http.Request) {
	items := []interface{}{}
	items = append(items, map[string]interface{}{})

	response := Response{}
	for _, field := range fields {
		response[field.title] = field.value
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, response)
}
