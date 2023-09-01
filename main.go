package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: subenum <domain> <concurrency>")
		return
	}

	domain := os.Args[1]
	concurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid concurrency value:", err)
		return
	}

	// List of common subdomains to check
	predefinedSubdomains := []string{
		"www",
		"mail",
		"ftp",
		"admin",
		"blog",
		"api",
		"app",
		"dev",
		"stage",
		"test",
		"secure",
		"support",
		"forum",
		// Add more predefined subdomains here
	}

	// Read custom subdomains from the text file
	customSubdomains, err := readSubdomainsFromFile("2m-subdomains.txt")
	if err != nil {
		fmt.Printf("Error reading custom subdomains file: %v\n", err)
		return
	}

	// Combine predefined and custom subdomains into a single list
	subdomains := append(predefinedSubdomains, customSubdomains...)

	// Create a channel for results
	resultChannel := make(chan map[string][]string)

	// Create a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit the number of concurrent workers
	sem := make(chan struct{}, concurrency)

	// Launch result processing Goroutine
	go func() {
		for result := range resultChannel {
			// Process the result, e.g., store it, print it, etc.
			// You can access the result map with subdomain and IP addresses here.
			fmt.Println(result)
		}
	}()

	// Launch Goroutines for subdomain enumeration
	for _, subdomain := range subdomains {
		sem <- struct{}{} // Acquire semaphore
		wg.Add(1)
		go func(subdomain string) {
			defer func() {
				<-sem // Release semaphore
				wg.Done()
			}()

			target := subdomain + "." + domain
			ips, err := resolveWithTimeout(target, 2*time.Second) // Set a timeout
			if err == nil && len(ips) > 0 {
				result := map[string][]string{target: ips}
				resultChannel <- result
			}
		}(subdomain)
	}

	// Wait for all workers to finish
	wg.Wait()
	close(resultChannel)
}

func readSubdomainsFromFile(filename string) ([]string, error) {
	var subdomains []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := strings.TrimSpace(scanner.Text())
		if subdomain != "" && !strings.HasPrefix(subdomain, "#") {
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains, scanner.Err()
}

func resolveWithTimeout(domain string, timeout time.Duration) ([]string, error) {
	var ips []string

	// Resolve the domain to an IP address with a timeout
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout: timeout,
			}
			return dialer.DialContext(ctx, network, address)
		},
	}

	addresses, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		// Return the error, but allow the enumeration to continue
		return nil, err
	}

	for _, address := range addresses {
		ips = append(ips, address)
	}

	return ips, nil
}
