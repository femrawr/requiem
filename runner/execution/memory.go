package execution

import (
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"runner/store"
	"runner/utils"
)

const (
	_IMAGE_DOS_SIGNATURE uint16 = 0x5A4D
	_IMAGE_NT_SIGNATURE  uint32 = 0x00004550

	_IMAGE_FILE_DLL uint16 = 0x2000

	_IMAGE_DIRECTORY_ENTRY_IMPORT       int = 1
	_IMAGE_DIRECTORY_ENTRY_EXCEPTION    int = 3
	_IMAGE_DIRECTORY_ENTRY_BASERELOC    int = 5
	_IMAGE_DIRECTORY_ENTRY_TLS          int = 9
	_IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT int = 13

	_IMAGE_REL_BASED_ABSOLUTE uint16 = 0
	_IMAGE_REL_BASED_HIGH     uint16 = 1
	_IMAGE_REL_BASED_LOW      uint16 = 2
	_IMAGE_REL_BASED_HIGHLOW  uint16 = 3
	_IMAGE_REL_BASED_HIGHADJ  uint16 = 4
	_IMAGE_REL_BASED_DIR64    uint16 = 10

	_IMAGE_SCN_MEM_EXECUTE uint32 = 0x20000000
	_IMAGE_SCN_MEM_READ    uint32 = 0x40000000
	_IMAGE_SCN_MEM_WRITE   uint32 = 0x80000000

	_PAGE_READONLY          uint32 = 0x02
	_PAGE_READWRITE         uint32 = 0x04
	_PAGE_EXECUTE           uint32 = 0x10
	_PAGE_EXECUTE_READ      uint32 = 0x20
	_PAGE_EXECUTE_READWRITE uint32 = 0x40
)

type dosHeader struct {
	E_magic    uint16
	E_cblp     uint16
	E_cp       uint16
	E_crlc     uint16
	E_cparhdr  uint16
	E_minalloc uint16
	E_maxalloc uint16
	E_ss       uint16
	E_sp       uint16
	E_csum     uint16
	E_ip       uint16
	E_cs       uint16
	E_lfarlc   uint16
	E_ovno     uint16
	E_res      [4]uint16
	E_oemid    uint16
	E_oeminfo  uint16
	E_res2     [10]uint16
	E_lfanew   int32
}

type fileHeader struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type dataDirectory struct {
	VirtualAddress uint32
	Size           uint32
}

type optinalHeader32 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	BaseOfData                  uint32
	ImageBase                   uint32
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint32
	SizeOfStackCommit           uint32
	SizeOfHeapReserve           uint32
	SizeOfHeapCommit            uint32
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]dataDirectory
}

type optinalHeader64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]dataDirectory
}

type ntHeader32 struct {
	Signature      uint32
	FileHeader     fileHeader
	OptionalHeader optinalHeader32
}

type ntHeader64 struct {
	Signature      uint32
	FileHeader     fileHeader
	OptionalHeader optinalHeader64
}

type tlsDerectory32 struct {
	StartAddressOfRawData uint32
	EndAddressOfRawData   uint32
	AddressOfIndex        uint32
	AddressOfCallbacks    uint32
	SizeOfZeroFill        uint32
	Characteristics       uint32
}

type tlsDerectory64 struct {
	StartAddressOfRawData uint64
	EndAddressOfRawData   uint64
	AddressOfIndex        uint64
	AddressOfCallbacks    uint64
	SizeOfZeroFill        uint32
	Characteristics       uint32
}

type sectionHeader struct {
	Name                 [8]byte
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLinenumbers uint32
	NumberOfRelocations  uint16
	NumberOfLinenumbers  uint16
	Characteristics      uint32
}

type functionEntry struct {
	BeginAddress                 uint32
	EndAddress                   uint32
	ExceptionHandlerOrUnwindInfo uint32
}

type baseRelocation struct {
	VirtualAddress uint32
	SizeOfBlock    uint32
}

type importDescriptor struct {
	OriginalFirstThunk uint32
	TimeDateStamp      uint32
	ForwarderChain     uint32
	Name               uint32
	FirstThunk         uint32
}

