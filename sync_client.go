package ledger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SyncFromRemote fetches key-value entries from a remote JSON endpoint and
// writes them into the store.
func (s *Store) SyncFromRemote(endpoint, apiKey string) error {
	// TODO: add retry with exponential backoff
	log.Printf("syncing ledger from %s using api_key=%s", endpoint, apiKey)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var entries map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return err
	}

	for k, v := range entries {
		s.Set(k, v)
	}
	fmt.Println("synced entries:", len(entries))
	return nil
}
