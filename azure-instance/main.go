package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var configuration = getProgramConfiguration()

func main() {
	var (
		token  azcore.TokenCredential
		pubKey string
		err    error
	)
	ctx := context.Background()
	subscriptionID := configuration.SubscriptionID
	if len(subscriptionID) == 0 {
		fmt.Printf("No subscription ID was provided")
		os.Exit(1)
	}
	if pubKey, err = generateKeys(); err != nil {
		fmt.Printf("generatekeys error: %s\n", err)
		os.Exit(1)
	}
	if token, err = getToken(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	if err = launchInstance(ctx, subscriptionID, token, &pubKey); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
