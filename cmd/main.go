package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/flygerian/assignment-system/assignmentsystem"
	"github.com/flygerian/assignment-system/loadtest"
)

func main() {
	fmt.Println("Starting assignment loop...")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	log.Printf("Generating work queues, this may take a while")
	agentWqs := loadtest.GenerateAgentWorkQueues()
	log.Printf("Making unique list of acccounts")

	accounts := loadtest.GetUniqueAccounts(agentWqs)
	log.Printf("Generating conversations")
	conversations := generateConversations(accounts, 10000)

	system := assignmentsystem.NewAssignmentSystem(agentWqs)

	assignmentInProgress := false
	tickCounter := 0

	for {
		select {
		case <-ticker.C:
			if !assignmentInProgress {
				// Mark assignment as in progress
				assignmentInProgress = true

				log.Printf("Starting assignment batch %d at %s", tickCounter+1, time.Now().Format("15:04:05"))

				start := time.Now()
				_, err := system.Assign(conversations[tickCounter*100 : (tickCounter+1)*100])
				if err != nil {
					log.Fatal(err)
					// We probably should have metrics here to measure failure ans alerting
				}
				elapsed := time.Since(start)

				log.Printf("Completed assignment batch %d in %v", tickCounter+1, elapsed)

				tickCounter++
				assignmentInProgress = false

				// Exit after 100th tick (tickCounter will be 100 after increment)
				if tickCounter >= 100 {
					log.Println("Completed 100 ticks, exiting...")
					return
				}

			} else {
				log.Println("Assignment in progress, skipping this tick")
			}
		}
	}
}

// generateConversations creates 100 conversations to assign based on the tick number
func generateConversations(accounts []string, numberToGenerate int) []assignmentsystem.ConversationToAssign {
	conversations := make([]assignmentsystem.ConversationToAssign, numberToGenerate)

	for i := range numberToGenerate {
		// Randomly select an account from the provided list
		accountIndex := rand.Intn(len(accounts))

		conversations[i] = assignmentsystem.ConversationToAssign{
			ConversationID: fmt.Sprintf("conversation-%d", i+1),
			Account:        accounts[accountIndex],
		}
	}

	return conversations
}
