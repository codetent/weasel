package plugin

import "github.com/codetent/weasel/pkg/weasel/plugin/plugins"

var availablePlugins = []Plugin{
	&plugins.GitPlugin{},
}

func EnterAll(name string) error {
	for _, plugin := range availablePlugins {
		err := plugin.Enter(name)
		if err != nil {
			return err
		}
	}

	return nil
}
