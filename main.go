package main

import(
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/mattermost/platform/model"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	url := "http://localhost:8065"
	username := "user-1"
	password := "SampleUs@r-1"
	execTime := time.Now().UTC()

	c := model.NewAPIv4Client(url)
	u, _ := c.Login(username, password)

	teams, _ := c.GetAllTeams("", 0, 100)

	fmt.Println("mmstats")
	fmt.Println("Executed", execTime.Format(time.RFC1123))
	for _, t := range teams {
		printTeamStats(c, execTime, u, t)
	}
}

func printTeamStats(c *model.Client4, execTime time.Time, u *model.User, t *model.Team) {
	fmt.Println("---")
	fmt.Println("Team:", t.Name)

	weeklyCountByType := map[string]int{} // Active channels past week by type
	weeklyNamesByType := map[string][]string{} // Active channel names past week by type
	monthlyCountByType := map[string]int{}
	monthlyNamesByType := map[string][]string{}

	chs, _ := c.GetChannelsForTeamForUser(t.Id, u.Id, "")

	for _, ch := range chs {
		// Find the best name
		name := ch.DisplayName
		if ch.Type == "D" {
			name = ch.Name
			name = strings.ReplaceAll(name, u.Id, "")
			name = strings.ReplaceAll(name, "__", "")

			// Ignore channels with myself
			if name == "" {
				continue
			}

			otherUser, _ := c.GetUser(name, "")
			name = otherUser.Username
		}

		// Past week
		if ch.LastPostAt >= execTime.AddDate(0, 0, -7).Unix() * 1000 {
			// Count
			if count, ok := weeklyCountByType[ch.Type]; ok {
				weeklyCountByType[ch.Type] = count + 1
			} else {
				weeklyCountByType[ch.Type] = 1
			}

			// Name
			if names, ok := weeklyNamesByType[ch.Type]; ok {
				weeklyNamesByType[ch.Type] = append(names, name)
			} else {
				weeklyNamesByType[ch.Type] = []string{name}
			}
		}

		// Past month
		if ch.LastPostAt >= execTime.AddDate(0, -1, 0).Unix() * 1000 {
			// Count
			if count, ok := monthlyCountByType[ch.Type]; ok {
				monthlyCountByType[ch.Type] = count + 1
			} else {
				monthlyCountByType[ch.Type] = 1
			}

			// Name
			if names, ok := monthlyNamesByType[ch.Type]; ok {
				monthlyNamesByType[ch.Type] = append(names, name)
			} else {
				monthlyNamesByType[ch.Type] = []string{name}
			}
		}
	}

	fmt.Println("Past week:")
	printNamesByType(weeklyNamesByType, weeklyCountByType)
	fmt.Println("Past month:")
	printNamesByType(monthlyNamesByType, monthlyCountByType)
}

func printNamesByType(names map[string][]string, counts map[string]int) {
	for k, v := range names {
		fmt.Println("  - ", k, counts[k])
		for _, n := range v {
			fmt.Println("    - ", n)
		}
	}
}

func credentials() (string, string) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter Username: ")
    username, _ := reader.ReadString('\n')

    fmt.Print("Enter Password: ")
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err == nil {
        fmt.Println("\nPassword typed: " + string(bytePassword))
    }
    password := string(bytePassword)

    return strings.TrimSpace(username), strings.TrimSpace(password)
}
