package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joel-thompson/my-go-service/storage"
	"github.com/spf13/cobra"
)

var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Manage items",
	Long:  "Commands for creating and managing items via the API",
}

var createItemCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new item",
	Long:  "Create a new item by providing name and optional description",
	RunE:  runCreateItem,
}

var (
	itemName        string
	itemDescription string
)

func init() {
	// Add flags for create command
	createItemCmd.Flags().StringVar(&itemName, "name", "", "Item name (required)")
	createItemCmd.Flags().StringVar(&itemDescription, "description", "", "Item description (optional)")
	createItemCmd.MarkFlagRequired("name")

	// Add subcommands to items
	itemsCmd.AddCommand(createItemCmd)
}

func runCreateItem(cmd *cobra.Command, args []string) error {
	// Prepare request payload
	req := storage.CreateItemRequest{
		Name: itemName,
	}
	if itemDescription != "" {
		req.Description = &itemDescription
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := serverURL + "/items"
	verboseLog(fmt.Sprintf("Making POST request to: %s", url))
	verboseLog(fmt.Sprintf("Request body: %s", string(jsonData)))

	// Make HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
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

	// Pretty format
	if resp.StatusCode == http.StatusCreated {
		var item storage.Item
		if err := json.Unmarshal(body, &item); err != nil {
			fmt.Printf("❌ API returned invalid response (not JSON)\n")
			if verbose {
				fmt.Printf("Response: %s\n", string(body))
			}
			return nil
		}

		fmt.Printf("✅ Item created successfully!\n")
		fmt.Printf("   ID: %s\n", item.ID)
		fmt.Printf("   Name: %s\n", item.Name)
		if item.Description != nil {
			fmt.Printf("   Description: %s\n", *item.Description)
		}
		fmt.Printf("   Created: %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("❌ Failed to create item (status: %s)\n", resp.Status)
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
	}

	return nil
}
