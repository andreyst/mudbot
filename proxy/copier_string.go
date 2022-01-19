// Code generated by "stringer -type=accumulationPolicy,copierState -output copier_string.go"; DO NOT EDIT.

package proxy

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ACCUMULATION_POLICY_DO-0]
	_ = x[ACCUMULATION_POLICY_DONT-1]
}

const _accumulationPolicy_name = "ACCUMULATION_POLICY_DOACCUMULATION_POLICY_DONT"

var _accumulationPolicy_index = [...]uint8{0, 22, 46}

func (i accumulationPolicy) String() string {
	if i < 0 || i >= accumulationPolicy(len(_accumulationPolicy_index)-1) {
		return "accumulationPolicy(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _accumulationPolicy_name[_accumulationPolicy_index[i]:_accumulationPolicy_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[STATE_START-0]
	_ = x[STATE_TELNET_IAC-1]
	_ = x[STATE_TELNET_SB-2]
	_ = x[STATE_TELNET_MCCPV2-3]
	_ = x[STATE_TELNET_MCCPV2_IAC-4]
}

const _copierState_name = "STATE_STARTSTATE_TELNET_IACSTATE_TELNET_SBSTATE_TELNET_MCCPV2STATE_TELNET_MCCPV2_IAC"

var _copierState_index = [...]uint8{0, 11, 27, 42, 61, 84}

func (i copierState) String() string {
	if i < 0 || i >= copierState(len(_copierState_index)-1) {
		return "copierState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _copierState_name[_copierState_index[i]:_copierState_index[i+1]]
}
