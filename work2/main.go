package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {
	keys := flag.String("keys", "", "keys to extract from JSON. Comma-separated. Use dot for nested keys, key.subkey")
	file := flag.String("file", "", "JSON file to read")
	help := flag.Bool("help", false, "show usage")
	flag.Parse()

	if *file == "" || *help {
		showUsage()
		return
	}

	data, err := os.ReadFile(*file)
	if err != nil {
		slog.Error("Error reading JSON file", "err", err)
		os.Exit(1)
	}

	var payload map[string]any
	err = json.Unmarshal(data, &payload)
	if err != nil {
		slog.Error("Error parsing JSON file:", "err", err)
		os.Exit(1)
	}

	if *keys == "" {
		printAll("", payload)
		return
	} else {
		requestedKeys := strings.Split(*keys, ",")
		keysPrint(payload, requestedKeys)
	}
}

func showUsage() {
	fmt.Println("JSON Parser")
	fmt.Println("Usage: go run main.go --file=order.json [--keys=key1,key2.subkey,...]")
	fmt.Println("Example: go run main.go --file=order.json --keys=id,customer.name")
}

func keysPrint(obj map[string]any, requestedKeys []string) {
	for _, key := range requestedKeys {
		nestedKeys := strings.Split(key, ".")
		value, exists := getNestedValue(obj, nestedKeys)

		if exists {
			fmt.Printf("%s: %v\n", key, value)
		}
	}
}

func getNestedValue(obj any, keys []string) (any, bool) {
	current := obj
	for _, key := range keys {
		switch v := current.(type) {
		case map[string]any:
			val, exists := v[key]
			if !exists {
				return nil, false
			}
			current = val
		default:
			return nil, false
		}
	}
	return current, true
}

func printAll(prefix string, obj any) {
	m, ok := obj.(map[string]any)
	if !ok {
		fmt.Printf("%s: %v\n", prefix, obj)
		return
	}
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		printAll(fullKey, v)
	}
}
