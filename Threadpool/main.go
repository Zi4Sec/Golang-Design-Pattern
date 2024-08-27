package main

import (
	"bytes"
	"fmt"
)

func main() {
	scan_worker := 4
	pool := NewPool(scan_worker) // Create a pool with 3 workers
	ruleset := [][]byte{{0x3C, 0x3E, 0x7e}, {0x22, 0x27}, {0x30, 0x12, 0x25, 0x19}, {0x27, 0x22}, {0x99}}
	payload := []byte{0x3C, 0x3E, 0x27, 0x22, 0x27}
	howmany := 6

	for i := 0; i < len(ruleset); i += howmany {
		// taskNum := i
		pool.AddTask(func() {

			if i+howmany-1 < len(ruleset) {
				scan(payload, ruleset[i:i+howmany])
			} else {
				scan(payload, ruleset[i:])
			}
		})
	}

	pool.Stop() // Wait for all tasks to finish
}

func scan(payload []byte, rules [][]byte) {
	for _, rule := range rules {
		if bytes.Contains(payload, rule) {
			fmt.Printf("One match has been found, %x\n", rule)
		}
	}

}
