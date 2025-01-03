//******************************************************************************************************
//  Ticks_test.go - Gbtc
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
//  09/11/2021 - J. Ritchie Carroll
//       Generated original version of source code.
//
//******************************************************************************************************
package ticks

import (
	"testing"
	"time"
)

func TestValidateTicksConstants(t *testing.T) {
	if LeapSecondFlag != 0x8000000000000000 {
		t.Fatalf("ValidateTicksConstants: unexpected ticks leap second flag value")
	}

	if ValueMask != 0x3FFFFFFFFFFFFFFF {
		t.Fatalf("ValidateTicksConstants: unexpected ticks value mask value")
	}

	if LeapSecondDirection != 0x4000000000000000 {
		t.Fatalf("ValidateTicksConstants: unexpected ticks leap second direction flag value")
	}
	if UnixBaseOffset.ToUnixNano() != 0 {
		t.Fatalf("unix base offset is incorrect")
	}
}

func TestTicksTimeConversions(t *testing.T) {
	timestamp := time.Date(2021, 9, 11, 14, 46, 39, 339127800, time.UTC)
	ticks := FromTime(timestamp)

	if ticks != 637669683993391278 {
		t.Fatalf("TicksToTimeConversions: unexpected ToTicks value conversion")
	}

	ticks = 637669698432643641
	timestamp = ticks.ToTime()

	if timestamp != time.Date(2021, 9, 11, 15, 10, 43, 264364100, time.UTC) {
		t.Fatalf("TicksToTimeConversions: unexpected ToTime value conversion")
	}
}