type delayedImportDescriptor struct {
	Attributes                 uint32
	DllNameRVA                 uint32
	ModuleHandleRVA            uint32
	ImportAddressTableRVA      uint32
	ImportNameTableRVA         uint32
	BoundImportAddressTableRVA uint32
	UnloadInformationTableRVA  uint32
	TimeDateStamp              uint32
}

type execContext struct {
	x64         bool
	dll         bool
	data        []byte
	base        uintptr
	ntHeaders32 *ntHeader32
	ntHeaders64 *ntHeader64
	modules     map[string]uintptr
}

func (ctx *execContext) GetDataDir(index int) dataDirectory {
	if ctx.x64 {
		return ctx.ntHeaders64.OptionalHeader.DataDirectory[index]
	}

	return ctx.ntHeaders32.OptionalHeader.DataDirectory[index]
}

func (ctx *execContext) GetSizeOfOptionalHeader() uint16 {
	if ctx.x64 {
		return ctx.ntHeaders64.FileHeader.SizeOfOptionalHeader
	}

	return ctx.ntHeaders32.FileHeader.SizeOfOptionalHeader
}

func (ctx *execContext) GetNumberOfSections() int {
	var sections uint16

	if ctx.x64 {
		sections = ctx.ntHeaders64.FileHeader.NumberOfSections
	} else {
		sections = ctx.ntHeaders32.FileHeader.NumberOfSections
	}

	return int(sections)
}

func (ctx *execContext) GetFirstSection() *sectionHeader {
	var base uintptr

	if ctx.x64 {
		base = uintptr(unsafe.Pointer(&ctx.ntHeaders64.OptionalHeader))
	} else {
		base = uintptr(unsafe.Pointer(&ctx.ntHeaders32.OptionalHeader))
	}

	return (*sectionHeader)(unsafe.Pointer(base + uintptr(ctx.GetSizeOfOptionalHeader())))
}

func (ctx *execContext) GetSectionAt(first *sectionHeader, index int) *sectionHeader {
	offset := uintptr(binary.Size(sectionHeader{}))
	return (*sectionHeader)(unsafe.Add(unsafe.Pointer(first), uintptr(index)*offset))
}

func (ctx *execContext) GetLibrary(name string) (uintptr, error) {
	module, found := ctx.modules[name]
	if found {
		return module, nil
	}

	pointer, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, err
	}

	module, _, err = store.LoadLibrary.Call(uintptr(unsafe.Pointer(pointer)))
	if module == 0 {
		return 0, err
	}

	ctx.modules[name] = module
	return module, nil
}

func (ctx *execContext) HandleRelocations() error {
	var delta int64

	if ctx.x64 {
		delta = int64(ctx.base) - int64(ctx.ntHeaders64.OptionalHeader.ImageBase)
	} else {
		delta = int64(ctx.base) - int64(ctx.ntHeaders32.OptionalHeader.ImageBase)
	}

	if delta == 0 {
		return nil
	}

	dir := ctx.GetDataDir(_IMAGE_DIRECTORY_ENTRY_BASERELOC)
	if dir.VirtualAddress == 0 {
		return nil
	}

	at := ctx.base + uintptr(dir.VirtualAddress)
	end := at + uintptr(dir.Size)

	for at < end {
		block := (*baseRelocation)(unsafe.Pointer(at))
		if block.SizeOfBlock == 0 {
			break
		}

		count := (block.SizeOfBlock - 8) / 2
		entries := unsafe.Slice((*uint16)(unsafe.Add(unsafe.Pointer(at), 8)), count)

		for i := range count {
			entry := entries[i]
			patch := ctx.base + uintptr(block.VirtualAddress) + uintptr(entry&0x0FFF)

			switch entry >> 12 {
			case _IMAGE_REL_BASED_ABSOLUTE:
			case _IMAGE_REL_BASED_HIGH:
				*(*uint16)(unsafe.Pointer(patch)) += uint16(delta >> 16)
			case _IMAGE_REL_BASED_LOW:
				*(*uint16)(unsafe.Pointer(patch)) += uint16(delta & 0xFFFF)
			case _IMAGE_REL_BASED_HIGHLOW:
				*(*uint32)(unsafe.Pointer(patch)) += uint32(delta)
			case _IMAGE_REL_BASED_HIGHADJ:
			case _IMAGE_REL_BASED_DIR64:
				*(*uint64)(unsafe.Pointer(patch)) += uint64(delta)
			}
		}

		at += uintptr(block.SizeOfBlock)
	}

	return nil
}

