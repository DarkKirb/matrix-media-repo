package config

import (
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-media-repo/common/globals"
)

func Watch() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}

	err = watcher.Add(Path)
	if err != nil {
		logrus.Fatal(err)
	}

	go func() {
		debounced := debounce.New(1 * time.Second)
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				debounced(onFileChanged)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Error("error in config watcher:", err)
			}
		}
	}()

	return watcher
}

func onFileChanged() {
	logrus.Info("Config file change detected - reloading")
	configNow := Get()
	configNew, err := reloadConfig()
	if err != nil {
		logrus.Error("Error reloading configuration - ignoring")
		logrus.Error(err)
		return
	}

	logrus.Info("Applying reloaded config live")
	instance = configNew

	bindAddressChange := configNew.General.BindAddress != configNow.General.BindAddress
	bindPortChange := configNew.General.Port != configNow.General.Port
	forwardAddressChange := configNew.General.TrustAnyForward != configNow.General.TrustAnyForward
	forwardedHostChange := configNew.General.UseForwardedHost != configNow.General.UseForwardedHost
	if bindAddressChange || bindPortChange || forwardAddressChange || forwardedHostChange {
		logrus.Warn("Webserver configuration changed - remounting")
		globals.WebReloadChan <- true
	}

	metricsEnableChange := configNew.Metrics.Enabled != configNow.Metrics.Enabled
	metricsBindAddressChange := configNew.Metrics.BindAddress != configNow.Metrics.BindAddress
	metricsBindPortChange := configNew.Metrics.Port != configNow.Metrics.Port
	if metricsEnableChange || metricsBindAddressChange || metricsBindPortChange {
		logrus.Warn("Metrics configuration changed - remounting")
		globals.MetricsReloadChan <- true
	}

	databaseChange := configNew.Database.Postgres != configNow.Database.Postgres
	poolConnsChange := configNew.Database.Pool.MaxConnections != configNow.Database.Pool.MaxConnections
	poolIdleChange := configNew.Database.Pool.MaxIdle != configNow.Database.Pool.MaxIdle
	if databaseChange || poolConnsChange || poolIdleChange {
		logrus.Warn("Database configuration changed - reconnecting")
		globals.DatabaseReloadChan <- true
	}

	logChange := configNew.General.LogDirectory != configNow.General.LogDirectory
	if logChange {
		logrus.Warn("Log configuration changed - restart the media repo to apply changes")
	}

	// Always update the datastores
	logrus.Warn("Updating datastores to ensure accuracy")
	globals.DatastoresReloadChan <- true
}
