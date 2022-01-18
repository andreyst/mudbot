package proxy

import (
	"bytes"
	"fmt"
	"log"
)

type telnetCommand int
type telnetOption int

const TELNET_IAC = 0xFF

const (
	TELNET_CMD_NOP               telnetCommand = 0xF1
	TELNET_CMD_SE                              = 0xF0
	TELNET_CMD_DATAMARK                        = 0xF2
	TELNET_CMD_BREAK                           = 0xF3
	TELNET_CMD_INTERRUPT_PROCESS               = 0xF4
	TELNET_CMD_ABORT_OUTPUT                    = 0xF5
	TELNET_CMD_ARE_YOU_THERE                   = 0xF6
	TELNET_CMD_ERASE_CHARACTER                 = 0xF7
	TELNET_CMD_ERASE_LINE                      = 0xF8
	TELNET_CMD_GO_AHEAD                        = 0xF9
	TELNET_CMD_SB                              = 0xFA
	TELNET_CMD_WILL                            = 0xFB
	TELNET_CMD_WONT                            = 0xFC
	TELNET_CMD_DO                              = 0xFD
	TELNET_CMD_DONT                            = 0xFE
	TELNET_CMD_IAC                             = 0xFF
)

const (
	TELNET_OPT_BINARY_TRANSMISSION                telnetOption = 0x0
	TELNET_OPT_ECHO                                            = 0x1
	TELNET_OPT_RECONNECTION                                    = 0x2
	TELNET_OPT_SUPPRESS_GO_AHEAD                               = 0x3
	TELNET_OPT_APPROX_MESSAGE_SIZE_NEGOTIATION                 = 0x4
	TELNET_OPT_STATUS                                          = 0x5
	TELNET_OPT_TIMING_MARK                                     = 0x6
	TELNET_OPT_REMOTE_CONTROLLED_TRANS_AND_ECHO                = 0x7
	TELNET_OPT_OUTPUT_LINE_WIDTH                               = 0x8
	TELNET_OPT_OUTPUT_PAGE_SIZE                                = 0x9
	TELNET_OPT_OUTPUT_CARRIAGE_RETURN_DISPOSITION              = 0xA
	TELNET_OPT_OUTPUT_HORIZONTAL_TAB_STOPS                     = 0xB
	TELNET_OPT_OUTPUT_HORIZONTAL_TAB_DISPOSITION               = 0xC
	TELNET_OPT_OUTPUT_FORMFEED_DISPOSITION                     = 0xD
	TELNET_OPT_OUTPUT_VERTICAL_TABSTOPS                        = 0xE
	TELNET_OPT_OUTPUT_VERTICAL_TAB_DISPOSITION                 = 0xF
	TELNET_OPT_OUTPUT_LINEFEED_DISPOSITION                     = 0x10
	TELNET_OPT_EXTENDED_ASCII                                  = 0x11
	TELNET_OPT_LOGOUT                                          = 0x12
	TELNET_OPT_BYTE_MACRO                                      = 0x13
	TELNET_OPT_DATA_ENTRY_TERMINAL                             = 0x14
	TELNET_OPT_SUPDUP                                          = 0x15
	TELNET_OPT_SUPDUP_OUTPUT                                   = 0x16
	TELNET_OPT_SEND_LOCATION                                   = 0x17
	TELNET_OPT_TERMINAL_TYPE                                   = 0x18
	TELNET_OPT_END_OF_RECORD                                   = 0x19
	TELNET_OPT_TACACS_USER_IDENTIFICATION                      = 0x1A
	TELNET_OPT_OUTPUT_MARKING                                  = 0x1B
	TELNET_OPT_TERMINAL_LOCATION_NUMBER                        = 0x1C
	TELNET_OPT_TELNET_3270_REGIME                              = 0x1D
	TELNET_OPT_X3_PAD                                          = 0x1E
	TELNET_OPT_NEGOTIATE_ABOUT_WINDOW_SIZE                     = 0x1F
	TELNET_OPT_TERMINAL_SPEED                                  = 0x20
	TELNET_OPT_REMOTE_FLOW_CONTROL                             = 0x21
	TELNET_OPT_LINEMODE                                        = 0x22
	TELNET_OPT_X_DISPLAY_LOCATION                              = 0x23
	TELNET_OPT_MSDP                                            = 0x45
	TELNET_OPT_MSSP                                            = 0x46
	TELNET_OPT_MCCPV1                                          = 0x55
	TELNET_OPT_MCCPV2                                          = 0x56
	TELNET_OPT_MSP                                             = 0x5A
	TELNET_OPT_MXP                                             = 0x5B
	TELNET_OPT_ZMP                                             = 0x5D
	TELNET_OPT_ATCP                                            = 0xC8
	TELNET_OPT_GMCP                                            = 0xC9
	TELNET_OPT_EXTENDED_OPTIONS_LIST                           = 0xFF
)

