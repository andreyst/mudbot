package telnet

import (
	"bytes"
	"fmt"
	"log"
)

const IAC byte = 0xFF

const (
	COMMAND_NOP               byte = 0xF1
	COMMAND_SE                byte = 0xF0
	COMMAND_DATAMARK          byte = 0xF2
	COMMAND_BREAK             byte = 0xF3
	COMMAND_INTERRUPT_PROCESS byte = 0xF4
	COMMAND_ABORT_OUTPUT      byte = 0xF5
	COMMAND_ARE_YOU_THERE     byte = 0xF6
	COMMAND_ERASE_CHARACTER   byte = 0xF7
	COMMAND_ERASE_LINE        byte = 0xF8
	COMMAND_GO_AHEAD          byte = 0xF9
	COMMAND_SB                byte = 0xFA
	COMMAND_WILL              byte = 0xFB
	COMMAND_WONT              byte = 0xFC
	COMMAND_DO                byte = 0xFD
	COMMAND_DONT              byte = 0xFE
	COMMAND_IAC               byte = 0xFF
)

const (
	OPTION_BINARY_TRANSMISSION                byte = 0x0
	OPTION_ECHO                               byte = 0x1
	OPTION_RECONNECTION                       byte = 0x2
	OPTION_SUPPRESS_GO_AHEAD                  byte = 0x3
	OPTION_APPROX_MESSAGE_SIZE_NEGOTIATION    byte = 0x4
	OPTION_STATUS                             byte = 0x5
	OPTION_TIMING_MARK                        byte = 0x6
	OPTION_REMOTE_CONTROLLED_TRANS_AND_ECHO   byte = 0x7
	OPTION_OUTPUT_LINE_WIDTH                  byte = 0x8
	OPTION_OUTPUT_PAGE_SIZE                   byte = 0x9
	OPTION_OUTPUT_CARRIAGE_RETURN_DISPOSITION byte = 0xA
	OPTION_OUTPUT_HORIZONTAL_TAB_STOPS        byte = 0xB
	OPTION_OUTPUT_HORIZONTAL_TAB_DISPOSITION  byte = 0xC
	OPTION_OUTPUT_FORMFEED_DISPOSITION        byte = 0xD
	OPTION_OUTPUT_VERTICAL_TABSTOPS           byte = 0xE
	OPTION_OUTPUT_VERTICAL_TAB_DISPOSITION    byte = 0xF
	OPTION_OUTPUT_LINEFEED_DISPOSITION        byte = 0x10
	OPTION_EXTENDED_ASCII                     byte = 0x11
	OPTION_LOGOUT                             byte = 0x12
	OPTION_BYTE_MACRO                         byte = 0x13
	OPTION_DATA_ENTRY_TERMINAL                byte = 0x14
	OPTION_SUPDUP                             byte = 0x15
	OPTION_SUPDUP_OUTPUT                      byte = 0x16
	OPTION_SEND_LOCATION                      byte = 0x17
	OPTION_TERMINAL_TYPE                      byte = 0x18
	OPTION_END_OF_RECORD                      byte = 0x19
	OPTION_TACACS_USER_IDENTIFICATION         byte = 0x1A
	OPTION_OUTPUT_MARKING                     byte = 0x1B
	OPTION_TERMINAL_LOCATION_NUMBER           byte = 0x1C
	OPTION_TELNET_3270_REGIME                 byte = 0x1D
	OPTION_X3_PAD                             byte = 0x1E
	OPTION_NEGOTIATE_ABOUT_WINDOW_SIZE        byte = 0x1F
	OPTION_TERMINAL_SPEED                     byte = 0x20
	OPTION_REMOTE_FLOW_CONTROL                byte = 0x21
	OPTION_LINEMODE                           byte = 0x22
	OPTION_X_DISPLAY_LOCATION                 byte = 0x23
	OPTION_MSDP                               byte = 0x45
	OPTION_MSSP                               byte = 0x46
	OPTION_MCCPV1                             byte = 0x55
	OPTION_MCCPV2                             byte = 0x56
	OPTION_MSP                                byte = 0x5A
	OPTION_MXP                                byte = 0x5B
	OPTION_ZMP                                byte = 0x5D
	OPTION_ATCP                               byte = 0xC8
	OPTION_GMCP                               byte = 0xC9
	OPTION_EXTENDED_OPTIONS_LIST              byte = 0xFF
)

var CompressionStartSequence = []byte{IAC, COMMAND_SB, OPTION_MCCPV2, IAC, COMMAND_SE}
var GaSequence = []byte{IAC, COMMAND_GO_AHEAD}

func isCommandWithOption(cmd byte) bool {
	if cmd == COMMAND_SB || cmd == COMMAND_WILL || cmd == COMMAND_WONT || cmd == COMMAND_DO || cmd == COMMAND_DONT {
		return true
	}
	return false
}