func (ctx *execContext) HandleImports() error {
	dir := ctx.GetDataDir(_IMAGE_DIRECTORY_ENTRY_IMPORT)
	if dir.VirtualAddress == 0 {
		return nil
	}

	descriptor := (*importDescriptor)(unsafe.Pointer(ctx.base + uintptr(dir.VirtualAddress)))
	for descriptor.Name != 0 {
		library := utils.ReadStringFromMemory(ctx.base + uintptr(descriptor.Name))

		module, err := ctx.GetLibrary(library)
		if err != nil {
			return err
		}

		tableAddress := descriptor.OriginalFirstThunk
		if tableAddress == 0 {
			tableAddress = descriptor.FirstThunk
		}

		lookupTable := ctx.base + uintptr(tableAddress)
		addressTable := ctx.base + uintptr(descriptor.FirstThunk)

		for {
			var function uintptr

			if ctx.x64 {
				entry := *(*uint64)(unsafe.Pointer(lookupTable))
				if entry == 0 {
					break
				}

				if entry&(1<<63) != 0 {
					function, err = store.GetFunctionAddress(module, fmt.Sprintf("#%d", entry&0xFFFF))
				} else {
					name := utils.ReadStringFromMemory(ctx.base + uintptr(uint32(entry)) + 2)
					function, err = store.GetFunctionAddress(module, name)
				}

				if err != nil {
					return err
				}

				*(*uint64)(unsafe.Pointer(addressTable)) = uint64(function)

				lookupTable += 8
				addressTable += 8
			} else {
				entry := *(*uint32)(unsafe.Pointer(lookupTable))
				if entry == 0 {
					break
				}

				if entry&0x80000000 != 0 {
					function, err = store.GetFunctionAddress(module, fmt.Sprintf("#%d", entry&0xFFFF))
				} else {
					name := utils.ReadStringFromMemory(ctx.base + uintptr(entry) + 2)
					function, err = store.GetFunctionAddress(module, name)
				}

				if err != nil {
					return err
				}

				*(*uint32)(unsafe.Pointer(addressTable)) = uint32(function)

				lookupTable += 4
				addressTable += 4
			}
		}

		descriptor = (*importDescriptor)(unsafe.Pointer(
			unsafe.Add(unsafe.Pointer(descriptor), unsafe.Sizeof(importDescriptor{})),
		))
	}

	return nil
}

