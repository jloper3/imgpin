package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Digest    string    `json:"digest"`
	FetchedAt time.Time `json:"fetched_at"`
}
type Cache struct {
	Path string
	Data map[string]Entry
}

func Load(p string) (*Cache, error) {
	c := &Cache{Path: p, Data: map[string]Entry{}}
	b, err := os.ReadFile(p)
	if err == nil {
		json.Unmarshal(b, &c.Data)
	}
	return c, nil
}
func (c *Cache) Save() error {
	os.MkdirAll(filepath.Dir(c.Path), 0755)
	b, _ := json.MarshalIndent(c.Data, "", "  ")
	return os.WriteFile(c.Path, b, 0644)
}
func (c *Cache) Get(r string) (Entry, bool) { e, ok := c.Data[r]; return e, ok }
func (c *Cache) Set(r, d string)            { c.Data[r] = Entry{Digest: d, FetchedAt: time.Now()} }
