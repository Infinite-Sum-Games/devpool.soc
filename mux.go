package main

import "sync"

type BotMux struct {
	IssueClaim  chan IssueAction
	Bounty      chan BountyAction
	Achievement chan Achievement
	Solution    chan Solution
	Mutex       sync.Mutex
}

func NewBotMux(bufferSize int) *BotMux {
	return &BotMux{
		IssueClaim:  make(chan IssueAction, bufferSize),
		Bounty:      make(chan BountyAction, bufferSize),
		Achievement: make(chan Achievement, bufferSize),
		Solution:    make(chan Solution, bufferSize),
		Mutex:       sync.Mutex{},
	}
}

// Initiates the bot multiplexer to list to incoming messages for various event
// streams and process them accordingly
func (bm *BotMux) Start() {
	for {
		select {
		case issue := <-bm.IssueClaim:
			bm.Mutex.Lock()
			if issue.Extend == true {
				ManageExtension(issue.ParticipantUsername, issue.Url)
			} else {
				ManageIssueClaim(issue.ParticipantUsername, issue.Claim, issue.Url)
			}
			bm.Mutex.Unlock()
		case cash := <-bm.Bounty:
			bm.Mutex.Lock()
			ManageBounty(cash.ParticipantUsername, cash.Amount, cash.Action, cash.Url)
			bm.Mutex.Unlock()
		case badge := <-bm.Achievement:
			bm.Mutex.Lock()
			ManageBadge(badge.ParticipantUsername, badge.Type, badge.Url)
			bm.Mutex.Unlock()
		case sol := <-bm.Solution:
			bm.Mutex.Lock()
			ManageSolution(sol.ParticipantUsername, sol.Merged, sol.Url)
			bm.Mutex.Unlock()
		}
	}
}

// Gracefully shuts down the Bot multiplexer
func (bm *BotMux) Shutdown() {
	close(bm.IssueClaim)
	close(bm.Bounty)
	close(bm.Achievement)
	close(bm.Solution)
}
