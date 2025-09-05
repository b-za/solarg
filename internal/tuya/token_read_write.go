package tuya

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const tokenFileName = "tuya_token.json"

func SaveToken(tokenResult TokenResult) {

	tokenFilePath, err := getTokenPath()
	if err != nil {
		log.Printf("Failed to get token path: %v", err)
		return
	}

	tokenData, err := json.MarshalIndent(tokenResult, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal token data: %v", err)
		return
	}

	err = os.WriteFile(tokenFilePath, tokenData, 0644)
	if err != nil {
		log.Printf("Failed to write token to file '%s': %v", tokenFilePath, err)
	} else {
		fmt.Printf("\nToken saved successfully to %s\n", tokenFilePath)
	}
}

func ReadToken() (TokenResult, error) {
	var tokenResult TokenResult

	tokenFilePath, err := getTokenPath()
	if err != nil {
		log.Printf("Failed to get token path: %v", err)
		return tokenResult, err
	}

	tokenData, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return tokenResult, fmt.Errorf("failed to read token file '%s': %w", tokenFilePath, err)
	}

	if err := json.Unmarshal(tokenData, &tokenResult); err != nil {
		return tokenResult, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return tokenResult, nil
}

func getTokenPath() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not get the current file path")
	}

	dir := filepath.Dir(currentFile)

	projectRoot := ""
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot = dir
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	if projectRoot == "" {
		return "", fmt.Errorf("could not find go.mod in any parent directory")
	}

	credsDir := filepath.Join(projectRoot, "creds")
	if err := os.MkdirAll(credsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create creds directory: %w", err)
	}

	filePath := filepath.Join(credsDir, tokenFileName)
	return filePath, nil
}
