package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/marcus-crane/wails-autoupdater/pkg/autoupdater"
)

// App struct
type App struct {
	version string
	ctx     context.Context
}

type UpdateStatus struct {
	UpdateAvailable bool   `json:"update_available"`
	RemoteVersion   string `json:"remote_version"`
}

// NewApp creates a new App application struct
func NewApp(version string) *App {
	return &App{version: version}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	if err := autoupdater.CleanupOldDarwinApp(); err != nil {
		log.Printf("Failed to complete darwin cleanup: %+v", err)
	}
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// CheckForUpdate polls Github to see if an update is available
func (a *App) CheckForUpdate() UpdateStatus {
	updateAvailable, remoteVersion := autoupdater.CheckForNewerVersion(a.version)
	return UpdateStatus{
		UpdateAvailable: updateAvailable,
		RemoteVersion:   remoteVersion,
	}
}

func (a *App) PerformUpdate() bool {
	var updaterFunc func() (bool, error)
	if runtime.GOOS == "windows" {
		updaterFunc = autoupdater.PerformUpdateWindows
	}
	if runtime.GOOS == "darwin" {
		updaterFunc = autoupdater.PerformUpdateDarwin
	}
	if runtime.GOOS == "linux" {
		updaterFunc = autoupdater.PerformUpdateWindows
	}
	success, err := updaterFunc()
	if err != nil {
		log.Print(err)
	}
	return success
}

func (a *App) GetCurrentVersion() string {
	return a.version
}
