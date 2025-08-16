package pokecache

import (
	"reflect"
	"slices"
	"testing"
)

func TestNewCache(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Errorf("NewCache() returned nil")
	}
	if reflect.TypeOf(c).Elem().Name() != "Cache" {
		t.Errorf("c is NOT cache type")
	}
}

func TestAdd(t *testing.T) {
	c := NewCache()
	key := "https://url1.com"
	key2 := "http://pokeapi.co/api/v2/pokemon/charizard"
	str := "{response:ABC123}"
	str2 := "{response:Charizard}"
	c.Add(key, []byte(str))
	c.Add(key2, []byte(str2))

	if string(c.cache_map[key].bytes) != str || string(c.cache_map[key2].bytes) != str2 {
		t.Errorf("Failed to add cache items")
	}
}

func TestGet(t *testing.T) {
	c := NewCache()

	cases := []struct {
		input    map[string]string
		expected []string
	}{
		{
			input: map[string]string{
				"https://url1.com": "{response:ABC123}",
				"https://pokeapi.co/api/v2/pokemon/charizard": "{response:Charizard}",
			},
			expected: []string{"{response:ABC123}", "{response:Charizard}"},
		},
	}

	for _, case_ := range cases {
		for key, val := range case_.input {
			c.Add(key, []byte(val))
		}
	}

	for _, case_ := range cases {
		for key, _ := range case_.input {
			res, ok := c.Get(key)
			if !ok {
				t.Errorf("key not found in cache.")
			}
			if !slices.Contains(case_.expected, string(res.bytes)) {
				t.Errorf("Value missing in cache")
			}
		}
	}

}
