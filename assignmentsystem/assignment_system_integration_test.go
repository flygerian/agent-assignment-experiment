package assignmentsystem

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// NewAssignmentSystemWithState creates a new AssignmentSystem with predefined internal state
// This constructor is useful for integration testing scenarios where you need specific initial conditions
func NewAssignmentSystemWithState(agentAssignments map[string]*AgentWorkQueue) AssignmentSystem {
	// Build the accountAgents map from the agentAssignments
	accountAgents := make(map[string][]string)
	for agentName, wq := range agentAssignments {
		if _, exists := accountAgents[wq.Account]; !exists {
			accountAgents[wq.Account] = make([]string, 0)
		}
		accountAgents[wq.Account] = append(accountAgents[wq.Account], agentName)
	}

	return AssignmentSystem{
		accountAgents:    accountAgents,
		agentAssignments: agentAssignments,
	}
}

func TestIntegrationAssignmentSystemBasicWorkflow(t *testing.T) {
	// Test basic assignment workflow with multiple agents and conversations
	initialData := []AgentNameAndAccount{
		{Name: "agent1", Account: "account1", Limit: 3},
		{Name: "agent2", Account: "account1", Limit: 2},
		{Name: "agent3", Account: "account2", Limit: 1},
	}

	system := NewAssignmentSystem(initialData)

	conversations := []ConversationToAssign{
		{ConversationID: "conv1", Account: "account1"},
		{ConversationID: "conv2", Account: "account1"},
		{ConversationID: "conv3", Account: "account2"},
	}

	assignedAgents, err := system.Assign(conversations)

	assert.NoError(t, err)
	assert.Len(t, assignedAgents, 3)

	// Verify assignments were made correctly
	assert.Contains(t, []string{"agent1", "agent2"}, assignedAgents[0]) // account1
	assert.Contains(t, []string{"agent1", "agent2"}, assignedAgents[1]) // account1
	assert.Equal(t, "agent3", assignedAgents[2])                        // account2
}

func TestIntegrationAssignmentSystemWithPreExistingState(t *testing.T) {
	// Test assignment with pre-existing state using the new constructor
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	preExistingState := map[string]*AgentWorkQueue{
		"agent1": {
			AgentName:          "agent1",
			Account:            "account1",
			Limit:              3,
			Queue:              []string{"existing1"},
			LastAssignmentTime: &oneHourAgo,
		},
		"agent2": {
			AgentName: "agent2",
			Account:   "account1",
			Limit:     3,
			Queue:     []string{},
		},
		"agent3": {
			AgentName: "agent3",
			Account:   "account2",
			Limit:     1,
			Queue:     []string{"existing2"},
		},
	}

	system := NewAssignmentSystemWithState(preExistingState)

	conversations := []ConversationToAssign{
		{ConversationID: "new1", Account: "account1"},
		{ConversationID: "new2", Account: "account1"},
		{ConversationID: "new3", Account: "account2"},
	}

	assignedAgents, err := system.Assign(conversations)

	// Should have error because agent3 is at capacity, but only 2 successful assignments
	assert.Error(t, err)
	assert.Len(t, assignedAgents, 2) // Only 2 successful assignments

	// agent2 should get the first assignment (empty queue)
	assert.Equal(t, "agent2", assignedAgents[0])
	// agent1 should get the second (has room, was assigned an hour ago)
	assert.Equal(t, "agent1", assignedAgents[1])
}

func TestIntegrationAssignmentSystemLoadBalancing(t *testing.T) {
	// Test that assignments are properly load balanced
	initialData := []AgentNameAndAccount{
		{Name: "agent1", Account: "account1", Limit: 5},
		{Name: "agent2", Account: "account1", Limit: 5},
		{Name: "agent3", Account: "account1", Limit: 5},
	}

	system := NewAssignmentSystem(initialData)

	// Assign 9 conversations to 3 agents (should distribute evenly)
	conversations := make([]ConversationToAssign, 9)
	for i := range 9 {
		conversations[i] = ConversationToAssign{
			ConversationID: fmt.Sprintf("conv%d", i+1),
			Account:        "account1",
		}
	}

	assignedAgents, err := system.Assign(conversations)

	assert.NoError(t, err)
	assert.Len(t, assignedAgents, 9)

	// Count assignments per agent
	agentCounts := make(map[string]int)
	for _, agent := range assignedAgents {
		agentCounts[agent]++
	}

	// Each agent should have 3 assignments (perfect distribution)
	assert.Equal(t, 3, agentCounts["agent1"])
	assert.Equal(t, 3, agentCounts["agent2"])
	assert.Equal(t, 3, agentCounts["agent3"])
}

