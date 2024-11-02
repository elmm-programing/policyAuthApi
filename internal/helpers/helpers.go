package helpers

import "policyAuth/internal/database"

type Helpers struct {
  DB   database.DatabaseService
}

func UniqueStrings(input []string) []string {
    uniqueMap := make(map[string]struct{})
    uniqueSlice := []string{}

    for _, elem := range input {
        if _, exists := uniqueMap[elem]; !exists {
            uniqueMap[elem] = struct{}{}
            uniqueSlice = append(uniqueSlice, elem)
        }
    }

    return uniqueSlice
}
