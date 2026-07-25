package main

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "moonBASIC IDE",
		Width:            1280,
		Height:           800,
		MinWidth:         900,
		MinHeight:        600,
		DisableResize:    false,
		Fullscreen:       false,
		Frameless:        false,
		StartHidden:      false,
		HideWindowOnClose: false,
		BackgroundColour: &options.RGBA{R: 10, G: 12, B: 16, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		Logger:           nil,
		LogLevel:         0,
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
		// Platform-specific options
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "moonBASIC IDE",
				Message: "IDE for moonBASIC .mb source\nUses moonbasic + moonrun toolchain",
			},
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			IsZoomControlEnabled:              false,
			EnableSwipeGestures:               false,
			BackdropType:                      windows.Mica,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               filepath.Join(os.TempDir(), "moonbasic-ide"),
			Theme: windows.Dark,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(10, 12, 16),
				DarkModeTitleText:  windows.RGB(0, 212, 255),
				DarkModeBorder:     windows.RGB(30, 58, 95),
				LightModeTitleBar:  windows.RGB(10, 12, 16),
				LightModeTitleText: windows.RGB(0, 212, 255),
				LightModeBorder:    windows.RGB(30, 58, 95),
			},
		},
		Linux: &linux.Options{
			ProgramName: "moonbasic-ide",
			WindowIsTranslucent: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
