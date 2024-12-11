package util

import (
	"fmt"
	"math/rand"
	"time"
)

type Account struct {
	Name     string
	Balance  float64
	Currency string
}

// Initialize the random seed only once
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomName generates a random full name
func GenerateRandomName() string {
	firstNames := []string{"John", "Alice", "Siddhesh", "Emma", "Liam", "Sophia", "Noah", "Ava", "Oliver", "Isabella"}
	lastNames := []string{"Smith", "Johnson", "Brown", "Taylor", "Anderson", "Thomas", "Jackson", "White", "Harris", "Martin"}

	// Generate random indices
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	return fmt.Sprintf("%s %s", firstName, lastName)
}

// GenerateRandomBalance generates a random balance between min and max
func GenerateRandomBalance(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// GenerateRandomCurrency generates a random currency
func GenerateRandomCurrency() string {
	currencies := []string{"INR", "USD"}
	return currencies[rand.Intn(len(currencies))]
}

// GenerateRandomAccount creates a random account with random values
func GenerateRandomAccount(minBalance, maxBalance float64) Account {
	return Account{
		Name:     GenerateRandomName(),
		Balance:  GenerateRandomBalance(minBalance, maxBalance),
		Currency: GenerateRandomCurrency(),
	}
}
