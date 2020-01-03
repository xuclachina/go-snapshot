package main

import (
	"sync"
	"testing"
)

func TestAdd1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Logio(".", &wg)
	wg.Wait()
}
