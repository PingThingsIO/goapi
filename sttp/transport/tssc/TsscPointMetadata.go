//******************************************************************************************************
//  TsscPointMetadata.go - Gbtc
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
//  12/02/2016 - Steven E. Chisholm
//       Generated original version of source code.
//  09/20/2021 - J. Ritchie Carroll
//       Migrated code to Go.
//
//******************************************************************************************************

package tssc

import "math"

var codeWords = struct {
	EndOfStream       byte
	PointIDXOR4       byte
	PointIDXOR8       byte
	PointIDXOR12      byte
	PointIDXOR16      byte
	PointIDXOR20      byte
	PointIDXOR24      byte
	PointIDXOR32      byte
	TimeDelta1Forward byte
	TimeDelta2Forward byte
	TimeDelta3Forward byte
	TimeDelta4Forward byte
	TimeDelta1Reverse byte
	TimeDelta2Reverse byte
	TimeDelta3Reverse byte
	TimeDelta4Reverse byte
	Timestamp2        byte
	TimeXOR7Bit       byte
	StateFlags2       byte
	StateFlags7Bit32  byte
	Value1            byte
	Value2            byte
	Value3            byte
	ValueZero         byte
	ValueXOR4         byte
	ValueXOR8         byte
	ValueXOR12        byte
	ValueXOR16        byte
	ValueXOR20        byte
	ValueXOR24        byte
	ValueXOR28        byte
	ValueXOR32        byte
}{
	EndOfStream:       0,
	PointIDXOR4:       1,
	PointIDXOR8:       2,
	PointIDXOR12:      3,
	PointIDXOR16:      4,
	PointIDXOR20:      5,
	PointIDXOR24:      6,
	PointIDXOR32:      7,
	TimeDelta1Forward: 8,
	TimeDelta2Forward: 9,
	TimeDelta3Forward: 10,
	TimeDelta4Forward: 11,
	TimeDelta1Reverse: 12,
	TimeDelta2Reverse: 13,
	TimeDelta3Reverse: 14,
	TimeDelta4Reverse: 15,
	Timestamp2:        16,
	TimeXOR7Bit:       17,
	StateFlags2:       18,
	StateFlags7Bit32:  19,
	Value1:            20,
	Value2:            21,
	Value3:            22,
	ValueZero:         23,
	ValueXOR4:         24,
	ValueXOR8:         25,
	ValueXOR12:        26,
	ValueXOR16:        27,
	ValueXOR20:        28,
	ValueXOR24:        29,
	ValueXOR28:        30,
	ValueXOR32:        31,
}

type pointMetadata struct {
	PrevNextPointID1 int32
	PrevStateFlags1  uint32
	PrevStateFlags2  uint32
	PrevValue1       uint32
	PrevValue2       uint32
	PrevValue3       uint32

	commandStats                [32]byte
	commandsSentSinceLastChange int32

	// Bit codes for the 4 modes of encoding
	mode byte

	// Mode 1 means no prefix
	mode21      byte
	mode31      byte
	mode301     byte
	mode41      byte
	mode401     byte
	mode4001    byte
	startupMode int32

	writeBits func(int32, int32)
	readBit   func() int32
	readBits5 func() int32
}

func newPointMetadata(writeBits func(int32, int32), readBit func() int32, readBits5 func() int32) *pointMetadata {
	return &pointMetadata{
		mode:      4,
		mode41:    codeWords.Value1,
		mode401:   codeWords.Value2,
		mode4001:  codeWords.Value3,
		writeBits: writeBits,
		readBit:   readBit,
		readBits5: readBits5,
	}
}

