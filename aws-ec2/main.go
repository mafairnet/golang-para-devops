package main

import (
	"context"
	"fmt"
	"os"
)

var configuration = getProgramConfiguration()

func main() {

	var (
		instanceId string
		err        error
	)

	if instanceId, err = createEC2(context.Background(), configuration.DefaulRegion); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("ID de la Instancia: %s\n", instanceId)
}
