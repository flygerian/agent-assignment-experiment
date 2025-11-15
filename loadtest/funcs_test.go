package loadtest

import (
	"testing"

	"github.com/flygerian/assignment-system/assignmentsystem"
)

func TestGenerateAgentWorkQueues(t *testing.T) {
	agents := GenerateAgentWorkQueues()

	// Test total number of agents
	if len(agents) != 1_000_000 {
		t.Errorf("Expected 1,000,000 agents, got %d", len(agents))
	}

	// Test limit ranges
	for _, agent := range agents {
		if agent.Limit < 5 || agent.Limit > 20 {
			t.Errorf("Agent %s has invalid limit %d, expected 5-20", agent.Name, agent.Limit)
		}
	}

	// Test account distribution
	accountMap := make(map[string]int)
	for _, agent := range agents {
		accountMap[agent.Account]++
	}

	var largeAccounts, smallAccounts int
	var largeAgents, smallAgents int

	for _, count := range accountMap {
		if count >= 1000 {
			largeAccounts++
			largeAgents += count
		} else {
			smallAccounts++
			smallAgents += count
		}
	}

	// Check distribution ratios (with some tolerance)
	largeRatio := float64(largeAgents) / float64(len(agents))
	smallRatio := float64(smallAgents) / float64(len(agents))

	if largeRatio < 0.15 || largeRatio > 0.25 {
		t.Errorf("Large account ratio %.2f outside expected range 0.15-0.25", largeRatio)
	}

	if smallRatio < 0.75 || smallRatio > 0.85 {
		t.Errorf("Small account ratio %.2f outside expected range 0.75-0.85", smallRatio)
	}

	// Test that all agents have unique names
	names := make(map[string]bool)
	for _, agent := range agents {
		if names[agent.Name] {
			t.Errorf("Duplicate agent name found: %s", agent.Name)
		}
		names[agent.Name] = true
	}

	// Test that agents are assigned to accounts correctly
	for _, agent := range agents {
		if agent.Account == "" {
			t.Errorf("Agent %s has empty account", agent.Name)
		}
		if agent.Name == "" {
			t.Errorf("Agent has empty name")
		}
	}
}

func TestGetUniqueAccounts(t *testing.T) {
	// Test case 1: Multiple agents with duplicate accounts
	agents := []assignmentsystem.AgentNameAndAccount{
		{Name: "agent1", Account: "accountA", Limit: 10},
		{Name: "agent2", Account: "accountB", Limit: 15},
		{Name: "agent3", Account: "accountA", Limit: 20},
		{Name: "agent4", Account: "accountC", Limit: 5},
		{Name: "agent5", Account: "accountB", Limit: 10},
	}

	uniqueAccounts := GetUniqueAccounts(agents)

	// Should return 3 unique accounts
	if len(uniqueAccounts) != 3 {
		t.Errorf("Expected 3 unique accounts, got %d", len(uniqueAccounts))
	}

	// Check that all expected accounts are present
	expectedAccounts := map[string]bool{"accountA": false, "accountB": false, "accountC": false}
	for _, account := range uniqueAccounts {
		if _, exists := expectedAccounts[account]; !exists {
			t.Errorf("Unexpected account: %s", account)
		}
		expectedAccounts[account] = true
	}

	// Verify all accounts were found
	for account, found := range expectedAccounts {
		if !found {
			t.Errorf("Missing expected account: %s", account)
		}
	}

	// Test case 2: Empty slice
	emptyAgents := []assignmentsystem.AgentNameAndAccount{}
	emptyResult := GetUniqueAccounts(emptyAgents)
	if len(emptyResult) != 0 {
		t.Errorf("Expected empty result for empty input, got %d accounts", len(emptyResult))
	}

	// Test case 3: Single agent
	singleAgent := []assignmentsystem.AgentNameAndAccount{
		{Name: "solo", Account: "soloAccount", Limit: 10},
	}
	singleResult := GetUniqueAccounts(singleAgent)
	if len(singleResult) != 1 || singleResult[0] != "soloAccount" {
		t.Errorf("Expected single account 'soloAccount', got %v", singleResult)
	}
}
