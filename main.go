package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"hyprtime/internal/gui/service"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	screenTimeService := service.NewScreenTimeService()
	defer screenTimeService.Close()

	app := application.New(application.Options{
		Name: "hyprtime",
		Services: []application.Service{
			application.NewService(screenTimeService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "hyprtime",
		BackgroundColour: application.NewRGB(255, 255, 255),
		Hidden:           true,
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
