package matrix

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// MaxRetries for operations
const MaxRetries = 10

// RetryDelay between retries
const RetryDelay = 10 * time.Second

var (
	client     *mautrix.Client
	ignored    map[id.UserID]struct{}
	retriables = []string{"429", "502", "504"}
)

// Init matrix client
func Init(hs, login, password, ignore string) error {
	var err error
	client, err = mautrix.NewClient(hs, "", "")
	if err != nil {
		return err
	}
	err = retry(func() error {
		log.Println("authorizing...")
		_, loginErr := client.Login(context.Background(), &mautrix.ReqLogin{
			Type: "m.login.password",
			Identifier: mautrix.UserIdentifier{
				Type: mautrix.IdentifierTypeUser,
				User: login,
			},
			Password:         password,
			StoreCredentials: true,
		})
		if loginErr != nil {
			fmt.Println("login error:", loginErr.Error())
		}
		return loginErr
	})
	if err != nil {
		return err
	}

	resolveIgnored(ignore)

	return nil
}

// Exit / stop matrix client
func Exit() {
	//nolint // nobody cares at that moment
	retry(func() error {
		log.Println("exiting...")
		_, logoutErr := client.Logout(context.Background())
		return logoutErr
	})
}

func retry(handler func() error) error {
	var err error
	for i := 0; i < MaxRetries; i++ {
		if err = handler(); err != nil {
			if isRetriable(err) {
				log.Println("error:", err.Error(), ", retyring in", RetryDelay.Seconds(), "seconds")
				time.Sleep(RetryDelay)
				continue
			}

			log.Println("error:", err.Error(), ", considered as non-retriable")
			break
		}
		break
	}

	return err
}

func isRetriable(err error) bool {
	for _, sub := range retriables {
		if strings.Contains(err.Error(), sub) {
			return true
		}
	}

	return false
}
