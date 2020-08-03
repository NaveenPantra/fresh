package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedExt(ev.Name) && !ev.IsAttrib() {
					if isDebug() {
						watcherLog("sending event %s", ev)
					}
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("Error: %s", err)
			}
		}
	}()

	ppath := path
	if p, err := filepath.Rel("./", path); err != nil {
		ppath = p
	}

	watcherLog("Watching %s", ppath)

	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watch() {
	root := root()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			if isIgnored(path) {
				if isDebug() {
					watcherLog("Ignoring %s", path)
				}
				// Not recursively ignoring anymore
				//return filepath.SkipDir
			} else {
				watchFolder(path)
			}
		}

		return err
	})
}
