package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Will open a link in the browser
func OpenURL(url string) error {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "linux":
        cmd = exec.Command("xdg-open", url)
    case "darwin":
        cmd = exec.Command("open", url)
    case "windows":
        cmd = exec.Command("cmd", "/c", "start","\"" + url + "\"")
    default:
        return fmt.Errorf("unsupported operating system")
    }

    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}