func GetCommandsStrings(buf []byte) []string {
	res := make([]string, 0)
	for {
		idx := bytes.Index(buf, []byte{IAC})
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
		case COMMAND_NOP:
			cmdStr = "NOP"
		case COMMAND_SE:
			cmdStr = "SE"
		case COMMAND_DATAMARK:
			cmdStr = "Data Mark"
		case COMMAND_BREAK:
			cmdStr = "Break"
		case COMMAND_INTERRUPT_PROCESS:
			cmdStr = "Interrupt Process"
		case COMMAND_ABORT_OUTPUT:
			cmdStr = "Abort output"
		case COMMAND_ARE_YOU_THERE:
			cmdStr = "Are You There"
		case COMMAND_ERASE_CHARACTER:
			cmdStr = "Erase character"
		case COMMAND_ERASE_LINE:
			cmdStr = "Erase Line"
		case COMMAND_GO_AHEAD:
			cmdStr = "Go ahead"
		case COMMAND_SB:
			cmdStr = "SB"
		case COMMAND_WILL:
			cmdStr = "WILL"
		case COMMAND_WONT:
			cmdStr = "WON'T"
		case COMMAND_DO:
			cmdStr = "DO"
		case COMMAND_DONT:
			cmdStr = "DON'T"
		case COMMAND_IAC:
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
			case OPTION_BINARY_TRANSMISSION:
				optStr = "Binary Transmission"
			case OPTION_ECHO:
				optStr = "Echo"
			case OPTION_RECONNECTION:
				optStr = "Reconnection"
			case OPTION_SUPPRESS_GO_AHEAD:
				optStr = "Suppress Go Ahead"
			case OPTION_APPROX_MESSAGE_SIZE_NEGOTIATION:
				optStr = "Approx Message Size Negotiation"
			case OPTION_STATUS:
				optStr = "Status"
			case OPTION_TIMING_MARK:
				optStr = "Timing Mark"
			case OPTION_REMOTE_CONTROLLED_TRANS_AND_ECHO:
				optStr = "Remote Controlled Trans and Echo"
			case OPTION_OUTPUT_LINE_WIDTH:
				optStr = "Output Line Width"
			case OPTION_OUTPUT_PAGE_SIZE:
				optStr = "Output Page Size"
			case OPTION_OUTPUT_CARRIAGE_RETURN_DISPOSITION:
				optStr = "Output Carriage-Return Disposition"
			case OPTION_OUTPUT_HORIZONTAL_TAB_STOPS:
				optStr = "Output Horizontal Tab Stops"
			case OPTION_OUTPUT_HORIZONTAL_TAB_DISPOSITION:
				optStr = "Output Horizontal Tab Disposition"
			case OPTION_OUTPUT_FORMFEED_DISPOSITION:
				optStr = "Output Formfeed Disposition"
			case OPTION_OUTPUT_VERTICAL_TABSTOPS:
				optStr = "Output Vertical Tabstops"
			case OPTION_OUTPUT_VERTICAL_TAB_DISPOSITION:
				optStr = "Output Vertical Tab Disposition"
			case OPTION_OUTPUT_LINEFEED_DISPOSITION:
				optStr = "Output Linefeed Disposition"
			case OPTION_EXTENDED_ASCII:
				optStr = "Extended ASCII"
			case OPTION_LOGOUT:
				optStr = "Logout"
			case OPTION_BYTE_MACRO:
				optStr = "Byte Macro"
			case OPTION_DATA_ENTRY_TERMINAL:
				optStr = "Data Entry Terminal"
			case OPTION_SUPDUP:
				optStr = "SUPDUP"
			case OPTION_SUPDUP_OUTPUT:
				optStr = "SUPDUP Output"
			case OPTION_SEND_LOCATION:
				optStr = "Send Location"
			case OPTION_TERMINAL_TYPE:
				optStr = "Terminal Type"
			case OPTION_END_OF_RECORD:
				optStr = "End of Record"
			case OPTION_TACACS_USER_IDENTIFICATION:
				optStr = "TACACS User Identification"
			case OPTION_OUTPUT_MARKING:
				optStr = "Output Marking"
			case OPTION_TERMINAL_LOCATION_NUMBER:
				optStr = "Terminal Location Number"
			case OPTION_TELNET_3270_REGIME:
				optStr = "Telnet 3270 Regime"
			case OPTION_X3_PAD:
				optStr = "X.3 PAD"
			case OPTION_NEGOTIATE_ABOUT_WINDOW_SIZE:
				optStr = "Negotiate About Window Size"
			case OPTION_TERMINAL_SPEED:
				optStr = "Terminal Speed"
			case OPTION_REMOTE_FLOW_CONTROL:
				optStr = "Remote Flow Control"
			case OPTION_LINEMODE:
				optStr = "Linemode"
			case OPTION_X_DISPLAY_LOCATION:
				optStr = "X Display Location"
			case OPTION_MSDP:
				optStr = "MSDP"
			case OPTION_MSSP:
				optStr = "MSSP"
			case OPTION_MCCPV1:
				optStr = "MCCPv1"
			case OPTION_MCCPV2:
				optStr = "MCCPv2"
			case OPTION_MSP:
				optStr = "MSP"
			case OPTION_MXP:
				optStr = "MXP"
			case OPTION_ZMP:
				optStr = "ZMP"
			case OPTION_ATCP:
				optStr = "ATCP"
			case OPTION_GMCP:
				optStr = "GMCP"
			case OPTION_EXTENDED_OPTIONS_LIST:
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

func StripTelnet(b []byte) []byte {
	var cursorFrom, cursorTo int
	for ; cursorFrom < len(b); cursorFrom++ {
		if b[cursorFrom] == IAC {
			if cursorFrom < len(b)-1 {
				if isCommandWithOption(b[cursorFrom+1]) {
					cursorFrom += 2
				} else {
					cursorFrom += 1
				}
			}
		} else {
			if cursorFrom > cursorTo {
				b[cursorTo] = b[cursorFrom]
			}
			cursorTo++
		}
	}

	return b[:cursorTo]
}
