package pokeapi

import (
	"io"
	"net/http"
	"time"

	"github.com/crisp-coder/go-pokedex/internal/pokecache"
)

type PokeClient struct {
	client http.Client
	cache  *pokecache.PokeCache
}

func NewPokeClient(cache_interval time.Duration) *PokeClient {
	p := &PokeClient{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		cache: pokecache.NewPokeCache(cache_interval),
	}
	return p
}

func (p *PokeClient) Get(url string) ([]byte, error) {
	if val, ok := p.cache.Get(url); ok {
		return val, nil
	}

	resp, err := p.client.Get(url)
	if err != nil {
		return []byte{}, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
