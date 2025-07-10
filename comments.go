package main

import (
	"context"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

const (
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

	Tester = `Alright, who's the genius dropping test cases into our code? You @%s! 
	Nice! Seriously, thanks for making sure this whole thing doesn't blow up. 
	You're the real MVP, you magnificent creature.`

	Helper = `Hot damn @%s! Someone's got their Good Samaritan pants on today! You 
	actually went out of your way to lend a hand, and for that, you get a gold star 
	(and maybe a new shiny badge, if you're lucky). Thanks for being so helpful!`

	PROpened = `Hey @%s! Thank you for opening a pull request. Make sure to tag the 
	maintainers and link the issue you are solving in your pull-request description 
	if you have not already and they'll review as soon as possible.`

	PRMerged = `Congratulations @%s! Our systems are going to start turning their 
	wheels and cogs to compute your scores. Make sure to check your profile for 
	any new achievements at [amsoc.vercel.app](https://amsoc.vercel.app/profile).`

	BountyDelivered = `Another day, another coin. Way to get that bounty @%s you 
	glorious keyboard-tapping, coffee-sipping, vibe-coding witcher!`

	PenaltyDelivered = `Looks like someone took an L today. Chin up, buttercup! 
	There is always a next time to, NOT get a penalty :')`
)

func postComment(ctx context.Context, installationToken, owner, repo string,
	number int, cmt string, sink string) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: installationToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if sink == "issue" {
		comment := &github.IssueComment{
			Body: github.String(cmt),
		}
		_, _, err := client.Issues.CreateComment(ctx, owner, repo, number, comment)
		if err != nil {
			return err
		}
	} else {
		// This handles pull request comments
		comment := &github.PullRequestComment{
			Body: github.String(cmt),
		}

		_, _, err := client.PullRequests.CreateComment(ctx, owner, repo, number, comment)
		if err != nil {
			return err
		}
	}
	return nil
}
