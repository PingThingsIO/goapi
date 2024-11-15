//******************************************************************************************************
//  Ticks.go - Gbtc
//
//  Copyright © 2021, Grid Protection Alliance.  All Rights Reserved.
//
//  Licensed to the Grid Protection Alliance (GPA) under one or more contributor license agreements. See
//  the NOTICE file distributed with this work for additional information regarding copyright ownership.
//  The GPA licenses this file to you under the MIT License (MIT), the "License"; you may not use this
//  file except in compliance with the License. You may obtain a copy of the License at:
//
//      http://opensource.org/licenses/MIT
//
//  Unless agreed to in writing, the subject software distributed under the License is distributed on an
//  "AS-IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. Refer to the
//  License for the specific language governing permissions and limitations.
//
//  Code Modification History:
//  ----------------------------------------------------------------------------------------------------
//  12.01.2024 - Noam Preil
//       Switch to more idiomatic implementation
//  09/09/2021 - J. Ritchie Carroll
//       Generated original version of source code.
//
//******************************************************************************************************

package ticks

import (
	"time"
)

// Ticks is a 64-bit integer used to designate time in STTP. The value represents the number of 100-nanosecond intervals
// that have elapsed since 12:00:00 midnight, January 1, 0001 UTC, in the Gregorian calendar. A single tick represents 100ns.
// Only bits 01 to 62 (0x3FFFFFFFFFFFFFFF) are used to represent the timestamp value. Bit 64 (0x8000000000000000) is used
// to denote a leap second, and bit 63 (0x4000000000000000) is used to denote the leap second's direction, 0 for add, 1 for delete.
// Leap seconds are exposed, but are silently discarded upon conversion to Go or Unix timestamps.
type Ticks uint64

// Equivalent to UTC time 01/01/0001 00:00:00.000.
const Min Ticks = 0
// TODO: this comment is definitely wrong.
// Max is the maximum value for Ticks. It represents UTC time 12/31/1999 11:59:59.999.
const Max Ticks = 3155378975999999999

// Ticks are every 100ns == 0.1us
const PerMicrosecond = 10
const PerMillisecond Ticks = PerMicrosecond*1000

const LeapSecondFlag Ticks = 1 << 63
const LeapSecondDirection Ticks = 1 << 62
const ValueMask Ticks = ^LeapSecondFlag & ^LeapSecondDirection

// Ticks representation of the Unix epoch, equivalent to FromUnixNs(0).
const UnixBaseOffset Ticks = 621355968000000000

const TimeFormat string = "2006-01-02 15:04:05.999999999"
const ShortTimeFormat string = "15:04:05.999"

func (t Ticks) TimestampValue() int64 {
	return int64(t&ValueMask)
}

// Converts a unix nanoseconds timestamp into a Ticks value.
func FromUnixNs(ns uint64) Ticks {
	return Ticks(ns / 100) + UnixBaseOffset
}

// FromTime converts a standard Go Time value to a Ticks value.
func FromTime(time time.Time) Ticks {
	return FromUnixNs(uint64(time.UnixNano()))
}

// Now gets the current local time as a Ticks value.
func Now() Ticks {
	return FromTime(time.Now())
}

// UtcNow gets the current time in UTC as a Ticks value.
func UtcNow() Ticks {
	return FromTime(time.Now().UTC())
}

// ToTime converts a Ticks value to standard Go Time value.
func (t Ticks) ToTime() time.Time {
	return time.Unix(0, int64((t-UnixBaseOffset)&ValueMask)*100).UTC()
}

// Converts the ticks value into a Unix nanoseconds timestamp
func (t Ticks) ToUnixNs() uint64 {
	return uint64(((t & ValueMask) - UnixBaseOffset) * 100)
}

func (t Ticks) String() string {
	return t.ToTime().Format(TimeFormat)
}

func (t Ticks) ShortTime() string {
	return t.ToTime().Format(ShortTimeFormat)
}

func (t Ticks) IsLeapSecond() bool {
	return (t&LeapSecondFlag) != 0
}

func (t *Ticks) SetLeapSecond(){
	*t |= LeapSecondFlag
}

func (t *Ticks) SetLeapSecondDirection(negative bool) {
	if negative {
		*t |= LeapSecondDirection
	} else {
		*t &= ^LeapSecondDirection
	}
}

func (t Ticks) IsNegativeLeapSecond() bool {
	return t.IsLeapSecond() && (t&LeapSecondDirection) != 0
}

