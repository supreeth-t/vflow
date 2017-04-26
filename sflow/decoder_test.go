//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder_test.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package sflow

import (
	"bytes"
	"testing"
)

var TestsFlowRawPacket = []byte{
	0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x1, 0xc0, 0xe5, 0xd6, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x6d, 0x3f, 0x61, 0x11, 0x57, 0x35, 0x0, 0x0, 0x0,
	0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x9c, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x2, 0x16, 0x0, 0x0, 0x7, 0xd0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x2, 0x28, 0x0, 0x0, 0x2, 0x16, 0x0, 0x0, 0x0, 0x2, 0x0,
	0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x5c, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0,
	0x4e, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x4a, 0xde, 0xad, 0x7a, 0x48,
	0xcc, 0x37, 0xd4, 0x4, 0xff, 0x1, 0x18, 0x1e, 0x81, 0x0, 0x0, 0x7, 0x8,
	0x0, 0x45, 0x0, 0x0, 0x38, 0x0, 0x0, 0x0, 0x0, 0xef, 0x1, 0xff, 0x3e,
	0xb5, 0x1e, 0x80, 0x6a, 0xc0, 0xe5, 0xd6, 0x17, 0xb, 0x0, 0xf4, 0xff,
	0x0, 0x0, 0x0, 0x0, 0x45, 0x0, 0x0, 0x40, 0x65, 0x2d, 0x0, 0x0, 0x1,
	0x1, 0xfc, 0x4d, 0xc0, 0xe5, 0xd6, 0x17, 0xc0, 0x10, 0x1, 0x35, 0x8,
	0x0, 0x9f, 0x7a, 0x34, 0x2, 0x24, 0x83, 0x0, 0x0, 0x0, 0x0, 0x3,
	0xe9, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
}

func TestSFHeaderDecode(t *testing.T) {
	filter := []uint32{DataCounterSample}
	reader := bytes.NewReader(TestsFlowRawPacket)
	d := NewSFDecoder(reader, filter)
	datagram, err := d.sfHeaderDecode()
	if err != nil {
		t.Error("unexpected error", err)
	}

	if datagram.IPAddress.String() != "192.229.214.0" {
		t.Error("expected agent ip address: 192.229.214.0, got",
			datagram.IPAddress.String())
	}
}

func TestGetSampleInfo(t *testing.T) {
	filter := []uint32{DataCounterSample}
	reader := bytes.NewReader(TestsFlowRawPacket)
	// skip sflow header
	skip := make([]byte, 4*7)
	reader.Read(skip)

	d := NewSFDecoder(reader, filter)
	sfTypeFormat, sfDataLength, err := d.getSampleInfo()
	if err != nil {
		t.Error("unexpected error", err)
	}
	if sfTypeFormat != 1 {
		t.Error("expected type format# 1, got", sfTypeFormat)
	}
	if sfDataLength != 156 {
		t.Error("expected data length: 156, got", sfDataLength)
	}
}

func TestSFDecode(t *testing.T) {
	filter := []uint32{DataCounterSample}
	reader := bytes.NewReader(TestsFlowRawPacket)
	d := NewSFDecoder(reader, filter)
	_, err := d.SFDecode()
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func BenchmarkSFDecode(b *testing.B) {
	filter := []uint32{DataCounterSample}
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(TestsFlowRawPacket)
		d := NewSFDecoder(reader, filter)
		d.SFDecode()
	}
}
