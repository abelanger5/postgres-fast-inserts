package main

import (
	"encoding/json"
	"math/rand"
	"strings"
)

const numFields = 10

// Generate a length value using random number generation
func generateLength() int {
	x := rand.Float64() * float64(maxPayloadSize)

	length := int(x / numFields)

	// Ensure minimum length of 1
	if length < 1 {
		length = 1
	}

	return length
}

// Generate a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Generate a field name combining a random adjective with a famous scientist's name
func generateFieldName() string {
	adjectives := []string{
		"curious", "brilliant", "eccentric", "visionary", "persistent",
		"innovative", "meticulous", "logical", "analytical", "creative",
		"theoretical", "practical", "quantum", "relativistic", "molecular",
		"astute", "pioneering", "dedicated", "revolutionary", "insightful",
	}

	scientists := []string{
		"Einstein", "Newton", "Curie", "Darwin", "Tesla",
		"Bohr", "Hawking", "Feynman", "Lovelace", "Turing",
		"Hopper", "Heisenberg", "Schrodinger", "Goodall", "Meitner",
		"Fermi", "Planck", "Franklin", "Faraday", "Maxwell",
	}

	adjective := adjectives[rand.Intn(len(adjectives))]
	scientist := scientists[rand.Intn(len(scientists))]

	return strings.ToLower(adjective + "-" + scientist)
}

// Generate a JSON payload with random field names and character array values
func generateJSONPayload() []byte {
	// Create a map to hold our field names and values
	payload := make(map[string]string)

	// Generate the specified number of fields
	for i := 0; i < numFields; i++ {
		fieldName := generateFieldName()

		// Ensure unique field names
		for _, exists := payload[fieldName]; exists; {
			fieldName = generateFieldName()
			_, exists = payload[fieldName]
		}

		length := generateLength()

		// Generate a random string of that length
		payload[fieldName] = generateRandomString(length)
	}

	// Convert the map to JSON
	resBytes, err := json.Marshal(payload)

	if err != nil {
		panic(err)
	}

	// Return the JSON payload
	return resBytes
}
