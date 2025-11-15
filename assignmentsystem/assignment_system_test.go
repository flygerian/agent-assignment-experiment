package assignmentsystem

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGettingEligibleAccounts(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]*AgentWorkQueue
		expectation []*AgentWorkQueue
	}{
		{
			name: "Finds agents below limit for account",
			input: map[string]*AgentWorkQueue{
				"agent1": {
					Account: "account1",
					Limit:   5,
					Queue:   []string{"item1"}, // Below limit
				},
				"agent2": {
					Account: "account2",
					Limit:   2,
					Queue:   []string{"item1"}, // Below limit
				},
				"agent3": {
					Account: "account1",
					Limit:   1,
					Queue:   []string{"item1"}, // At limit
				},
			},
			expectation: []*AgentWorkQueue{
				{
					Account: "account2",
					Limit:   2,
					Queue:   []string{"item1"},
				},
			},
		},
		{
			name:        "Empty input map",
			input:       map[string]*AgentWorkQueue{},
			expectation: []*AgentWorkQueue{},
		},
		{
			name: "Agents below limit with matching account",
			input: map[string]*AgentWorkQueue{
				"agent1": {
					Account: "account1",
					Limit:   5,
					Queue:   []string{"item1", "item2"},
				},
				"agent2": {
					Account: "account1",
					Limit:   3,
					Queue:   []string{"item1"},
				},
			},
			expectation: []*AgentWorkQueue{
				{
					Account: "account1",
					Limit:   5,
					Queue:   []string{"item1", "item2"},
				},
				{
					Account: "account1",
					Limit:   3,
					Queue:   []string{"item1"},
				},
			},
		},
		{
			name: "Agents at limit should be excluded",
			input: map[string]*AgentWorkQueue{
				"agent1": {
					Account: "account2",
					Limit:   2,
					Queue:   []string{"item1", "item2"},
				},
				"agent2": {
					Account: "account2",
					Limit:   1,
					Queue:   []string{},
				},
			},
			expectation: []*AgentWorkQueue{
				{
					Account: "account2",
					Limit:   1,
					Queue:   []string{},
				},
			},
		},
		{
			name: "Mixed scenarios - some below limit, some at limit",
			input: map[string]*AgentWorkQueue{
				"agent1": {
					Account: "account3",
					Limit:   3,
					Queue:   []string{"item1", "item2", "item3"},
				},
				"agent2": {
					Account: "account3",
					Limit:   2,
					Queue:   []string{"item1"},
				},
				"agent3": {
					Account: "account1",
					Limit:   1,
					Queue:   []string{"item1"},
				},
			},
			expectation: []*AgentWorkQueue{
				{
					Account: "account3",
					Limit:   2,
					Queue:   []string{"item1"},
				},
			},
		},
		{
			name: "Multiple agents below limit with same account",
			input: map[string]*AgentWorkQueue{
				"agent1": {
					Account: "account4",
					Limit:   3,
					Queue:   []string{"item1", "item2"},
				},
				"agent2": {
					Account: "account1",
					Limit:   1,
					Queue:   []string{"item1"},
				},
				"agent3": {
					Account: "account4",
					Limit:   3,
					Queue:   []string{"item1"},
				},
			},
			expectation: []*AgentWorkQueue{
				{
					Account: "account4",
					Limit:   3,
					Queue:   []string{"item1", "item2"},
				},
				{
					Account: "account4",
					Limit:   3,
					Queue:   []string{"item1"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Determine the account to search for based on the test case
			accountToSearch := "account2"
			switch test.name {
			case "Agents below limit with matching account", "No agents at limit":
				accountToSearch = "account1"
			case "Agents at limit should be excluded":
				accountToSearch = "account2"
			case "Mixed scenarios - some below limit, some at limit":
				accountToSearch = "account3"
			case "Multiple agents below limit with same account":
				accountToSearch = "account4"
			}

			// Create accountAgents map from the input data
			accountAgents := make(map[string][]string)
			for agentName, wq := range test.input {
				if _, exists := accountAgents[wq.Account]; !exists {
					accountAgents[wq.Account] = make([]string, 0)
				}
				accountAgents[wq.Account] = append(accountAgents[wq.Account], agentName)
			}

			accounts := getEligibleAgentWorkQueues(accountAgents, test.input, accountToSearch)
			assert.ElementsMatch(t, accounts, test.expectation)
		})
	}
}

func TestGetWorkQueueWithTheLeastRecentAssignment(t *testing.T) {
	// Create test times for consistent testing
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	thirtyMinutesAgo := now.Add(-30 * time.Minute)
	twoHoursAgo := now.Add(-2 * time.Hour)
	threeHoursAgo := now.Add(-3 * time.Hour)

	tests := []struct {
		name        string
		input       []*AgentWorkQueue
		expectation *AgentWorkQueue
	}{
		{
			name: "Multiple queues with different assignment times",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &thirtyMinutesAgo},
				{LastAssignmentTime: &twoHoursAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &twoHoursAgo},
		},
		{
			name: "Single work queue",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &oneHourAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &oneHourAgo},
		},
		{
			name: "All queues with same assignment time",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &oneHourAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &oneHourAgo},
		},
		{
			name:        "Empty input",
			input:       []*AgentWorkQueue{},
			expectation: nil,
		},
		{
			name: "Mix of nil and non-nil assignment times - nil should be prioritized",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: nil},
				{LastAssignmentTime: &thirtyMinutesAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: nil},
		},
		{
			name: "All nil assignment times",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: nil},
				{LastAssignmentTime: nil},
				{LastAssignmentTime: nil},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: nil},
		},
		{
			name: "Very recent assignment should not be selected",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &thirtyMinutesAgo},
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &twoHoursAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &twoHoursAgo},
		},
		{
			name: "Very old assignment should be selected",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &thirtyMinutesAgo},
				{LastAssignmentTime: &threeHoursAgo},
				{LastAssignmentTime: &twoHoursAgo},
				{LastAssignmentTime: &oneHourAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &threeHoursAgo},
		},
		{
			name: "Assignment times in chronological order",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &threeHoursAgo},
				{LastAssignmentTime: &twoHoursAgo},
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &thirtyMinutesAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &threeHoursAgo},
		},
		{
			name: "Assignment times in reverse order",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &thirtyMinutesAgo},
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &twoHoursAgo},
				{LastAssignmentTime: &threeHoursAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &threeHoursAgo},
		},
		{
			name: "Mixed recent and old assignments",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &thirtyMinutesAgo},
				{LastAssignmentTime: &threeHoursAgo},
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: &twoHoursAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &threeHoursAgo},
		},
		{
			name: "Single queue with nil assignment time",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: nil},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: nil},
		},
		{
			name: "Single queue with assignment time",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &oneHourAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: &oneHourAgo},
		},
		{
			name: "Nil assignment time in middle of list",
			input: []*AgentWorkQueue{
				{LastAssignmentTime: &oneHourAgo},
				{LastAssignmentTime: nil},
				{LastAssignmentTime: &twoHoursAgo},
				{LastAssignmentTime: &thirtyMinutesAgo},
			},
			expectation: &AgentWorkQueue{LastAssignmentTime: nil},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getWorkQueueWithTheLeastRecentAssignment(test.input)

			// For nil comparison, we need to handle it specially
			if test.expectation == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)

			// Compare the LastAssignmentTime fields
			if test.expectation.LastAssignmentTime == nil {
				assert.Nil(t, result.LastAssignmentTime)
			} else {
				assert.NotNil(t, result.LastAssignmentTime)
				// Compare the actual time values
				assert.Equal(t, test.expectation.LastAssignmentTime.Unix(), result.LastAssignmentTime.Unix())
			}
		})
	}
}

