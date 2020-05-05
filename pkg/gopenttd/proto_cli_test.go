package gopenttd

import (
	"bytes"
	"encoding/hex"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/global"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDecodeInfoPacketNoGRF(t *testing.T) {
	// real packet captured from Reddit Server #1 28/04/2020
	// this has the first 3 bytes removed because they interfere with the populateServerState func
	testPacket, _ := hex.DecodeString("040090260b00dacf0a000f0e2d7265646469742e636f6d2f722f4f70656e545444202331202d2056616e696c6c6100312e31302e3100000032090152616e646f6d204d617000000400080001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	expectObject := global.OpenttdServerState{
		Status:        true,
		Error:         nil,
		Dedicated:     true,
		Name:          "reddit.com/r/OpenTTD #1 - Vanilla",
		Version:       "1.10.1",
		NeedPass:      false,
		Language:      0,
		Environment:   0,
		Map:           "Random Map",
		MapHeight:     2048,
		MapWidth:      1024,
		DateStart:     time.Date(1940, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		DateCurrent:   time.Date(2000, time.Month(10), 10, 0, 0, 0, 0, time.UTC),
		NumClients:    9,
		MaxClients:    50,
		NumSpectators: 1,
		MaxSpectators: 45,
		NumCompanies:  14,
		MaxCompanies:  15,
		NewgrfCount:   0,
		NewgrfActive:  nil,
	}

	output := global.OpenttdServerState{}
	output.PopulateServerState(bytes.NewBuffer(testPacket))

	assert.EqualValues(t, expectObject, output)
}

func TestDecodeCompanyPacket(t *testing.T) {
	// real packet captured from Reddit Server #1 30/04/2020
	// this has the first 3 bytes removed because they interfere with the populateServerState func
	testPacket, _ := hex.DecodeString("060e00656c20636861706f2032303230205472616e73706f727400aa070000704e550000000000f29159000000000067b5030000000000fc0001010000000000020006000200000000000200030000015361696e745f616e7269205472616e73706f727400b7070000d997b6060000000076d4b10600000000aef47700000000001702015a00000000000000000009000100000000000000000263696361646133333031205472616e73706f7274009407000078896b470000000059f9e04600000000ccb5b00300000000570301b80000000000120000002200000009000a00000000034e65772043686164696e6762757267207472616e7300ad070000915f060200000000abd80402000000005004170000000000290201230000000000000000000f0000000000000000000004477265676f7279202620436f2e00c607000042250b00000000005a700a00000000002311020000000000110101000000001400000000000000000002000000000000054d69535a55205472616e73706f7274009d07000004f2d40a000000003635ba0a000000001fdf8700000000005403014b00000000001c0000002c000000000003000000000676756e646f205472616e73706f727400bf0700000f3f5a0000000000cd52530000000000dec20000000000005a01010d0005000700000000000500020008000000000000075472616e73706f7465205472616e73706f72740099070000260fe74a0000000078e7bc49000000003823960300000000ff02015a0019002100570000000f000c000c000b0000000008537574747964726f697477696368205472616e73706f727400b907000001fbe700000000000042e60000000000f5f90c000000000028020109000b000000060000000e000600050003000000000946656c6978205472616e73706f727400ba07000094fdd400000000005165a00000000000ce570f00000000001502011e0000002200000004000b000200230000000900000a6b616c616d697469206b6f7270009d0700002f0a270200000000274723020000000096e010000000000079020120000c00300002000400090006000f0004000100000b45724b6f576f205472616e73706f727400c00700006d7eaf01000000002d5d4901000000008f61f3fffffffffffd02010b00000000006400000004000200060005000000000d52756e6172205472616e73706f727400ce070000bd9a19000000000028431500000000003cca0900000000001a01010000000000000800000002000000000002000000000e544275735f45787072657373205472616e73706f727400c7070000bbe6020000000000bbdd000000000000352a000000000000900001000001001200000000000000020012000000000000000000000000000000000000")
	expectObject := global.OpenttdServerState{
		Companies: []global.OpenttdCompany{
			{
				Id:          0,
				Name:        "el chapo 2020 Transport",
				YearStart:   1962,
				Value:       5590640,
				Money:       5870066,
				Income:      243047,
				Performance: 252,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    1,
					Truck:    0,
					Bus:      0,
					Aircraft: 2,
					Ship:     6,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    2,
					Truck:    0,
					Bus:      0,
					Aircraft: 2,
					Ship:     3,
				},
			},
			{
				Id:          1,
				Name:        "Saint_anri Transport",
				YearStart:   1975,
				Value:       112629721,
				Money:       112317558,
				Income:      7861422,
				Performance: 535,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    90,
					Truck:    0,
					Bus:      0,
					Aircraft: 0,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    9,
					Truck:    1,
					Bus:      0,
					Aircraft: 0,
					Ship:     0,
				},
			},
			{
				Id:          2,
				Name:        "cicada3301 Transport",
				YearStart:   1940,
				Value:       1198229880,
				Money:       1189149017,
				Income:      61912524,
				Performance: 855,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    184,
					Truck:    0,
					Bus:      0,
					Aircraft: 18,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    34,
					Truck:    0,
					Bus:      9,
					Aircraft: 10,
					Ship:     0,
				},
			},
			{
				Id:          3,
				Name:        "New Chadingburg trans",
				YearStart:   1965,
				Value:       33972113,
				Money:       33872043,
				Income:      1508432,
				Performance: 553,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    35,
					Truck:    0,
					Bus:      0,
					Aircraft: 0,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    15,
					Truck:    0,
					Bus:      0,
					Aircraft: 0,
					Ship:     0,
				},
			},
			{
				Id:          4,
				Name:        "Gregory & Co.",
				YearStart:   1990,
				Value:       730434,
				Money:       684122,
				Income:      135459,
				Performance: 273,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    0,
					Truck:    0,
					Bus:      20,
					Aircraft: 0,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    0,
					Truck:    0,
					Bus:      2,
					Aircraft: 0,
					Ship:     0,
				},
			},
			{
				Id:          5,
				Name:        "MiSZU Transport",
				YearStart:   1949,
				Value:       181727748,
				Money:       179975478,
				Income:      8904479,
				Performance: 852,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    75,
					Truck:    0,
					Bus:      0,
					Aircraft: 28,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    44,
					Truck:    0,
					Bus:      0,
					Aircraft: 3,
					Ship:     0,
				},
			},
			{
				Id:          6,
				Name:        "vundo Transport",
				YearStart:   1983,
				Value:       5914383,
				Money:       5460685,
				Income:      49886,
				Performance: 346,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    13,
					Truck:    5,
					Bus:      7,
					Aircraft: 0,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    5,
					Truck:    2,
					Bus:      8,
					Aircraft: 0,
					Ship:     0,
				},
			},
			{
				Id:          7,
				Name:        "Transpote Transport",
				YearStart:   1945,
				Value:       1256656678,
				Money:       1237116792,
				Income:      60171064,
				Performance: 767,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    90,
					Truck:    25,
					Bus:      33,
					Aircraft: 87,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    15,
					Truck:    12,
					Bus:      12,
					Aircraft: 11,
					Ship:     0,
				},
			},
			{
				Id:          8,
				Name:        "Suttydroitwich Transport",
				YearStart:   1977,
				Value:       15203073,
				Money:       15090176,
				Income:      850421,
				Performance: 552,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    9,
					Truck:    11,
					Bus:      0,
					Aircraft: 6,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    14,
					Truck:    6,
					Bus:      5,
					Aircraft: 3,
					Ship:     0,
				},
			},
			{
				Id:          9,
				Name:        "Felix Transport",
				YearStart:   1978,
				Value:       13958548,
				Money:       10511697,
				Income:      1005518,
				Performance: 533,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    30,
					Truck:    0,
					Bus:      34,
					Aircraft: 0,
					Ship:     4,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    11,
					Truck:    2,
					Bus:      35,
					Aircraft: 0,
					Ship:     9,
				},
			},
			{
				Id:          10,
				Name:        "kalamiti korp",
				YearStart:   1949,
				Value:       36112943,
				Money:       35866407,
				Income:      1106070,
				Performance: 633,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    32,
					Truck:    12,
					Bus:      48,
					Aircraft: 2,
					Ship:     4,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    9,
					Truck:    6,
					Bus:      15,
					Aircraft: 4,
					Ship:     1,
				},
			},
			{
				Id:          11,
				Name:        "ErKoWo Transport",
				YearStart:   1984,
				Value:       28278381,
				Money:       21585197,
				Income:      18446744073708724623,
				Performance: 765,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    11,
					Truck:    0,
					Bus:      0,
					Aircraft: 100,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    4,
					Truck:    2,
					Bus:      6,
					Aircraft: 5,
					Ship:     0,
				},
			},
			{
				Id:          13,
				Name:        "Runar Transport",
				YearStart:   1998,
				Value:       1678013,
				Money:       1393448,
				Income:      641596,
				Performance: 282,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    0,
					Truck:    0,
					Bus:      0,
					Aircraft: 8,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    2,
					Truck:    0,
					Bus:      0,
					Aircraft: 2,
					Ship:     0,
				},
			},
			{
				Id:          14,
				Name:        "TBus_Express Transport",
				YearStart:   1991,
				Value:       190139,
				Money:       56763,
				Income:      10805,
				Performance: 144,
				Passworded:  true,
				Vehicles: global.OpenttdTypeCounts{
					Train:    0,
					Truck:    1,
					Bus:      18,
					Aircraft: 0,
					Ship:     0,
				},
				Stations: global.OpenttdTypeCounts{
					Train:    0,
					Truck:    2,
					Bus:      18,
					Aircraft: 0,
					Ship:     0,
				},
			},
		},
	}

	output := global.OpenttdServerState{}
	output.PopulateCompanyState(bytes.NewBuffer(testPacket))

	assert.EqualValues(t, expectObject, output)
}
