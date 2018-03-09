package dumputils

import "fmt"

// Go routine to dump summary
func DumpSummary(done chan bool, m map[string]int) {
	fmt.Println("Show Summary")
	fmt.Println("Prefix		", "NumRcvd")
	fmt.Println("------		", "---------------")
	for key, value := range m {
		fmt.Println(key, "	", value)
	}
	done <- true
}
