package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	walletAddress = "[YOUR ADDRESS HERE]"
	airdropAmount = "2.5"
	network       = "devnet"
	retryDelay    = 30 * time.Second
	waitBetween   = 3 * time.Hour
	cycleInterval = 4*time.Hour
)

func requestAirdrop() bool {
	cmd := exec.Command("solana", "airdrop", airdropAmount, walletAddress, "-u", network)
	output, err := cmd.CombinedOutput()
	
	outputStr := string(output)
	log.Printf("Airdrop attempt output: %s", outputStr)
	
	if err != nil {
		log.Printf("Command error: %v", err)
		return false
	}
	
	// Check if successful by looking for "Signature:" in output
	if strings.Contains(outputStr, "Signature:") {
		log.Println("✓ Airdrop successful!")
		return true
	}
	
	if strings.Contains(outputStr, "rate limit") || strings.Contains(outputStr, "Error") {
		log.Println("✗ Airdrop failed (rate limit or error)")
		return false
	}
	
	return false
}

func retryUntilSuccess() {
	log.Println("Starting airdrop request cycle...")
	
	for {
		if requestAirdrop() {
			return
		}
		
		log.Printf("Retrying in %v...", retryDelay)
		time.Sleep(retryDelay)
	}
}

func main() {
	log.Println("=== Solana Devnet Airdrop Bot Started ===")
	log.Printf("Wallet: %s", walletAddress)
	log.Printf("Amount: %s SOL per request", airdropAmount)
	log.Printf("Network: %s", network)
	log.Println("=========================================")
	
	for {
		// First airdrop attempt
		log.Println("\n[Attempt 1/2] Requesting first airdrop...")
		retryUntilSuccess()
		
		// Wait 3 hours
		log.Printf("\nWaiting %v before next attempt...\n", waitBetween)
		time.Sleep(waitBetween)
		
		// Second airdrop attempt
		log.Println("\n[Attempt 2/2] Requesting second airdrop...")
		retryUntilSuccess()
		
		// Wait for next cycle (4 hours 1 minute)
		log.Printf("\nCycle complete! Waiting %v before starting next cycle...\n", cycleInterval)
		time.Sleep(cycleInterval)
	}
}