var compressionStartSequence = []byte{TELNET_IAC, TELNET_CMD_SB, TELNET_OPT_MCCPV2, TELNET_IAC, TELNET_CMD_SE}
var gaSequence = []byte{TELNET_IAC, TELNET_CMD_GO_AHEAD}

func isCommandWithOption(cmd byte) bool {
	if cmd == 0xFA || cmd == 0xFB || cmd == 0xFC || cmd == 0xFD || cmd == 0xFE {
		return true
	}
	return false
}

func GetTelnetCommandsStrings(buf []byte) []string {
	res := make([]string, 0)
	for {
		idx := bytes.Index(buf, []byte{0xff})
		if idx == -1 {
			break
		}
		nextIdx := idx + 2

		if len(buf) < idx+2 {
			log.Fatalf("Buffer does not contain command after IAC: %v\n", buf)
		}

		cmd := buf[idx+1]
		var cmdStr string
		switch cmd {
		case byte(TELNET_CMD_NOP):
			cmdStr = "NOP"
		case byte(TELNET_CMD_SE):
			cmdStr = "SE"
		case byte(TELNET_CMD_DATAMARK):
			cmdStr = "Data Mark"
		case byte(TELNET_CMD_BREAK):
			cmdStr = "Break"
		case byte(TELNET_CMD_INTERRUPT_PROCESS):
			cmdStr = "Interrupt Process"
		case byte(TELNET_CMD_ABORT_OUTPUT):
			cmdStr = "Abort output"
		case byte(TELNET_CMD_ARE_YOU_THERE):
			cmdStr = "Are You There"
		case byte(TELNET_CMD_ERASE_CHARACTER):
			cmdStr = "Erase character"
		case byte(TELNET_CMD_ERASE_LINE):
			cmdStr = "Erase Line"
		case byte(TELNET_CMD_GO_AHEAD):
			cmdStr = "Go ahead"
		case byte(TELNET_CMD_SB):
			cmdStr = "SB"
		case byte(TELNET_CMD_WILL):
			cmdStr = "WILL"
		case byte(TELNET_CMD_WONT):
			cmdStr = "WON'T"
		case byte(TELNET_CMD_DO):
			cmdStr = "DO"
		case byte(TELNET_CMD_DONT):
			cmdStr = "DON'T"
		case byte(TELNET_CMD_IAC):
			cmdStr = "IAC"
		default:
			cmdStr = "!UNKNOWN"
		}

		s := fmt.Sprintf("IAC %s", cmdStr)

		if isCommandWithOption(cmd) {
			if len(buf) < idx+3 {
				log.Fatalf("Buffer does not contain option after command: %v\n", buf)
			}

			opt := buf[idx+2]
			var optStr string
			switch opt {
			case byte(TELNET_OPT_BINARY_TRANSMISSION):
				optStr = "Binary Transmission"
			case byte(TELNET_OPT_ECHO):
				optStr = "Echo"
			case byte(TELNET_OPT_RECONNECTION):
				optStr = "Reconnection"
			case byte(TELNET_OPT_SUPPRESS_GO_AHEAD):
				optStr = "Suppress Go Ahead"
			case byte(TELNET_OPT_APPROX_MESSAGE_SIZE_NEGOTIATION):
				optStr = "Approx Message Size Negotiation"
			case byte(TELNET_OPT_STATUS):
				optStr = "Status"
			case byte(TELNET_OPT_TIMING_MARK):
				optStr = "Timing Mark"
			case byte(TELNET_OPT_REMOTE_CONTROLLED_TRANS_AND_ECHO):
				optStr = "Remote Controlled Trans and Echo"
			case byte(TELNET_OPT_OUTPUT_LINE_WIDTH):
				optStr = "Output Line Width"
			case byte(TELNET_OPT_OUTPUT_PAGE_SIZE):
				optStr = "Output Page Size"
			case byte(TELNET_OPT_OUTPUT_CARRIAGE_RETURN_DISPOSITION):
				optStr = "Output Carriage-Return Disposition"
			case byte(TELNET_OPT_OUTPUT_HORIZONTAL_TAB_STOPS):
				optStr = "Output Horizontal Tab Stops"
			case byte(TELNET_OPT_OUTPUT_HORIZONTAL_TAB_DISPOSITION):
				optStr = "Output Horizontal Tab Disposition"
			case byte(TELNET_OPT_OUTPUT_FORMFEED_DISPOSITION):
				optStr = "Output Formfeed Disposition"
			case byte(TELNET_OPT_OUTPUT_VERTICAL_TABSTOPS):
				optStr = "Output Vertical Tabstops"
			case byte(TELNET_OPT_OUTPUT_VERTICAL_TAB_DISPOSITION):
				optStr = "Output Vertical Tab Disposition"
			case byte(TELNET_OPT_OUTPUT_LINEFEED_DISPOSITION):
				optStr = "Output Linefeed Disposition"
			case byte(TELNET_OPT_EXTENDED_ASCII):
				optStr = "Extended ASCII"
			case byte(TELNET_OPT_LOGOUT):
				optStr = "Logout"
			case byte(TELNET_OPT_BYTE_MACRO):
				optStr = "Byte Macro"
			case byte(TELNET_OPT_DATA_ENTRY_TERMINAL):
				optStr = "Data Entry Terminal"
			case byte(TELNET_OPT_SUPDUP):
				optStr = "SUPDUP"
			case byte(TELNET_OPT_SUPDUP_OUTPUT):
				optStr = "SUPDUP Output"
			case byte(TELNET_OPT_SEND_LOCATION):
				optStr = "Send Location"
			case byte(TELNET_OPT_TERMINAL_TYPE):
				optStr = "Terminal Type"
			case byte(TELNET_OPT_END_OF_RECORD):
				optStr = "End of Record"
			case byte(TELNET_OPT_TACACS_USER_IDENTIFICATION):
				optStr = "TACACS User Identification"
			case byte(TELNET_OPT_OUTPUT_MARKING):
				optStr = "Output Marking"
			case byte(TELNET_OPT_TERMINAL_LOCATION_NUMBER):
				optStr = "Terminal Location Number"
			case byte(TELNET_OPT_TELNET_3270_REGIME):
				optStr = "Telnet 3270 Regime"
			case byte(TELNET_OPT_X3_PAD):
				optStr = "X.3 PAD"
			case byte(TELNET_OPT_NEGOTIATE_ABOUT_WINDOW_SIZE):
				optStr = "Negotiate About Window Size"
			case byte(TELNET_OPT_TERMINAL_SPEED):
				optStr = "Terminal Speed"
			case byte(TELNET_OPT_REMOTE_FLOW_CONTROL):
				optStr = "Remote Flow Control"
			case byte(TELNET_OPT_LINEMODE):
				optStr = "Linemode"
			case byte(TELNET_OPT_X_DISPLAY_LOCATION):
				optStr = "X Display Location"
			case byte(TELNET_OPT_MSDP):
				optStr = "MSDP"
			case byte(TELNET_OPT_MSSP):
				optStr = "MSSP"
			case byte(TELNET_OPT_MCCPV1):
				optStr = "MCCPv1"
			case byte(TELNET_OPT_MCCPV2):
				optStr = "MCCPv2"
			case byte(TELNET_OPT_MSP):
				optStr = "MSP"
			case byte(TELNET_OPT_MXP):
				optStr = "MXP"
			case byte(TELNET_OPT_ZMP):
				optStr = "ZMP"
			case byte(TELNET_OPT_ATCP):
				optStr = "ATCP"
			case byte(TELNET_OPT_GMCP):
				optStr = "GMCP"
			case byte(TELNET_OPT_EXTENDED_OPTIONS_LIST):
				optStr = "Extended-Options-List"
			default:
				optStr = "!UNKNOWN"
			}

			s = fmt.Sprintf("%s %s", s, optStr)
			nextIdx += 1
		}

		res = append(res, s)

		buf = buf[nextIdx:]
	}

	return res
}
