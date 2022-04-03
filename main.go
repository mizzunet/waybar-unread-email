package main

import (
    "os"
//	"log"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
    "flag"
    "fmt"
)
var (
    SERVER *string
    USER *string
    PASS *string
    ICON *string
)
func init() {
    SERVER   = flag.String("S","","Server")
    USER     = flag.String("U","","Username")
    PASS     = flag.String("P","","Password")
    ICON     = flag.String("I","󰗰","Icon")

}

func main() {
    //const ICON_UNREAD="󰗰"

    flag.Parse()

	// Connect to server
	c, err := client.Dial(*SERVER)
	if err != nil {
//		log.Fatal(err)
        fmt.Println(err)
        os.Exit(0)
	}
    
	// Login
	if err := c.Login(*USER, *PASS); err != nil {
//		log.Fatal(err)
        fmt.Println(err)
        os.Exit(0)
	}

    _, err = c.Select("INBOX", false)
    if err !=nil {
//      log.Fatal(err)
        fmt.Println(err)
        os.Exit(0)
    }

    // Set search criteria
    criteria := imap.NewSearchCriteria()
    criteria.WithoutFlags = []string{imap.SeenFlag}
    COUNT, err := c.Search(criteria)
    if err != nil {
//     log.Fatal(err)
        fmt.Println(err)
        os.Exit(0)
    }

    if len(COUNT) > 0 {
    fmt.Println(`{"text":" `,len(COUNT), ICON ,`", "class": ["unread"]}` )
    } else {
        os.Exit(0)
    }

	defer c.Logout()
}
