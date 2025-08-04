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

var listItemsCmd = &cobra.Command{
	Use:   "list",
	Short: "List items",
	Long:  "List items with optional pagination",
	RunE:  runListItems,
}

var (
	itemName        string
	itemDescription string
	listLimit       int
	listOffset      int
)

func init() {
	// Add flags for create command
	createItemCmd.Flags().StringVar(&itemName, "name", "", "Item name (required)")
	createItemCmd.Flags().StringVar(&itemDescription, "description", "", "Item description (optional)")
	createItemCmd.MarkFlagRequired("name")

	// Add flags for list command
	listItemsCmd.Flags().IntVar(&listLimit, "limit", 10, "Number of items to retrieve (max 100)")
	listItemsCmd.Flags().IntVar(&listOffset, "offset", 0, "Number of items to skip")

	// Add subcommands to items
	itemsCmd.AddCommand(createItemCmd)
	itemsCmd.AddCommand(listItemsCmd)
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

func runListItems(cmd *cobra.Command, args []string) error {
	// Build URL with query parameters
	url := fmt.Sprintf("%s/items?limit=%d&offset=%d", serverURL, listLimit, listOffset)
	verboseLog(fmt.Sprintf("Making GET request to: %s", url))

	// Make HTTP request
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
		fmt.Printf("❌ Failed to list items (status: %s)\n", resp.Status)
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	// Pretty format - try to parse JSON
	var response storage.ListItemsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ API returned invalid response (not JSON)\n")
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	// Display results
	if len(response.Items) == 0 {
		fmt.Println("📭 No items found")
		fmt.Printf("   Total: %d items\n", response.Total)
		return nil
	}

	fmt.Printf("📋 Found %d items (showing %d-%d of %d total)\n",
		len(response.Items),
		response.Offset+1,
		response.Offset+len(response.Items),
		response.Total)
	fmt.Println()

	for i, item := range response.Items {
		fmt.Printf("%d. %s\n", i+1, item.Name)
		fmt.Printf("   ID: %s\n", item.ID)
		if item.Description != nil {
			fmt.Printf("   Description: %s\n", *item.Description)
		}
		fmt.Printf("   Created: %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
		if i < len(response.Items)-1 {
			fmt.Println()
		}
	}

	// Show pagination info
	if response.Total > response.Limit {
		fmt.Println()
		nextOffset := response.Offset + response.Limit
		if nextOffset < response.Total {
			fmt.Printf("💡 To see more items, use: --offset %d\n", nextOffset)
		}
	}

	return nil
}
