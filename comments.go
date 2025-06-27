package main

import (
	"context"
	"github.com/google/go-github/v62/github"
)

const (

	// For managing all issues
	IssueOnboard = `Howdy contributors, listen up! This issue is under the 
	watchful gaze of DevPool (official bot of Summer of Code, 2025). Wanna tackle
	this one ? Just shout "/assign" in the comments and you can fork the repository
	and start working AFTER my confirmation.`

	IssueClaimed = `Alright @%s! Consider yourself officially assigned as of 
	%s. If you are thinking, "Whoa, hold up, I didn't sign up for this," no 
	worries, just yell "/unassign" and we can pretend this never happened. But 
	if you are stickign around, just know your clock for this ticks out on 
	%s. Don't screw it up :)`

	IssueUnclaimed = `Adios @%s! You are officially off the hook for this 
	issue. If you had a change of heart or just needed a break, no sweat. However 
	if this digital adventure still calls you (and if it's still open obv), just 
	shout "/assign" again. I am not picky! Or am I ?`

	DocSubmissions = `Someone actually spent time on ... documentation ? Well, you
	have made the world a slightly less confusion place for now @%s. Hopefully this 
	makes things easier for the poor schmucks.(unlike me, the a 100x dev)`

	HighImpact = `Consider my tiny, shriveled heart officially impressed. Now, go 
	and grab yourself a latte. You have earned it!`

	BugReport = `Well, well, @%s been paying attention! A bug, you say? 
	Is it squishy? Does it glow? Anyway, good job finding it! Bug report accepted!`

	// For managing pull requests
	PROpened = `Hey @%s! Thank you for opening a pull request. Make sure to tag the 
	maintainers and link the issue you are solving in your pull-request description 
	if you have not already and they'll review as soon as possible.`

	PRMerged = `Congratulations @%s! Our systems are going to start turning their 
	wheels and cogs to compute your scores. Make sure to check your profile for 
	any new achievements at [amsoc.vercel.app](https://amsoc.vercel.app/profile).`

	// Bounties and penalties
	BountyDelivered = `Another day, another coin. Way to get that bounty @%s you 
	glorious keyboard-tapping, coffee-sipping, vibe-coding witcher!`

	PenaltyDelivered = `Looks like someone took an L today. Chin up, buttercup! 
	There is always a next time to, NOT get a penalty :')`

	// Unauthorized comments (participants trying to run maintainer bot-commands)
	PretendMaintainer = `Buddy, that's a maintainers ONLY command! Your 'access' is 
	as real as my long-term relationship goals. Which is completely non-existant :)
	Bless your heart for trying LOL!`

	UnregisteredParticipant = `Oopsie @%s! Someone forgot to register for the 
	Amrita Summer of Code, 2025 program aye ? Register then login via GitHub first.
	Head over to the [website](https://amsoc.vercel.app/register).`
)

func PostComment(ctx context.Context, client *github.Client, owner, repo string, issueNumber int, comment string) error {

}
