package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Test the hello endpoint",
	Long:  "Calls the /hello endpoint to test basic API functionality",
	RunE:  runHelloCheck,
}

func runHelloCheck(cmd *cobra.Command, args []string) error {
	url := serverURL + "/hello"
	verboseLog(fmt.Sprintf("Making request to: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Cannot connect to API server at %s\n", serverURL)
		if verbose {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("💡 Make sure the server is running with: ./do start")
		return nil // Don't exit with error for connection issues
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	verboseLog(fmt.Sprintf("Response status: %s", resp.Status))

	if format == "json" {
		fmt.Println(string(body))
		return nil
	}

	// Check if response is successful before parsing
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Hello endpoint failed (status: %s)\n", resp.Status)
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	// Pretty format - try to parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("❌ API returned invalid response (not JSON)\n")
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	fmt.Printf("👋 %v\n", result["message"])
	return nil
}
