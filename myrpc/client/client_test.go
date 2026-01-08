package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewClient(t *testing.T) {
	client := NewClient("")
	assert.NotNil(t, client)
}
