package main

import (
	"context"
	"fmt"
	"time"
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
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// CheckForUpdate polls Github to see if an update is available
func (a *App) CheckForUpdate() UpdateStatus {
	return UpdateStatus{
		UpdateAvailable: true,
		RemoteVersion:   "2.0.0-alpha1",
	}
}

func (a *App) PerformUpdate() bool {
	time.Sleep(time.Second * 3)
	return true
}

func (a *App) GetCurrentVersion() string {
	return a.version
}
