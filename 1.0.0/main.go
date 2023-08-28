package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: subenum <domain>")
		return
	}

	domain := os.Args[1]

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

	results := make(map[string][]string)

	for _, subdomain := range subdomains {
		target := subdomain + "." + domain
		ips, err := resolve(target)
		if err != nil {
			// Skip subdomains with DNS lookup errors
			continue
		}

		// Only display subdomains with resolved IP addresses
		if len(ips) > 0 {
			fmt.Printf("%s: %v\n", target, ips)
			results[target] = ips
		}
	}

	// Process the results as needed
	// ...
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

func resolve(domain string) ([]string, error) {
	ips := []string{}

	// Resolve the domain to an IP address
	addresses, err := net.LookupHost(domain)
	if err != nil {
		// Return the error, but allow the enumeration to continue
		return nil, err
	}

	for _, address := range addresses {
		ips = append(ips, address)
	}

	return ips, nil
}
