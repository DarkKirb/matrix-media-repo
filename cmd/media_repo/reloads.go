package main

import (
	"github.com/turt2live/matrix-media-repo/api/webserver"
	"github.com/turt2live/matrix-media-repo/common/globals"
	"github.com/turt2live/matrix-media-repo/metrics"
	"github.com/turt2live/matrix-media-repo/storage"
)

func setupReloads() {
	reloadWebOnChan(globals.WebReloadChan)
	reloadMetricsOnChan(globals.MetricsReloadChan)
	reloadDatabaseOnChan(globals.DatabaseReloadChan)
	reloadDatastoresOnChan(globals.DatastoresReloadChan)
}

func stopReloads() {
	// send stop signal to reload fns
	globals.WebReloadChan <- false
	globals.MetricsReloadChan <- false
	globals.DatabaseReloadChan <- false
	globals.DatastoresReloadChan <- false
}

func reloadWebOnChan(reloadChan chan bool) {
	go func() {
		for {
			shouldReload := <-reloadChan
			if shouldReload {
				webserver.Reload()
			} else {
				return // received stop
			}
		}
	}()
}

func reloadMetricsOnChan(reloadChan chan bool) {
	go func() {
		for {
			shouldReload := <-reloadChan
			if shouldReload {
				metrics.Reload()
			} else {
				return // received stop
			}
		}
	}()
}

func reloadDatabaseOnChan(reloadChan chan bool) {
	go func() {
		for {
			shouldReload := <-reloadChan
			if shouldReload {
				storage.ReloadDatabase()
				loadDatabase()
				globals.DatastoresReloadChan <- true
			} else {
				return // received stop
			}
		}
	}()
}

func reloadDatastoresOnChan(reloadChan chan bool) {
	go func() {
		for {
			shouldReload := <-reloadChan
			if shouldReload {
				loadDatastores()
			} else {
				return // received stop
			}
		}
	}()
}
