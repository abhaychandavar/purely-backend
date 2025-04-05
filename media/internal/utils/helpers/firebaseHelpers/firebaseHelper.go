package firebaseHelper

import (
	"context"
	"fmt"
	"log"
	"media/internal/config"
	"sync"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	app  *firebase.App
	once sync.Once
)

func App() *firebase.App {
	once.Do(func() {
		opt := option.WithCredentialsFile(config.GetConfig().FirebaseConfigPath)
		firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)

		if err != nil {
			panic(fmt.Errorf("error initializing app: %v", err))
		}
		log.Default().Printf("Firebase app initialized")
		app = firebaseApp
	})
	return app
}