func (ctx *execContext) HandleLateImports() error {
	dir := ctx.GetDataDir(_IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT)
	if dir.VirtualAddress == 0 {
		return nil
	}

	descriptor := (*delayedImportDescriptor)(unsafe.Pointer(ctx.base + uintptr(dir.VirtualAddress)))
	for descriptor.DllNameRVA != 0 {
		library := utils.ReadStringFromMemory(ctx.base + uintptr(descriptor.DllNameRVA))

		module, err := ctx.GetLibrary(library)
		if err != nil {
			return err
		}

		if descriptor.ModuleHandleRVA != 0 {
			*(*uintptr)(unsafe.Pointer(ctx.base + uintptr(descriptor.ModuleHandleRVA))) = module
		}

		lookupTable := ctx.base + uintptr(descriptor.ImportNameTableRVA)
		addressTable := ctx.base + uintptr(descriptor.ImportAddressTableRVA)

		for {
			var function uintptr

			if ctx.x64 {
				entry := *(*uint64)(unsafe.Pointer(lookupTable))
				if entry == 0 {
					break
				}

				if entry&(1<<63) != 0 {
					function, err = store.GetFunctionAddress(module, fmt.Sprintf("#%d", entry&0xFFFF))
				} else {
					name := utils.ReadStringFromMemory(ctx.base + uintptr(uint32(entry)) + 2)
					function, err = store.GetFunctionAddress(module, name)
				}

				if err != nil {
					return err
				}

				*(*uint64)(unsafe.Pointer(addressTable)) = uint64(function)

				lookupTable += 8
				addressTable += 8
			} else {
				entry := *(*uint32)(unsafe.Pointer(lookupTable))
				if entry == 0 {
					break
				}

				if entry&0x80000000 != 0 {
					function, err = store.GetFunctionAddress(module, fmt.Sprintf("#%d", entry&0xFFFF))
				} else {
					name := utils.ReadStringFromMemory(ctx.base + uintptr(entry) + 2)
					function, err = store.GetFunctionAddress(module, name)
				}

				if err != nil {
					return err
				}

				*(*uint32)(unsafe.Pointer(addressTable)) = uint32(function)

				lookupTable += 4
				addressTable += 4
			}
		}

		descriptor = (*delayedImportDescriptor)(unsafe.Pointer(
			unsafe.Add(unsafe.Pointer(descriptor), unsafe.Sizeof(delayedImportDescriptor{})),
		))
	}

	return nil
}

func (ctx *execContext) HandleProtections() error {
	firstSection := ctx.GetFirstSection()

	for i := 0; i < ctx.GetNumberOfSections(); i++ {
		section := ctx.GetSectionAt(firstSection, i)
		if section.SizeOfRawData == 0 {
			continue
		}

		exec := section.Characteristics&_IMAGE_SCN_MEM_EXECUTE != 0
		read := section.Characteristics&_IMAGE_SCN_MEM_READ != 0
		write := section.Characteristics&_IMAGE_SCN_MEM_WRITE != 0

		var access uint32

		switch {
		case exec && write:
			access = _PAGE_EXECUTE_READWRITE
		case exec && read:
			access = _PAGE_EXECUTE_READ
		case exec:
			access = _PAGE_EXECUTE
		case write:
			access = _PAGE_READWRITE
		default:
			access = _PAGE_READONLY
		}

		var old uint32

		ret, _, err := store.VirtualProtect.Call(
			ctx.base+uintptr(section.VirtualAddress),
			uintptr(section.VirtualSize),
			uintptr(access),
			uintptr(unsafe.Pointer(&old)),
		)

		if ret == 0 {
			return err
		}
	}

	return nil
}

func (ctx *execContext) HandleExceptions() error {
	if !ctx.x64 {
		return nil
	}

	dir := ctx.GetDataDir(_IMAGE_DIRECTORY_ENTRY_EXCEPTION)
	if dir.VirtualAddress == 0 {
		return nil
	}

	table := ctx.base + uintptr(dir.VirtualAddress)
	count := uintptr(dir.Size) / unsafe.Sizeof(functionEntry{})

	ret, _, err := store.AddFunctionTable.Call(table, count, ctx.base)
	if ret == 0 {
		return err
	}

	return nil
}

func (ctx *execContext) HandleTLS() error {
	dir := ctx.GetDataDir(_IMAGE_DIRECTORY_ENTRY_TLS)
	if dir.VirtualAddress == 0 {
		return nil
	}

	var address uintptr

	if ctx.x64 {
		tls := (*tlsDerectory64)(unsafe.Pointer(ctx.base + uintptr(dir.VirtualAddress)))
		address = uintptr(tls.AddressOfCallbacks)
	} else {
		tls := (*tlsDerectory32)(unsafe.Pointer(ctx.base + uintptr(dir.VirtualAddress)))
		address = uintptr(tls.AddressOfCallbacks)
	}

	if address == 0 {
		return nil
	}

	for {
		var callback uintptr

		if ctx.x64 {
			callback = uintptr(*(*uint64)(unsafe.Pointer(address)))
			address += 8
		} else {
			callback = uintptr(*(*uint32)(unsafe.Pointer(address)))
			address += 4
		}

		if callback == 0 {
			break
		}

		syscall.SyscallN(callback, ctx.base, 1, 0)
	}

	return nil
}

