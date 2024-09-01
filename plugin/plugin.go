package plugin

import (
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"os/exec"
	"path/filepath"
)

// Load all plugins from the plugin directory
func Load() error {
	dir := config.PluginsDir.Value()
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range dirEntries {
		pluginPath := filepath.Join(dir, entry.Name())
		if !entry.Type().IsRegular() {
			logging.Debug.Printf("Skipping non-regular file in plugins directory: %s", pluginPath)
			continue
		}
		go run(pluginPath)
	}

	return nil
}

func run(pluginPath string) {
	plugin := exec.Command(pluginPath)
	plugin.Stderr = os.Stderr
	writer, err := plugin.StdinPipe()
	if err != nil {
		return
	}
	reader, err := plugin.StdoutPipe()
	if err != nil {
		return
	}
	if err := plugin.Start(); err != nil {
		logging.Error.Printf("Error starting plugin %s with error: %s", pluginPath, err)
		return
	}
	// send version
	versionHeader := header{
		sequence:    0,
		topicLength: uint16(len(TopicVersion)),
		bodyLength:  0,
		reserved:    0,
		topic:       TopicVersion,
		body:        nil,
	}
	if err := versionHeader.write(writer); err != nil {
		return
	}
	versionResponse, err := read(reader)
	if err != nil {
		return
	}
	var version Version
	if _, err := version.UnmarshalMsg(versionResponse.body); err != nil {
		return
	}

	logging.Debug.Printf("Version Response for %s: %s", pluginPath, version.Version)
	// send register
	if err := plugin.Wait(); err != nil {
		return
	}
}
