package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// create empty directory to put build outputs
	outputs := client.Directory()

    m2 := client.CacheVolume("m2")

	// get `maven` image
	maven := client.Container().From("maven:3-eclipse-temurin-17").WithMountedCache("~/.m2", m2)

	// mount cloned repository into `maven` image
	maven = maven.WithDirectory("/src", src).WithWorkdir("/src")

	out, err := maven.WithExec([]string{"mvn", "--batch-mode", "verify"}).Stdout(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)

	// write build artifacts to host
	_, err = outputs.Export(ctx, ".")
	if err != nil {
		return err
	}

	return nil
}
