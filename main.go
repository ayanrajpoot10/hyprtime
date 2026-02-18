package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {

	screenTimeService := NewScreenTimeService()
	defer screenTimeService.Close()

	app := application.New(application.Options{
		Name:        "hyprtime",
		Description: "Hyprland Screen Time Viewer",
		Services: []application.Service{
			application.NewService(screenTimeService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Hypr Screen Time",
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
		Width:            1200,
		Height:           800,
		Hidden:           true,
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
