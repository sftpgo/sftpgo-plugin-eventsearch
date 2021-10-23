package db

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	driver := os.Getenv("SFTPGO_PLUGIN_EVENTSEARCH_DRIVER")
	dsn := os.Getenv("SFTPGO_PLUGIN_EVENTSEARCH_DSN")
	if driver == "" || dsn == "" {
		fmt.Println("Driver and/or DSN not set, unable to execute test")
		os.Exit(1)
	}
	if err := Initialize(driver, dsn); err != nil {
		fmt.Printf("unable to initialize database: %v\n", err)
		os.Exit(1)
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	err := sess.AutoMigrate(&ProviderEvent{}, &FsEvent{})
	if err != nil {
		fmt.Printf("unable to migrate database: %v\n", err)
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
