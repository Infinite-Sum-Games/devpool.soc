package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IAmRiteshKoushik/devpool/db"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

type GitHubInfo struct {
	RepoOwner string
	RepoName  string
	Type      string
	Number    int
}

type InstallationTokenResp struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

const (
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
	and grab yourself a latte @%s. You have earned it!`

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

	PenaltyDelivered = `Looks like someone took an L today. Chin up, @%s! 
	There is always a next time to, NOT get a penalty :')`

	ExtensionGranted = `Alright @%s! I'm giving you an extension. Don't screw it 
	up, or I'm telling Cable you need a dial-up.`
)

func postComment(ctx context.Context, installationToken, owner, repo string,
	number int, cmt string, sink string) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: installationToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if sink == "issues" {
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

func ParseGitHubURL(rawURL string) (*GitHubInfo, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	if parsedURL.Host != "github.com" {
		return nil, fmt.Errorf("not a github.com URL: %s", rawURL)
	}

	// Split the path into components
	// ["Infinite-Sum-Games", "platform.soc", "issues", "62"]
	pathParts := strings.Split(parsedURL.Path, "/")

	// Format: /owner/repo/issues(or pull)/number
	if len(pathParts) < 4 {
		return nil, fmt.Errorf("invalid GitHub issue/pull URL format: not enough path components in %s", rawURL)
	}

	// Validate the "issues" and "pull" part
	if pathParts[2] != "issues" && pathParts[2] != "pull" {
		return nil, fmt.Errorf("invalid GitHub issue URL format: expected 'issues/pull' but got '%s' in %s", pathParts[2], rawURL)
	}

	// Extract owner, repo, and number
	owner := pathParts[0]
	repo := pathParts[1]
	issueNumStr := pathParts[3]

	// Convert issue number string to integer
	issueNum, err := strconv.Atoi(issueNumStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert number '%s' to integer: %w",
			issueNumStr,
			err,
		)
	}

	return &GitHubInfo{
		RepoOwner: owner,
		RepoName:  repo,
		Type:      pathParts[2],
		Number:    issueNum,
	}, nil
}

func NewInstallationToken(repoUrl string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get connection from pool: %w", err)
	}
	defer conn.Release()

	q := db.New()
	installationID, err := q.GetInstallationIdQuery(ctx, conn, repoUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get installation token from db: %w", err)
	}

	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.RegisteredClaims{
			Issuer:    fmt.Sprintf("%d", App.AppID),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-60 * time.Second)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	)
	jwtString, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Create an API request
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens",
		installationID.Int64)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtString)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request installation token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d: %s",
			resp.StatusCode, string(body))
	}

	var tokenResp InstallationTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Re-populate cache in-case is does not exist
	hash := "repo-token:" + repoUrl
	value := tokenResp.Token
	ttl := time.Until(tokenResp.ExpiresAt)
	if err := Valkey.Set(ctx, hash, value, ttl).Err(); err != nil {
		// Don't cancel the operation because Redis cache failed
		Log.Error("failed to store token in cache", err)
	}

	return tokenResp.Token, nil
}

func FetchInstallationToken(repoUrl string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	token, err := Valkey.Get(ctx, "repo-token:"+repoUrl).Result()
	if err != nil {
		return "", err // Might be Redis.Nil if token does not exist
	}
	return token, nil
}
