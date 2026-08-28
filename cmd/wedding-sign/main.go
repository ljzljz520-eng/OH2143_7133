package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"wedding-sign/internal/config"
	"wedding-sign/internal/service"
	"wedding-sign/internal/store"
	"wedding-sign/internal/workflow"
)

func main() {
	base := flag.String("data", ".data", "data directory")
	demo := flag.Bool("demo", false, "create a demo welcome record")
	flag.Parse()
	settings := config.DefaultSettings(*base)
	if err := settings.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := config.EnsureDataDirectory(settings); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	db, err := store.Open(filepath.Clean(settings.DataPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()
	svc, err := service.New(db, settings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if *demo {
		result, runErr := workflow.AcceptWelcome(svc, "ceremony", "Lin & Kai", "2026-08-24", "Jasmine Hall", "")
		if runErr != nil {
			fmt.Fprintln(os.Stderr, runErr)
			return
		}
		fmt.Printf("%s | %s | %s\n", result.Record.DisplayTitle(), result.Record.Date, result.Record.Venue)
		return
	}
	fmt.Printf("wedding welcome sign ready at %s\n", settings.DataPath)
}
