package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andygrunwald/go-jira"

	"github.com/urfave/cli/v2"
)

var (
	flagToken = cli.StringFlag{
		Name:     "token",
		EnvVars:  []string{"JIRA_TOKEN"},
		Usage:    "JIRA API Token",
		Required: true,
	}
	flagURL = cli.StringFlag{
		Name:     "url",
		EnvVars:  []string{"JIRA_URL"},
		Usage:    "JIRA base url",
		Required: true,
	}
	flagUser = cli.StringFlag{
		Name:     "user",
		EnvVars:  []string{"JIRA_USER"},
		Usage:    "JIRA user name",
		Required: true,
	}
	flagJQL = cli.StringFlag{
		Name:     "jql",
		Aliases:  []string{"q"},
		Usage:    "jira query expression",
		Required: true,
	}
)

func main() {

	app := &cli.App{
		Name:  "jira client",
		Usage: "allows querying jira using JQL via commandline",
		Flags: []cli.Flag{
			&flagJQL,
			&flagToken,
			&flagURL,
			&flagUser,
		},
		Action: handler,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func handler(c *cli.Context) error {
	tp := jira.BasicAuthTransport{
		Username: c.String(flagUser.Name),
		Password: c.String(flagToken.Name),
	}

	client, err := jira.NewClient(tp.Client(), "https://livesport.atlassian.net")
	if err != nil {
		return fmt.Errorf("as error occurred initializing jira client: %s", err)

	}

	opts := &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 10,
	}

	jql := c.String(flagJQL.Name)
	issues, _, err := client.Issue.Search(jql, opts)
	if err != nil {
		return fmt.Errorf("as error occurred runnig JQL '%s': %s", jql, err)
	}

	for _, issue := range issues {
		fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
		fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
		fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	}

	return nil
}