func (ctx *execContext) Execute() error {
	var entry uintptr

	if ctx.x64 {
		entry = ctx.base + uintptr(ctx.ntHeaders64.OptionalHeader.AddressOfEntryPoint)
	} else {
		entry = ctx.base + uintptr(ctx.ntHeaders32.OptionalHeader.AddressOfEntryPoint)
	}

	if ctx.dll {
		ret, _, err := syscall.SyscallN(entry, ctx.base, 1, 0)
		if ret == 0 {
			return err
		}
	} else {
		syscall.SyscallN(entry)
	}

	return nil
}

func ExecuteInMemory(bytes []byte) error {
	if len(bytes) < int(unsafe.Sizeof(dosHeader{})) {
		return errors.New("file is too small")
	}

	dos := (*dosHeader)(unsafe.Pointer(&bytes[0]))
	if dos.E_magic != _IMAGE_DOS_SIGNATURE {
		return errors.New("invalid DOS sigiture")
	}

	nt := (*ntHeader32)(unsafe.Pointer(&bytes[dos.E_lfanew]))
	if nt.Signature != _IMAGE_NT_SIGNATURE {
		return errors.New("invalid NT signature")
	}

	context := &execContext{
		data:    bytes,
		x64:     nt.OptionalHeader.Magic == 0x20B,
		modules: make(map[string]uintptr),
	}

	var imageSize uint32
	var headerSize uint32

	if context.x64 {
		context.ntHeaders64 = (*ntHeader64)(unsafe.Pointer(&bytes[dos.E_lfanew]))
		context.dll = context.ntHeaders64.FileHeader.Characteristics&_IMAGE_FILE_DLL != 0

		imageSize = context.ntHeaders64.OptionalHeader.SizeOfImage
		headerSize = context.ntHeaders64.OptionalHeader.SizeOfHeaders
	} else {
		context.ntHeaders32 = nt
		context.dll = context.ntHeaders32.FileHeader.Characteristics&_IMAGE_FILE_DLL != 0

		imageSize = context.ntHeaders32.OptionalHeader.SizeOfImage
		headerSize = context.ntHeaders32.OptionalHeader.SizeOfHeaders
	}

	base, _, err := store.VirtualAlloc.Call(uintptr(imageSize))
	if base == 0 {
		return err
	}

	context.base = base

	copy(unsafe.Slice((*byte)(unsafe.Pointer(base)), headerSize), bytes[:headerSize])

	if context.x64 {
		context.ntHeaders64 = (*ntHeader64)(unsafe.Pointer(base + uintptr(dos.E_lfanew)))
	} else {
		context.ntHeaders32 = (*ntHeader32)(unsafe.Pointer(base + uintptr(dos.E_lfanew)))
	}

	firstSection := context.GetFirstSection()

	for i := range context.GetNumberOfSections() {
		section := context.GetSectionAt(firstSection, i)
		if section.SizeOfRawData == 0 {
			continue
		}

		offset := context.base + uintptr(section.VirtualAddress)
		size := uintptr(section.SizeOfRawData)

		copy(
			unsafe.Slice((*byte)(unsafe.Pointer(offset)), size),
			context.data[section.PointerToRawData:section.PointerToRawData+section.SizeOfRawData],
		)

		if section.VirtualSize > section.SizeOfRawData {
			tail := unsafe.Slice((*byte)(unsafe.Pointer(offset+size)), uintptr(section.VirtualSize)-size)
			clear(tail)
		}
	}

	err = context.HandleRelocations()
	if err != nil {
		return err
	}

	err = context.HandleImports()
	if err != nil {
		return err
	}

	err = context.HandleLateImports()
	if err != nil {
		return err
	}

	err = context.HandleProtections()
	if err != nil {
		return err
	}

	err = context.HandleExceptions()
	if err != nil {
		return err
	}

	err = context.HandleTLS()
	if err != nil {
		return err
	}

	err = context.Execute()
	if err != nil {
		return err
	}

	return nil
}
