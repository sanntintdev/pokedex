package api

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_BasicOperations(t *testing.T) {
	cache := NewCache(1 * time.Minute)

	key := "pokemon-25"
	value := []byte(`{"name": "pikachu", "id": 25}`)

	cache.Add(key, value)
	fmt.Printf("Added to cache: %s\n", key)

	cachedData, found := cache.Get(key)

	if found == false {
		t.Errorf("Expected to find the key in cache, but it was not found")
	}

	fmt.Printf("Retrieved from cache: %s\n", string(cachedData))
	fmt.Println("Basic operations test PASSED!")
}

func TestCache_Expire(t *testing.T) {
	cache := NewCache(200 * time.Millisecond)

	key := "pokemon-25"
	value := []byte(`{"name": "pikachu", "id": 25}`)

	cache.Add(key, value)

	time.Sleep(400 * time.Millisecond)

	_, found := cache.Get(key)

	if found == true {
		t.Errorf("Expected key to be expired, but it was still in cache")
	}

	fmt.Println("Expiration test PASSED!")

}
