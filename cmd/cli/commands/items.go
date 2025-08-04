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

var getItemCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single item by ID",
	Long:  "Retrieve a single item using its UUID",
	RunE:  runGetItem,
}

var updateItemCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing item",
	Long:  "Update an item's name and/or description",
	RunE:  runUpdateItem,
}

var deleteItemCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an item",
	Long:  "Delete an item by its UUID",
	RunE:  runDeleteItem,
}

var (
	itemName        string
	itemDescription string
	listLimit       int
	listOffset      int
	itemID          string
	updateName      string
	updateDesc      string
)

func init() {
	// Add flags for create command
	createItemCmd.Flags().StringVar(&itemName, "name", "", "Item name (required)")
	createItemCmd.Flags().StringVar(&itemDescription, "description", "", "Item description (optional)")
	createItemCmd.MarkFlagRequired("name")

	// Add flags for list command
	listItemsCmd.Flags().IntVar(&listLimit, "limit", 10, "Number of items to retrieve (max 100)")
	listItemsCmd.Flags().IntVar(&listOffset, "offset", 0, "Number of items to skip")

	// Add flags for get command
	getItemCmd.Flags().StringVar(&itemID, "id", "", "Item ID (required)")
	getItemCmd.MarkFlagRequired("id")

	// Add flags for update command
	updateItemCmd.Flags().StringVar(&itemID, "id", "", "Item ID (required)")
	updateItemCmd.Flags().StringVar(&updateName, "name", "", "New item name")
	updateItemCmd.Flags().StringVar(&updateDesc, "description", "", "New item description")
	updateItemCmd.MarkFlagRequired("id")

	// Add flags for delete command
	deleteItemCmd.Flags().StringVar(&itemID, "id", "", "Item ID (required)")
	deleteItemCmd.MarkFlagRequired("id")

	// Add subcommands to items
	itemsCmd.AddCommand(createItemCmd)
	itemsCmd.AddCommand(listItemsCmd)
	itemsCmd.AddCommand(getItemCmd)
	itemsCmd.AddCommand(updateItemCmd)
	itemsCmd.AddCommand(deleteItemCmd)
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

func runGetItem(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/items/%s", serverURL, itemID)
	verboseLog(fmt.Sprintf("Making GET request to: %s", url))

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Cannot connect to API server at %s\n", serverURL)
		if verbose {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("💡 Make sure the server is running with: ./do start")
		return nil
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

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("❌ Item not found (ID: %s)\n", itemID)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to get item (status: %s)\n", resp.Status)
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	// Pretty format
	var item storage.Item
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Printf("❌ API returned invalid response (not JSON)\n")
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	fmt.Printf("📄 Item Details\n")
	fmt.Printf("   ID: %s\n", item.ID)
	fmt.Printf("   Name: %s\n", item.Name)
	if item.Description != nil {
		fmt.Printf("   Description: %s\n", *item.Description)
	} else {
		fmt.Printf("   Description: (none)\n")
	}
	fmt.Printf("   Created: %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Updated: %s\n", item.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func runUpdateItem(cmd *cobra.Command, args []string) error {
	// Build update request
	reqData := make(map[string]interface{})
	if updateName != "" {
		reqData["name"] = updateName
	}
	if updateDesc != "" {
		reqData["description"] = updateDesc
	}

	// Check if at least one field is being updated
	if len(reqData) == 0 {
		fmt.Println("❌ At least one field (--name or --description) must be provided for update")
		return nil
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/items/%s", serverURL, itemID)
	verboseLog(fmt.Sprintf("Making PUT request to: %s", url))
	verboseLog(fmt.Sprintf("Request body: %s", string(jsonData)))

	// Create PUT request
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Cannot connect to API server at %s\n", serverURL)
		if verbose {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("💡 Make sure the server is running with: ./do start")
		return nil
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

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("❌ Item not found (ID: %s)\n", itemID)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to update item (status: %s)\n", resp.Status)
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	// Pretty format
	var item storage.Item
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Printf("❌ API returned invalid response (not JSON)\n")
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	fmt.Printf("✅ Item updated successfully!\n")
	fmt.Printf("   ID: %s\n", item.ID)
	fmt.Printf("   Name: %s\n", item.Name)
	if item.Description != nil {
		fmt.Printf("   Description: %s\n", *item.Description)
	}
	fmt.Printf("   Updated: %s\n", item.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func runDeleteItem(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/items/%s", serverURL, itemID)
	verboseLog(fmt.Sprintf("Making DELETE request to: %s", url))

	// Create DELETE request
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Cannot connect to API server at %s\n", serverURL)
		if verbose {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("💡 Make sure the server is running with: ./do start")
		return nil
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

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("❌ Item not found (ID: %s)\n", itemID)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to delete item (status: %s)\n", resp.Status)
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	// Pretty format
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ API returned invalid response (not JSON)\n")
		if verbose {
			fmt.Printf("Response: %s\n", string(body))
		}
		return nil
	}

	fmt.Printf("✅ Item deleted successfully!\n")
	if item, ok := response["item"].(map[string]interface{}); ok {
		if name, ok := item["name"].(string); ok {
			fmt.Printf("   Deleted: %s (ID: %s)\n", name, itemID)
		}
	}

	return nil
}
