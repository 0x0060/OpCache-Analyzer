package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Script struct {
	FullPath  string `json:"full_path"`
	Timestamp int64  `json:"timestamp"`
	Memory    struct {
		UsedMemory    int64 `json:"used_memory"`
		FreeMemory    int64 `json:"free_memory"`
		WastedMemory  int64 `json:"wasted_memory"`
	} `json:"memory"`
	Opcode struct {
		Count int `json:"opcodes_count"`
	} `json:"opcodes"`
}

type OpcacheStatus struct {
	Scripts map[string]Script `json:"scripts"`
}

func fetchOpcacheData(url string) (*OpcacheStatus, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch opcache data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var opcacheStatus OpcacheStatus
	if err := json.Unmarshal(body, &opcacheStatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &opcacheStatus, nil
}

func analyzeOpcache(opcacheStatus *OpcacheStatus) {
	for _, script := range opcacheStatus.Scripts {
		fmt.Printf("Script: %s\n", script.FullPath)
		fmt.Printf("  Timestamp: %d\n", script.Timestamp)
		fmt.Printf("  Memory:\n")
		fmt.Printf("    Used: %d bytes\n", script.Memory.UsedMemory)
		fmt.Printf("    Free: %d bytes\n", script.Memory.FreeMemory)
		fmt.Printf("    Wasted: %d bytes\n", script.Memory.WastedMemory)
		fmt.Printf("  Opcodes Count: %d\n", script.Opcode.Count)
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: analyze_opcache <opcache_data_url>")
		os.Exit(1)
	}

	url := os.Args[1]
	opcacheStatus, err := fetchOpcacheData(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	analyzeOpcache(opcacheStatus)
}