func TestIntegrationAssignmentSystemAccountIsolation(t *testing.T) {
	// Test that agents are properly isolated by account
	initialData := []AgentNameAndAccount{
		{Name: "agent1", Account: "account1", Limit: 3},
		{Name: "agent2", Account: "account2", Limit: 3},
		{Name: "agent3", Account: "account1", Limit: 3},
	}

	system := NewAssignmentSystem(initialData)

	conversations := []ConversationToAssign{
		{ConversationID: "conv1", Account: "account1"},
		{ConversationID: "conv2", Account: "account2"},
		{ConversationID: "conv3", Account: "account1"},
	}

	assignedAgents, err := system.Assign(conversations)

	assert.NoError(t, err)
	assert.Len(t, assignedAgents, 3)

	// Verify account isolation
	assert.Contains(t, []string{"agent1", "agent3"}, assignedAgents[0]) // account1
	assert.Equal(t, "agent2", assignedAgents[1])                        // account2
	assert.Contains(t, []string{"agent1", "agent3"}, assignedAgents[2]) // account1
}

func TestIntegrationAssignmentSystemAtCapacity(t *testing.T) {
	// Test behavior when all agents are at capacity
	initialData := []AgentNameAndAccount{
		{Name: "agent1", Account: "account1", Limit: 1},
		{Name: "agent2", Account: "account1", Limit: 1},
	}

	system := NewAssignmentSystem(initialData)

	// Fill agents to capacity
	conversations1 := []ConversationToAssign{
		{ConversationID: "conv1", Account: "account1"},
		{ConversationID: "conv2", Account: "account1"},
	}

	assignedAgents, err := system.Assign(conversations1)
	assert.NoError(t, err)
	assert.Len(t, assignedAgents, 2)

	// Try to assign more when at capacity
	conversations2 := []ConversationToAssign{
		{ConversationID: "conv3", Account: "account1"},
	}

	assignedAgents2, err := system.Assign(conversations2)
	assert.Error(t, err)
	// When no agents are available, the system returns empty slice (no successful assignments)
	assert.Len(t, assignedAgents2, 0) // No successful assignments
}

func TestIntegrationAssignmentSystemTimeBasedAssignment(t *testing.T) {
	// Test that time-based assignment logic works correctly
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	preExistingState := map[string]*AgentWorkQueue{
		"agent1": {
			AgentName:          "agent1",
			Account:            "account1",
			Limit:              3,
			Queue:              []string{"existing1"},
			LastAssignmentTime: &oneHourAgo,
		},
		"agent2": {
			AgentName:          "agent2",
			Account:            "account1",
			Limit:              3,
			Queue:              []string{"existing2"},
			LastAssignmentTime: &twoHoursAgo,
		},
	}

	system := NewAssignmentSystemWithState(preExistingState)

	// Both agents have same queue length, but agent2 was assigned longer ago
	conversations := []ConversationToAssign{
		{ConversationID: "new1", Account: "account1"},
	}

	assignedAgents, err := system.Assign(conversations)

	assert.NoError(t, err)
	assert.Len(t, assignedAgents, 1)
	assert.Equal(t, "agent2", assignedAgents[0]) // Should choose agent2 (older assignment)
}

func TestIntegrationAssignmentSystemMixedScenarios(t *testing.T) {
	// Test complex scenario with mixed states, different accounts, and various limits
	now := time.Now()
	thirtyMinutesAgo := now.Add(-30 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)

	preExistingState := map[string]*AgentWorkQueue{
		"agent1": {
			AgentName:          "agent1",
			Account:            "account1",
			Limit:              5,
			Queue:              []string{"existing1", "existing2"},
			LastAssignmentTime: &thirtyMinutesAgo,
		},
		"agent2": {
			AgentName:          "agent2",
			Account:            "account1",
			Limit:              3,
			Queue:              []string{"existing3"},
			LastAssignmentTime: &oneHourAgo,
		},
		"agent3": {
			AgentName: "agent3",
			Account:   "account2",
			Limit:     1,
			Queue:     []string{},
		},
		"agent4": {
			AgentName: "agent4",
			Account:   "account2",
			Limit:     1,
			Queue:     []string{"existing4"},
		},
	}

	system := NewAssignmentSystemWithState(preExistingState)

	conversations := []ConversationToAssign{
		{ConversationID: "new1", Account: "account1"}, // Should go to agent2 (less work, older assignment)
		{ConversationID: "new2", Account: "account1"}, // Should go to agent1 (has room)
		{ConversationID: "new3", Account: "account2"}, // Should go to agent3 (empty queue)
		{ConversationID: "new4", Account: "account2"}, // Should fail - both agents at capacity
	}

	assignedAgents, err := system.Assign(conversations)

	assert.Error(t, err)             // Should have error due to failed assignment
	assert.Len(t, assignedAgents, 3) // Only 3 successful assignments

	// Verify assignments based on complex logic
	assert.Equal(t, "agent2", assignedAgents[0]) // Least work + older assignment
	assert.Equal(t, "agent1", assignedAgents[1]) // Has room
	assert.Equal(t, "agent3", assignedAgents[2]) // Empty queue
}
