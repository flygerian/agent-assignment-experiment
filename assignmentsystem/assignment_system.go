package assignmentsystem

import (
	"fmt"
	"log"
	"time"
)

type AgentWorkQueue struct {
	Limit              int
	AgentName          string
	LastAssignmentTime *time.Time
	Queue              []string
	Account            string
}

type AssignmentSystem struct {
	accountAgents    map[string][]string
	agentAssignments map[string]*AgentWorkQueue // Pointer to AgentWorkQueue because map access in golang yields a copy
}

type AgentNameAndAccount struct {
	Name    string
	Account string
	Limit   int
}

type ConversationToAssign struct {
	ConversationID string
	Account        string
}

type conversationAssignmentError struct {
	ConversationToAssign
	Err error
}

func NewAssignmentSystem(initData []AgentNameAndAccount) AssignmentSystem {
	assignmentsystem := AssignmentSystem{
		accountAgents:    make(map[string][]string),
		agentAssignments: make(map[string]*AgentWorkQueue),
	}

	for _, nameAndAccount := range initData {
		assignmentsystem.agentAssignments[nameAndAccount.Name] = &AgentWorkQueue{
			AgentName: nameAndAccount.Name,
			Limit:     nameAndAccount.Limit,
			Queue:     make([]string, 0),
			Account:   nameAndAccount.Account,
		}

		if _, ok := assignmentsystem.accountAgents[nameAndAccount.Account]; !ok {
			assignmentsystem.accountAgents[nameAndAccount.Account] = make([]string, 0)
		}

		assignmentsystem.accountAgents[nameAndAccount.Account] = append(assignmentsystem.accountAgents[nameAndAccount.Account], nameAndAccount.Name)
	}

	return assignmentsystem
}

func (as *AssignmentSystem) SetLimit(agentName string, limit int) {
	as.agentAssignments[agentName].Limit = limit
}

func (as *AssignmentSystem) Assign(conversationsToAssign []ConversationToAssign) ([]string, error) {
	log.Printf("Assigning %d conversatons", len(conversationsToAssign))
	assignedAgents := make([]string, 0)
	failedAssignments := make([]conversationAssignmentError, 0)
	for _, conversation := range conversationsToAssign {
		assignment, err := as.assign(conversation)
		if err != nil {
			failedAssignments = append(failedAssignments, conversationAssignmentError{
				conversation,
				err,
			})

			continue
		}

		assignedAgents = append(assignedAgents, assignment)
	}

	return assignedAgents, as.constructError(failedAssignments)
}

func (as *AssignmentSystem) constructError(failedAssignments []conversationAssignmentError) error {
	if len(failedAssignments) == 0 {
		return nil
	}
	return fmt.Errorf("failed to assign %d conversations", len(failedAssignments))
}

func (as *AssignmentSystem) assign(conversation ConversationToAssign) (string, error) {
	// Get all the AgentWorkQueue(s) that belong to this account and are not at their limit
	eligibleWorkQueues := getEligibleAgentWorkQueues(as.accountAgents, as.agentAssignments, conversation.Account)
	// If no agents are available either reject maybe use a fall back queue
	if len(eligibleWorkQueues) == 0 {
		// we need to make a decision here
		return "", fmt.Errorf("no available agents to take on work")
	}
	// Get the agents with least amount of work
	workQueueWithLeastAmountOfWork := getWorkqueuesWithLeastAmountOfWork(eligibleWorkQueues)
	// If one assign the the current case to then return
	if len(workQueueWithLeastAmountOfWork) == 1 {
		// Assign
		return as.assignToWorkQueue(workQueueWithLeastAmountOfWork[0], conversation.ConversationID)
	}
	// more than one, pick the one with longest time.Now() - assignmentTime
	withLeastRecentAssignment := getWorkQueueWithTheLeastRecentAssignment(workQueueWithLeastAmountOfWork)

	return as.assignToWorkQueue(withLeastRecentAssignment, conversation.ConversationID)
}

func (as *AssignmentSystem) assignToWorkQueue(wq *AgentWorkQueue, conversationID string) (string, error) {
	wq.Queue = append(wq.Queue, conversationID)
	assignmentTime := time.Now()
	wq.LastAssignmentTime = &assignmentTime
	return wq.AgentName, nil
}

func getEligibleAgentWorkQueues(accountAgents map[string][]string, agentAssignments map[string]*AgentWorkQueue, account string) []*AgentWorkQueue {
	availableWorkQueues := make([]*AgentWorkQueue, 0)
	agentsForAccount := accountAgents[account]

	agentWqs := make([]*AgentWorkQueue, 0)
	for _, agent := range agentsForAccount {
		wq := agentAssignments[agent]
		agentWqs = append(agentWqs, wq)
	}

	for _, wq := range agentWqs {
		if len(wq.Queue) == wq.Limit {
			continue
		}

		availableWorkQueues = append(availableWorkQueues, wq)
	}

	return availableWorkQueues
}

func getWorkqueuesWithLeastAmountOfWork(workQueues []*AgentWorkQueue) []*AgentWorkQueue {
	if len(workQueues) == 0 {
		return []*AgentWorkQueue{}
	}

	lowsestWorkload := len(workQueues[0].Queue)
	workQueuesFoundSoFar := make([]*AgentWorkQueue, 0)

	for _, wq := range workQueues {
		if len(wq.Queue) < lowsestWorkload {
			workQueuesFoundSoFar = make([]*AgentWorkQueue, 0)
			workQueuesFoundSoFar = append(workQueuesFoundSoFar, wq)
			lowsestWorkload = len(wq.Queue)
			continue
		}

		if len(wq.Queue) == lowsestWorkload {
			workQueuesFoundSoFar = append(workQueuesFoundSoFar, wq)
		}
	}

	return workQueuesFoundSoFar
}

func getWorkQueueWithTheLeastRecentAssignment(workQueues []*AgentWorkQueue) *AgentWorkQueue {
	var highestDuration time.Duration
	var workQueueWithHighestDuration *AgentWorkQueue

	for _, wq := range workQueues {
		if wq.LastAssignmentTime == nil { // Nothing has been assigned to it break
			workQueueWithHighestDuration = wq
			break
		}

		if time.Since(*wq.LastAssignmentTime) > highestDuration {
			workQueueWithHighestDuration = wq
			highestDuration = time.Since(*wq.LastAssignmentTime)
			continue
		}
	}

	return workQueueWithHighestDuration
}
