package model

import (
	"fmt"
	"strings"
)

type Message struct {
	ResetTime    string `json:"reset_time"`
	ResetReason  string `json:"reset_reason"`
	FwVer        string `json:"fw_ver"`
	ElfSha256    string `json:"elf_sha256"`
	FreeHeap     int    `json:"free_heap"`
	MinHeapBlock int    `json:"min_heap_block"`
	PanicReason  string `json:"panic_reason"`
	StackOvfTask string `json:"stack_ovf_task"`
	Core         int    `json:"core"`
	PanicDesc    string `json:"panic_desc"`
	Pc           string `json:"PC"`
	Backtrace    string `json:"backtrace"`
}

type Report struct {
	Message  Message `json:"message"`
	DeviceID string  `json:"deviceID"`
}

func (r *Report) GetReport() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "<b>DeviceID:</b> %s\n", r.DeviceID)

	return sb.String()
}