func TestGetWorkqueuesWithLeastAmountOfWork(t *testing.T) {
	tests := []struct {
		name        string
		input       []*AgentWorkQueue
		expectation []*AgentWorkQueue
	}{
		{
			name: "Multiple queues with different lengths",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1", "item2", "item3"}}, // Length 3
				{Queue: []string{"item1"}},                   // Length 1
				{Queue: []string{"item1", "item2"}},          // Length 2
				{Queue: []string{"item1"}},                   // Length 1
			},
			expectation: []*AgentWorkQueue{
				{Queue: []string{"item1"}},
				{Queue: []string{"item1"}},
			},
		},
		{
			name: "Single work queue",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1", "item2"}},
			},
			expectation: []*AgentWorkQueue{
				{Queue: []string{"item1", "item2"}},
			},
		},
		{
			name: "All queues same length",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1"}},
				{Queue: []string{"item2"}},
				{Queue: []string{"item3"}},
			},
			expectation: []*AgentWorkQueue{
				{Queue: []string{"item1"}},
				{Queue: []string{"item2"}},
				{Queue: []string{"item3"}},
			},
		},
		{
			name:        "Empty input",
			input:       []*AgentWorkQueue{},
			expectation: []*AgentWorkQueue{},
		},
		{
			name: "Mix of zero and non-zero lengths",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1", "item2"}},
				{Queue: []string{}}, // Empty queue
				{Queue: []string{"item1"}},
				{Queue: []string{}}, // Empty queue
			},
			expectation: []*AgentWorkQueue{
				{Queue: []string{}},
				{Queue: []string{}},
			},
		},
		{
			name: "Multiple queues with same minimum length",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1", "item2", "item3"}},
				{Queue: []string{"item4"}},
				{Queue: []string{"item5", "item6"}},
				{Queue: []string{"item7"}},
				{Queue: []string{"item8", "item9", "item10"}},
			},
			expectation: []*AgentWorkQueue{
				{Queue: []string{"item4"}},
				{Queue: []string{"item7"}},
			},
		},
		{
			name: "All queues at different lengths",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1", "item2"}},                   // Length 2
				{Queue: []string{"item1", "item2", "item3", "item4"}}, // Length 4
				{Queue: []string{"item1"}},                            // Length 1
				{Queue: []string{"item1", "item2", "item3"}},          // Length 3
			},
			expectation: []*AgentWorkQueue{
				{Queue: []string{"item1"}},
			},
		},
		{
			name: "Nil queue fields treated as empty",
			input: []*AgentWorkQueue{
				{Queue: []string{"item1"}},
				{Queue: nil}, // Nil queue
				{Queue: []string{"item1", "item2"}},
				{Queue: nil}, // Nil queue
			},
			expectation: []*AgentWorkQueue{
				{Queue: nil},
				{Queue: nil},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getWorkqueuesWithLeastAmountOfWork(test.input)
			assert.ElementsMatch(t, result, test.expectation)
		})
	}
}
