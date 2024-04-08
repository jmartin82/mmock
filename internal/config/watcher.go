package config

import (
	"github.com/radovskyb/watcher"
	"time"
)

// Watcher interface contains the function watching process
type Watcher interface {
	Bind()
	UnBind()
}

type FileWatcher struct {
	watcher  *watcher.Watcher
	fsUpdate chan struct{}
	path     string
}

func NewFileWatcher(path string, fsUpdate chan struct{}) *FileWatcher {
	return &FileWatcher{path: path, fsUpdate: fsUpdate}
}

func (fw *FileWatcher) UnBind() {
	if fw.watcher != nil {
		fw.watcher.Close()
	}
}

// Bind start the watching process to detect any change on defintions
func (fw *FileWatcher) Bind() {

	fw.watcher = watcher.New()

	// SetMaxEvents to 1 to allow at most 1 Event to be received
	fw.watcher.SetMaxEvents(1)
	go func() {
		for {
			select {
			case event := <-fw.watcher.Event:
				log.Infof("Changes detected in mock definitions ", event.String())
				fw.fsUpdate <- struct{}{}
			case err := <-fw.watcher.Error:
				log.Errorf("File monitor error", err)
			}
		}
	}()

	// Watch dir recursively for changes.
	if err := fw.watcher.AddRecursive(fw.path); err != nil {
		log.Errorf("Impossible bind the config folder to the files monitor: ", err)
		return
	}

	go func() {
		log.Info("File monitor started")
		if err := fw.watcher.Start(time.Millisecond * 100); err != nil {
			log.Errorf("Impossible to start the config files monitor: ", err)
		}
	}()

}
