package telnet

import (
	"bytes"
	"fmt"
	"log"
)

type telnetCommand int
type telnetOption int

const IAC = 0xFF

const (
	COMMAND_NOP               telnetCommand = 0xF1
	COMMAND_SE                              = 0xF0
	COMMAND_DATAMARK                        = 0xF2
	COMMAND_BREAK                           = 0xF3
	COMMAND_INTERRUPT_PROCESS               = 0xF4
	COMMAND_ABORT_OUTPUT                    = 0xF5
	COMMAND_ARE_YOU_THERE                   = 0xF6
	COMMAND_ERASE_CHARACTER                 = 0xF7
	COMMAND_ERASE_LINE                      = 0xF8
	COMMAND_GO_AHEAD                        = 0xF9
	COMMAND_SB                              = 0xFA
	COMMAND_WILL                            = 0xFB
	COMMAND_WONT                            = 0xFC
	COMMAND_DO                              = 0xFD
	COMMAND_DONT                            = 0xFE
	COMMAND_IAC                             = 0xFF
)

const (
	OPTION_BINARY_TRANSMISSION                telnetOption = 0x0
	OPTION_ECHO                                            = 0x1
	OPTION_RECONNECTION                                    = 0x2
	OPTION_SUPPRESS_GO_AHEAD                               = 0x3
	OPTION_APPROX_MESSAGE_SIZE_NEGOTIATION                 = 0x4
	OPTION_STATUS                                          = 0x5
	OPTION_TIMING_MARK                                     = 0x6
	OPTION_REMOTE_CONTROLLED_TRANS_AND_ECHO                = 0x7
	OPTION_OUTPUT_LINE_WIDTH                               = 0x8
	OPTION_OUTPUT_PAGE_SIZE                                = 0x9
	OPTION_OUTPUT_CARRIAGE_RETURN_DISPOSITION              = 0xA
	OPTION_OUTPUT_HORIZONTAL_TAB_STOPS                     = 0xB
	OPTION_OUTPUT_HORIZONTAL_TAB_DISPOSITION               = 0xC
	OPTION_OUTPUT_FORMFEED_DISPOSITION                     = 0xD
	OPTION_OUTPUT_VERTICAL_TABSTOPS                        = 0xE
	OPTION_OUTPUT_VERTICAL_TAB_DISPOSITION                 = 0xF
	OPTION_OUTPUT_LINEFEED_DISPOSITION                     = 0x10
	OPTION_EXTENDED_ASCII                                  = 0x11
	OPTION_LOGOUT                                          = 0x12
	OPTION_BYTE_MACRO                                      = 0x13
	OPTION_DATA_ENTRY_TERMINAL                             = 0x14
	OPTION_SUPDUP                                          = 0x15
	OPTION_SUPDUP_OUTPUT                                   = 0x16
	OPTION_SEND_LOCATION                                   = 0x17
	OPTION_TERMINAL_TYPE                                   = 0x18
	OPTION_END_OF_RECORD                                   = 0x19
	OPTION_TACACS_USER_IDENTIFICATION                      = 0x1A
	OPTION_OUTPUT_MARKING                                  = 0x1B
	OPTION_TERMINAL_LOCATION_NUMBER                        = 0x1C
	OPTION_TELNET_3270_REGIME                              = 0x1D
	OPTION_X3_PAD                                          = 0x1E
	OPTION_NEGOTIATE_ABOUT_WINDOW_SIZE                     = 0x1F
	OPTION_TERMINAL_SPEED                                  = 0x20
	OPTION_REMOTE_FLOW_CONTROL                             = 0x21
	OPTION_LINEMODE                                        = 0x22
	OPTION_X_DISPLAY_LOCATION                              = 0x23
	OPTION_MSDP                                            = 0x45
	OPTION_MSSP                                            = 0x46
	OPTION_MCCPV1                                          = 0x55
	OPTION_MCCPV2                                          = 0x56
	OPTION_MSP                                             = 0x5A
	OPTION_MXP                                             = 0x5B
	OPTION_ZMP                                             = 0x5D
	OPTION_ATCP                                            = 0xC8
	OPTION_GMCP                                            = 0xC9
	OPTION_EXTENDED_OPTIONS_LIST                           = 0xFF
)

var CompressionStartSequence = []byte{IAC, COMMAND_SB, OPTION_MCCPV2, IAC, COMMAND_SE}
var GaSequence = []byte{IAC, COMMAND_GO_AHEAD}

func isCommandWithOption(cmd byte) bool {
	if cmd == 0xFA || cmd == 0xFB || cmd == 0xFC || cmd == 0xFD || cmd == 0xFE {
		return true
	}
	return false
}