func (pm *pointMetadata) WriteCode(code int32) {
	switch pm.mode {
	case 1:
		pm.writeBits(code, 5)
	case 2:
		if code == int32(pm.mode21) {
			pm.writeBits(1, 1)
		} else {
			pm.writeBits(code, 6)
		}
	case 3:
		if code == int32(pm.mode31) {
			pm.writeBits(1, 1)
		} else if code == int32(pm.mode301) {
			pm.writeBits(1, 2)
		} else {
			pm.writeBits(code, 7)
		}
	case 4:
		if code == int32(pm.mode41) {
			pm.writeBits(1, 1)
		} else if code == int32(pm.mode401) {
			pm.writeBits(1, 2)
		} else if code == int32(pm.mode4001) {
			pm.writeBits(1, 3)
		} else {
			pm.writeBits(code, 8)
		}
	default:
		panic("Coding Error")
	}

	pm.updatedCodeStatistics(code)
}

func (pm *pointMetadata) ReadCode() int32 {
	var code int32

	switch pm.mode {
	case 1:
		code = pm.readBits5()
	case 2:
		if pm.readBit() == 1 {
			code = int32(pm.mode21)
		} else {
			code = pm.readBits5()
		}
	case 3:
		if pm.readBit() == 1 {
			code = int32(pm.mode31)
		} else if pm.readBit() == 1 {
			code = int32(pm.mode301)
		} else {
			code = pm.readBits5()
		}
	case 4:
		if pm.readBit() == 1 {
			code = int32(pm.mode41)
		} else if pm.readBit() == 1 {
			code = int32(pm.mode401)
		} else if pm.readBit() == 1 {
			code = int32(pm.mode4001)
		} else {
			code = pm.readBits5()
		}
	default:
		panic("Unsupported compression mode")
	}

	pm.updatedCodeStatistics(code)
	return code
}

func (pm *pointMetadata) updatedCodeStatistics(code int32) {
	pm.commandsSentSinceLastChange++
	pm.commandStats[code]++

	if pm.startupMode == 0 && pm.commandsSentSinceLastChange > 5 {
		pm.startupMode++
		pm.adaptCommands()
	} else if pm.startupMode == 1 && pm.commandsSentSinceLastChange > 20 {
		pm.startupMode++
		pm.adaptCommands()
	} else if pm.startupMode == 2 && pm.commandsSentSinceLastChange > 100 {
		pm.adaptCommands()
	}
}

func (pm *pointMetadata) adaptCommands() {
	var code1 byte = 0
	var count1 int32 = 0

	var code2 byte = 1
	var count2 int32 = 0

	var code3 byte = 2
	var count3 int32 = 0

	var total int32 = 0

	for i := 0; i < len(pm.commandStats); i++ {
		var count int32 = int32(pm.commandStats[i])
		pm.commandStats[i] = 0

		total += count

		if count > count3 {
			if count > count1 {
				code3 = code2
				count3 = count2

				code2 = code1
				count2 = count1

				code1 = byte(i)
				count1 = count
			} else if count > count2 {
				code3 = code2
				count3 = count2

				code2 = byte(i)
				count2 = count
			} else {
				code3 = byte(i)
				count3 = count
			}
		}
	}

	var mode1Size int32 = total * 5
	var mode2Size int32 = count1*1 + (total-count1)*6
	var mode3Size int32 = count1*1 + count2*2 + (total-count1-count2)*7
	var mode4Size int32 = count1*1 + count2*2 + count3*3 + (total-count1-count2-count3)*8

	var minSize int32 = math.MaxInt32

	minSize = min(minSize, mode1Size)
	minSize = min(minSize, mode2Size)
	minSize = min(minSize, mode3Size)
	minSize = min(minSize, mode4Size)

	if minSize == mode1Size {
		pm.mode = 1
	} else if minSize == mode2Size {
		pm.mode = 2
		pm.mode21 = code1
	} else if minSize == mode3Size {
		pm.mode = 3
		pm.mode31 = code1
		pm.mode301 = code2
	} else if minSize == mode4Size {
		pm.mode = 4
		pm.mode41 = code1
		pm.mode401 = code2
		pm.mode4001 = code3
	} else {
		if pm.writeBits == nil {
			panic("Subscriber Coding Error")
		}

		panic("Publisher Coding Error")
	}

	pm.commandsSentSinceLastChange = 0
}

func min(lv int32, rv int32) int32 {
	if lv < rv {
		return lv
	}

	if rv < lv {
		return rv
	}

	return lv
}
