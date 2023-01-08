package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"sigs.k8s.io/yaml"
)

type Config struct {
	ShowZero bool           `json:"showZero"`
	Servers  []ServerConfig `json:"servers,omitempty"`
}

type ServerConfig struct {
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Encryption string `json:"encryption,omitempty"`
	SkipVerify bool   `json:"skipVerify,omitempty"`
}

type WaybarOutput struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip,omitempty"`
	Class      string `json:"class,omitempty"`
	Percentage int    `json:"percentage"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error finding user's home directory: %v", err)
		os.Exit(1)
	}

	var confPath, outputFmt string
	flag.StringVar(&confPath, "config", homeDir+"/.config/waybar-unread-email/config.yaml", "path to YAML file containing configuration")
	flag.StringVar(&outputFmt, "output", "json", "output format, one of: \"json\", \"yaml\", \"num\"")
	flag.Parse()

	content, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		os.Exit(1)
	}

	config := &Config{}
	yaml.Unmarshal(content, config)
	if err != nil {
		fmt.Printf("error unmarshaling config file: %v", err)
		os.Exit(1)
	}

	tooltips := []string{}
	totalUnread := 0
	for _, srv := range config.Servers {
		enc := strings.ToLower(srv.Encryption)
		if enc == "" {
			enc = "no"
		}

		var err error
		var imapClt *client.Client
		if enc == "tls" {
			imapClt, err = client.DialTLS(srv.Address, &tls.Config{InsecureSkipVerify: srv.SkipVerify})
		} else {
			imapClt, err = client.Dial(srv.Address)
		}
		if err != nil {
			fmt.Printf("error dialing server '%s': %v", srv.Name, err)
			os.Exit(1)
		}
		if enc == "starttls" {
			err = imapClt.StartTLS(&tls.Config{InsecureSkipVerify: srv.SkipVerify})
			if err != nil {
				fmt.Printf("error dialing server '%s': %v", srv.Name, err)
				os.Exit(1)
			}
		}
		defer imapClt.Logout()

		if err := imapClt.Login(srv.Username, srv.Password); err != nil {
			fmt.Printf("error logging into server '%s': %v", srv.Name, err)
			os.Exit(1)
		}

		_, err = imapClt.Select("INBOX", true)
		if err != nil {
			fmt.Printf("error selecting INBOX on server '%s': %v", srv.Name, err)
			os.Exit(1)
		}

		// Criteria for unread emails.
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}
		seqNums, err := imapClt.Search(criteria)
		if err != nil {
			fmt.Printf("error searching INBOX for unread messages on server '%s': %v", srv.Name, err)
			os.Exit(1)
		}

		count := len(seqNums)
		tooltips = append(tooltips, fmt.Sprintf("%s: %d unread", srv.Name, count))
		totalUnread += count
	}

	wbOut := &WaybarOutput{Tooltip: strings.Join(tooltips, "\n")}
	if totalUnread > 0 {
		wbOut.Text = strconv.Itoa(totalUnread)
		wbOut.Class = "unread"
		wbOut.Percentage = 100
	} else if config.ShowZero {
		wbOut.Text = "0"
	}

	out := ""
	switch strings.ToLower(outputFmt) {
	case "num":
		out = wbOut.Text
	case "yaml":
		str, _ := yaml.Marshal(wbOut)
		out = string(str)
	default:
		str, _ := json.Marshal(wbOut)
		out = string(str)
	}

	fmt.Println(out)
	os.Exit(0)
}
