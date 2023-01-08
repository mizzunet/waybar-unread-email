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
	Text       string   `json:"text"`
	Tooltip    string   `json:"tooltip,omitempty"`
	Class      []string `json:"class,omitempty"`
	Percentage int      `json:"percentage"`
}

var (
	OutputFmt string

	TotalUnread int
	Tooltips    []string
	Errors      []string

	Cfg = &Config{}
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		Errors = append(Errors, fmt.Sprintf("error finding user's home directory: %v", err))
		printAndExit()
	}

	var confPath string
	flag.StringVar(&confPath, "config", homeDir+"/.config/waybar-unread-email/config.yaml", "path to YAML file containing configuration")
	flag.StringVar(&OutputFmt, "output", "json", "output format, one of: \"json\", \"yaml\", \"num\"")
	flag.Parse()

	content, err := os.ReadFile(confPath)
	if err != nil {
		Errors = append(Errors, fmt.Sprintf("error reading config file: %v", err))
		printAndExit()
	}

	err = yaml.Unmarshal(content, Cfg)
	if err != nil {
		Errors = append(Errors, fmt.Sprintf("error unmarshaling config file: %v", err))
		printAndExit()
	}

	for _, srv := range Cfg.Servers {
		var err error
		var imapClt *client.Client
		enc := strings.ToLower(srv.Encryption)

		if enc == "tls" {
			imapClt, err = client.DialTLS(srv.Address, &tls.Config{InsecureSkipVerify: srv.SkipVerify})
		} else {
			imapClt, err = client.Dial(srv.Address)
		}
		if err != nil {
			Errors = append(Errors, fmt.Sprintf("error dialing server '%s': %v", srv.Name, err))
			continue
		}
		defer imapClt.Logout()

		if enc == "starttls" {
			err = imapClt.StartTLS(&tls.Config{InsecureSkipVerify: srv.SkipVerify})
			if err != nil {
				Errors = append(Errors, fmt.Sprintf("error starting TLS for server '%s': %v", srv.Name, err))
				continue
			}
		}

		if err := imapClt.Login(srv.Username, srv.Password); err != nil {
			Errors = append(Errors, fmt.Sprintf("error logging into server '%s': %v", srv.Name, err))
			continue
		}

		_, err = imapClt.Select("INBOX", true)
		if err != nil {
			Errors = append(Errors, fmt.Sprintf("error selecting INBOX on server '%s': %v", srv.Name, err))
			continue
		}

		// Criteria for unread emails.
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}
		seqNums, err := imapClt.Search(criteria)
		if err != nil {
			Errors = append(Errors, fmt.Sprintf("error searching INBOX on server '%s': %v", srv.Name, err))
			continue
		}

		count := len(seqNums)
		Tooltips = append(Tooltips, fmt.Sprintf("%s: %d unread", srv.Name, count))
		TotalUnread += count
	}

	printAndExit()
}

func printAndExit() {
	wbOut := &WaybarOutput{}
	wbOut.Tooltip = strings.Join(append(Tooltips, Errors...), "\n")
	if TotalUnread > 0 {
		wbOut.Text = strconv.Itoa(TotalUnread)
		wbOut.Class = []string{"unread"}
		wbOut.Percentage = 100
	} else if Cfg.ShowZero {
		wbOut.Text = "0"
	}
	if len(Errors) > 0 {
		wbOut.Class = append(wbOut.Class, "error")
		os.Stderr.WriteString(strings.Join(Errors, "\n") + "\n")
	}

	out := ""
	switch strings.ToLower(OutputFmt) {
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
