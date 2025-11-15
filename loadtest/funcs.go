package loadtest

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/flygerian/assignment-system/assignmentsystem"
)

// GenerateAgentWorkQueues creates 1 million agents with varying limits from 5-20,
// distributed unevenly across accounts (large accounts and small accounts)
func GenerateAgentWorkQueues() []assignmentsystem.AgentNameAndAccount {
	const totalAgents = 1_000_000
	const minLimit = 5
	const maxLimit = 20

	// Create account distribution: 80% small accounts, 20% large accounts
	const largeAccountRatio = 0.2
	const smallAccountRatio = 0.8

	// Large accounts will have 1000-5000 agents, small accounts will have 10-100 agents
	const minLargeAccountAgents = 1000
	const maxLargeAccountAgents = 5000
	const minSmallAccountAgents = 10
	const maxSmallAccountAgents = 100

	var agents []assignmentsystem.AgentNameAndAccount
	agentsCreated := 0
	accountID := 1

	// Create large accounts first (20% of agents)
	targetLargeAgents := int(float64(totalAgents) * largeAccountRatio)
	largeAgentsCreated := 0

	for largeAgentsCreated < targetLargeAgents {
		agentsInThisAccount := rand.Intn(maxLargeAccountAgents-minLargeAccountAgents+1) + minLargeAccountAgents
		if largeAgentsCreated+agentsInThisAccount > targetLargeAgents {
			agentsInThisAccount = targetLargeAgents - largeAgentsCreated
		}

		accountName := fmt.Sprintf("large_account_%d", accountID)
		for i := 0; i < agentsInThisAccount; i++ {
			agent := assignmentsystem.AgentNameAndAccount{
				Name:    fmt.Sprintf("agent_%s_%d", accountName, i+1),
				Account: accountName,
				Limit:   rand.Intn(maxLimit-minLimit+1) + minLimit,
			}
			agents = append(agents, agent)
		}

		largeAgentsCreated += agentsInThisAccount
		agentsCreated += agentsInThisAccount
		accountID++
	}

	// Create small accounts (remaining 80% of agents)
	targetSmallAgents := totalAgents - targetLargeAgents
	smallAgentsCreated := 0

	for smallAgentsCreated < targetSmallAgents {
		agentsInThisAccount := rand.Intn(maxSmallAccountAgents-minSmallAccountAgents+1) + minSmallAccountAgents
		if smallAgentsCreated+agentsInThisAccount > targetSmallAgents {
			agentsInThisAccount = targetSmallAgents - smallAgentsCreated
		}

		accountName := fmt.Sprintf("small_account_%d", accountID)
		for i := 0; i < agentsInThisAccount; i++ {
			agent := assignmentsystem.AgentNameAndAccount{
				Name:    fmt.Sprintf("agent_%s_%d", accountName, i+1),
				Account: accountName,
				Limit:   rand.Intn(maxLimit-minLimit+1) + minLimit,
			}
			agents = append(agents, agent)
		}

		smallAgentsCreated += agentsInThisAccount
		agentsCreated += agentsInThisAccount
		accountID++
	}

	// Shuffle the agents to ensure random distribution
	rand.Shuffle(len(agents), func(i, j int) {
		agents[i], agents[j] = agents[j], agents[i]
	})

	return agents
}

func GetUniqueAccounts(agents []assignmentsystem.AgentNameAndAccount) []string {
	var uniqueAccounts []string

	for _, agent := range agents {
		if !slices.Contains(uniqueAccounts, agent.Account) {
			uniqueAccounts = append(uniqueAccounts, agent.Account)
		}
	}

	return uniqueAccounts
}
