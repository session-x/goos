package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Function to handle --install flag
func main() {
	installFlag := flag.Bool("install", false, "Install required tools")
	flag.Parse()

	if *installFlag {
		err := installTools()
		if err != nil {
			fmt.Println("[!] Error during installation:", err)
			return
		}
		fmt.Println("[+] Installation completed successfully!")
		return
	}

	// Placeholder for the main logic of the tool
	fmt.Println("[+] GOOS tool is running...")

	// Example: Run your main logic here (URL processing, etc.)
	// This could be your actual tool functionality (e.g., calling waybackurls, gau, etc.)
}

// Install Go-based and non-Go-based tools
func installTools() error {
	// Install Go tools using go install
	tools := []string{
		"github.com/tomnomnom/waybackurls@latest",
		"github.com/lc/gau/v2/cmd/gau@latest",
		"github.com/ffuf/ffuf@latest", // For example, if you're using qsreplace
	}

	for _, tool := range tools {
		fmt.Printf("[+] Installing Go tool: %s\n", tool)
		err := goInstall(tool)
		if err != nil {
			return fmt.Errorf("failed to install Go tool %s: %v", tool, err)
		}
	}

	// Install non-Go tools using curl
	fmt.Println("[+] Installing non-Go tools...")

	// Install freq using curl
	err := installFreq()
	if err != nil {
		return fmt.Errorf("failed to install freq: %v", err)
	}

	// Install uro using curl
	err = installUro()
	if err != nil {
		return fmt.Errorf("failed to install uro: %v", err)
	}

	// Verify installation
	err = verifyInstallation()
	if err != nil {
		return fmt.Errorf("installation verification failed: %v", err)
	}

	return nil
}

// Install a Go tool using go install
func goInstall(tool string) error {
	cmd := exec.Command("go", "install", tool)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running go install: %s\nOutput: %s", err, output)
	}
	return nil
}

// Install freq using curl
func installFreq() error {
	cmd := exec.Command("curl", "-LO", "https://github.com/s0md3v/Arjun/releases/download/v0.2/freq-linux-amd64")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error installing freq: %s\nOutput: %s", err, output)
	}
	cmd = exec.Command("chmod", "+x", "freq-linux-amd64")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error making freq executable: %s\nOutput: %s", err, output)
	}
	cmd = exec.Command("mv", "freq-linux-amd64", "/usr/local/bin/freq")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error moving freq to /usr/local/bin: %s\nOutput: %s", err, output)
	}
	return nil
}

// Install uro using curl
func installUro() error {
	cmd := exec.Command("curl", "-LO", "https://github.com/s0md3v/uro/releases/download/v0.2/uro-linux-amd64")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error installing uro: %s\nOutput: %s", err, output)
	}
	cmd = exec.Command("chmod", "+x", "uro-linux-amd64")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error making uro executable: %s\nOutput: %s", err, output)
	}
	cmd = exec.Command("mv", "uro-linux-amd64", "/usr/local/bin/uro")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error moving uro to /usr/local/bin: %s\nOutput: %s", err, output)
	}
	return nil
}

// Verify that the tools were installed successfully
func verifyInstallation() error {
	tools := []string{"waybackurls", "gau", "ffuf", "freq", "uro"}
	for _, tool := range tools {
		cmd := exec.Command(tool, "--version")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error verifying %s: %v\nOutput: %s", tool, err, output)
		}
		fmt.Printf("%s version:\n%s\n", tool, string(output))
	}
	return nil
}
