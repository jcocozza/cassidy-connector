package auth

import (
	"fmt"
	"os/exec"
	"runtime"

	//config "github.com/jcocozza/cassidy-connector/strava/internal"
	//"github.com/jcocozza/cassidy-connector/strava/internal/auth"
)

// Will open a link in the browser
func openURL(url string) error {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "linux":
        cmd = exec.Command("xdg-open", url)
    case "darwin":
        cmd = exec.Command("open", url)
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", url)
    default:
        return fmt.Errorf("unsupported operating system")
    }

    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}
/*
// open the approval url in browser
func InitialAuthorizationDirect() {
	approvalUrl := auth.GenerateApprovalUrl(config.ClientId, "http://localhost/exchange_token", "activity:read_all")
	openURL(approvalUrl)
}
*/