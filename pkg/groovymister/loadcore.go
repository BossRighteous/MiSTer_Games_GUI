package groovymister

import (
	"fmt"
	"os"
	"strings"
)

func LaunchGroovyCore(path string) error {
	// path = /media/fat/_Utility/Groovy_x.rbf
	cmdInterface := "/dev/MiSTer_cmd"

	_, err := os.Stat(cmdInterface)
	if err != nil {
		return fmt.Errorf("command interface not accessible: %s", err)
	}

	if !(strings.HasSuffix(strings.ToLower(path), ".rbf")) {
		return fmt.Errorf("not a valid launch file: %s", path)
	}

	cmd, err := os.OpenFile(cmdInterface, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer cmd.Close()

	cmd.WriteString(fmt.Sprintf("load_core %s\n", path))
	return nil
}
