package gemini

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Cache はプロンプトの SHA256 ハッシュをキーにレスポンスをファイルキャッシュする。
// GEMINI_CACHE_ENABLED=false の場合はすべての操作が no-op になる。
type Cache struct {
	dir     string
	enabled bool
}

func NewCache(dir string, enabled bool) *Cache {
	if enabled {
		_ = os.MkdirAll(dir, 0755)
	}
	return &Cache{dir: dir, enabled: enabled}
}

func (c *Cache) Get(prompt string) (string, bool) {
	if !c.enabled {
		return "", false
	}
	data, err := os.ReadFile(c.cachePath(prompt))
	if err != nil {
		return "", false
	}
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}
	return entry.Response, true
}

func (c *Cache) Set(prompt, response string) {
	if !c.enabled {
		return
	}
	data, _ := json.Marshal(cacheEntry{Response: response})
	_ = os.WriteFile(c.cachePath(prompt), data, 0644)
}

func (c *Cache) cachePath(prompt string) string {
	h := sha256.Sum256([]byte(prompt))
	return filepath.Join(c.dir, fmt.Sprintf("%x.json", h))
}

type cacheEntry struct {
	Response string `json:"response"`
}
