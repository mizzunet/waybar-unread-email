package main

import (
    "os"
	"log"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
    "flag"
    "fmt"
)
var (
    SERVER *string
    USER *string
    PASS *string
)
func init() {
    SERVER   = flag.String("S","","Server")
    USER     = flag.String("U","","Username")
    PASS     = flag.String("P","","Password")
}

func main() {
    const ICON_UNREAD="󰗰"

    flag.Parse()

	// Connect to server
	c, err := client.Dial(*SERVER)
	if err != nil {
		log.Fatal(err)
	}
    
	// Login
	if err := c.Login(*USER, *PASS); err != nil {
		log.Fatal(err)
	}

    _, err = c.Select("INBOX", false)
    if err !=nil {
        log.Fatal(err)
    }

    // Set search criteria
    criteria := imap.NewSearchCriteria()
    criteria.WithoutFlags = []string{imap.SeenFlag}
    COUNT, err := c.Search(criteria)
    if err != nil {
        log.Fatal(err)
    }

    if len(COUNT) > 0 {
    fmt.Println(`{"text":" `,len(COUNT), ICON_UNREAD ,`", "class": ["unread"]}` )
    } else {
        os.Exit(0)
    }

	defer c.Logout()
}