func GetCommandsStrings(buf []byte) []string {
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
		case byte(COMMAND_NOP):
			cmdStr = "NOP"
		case byte(COMMAND_SE):
			cmdStr = "SE"
		case byte(COMMAND_DATAMARK):
			cmdStr = "Data Mark"
		case byte(COMMAND_BREAK):
			cmdStr = "Break"
		case byte(COMMAND_INTERRUPT_PROCESS):
			cmdStr = "Interrupt Process"
		case byte(COMMAND_ABORT_OUTPUT):
			cmdStr = "Abort output"
		case byte(COMMAND_ARE_YOU_THERE):
			cmdStr = "Are You There"
		case byte(COMMAND_ERASE_CHARACTER):
			cmdStr = "Erase character"
		case byte(COMMAND_ERASE_LINE):
			cmdStr = "Erase Line"
		case byte(COMMAND_GO_AHEAD):
			cmdStr = "Go ahead"
		case byte(COMMAND_SB):
			cmdStr = "SB"
		case byte(COMMAND_WILL):
			cmdStr = "WILL"
		case byte(COMMAND_WONT):
			cmdStr = "WON'T"
		case byte(COMMAND_DO):
			cmdStr = "DO"
		case byte(COMMAND_DONT):
			cmdStr = "DON'T"
		case byte(COMMAND_IAC):
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
			case byte(OPTION_BINARY_TRANSMISSION):
				optStr = "Binary Transmission"
			case byte(OPTION_ECHO):
				optStr = "Echo"
			case byte(OPTION_RECONNECTION):
				optStr = "Reconnection"
			case byte(OPTION_SUPPRESS_GO_AHEAD):
				optStr = "Suppress Go Ahead"
			case byte(OPTION_APPROX_MESSAGE_SIZE_NEGOTIATION):
				optStr = "Approx Message Size Negotiation"
			case byte(OPTION_STATUS):
				optStr = "Status"
			case byte(OPTION_TIMING_MARK):
				optStr = "Timing Mark"
			case byte(OPTION_REMOTE_CONTROLLED_TRANS_AND_ECHO):
				optStr = "Remote Controlled Trans and Echo"
			case byte(OPTION_OUTPUT_LINE_WIDTH):
				optStr = "Output Line Width"
			case byte(OPTION_OUTPUT_PAGE_SIZE):
				optStr = "Output Page Size"
			case byte(OPTION_OUTPUT_CARRIAGE_RETURN_DISPOSITION):
				optStr = "Output Carriage-Return Disposition"
			case byte(OPTION_OUTPUT_HORIZONTAL_TAB_STOPS):
				optStr = "Output Horizontal Tab Stops"
			case byte(OPTION_OUTPUT_HORIZONTAL_TAB_DISPOSITION):
				optStr = "Output Horizontal Tab Disposition"
			case byte(OPTION_OUTPUT_FORMFEED_DISPOSITION):
				optStr = "Output Formfeed Disposition"
			case byte(OPTION_OUTPUT_VERTICAL_TABSTOPS):
				optStr = "Output Vertical Tabstops"
			case byte(OPTION_OUTPUT_VERTICAL_TAB_DISPOSITION):
				optStr = "Output Vertical Tab Disposition"
			case byte(OPTION_OUTPUT_LINEFEED_DISPOSITION):
				optStr = "Output Linefeed Disposition"
			case byte(OPTION_EXTENDED_ASCII):
				optStr = "Extended ASCII"
			case byte(OPTION_LOGOUT):
				optStr = "Logout"
			case byte(OPTION_BYTE_MACRO):
				optStr = "Byte Macro"
			case byte(OPTION_DATA_ENTRY_TERMINAL):
				optStr = "Data Entry Terminal"
			case byte(OPTION_SUPDUP):
				optStr = "SUPDUP"
			case byte(OPTION_SUPDUP_OUTPUT):
				optStr = "SUPDUP Output"
			case byte(OPTION_SEND_LOCATION):
				optStr = "Send Location"
			case byte(OPTION_TERMINAL_TYPE):
				optStr = "Terminal Type"
			case byte(OPTION_END_OF_RECORD):
				optStr = "End of Record"
			case byte(OPTION_TACACS_USER_IDENTIFICATION):
				optStr = "TACACS User Identification"
			case byte(OPTION_OUTPUT_MARKING):
				optStr = "Output Marking"
			case byte(OPTION_TERMINAL_LOCATION_NUMBER):
				optStr = "Terminal Location Number"
			case byte(OPTION_TELNET_3270_REGIME):
				optStr = "Telnet 3270 Regime"
			case byte(OPTION_X3_PAD):
				optStr = "X.3 PAD"
			case byte(OPTION_NEGOTIATE_ABOUT_WINDOW_SIZE):
				optStr = "Negotiate About Window Size"
			case byte(OPTION_TERMINAL_SPEED):
				optStr = "Terminal Speed"
			case byte(OPTION_REMOTE_FLOW_CONTROL):
				optStr = "Remote Flow Control"
			case byte(OPTION_LINEMODE):
				optStr = "Linemode"
			case byte(OPTION_X_DISPLAY_LOCATION):
				optStr = "X Display Location"
			case byte(OPTION_MSDP):
				optStr = "MSDP"
			case byte(OPTION_MSSP):
				optStr = "MSSP"
			case byte(OPTION_MCCPV1):
				optStr = "MCCPv1"
			case byte(OPTION_MCCPV2):
				optStr = "MCCPv2"
			case byte(OPTION_MSP):
				optStr = "MSP"
			case byte(OPTION_MXP):
				optStr = "MXP"
			case byte(OPTION_ZMP):
				optStr = "ZMP"
			case byte(OPTION_ATCP):
				optStr = "ATCP"
			case byte(OPTION_GMCP):
				optStr = "GMCP"
			case byte(OPTION_EXTENDED_OPTIONS_LIST):
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
