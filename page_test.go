package dexcom

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestUnmarshalPage(t *testing.T) {
	cases := []struct {
		pageType   PageType
		pageNumber int
		page       string
		records    []string
	}{

		{
			ManufacturingData, 0,
			"00 00 00 00 01 00 00 00 00 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 3A 7D B8 A7 2B 0B 37 37 2B 0B 3C 4D 61 6E 75 66 61 63 74 75 72 69 6E 67 50 61 72 61 6D 65 74 65 72 73 20 53 65 72 69 61 6C 4E 75 6D 62 65 72 3D 22 53 4D 34 34 37 39 32 36 37 35 22 20 48 61 72 64 77 61 72 65 50 61 72 74 4E 75 6D 62 65 72 3D 22 4D 54 32 30 36 34 39 22 20 48 61 72 64 77 61 72 65 52 65 76 69 73 69 6F 6E 3D 22 32 33 22 20 44 61 74 65 54 69 6D 65 43 72 65 61 74 65 64 3D 22 32 30 31 34 2D 31 32 2D 30 39 20 31 38 3A 32 35 3A 35 39 2E 32 34 36 20 2D 30 38 3A 30 30 22 20 48 61 72 64 77 61 72 65 49 64 3D 22 7B 30 38 35 38 45 45 46 39 2D 46 46 32 46 2D 34 41 45 31 2D 42 41 35 46 2D 45 34 46 46 39 30 46 36 39 36 30 42 7D 22 20 2F 3E 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 A7 FB",
			[]string{"B8 A7 2B 0B 37 37 2B 0B 3C 4D 61 6E 75 66 61 63 74 75 72 69 6E 67 50 61 72 61 6D 65 74 65 72 73 20 53 65 72 69 61 6C 4E 75 6D 62 65 72 3D 22 53 4D 34 34 37 39 32 36 37 35 22 20 48 61 72 64 77 61 72 65 50 61 72 74 4E 75 6D 62 65 72 3D 22 4D 54 32 30 36 34 39 22 20 48 61 72 64 77 61 72 65 52 65 76 69 73 69 6F 6E 3D 22 32 33 22 20 44 61 74 65 54 69 6D 65 43 72 65 61 74 65 64 3D 22 32 30 31 34 2D 31 32 2D 30 39 20 31 38 3A 32 35 3A 35 39 2E 32 34 36 20 2D 30 38 3A 30 30 22 20 48 61 72 64 77 61 72 65 49 64 3D 22 7B 30 38 35 38 45 45 46 39 2D 46 46 32 46 2D 34 41 45 31 2D 42 41 35 46 2D 45 34 46 46 39 30 46 36 39 36 30 42 7D 22 20 2F 3E 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"},
		},
		{
			SensorData, 469,
			"CD 2D 00 00 15 00 00 00 03 01 D5 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 C9 B6 E0 5D 4D 0C 55 D9 71 0D D0 8B 01 00 60 5E 01 00 B1 00 D0 02 0C 5F 4D 0C 81 DA 71 0D 40 C1 01 00 C0 7E 01 00 AB 00 72 86 38 60 4D 0C AD DB 71 0D A0 D3 01 00 50 A3 01 00 AA 00 9D B9 64 61 4D 0C D9 DC 71 0D 20 05 02 00 A0 C8 01 00 AE 00 FA 8B 90 62 4D 0C 05 DE 71 0D C0 31 02 00 40 EF 01 00 AC 00 11 32 BC 63 4D 0C 31 DF 71 0D 80 44 02 00 00 15 02 00 B3 00 6E 00 E8 64 4D 0C 5C E0 71 0D 00 46 02 00 80 35 02 00 B2 00 A0 67 14 66 4D 0C 88 E1 71 0D E0 5C 02 00 60 4B 02 00 A4 00 7A BD 40 67 4D 0C B5 E2 71 0D 60 69 02 00 E0 57 02 00 A9 01 BB F7 6C 68 4D 0C E0 E3 71 0D E0 A2 02 00 40 66 02 00 A8 00 9E BC 98 69 4D 0C 0C E5 71 0D C0 D2 02 00 A0 82 02 00 A3 00 CE 48 C4 6A 4D 0C 38 E6 71 0D 80 FB 02 00 20 AE 02 00 A4 00 6A C1 F0 6B 4D 0C 64 E7 71 0D 40 21 03 00 E0 E0 02 00 A1 00 1F 8A 1C 6D 4D 0C 90 E8 71 0D E0 42 03 00 60 0F 03 00 A5 00 37 23 48 6E 4D 0C BC E9 71 0D C0 5B 03 00 20 34 03 00 AD 00 33 8E 74 6F 4D 0C E8 EA 71 0D 80 72 03 00 60 50 03 00 A7 00 4A 45 A0 70 4D 0C 14 EC 71 0D 20 6C 03 00 80 64 03 00 AD 00 DA 12 CC 71 4D 0C 40 ED 71 0D 20 6F 03 00 C0 6F 03 00 AF 00 F7 95 F8 72 4D 0C 6C EE 71 0D 80 6B 03 00 60 72 03 00 AA 00 D9 6D 24 74 4D 0C 98 EF 71 0D 00 5A 03 00 40 6D 03 00 AC 00 05 7A 50 75 4D 0C C4 F0 71 0D A0 47 03 00 A0 62 03 00 B0 00 B4 32 FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF",
			[]string{
				"50 75 4D 0C C4 F0 71 0D A0 47 03 00 A0 62 03 00 B0 00",
				"24 74 4D 0C 98 EF 71 0D 00 5A 03 00 40 6D 03 00 AC 00",
				"F8 72 4D 0C 6C EE 71 0D 80 6B 03 00 60 72 03 00 AA 00",
				"CC 71 4D 0C 40 ED 71 0D 20 6F 03 00 C0 6F 03 00 AF 00",
				"A0 70 4D 0C 14 EC 71 0D 20 6C 03 00 80 64 03 00 AD 00",
				"74 6F 4D 0C E8 EA 71 0D 80 72 03 00 60 50 03 00 A7 00",
				"48 6E 4D 0C BC E9 71 0D C0 5B 03 00 20 34 03 00 AD 00",
				"1C 6D 4D 0C 90 E8 71 0D E0 42 03 00 60 0F 03 00 A5 00",
				"F0 6B 4D 0C 64 E7 71 0D 40 21 03 00 E0 E0 02 00 A1 00",
				"C4 6A 4D 0C 38 E6 71 0D 80 FB 02 00 20 AE 02 00 A4 00",
				"98 69 4D 0C 0C E5 71 0D C0 D2 02 00 A0 82 02 00 A3 00",
				"6C 68 4D 0C E0 E3 71 0D E0 A2 02 00 40 66 02 00 A8 00",
				"40 67 4D 0C B5 E2 71 0D 60 69 02 00 E0 57 02 00 A9 01",
				"14 66 4D 0C 88 E1 71 0D E0 5C 02 00 60 4B 02 00 A4 00",
				"E8 64 4D 0C 5C E0 71 0D 00 46 02 00 80 35 02 00 B2 00",
				"BC 63 4D 0C 31 DF 71 0D 80 44 02 00 00 15 02 00 B3 00",
				"90 62 4D 0C 05 DE 71 0D C0 31 02 00 40 EF 01 00 AC 00",
				"64 61 4D 0C D9 DC 71 0D 20 05 02 00 A0 C8 01 00 AE 00",
				"38 60 4D 0C AD DB 71 0D A0 D3 01 00 50 A3 01 00 AA 00",
				"0C 5F 4D 0C 81 DA 71 0D 40 C1 01 00 C0 7E 01 00 AB 00",
				"E0 5D 4D 0C 55 D9 71 0D D0 8B 01 00 60 5E 01 00 B1 00",
			},
		},
		{
			EGVData, 312,
			"50 2E 00 00 17 00 00 00 04 02 38 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 8B FA 89 5B 4D 0C FD D6 71 0D 4F 00 14 87 FF B5 5C 4D 0C 29 D8 71 0D 58 00 13 EC 59 E1 5D 4D 0C 55 D9 71 0D 5F 00 13 3C 60 0D 5F 4D 0C 81 DA 71 0D 6E 00 12 86 86 39 60 4D 0C AD DB 71 0D 74 00 13 12 3E 65 61 4D 0C D9 DC 71 0D 82 00 12 43 7F 91 62 4D 0C 05 DE 71 0D 8F 00 12 E7 77 BD 63 4D 0C 31 DF 71 0D 95 00 12 58 18 E9 64 4D 0C 5D E0 71 0D 95 00 13 D2 5F 14 66 4D 0C 89 E1 71 0D 9C 00 14 9E B3 41 67 4D 0C B5 E2 71 0D 9F 00 14 6B D8 6C 68 4D 0C E1 E3 71 0D B0 00 13 C1 BB 98 69 4D 0C 0D E5 71 0D BE 00 12 B2 7B C5 6A 4D 0C 39 E6 71 0D CA 00 12 76 A2 F0 6B 4D 0C 65 E7 71 0D D5 00 12 74 29 1C 6D 4D 0C 91 E8 71 0D DF 00 12 9B 22 48 6E 4D 0C BD E9 71 0D E6 00 13 D3 E1 74 6F 4D 0C E9 EA 71 0D ED 00 13 58 97 A0 70 4D 0C 15 EC 71 0D EB 00 14 2F AE CC 71 4D 0C 41 ED 71 0D EC 00 14 6D 37 F8 72 4D 0C 6D EE 71 0D EB 00 14 4F 5F 24 74 4D 0C 99 EF 71 0D E6 00 14 D4 AE 50 75 4D 0C C5 F0 71 0D E0 00 14 9F FE FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF",
			[]string{
				"50 75 4D 0C C5 F0 71 0D E0 00 14",
				"24 74 4D 0C 99 EF 71 0D E6 00 14",
				"F8 72 4D 0C 6D EE 71 0D EB 00 14",
				"CC 71 4D 0C 41 ED 71 0D EC 00 14",
				"A0 70 4D 0C 15 EC 71 0D EB 00 14",
				"74 6F 4D 0C E9 EA 71 0D ED 00 13",
				"48 6E 4D 0C BD E9 71 0D E6 00 13",
				"1C 6D 4D 0C 91 E8 71 0D DF 00 12",
				"F0 6B 4D 0C 65 E7 71 0D D5 00 12",
				"C5 6A 4D 0C 39 E6 71 0D CA 00 12",
				"98 69 4D 0C 0D E5 71 0D BE 00 12",
				"6C 68 4D 0C E1 E3 71 0D B0 00 13",
				"41 67 4D 0C B5 E2 71 0D 9F 00 14",
				"14 66 4D 0C 89 E1 71 0D 9C 00 14",
				"E9 64 4D 0C 5D E0 71 0D 95 00 13",
				"BD 63 4D 0C 31 DF 71 0D 95 00 12",
				"91 62 4D 0C 05 DE 71 0D 8F 00 12",
				"65 61 4D 0C D9 DC 71 0D 82 00 12",
				"39 60 4D 0C AD DB 71 0D 74 00 13",
				"0D 5F 4D 0C 81 DA 71 0D 6E 00 12",
				"E1 5D 4D 0C 55 D9 71 0D 5F 00 13",
				"B5 5C 4D 0C 29 D8 71 0D 58 00 13",
				"89 5B 4D 0C FD D6 71 0D 4F 00 14",
			},
		},
		{
			CalibrationData, 1432,
			"30 0B 00 00 02 00 00 00 05 03 98 05 00 00 00 00 00 00 00 00 00 00 00 00 00 00 60 8F 49 3B C1 0F 32 CF C0 0F 0E 86 1B 91 2A F5 8E 40 E4 52 D5 9E 9C F8 E5 40 00 00 00 00 00 00 F0 3F 03 06 03 00 00 00 00 00 00 00 00 08 D6 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 E2 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 50 57 BD 0F B0 00 00 00 BF 93 03 00 67 58 BD 0F 00 FE 12 BE 0F 7E 00 00 00 71 96 02 00 E5 13 BE 0F 00 3E A7 BE 0F D0 00 00 00 82 AB 03 00 B7 A8 BE 0F 00 61 53 BF 0F 7A 00 00 00 80 84 02 00 F9 54 BF 0F 00 9C 70 C0 0F 50 00 00 00 A0 ED 01 00 BA 71 C0 0F 00 C3 39 C1 0F 89 00 00 00 C0 C6 02 00 47 3B C1 0F 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 33 33 33 33 33 33 D3 3F 55 55 55 55 55 55 E5 3F 63 FF 90 FE C1 0F 79 92 C1 0F 7E EC 2A 95 64 0D 90 40 8E 65 F7 FE 13 83 E3 40 00 00 00 00 00 00 F0 3F 03 06 08 00 00 00 00 00 00 00 00 08 E2 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 50 57 BD 0F B0 00 00 00 BF 93 03 00 67 58 BD 0F 00 FE 12 BE 0F 7E 00 00 00 71 96 02 00 E5 13 BE 0F 00 3E A7 BE 0F D0 00 00 00 82 AB 03 00 B7 A8 BE 0F 00 61 53 BF 0F 7A 00 00 00 80 84 02 00 F9 54 BF 0F 00 9C 70 C0 0F 50 00 00 00 A0 ED 01 00 BA 71 C0 0F 00 C3 39 C1 0F 89 00 00 00 C0 C6 02 00 47 3B C1 0F 00 6D FE C1 0F 83 00 00 00 A0 AA 02 00 CD FD C1 0F 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 33 33 33 33 33 33 D3 3F 55 55 55 55 55 55 E5 3F 97 BB FF FF",
			[]string{
				"90 FE C1 0F 79 92 C1 0F 7E EC 2A 95 64 0D 90 40 8E 65 F7 FE 13 83 E3 40 00 00 00 00 00 00 F0 3F 03 06 08 00 00 00 00 00 00 00 00 08 E2 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 50 57 BD 0F B0 00 00 00 BF 93 03 00 67 58 BD 0F 00 FE 12 BE 0F 7E 00 00 00 71 96 02 00 E5 13 BE 0F 00 3E A7 BE 0F D0 00 00 00 82 AB 03 00 B7 A8 BE 0F 00 61 53 BF 0F 7A 00 00 00 80 84 02 00 F9 54 BF 0F 00 9C 70 C0 0F 50 00 00 00 A0 ED 01 00 BA 71 C0 0F 00 C3 39 C1 0F 89 00 00 00 C0 C6 02 00 47 3B C1 0F 00 6D FE C1 0F 83 00 00 00 A0 AA 02 00 CD FD C1 0F 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 33 33 33 33 33 33 D3 3F 55 55 55 55 55 55 E5 3F",
				"49 3B C1 0F 32 CF C0 0F 0E 86 1B 91 2A F5 8E 40 E4 52 D5 9E 9C F8 E5 40 00 00 00 00 00 00 F0 3F 03 06 03 00 00 00 00 00 00 00 00 08 D6 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 E2 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 50 57 BD 0F B0 00 00 00 BF 93 03 00 67 58 BD 0F 00 FE 12 BE 0F 7E 00 00 00 71 96 02 00 E5 13 BE 0F 00 3E A7 BE 0F D0 00 00 00 82 AB 03 00 B7 A8 BE 0F 00 61 53 BF 0F 7A 00 00 00 80 84 02 00 F9 54 BF 0F 00 9C 70 C0 0F 50 00 00 00 A0 ED 01 00 BA 71 C0 0F 00 C3 39 C1 0F 89 00 00 00 C0 C6 02 00 47 3B C1 0F 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 33 33 33 33 33 33 D3 3F 55 55 55 55 55 55 E5 3F",
			},
		},
		{
			CalibrationData, 1432,
			// The one record in this page has CRC = 63 FF,
			// followed by FF bytes for padding.
			"30 0B 00 00 01 00 00 00 05 03 98 05 00 00 00 00 00 00 00 00 00 00 00 00 00 00 08 39 49 3B C1 0F 32 CF C0 0F 0E 86 1B 91 2A F5 8E 40 E4 52 D5 9E 9C F8 E5 40 00 00 00 00 00 00 F0 3F 03 06 03 00 00 00 00 00 00 00 00 08 D6 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 E2 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 50 57 BD 0F B0 00 00 00 BF 93 03 00 67 58 BD 0F 00 FE 12 BE 0F 7E 00 00 00 71 96 02 00 E5 13 BE 0F 00 3E A7 BE 0F D0 00 00 00 82 AB 03 00 B7 A8 BE 0F 00 61 53 BF 0F 7A 00 00 00 80 84 02 00 F9 54 BF 0F 00 9C 70 C0 0F 50 00 00 00 A0 ED 01 00 BA 71 C0 0F 00 C3 39 C1 0F 89 00 00 00 C0 C6 02 00 47 3B C1 0F 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 33 33 33 33 33 33 D3 3F 55 55 55 55 55 55 E5 3F 63 FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF",
			[]string{
				"49 3B C1 0F 32 CF C0 0F 0E 86 1B 91 2A F5 8E 40 E4 52 D5 9E 9C F8 E5 40 00 00 00 00 00 00 F0 3F 03 06 03 00 00 00 00 00 00 00 00 08 D6 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 E2 92 BC 0F 62 00 00 00 B7 1D 02 00 89 93 BC 0F 00 50 57 BD 0F B0 00 00 00 BF 93 03 00 67 58 BD 0F 00 FE 12 BE 0F 7E 00 00 00 71 96 02 00 E5 13 BE 0F 00 3E A7 BE 0F D0 00 00 00 82 AB 03 00 B7 A8 BE 0F 00 61 53 BF 0F 7A 00 00 00 80 84 02 00 F9 54 BF 0F 00 9C 70 C0 0F 50 00 00 00 A0 ED 01 00 BA 71 C0 0F 00 C3 39 C1 0F 89 00 00 00 C0 C6 02 00 47 3B C1 0F 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 33 33 33 33 33 33 D3 3F 55 55 55 55 55 55 E5 3F",
			},
		},
	}
	for _, c := range cases {
		page, err := unmarshalPage(hexdata(c.page))
		if err != nil {
			t.Errorf("UnmarshalPage returned %v", err)
		}
		if page.Type != c.pageType {
			t.Errorf("UnmarshalPage: page type == %d, want %d", page.Type, c.pageType)
		}
		if page.Number != c.pageNumber {
			t.Errorf("UnmarshalPage: page number == %d, want %d", page.Number, c.pageNumber)
		}
		if len(page.Records) != len(c.records) {
			t.Errorf("UnmarshalPage: #records == %d, want %d", len(page.Records), len(c.records))
		}
		for i, v := range page.Records {
			r := hexdata(c.records[i])
			if !bytes.Equal(v, r) {
				t.Errorf("UnmarshalPage: record #%d == % X, want % X", i, v, r)
			}
		}
	}
}

func hexdata(str string) []byte {
	fields := strings.Fields(str)
	data := make([]byte, len(fields))
	for i, s := range fields {
		b, err := strconv.ParseUint(string(s), 16, 8)
		if err != nil {
			panic(err)
		}
		data[i] = byte(b)
	}
	return data
}

/*
ManufacturingData page 0


[
  {
    "Timestamp": {
      "SystemTime": "2014-12-10T02:26:00-05:00",
      "DisplayTime": "2014-12-09T18:25:59-05:00"
    },
    "XML": {
      "DateTimeCreated": "2014-12-09 18:25:59.246 -08:00",
      "HardwareId": "{0858EEF9-FF2F-4AE1-BA5F-E4FF90F6960B}",
      "HardwarePartNumber": "MT20649",
      "HardwareRevision": "23",
      "SerialNumber": "SM44792675"
    }
  }
]

[
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T22:08:16-04:00",
      "DisplayTime": "2016-02-24T18:36:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 214944,
      "Filtered": 221856,
      "RSSI": -80,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T22:03:16-04:00",
      "DisplayTime": "2016-02-24T18:31:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 219648,
      "Filtered": 224576,
      "RSSI": -84,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:58:16-04:00",
      "DisplayTime": "2016-02-24T18:26:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 224128,
      "Filtered": 225888,
      "RSSI": -86,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:53:16-04:00",
      "DisplayTime": "2016-02-24T18:21:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 225056,
      "Filtered": 225216,
      "RSSI": -81,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:48:16-04:00",
      "DisplayTime": "2016-02-24T18:16:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 224288,
      "Filtered": 222336,
      "RSSI": -83,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:43:16-04:00",
      "DisplayTime": "2016-02-24T18:11:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 225920,
      "Filtered": 217184,
      "RSSI": -89,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:38:16-04:00",
      "DisplayTime": "2016-02-24T18:06:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 220096,
      "Filtered": 209952,
      "RSSI": -83,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:33:16-04:00",
      "DisplayTime": "2016-02-24T18:01:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 213728,
      "Filtered": 200544,
      "RSSI": -91,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:28:16-04:00",
      "DisplayTime": "2016-02-24T17:56:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 205120,
      "Filtered": 188640,
      "RSSI": -95,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:23:16-04:00",
      "DisplayTime": "2016-02-24T17:51:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 195456,
      "Filtered": 175648,
      "RSSI": -92,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:18:16-04:00",
      "DisplayTime": "2016-02-24T17:46:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 185024,
      "Filtered": 164512,
      "RSSI": -93,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:13:16-04:00",
      "DisplayTime": "2016-02-24T17:41:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 172768,
      "Filtered": 157248,
      "RSSI": -88,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:08:16-04:00",
      "DisplayTime": "2016-02-24T17:36:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 158048,
      "Filtered": 153568,
      "RSSI": -87,
      "Unknown": 1
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:03:16-04:00",
      "DisplayTime": "2016-02-24T17:31:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 154848,
      "Filtered": 150368,
      "RSSI": -92,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:58:16-04:00",
      "DisplayTime": "2016-02-24T17:26:52-05:00"
    },
    "Sensor": {
      "Unfiltered": 148992,
      "Filtered": 144768,
      "RSSI": -78,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:53:16-04:00",
      "DisplayTime": "2016-02-24T17:21:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 148608,
      "Filtered": 136448,
      "RSSI": -77,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:48:16-04:00",
      "DisplayTime": "2016-02-24T17:16:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 143808,
      "Filtered": 126784,
      "RSSI": -84,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:43:16-04:00",
      "DisplayTime": "2016-02-24T17:11:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 132384,
      "Filtered": 116896,
      "RSSI": -82,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:38:16-04:00",
      "DisplayTime": "2016-02-24T17:06:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 119712,
      "Filtered": 107344,
      "RSSI": -86,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:33:16-04:00",
      "DisplayTime": "2016-02-24T17:01:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 115008,
      "Filtered": 97984,
      "RSSI": -85,
      "Unknown": 0
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:28:16-04:00",
      "DisplayTime": "2016-02-24T16:56:53-05:00"
    },
    "Sensor": {
      "Unfiltered": 101328,
      "Filtered": 89696,
      "RSSI": -79,
      "Unknown": 0
    }
  }
]



[
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T22:08:16-04:00",
      "DisplayTime": "2016-02-24T18:36:53-05:00"
    },
    "EGV": {
      "Glucose": 224,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T22:03:16-04:00",
      "DisplayTime": "2016-02-24T18:31:53-05:00"
    },
    "EGV": {
      "Glucose": 230,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:58:16-04:00",
      "DisplayTime": "2016-02-24T18:26:53-05:00"
    },
    "EGV": {
      "Glucose": 235,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:53:16-04:00",
      "DisplayTime": "2016-02-24T18:21:53-05:00"
    },
    "EGV": {
      "Glucose": 236,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:48:16-04:00",
      "DisplayTime": "2016-02-24T18:16:53-05:00"
    },
    "EGV": {
      "Glucose": 235,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:43:16-04:00",
      "DisplayTime": "2016-02-24T18:11:53-05:00"
    },
    "EGV": {
      "Glucose": 237,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:38:16-04:00",
      "DisplayTime": "2016-02-24T18:06:53-05:00"
    },
    "EGV": {
      "Glucose": 230,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:33:16-04:00",
      "DisplayTime": "2016-02-24T18:01:53-05:00"
    },
    "EGV": {
      "Glucose": 223,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:28:16-04:00",
      "DisplayTime": "2016-02-24T17:56:53-05:00"
    },
    "EGV": {
      "Glucose": 213,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:23:17-04:00",
      "DisplayTime": "2016-02-24T17:51:53-05:00"
    },
    "EGV": {
      "Glucose": 202,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:18:16-04:00",
      "DisplayTime": "2016-02-24T17:46:53-05:00"
    },
    "EGV": {
      "Glucose": 190,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:13:16-04:00",
      "DisplayTime": "2016-02-24T17:41:53-05:00"
    },
    "EGV": {
      "Glucose": 176,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:08:17-04:00",
      "DisplayTime": "2016-02-24T17:36:53-05:00"
    },
    "EGV": {
      "Glucose": 159,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T21:03:16-04:00",
      "DisplayTime": "2016-02-24T17:31:53-05:00"
    },
    "EGV": {
      "Glucose": 156,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:58:17-04:00",
      "DisplayTime": "2016-02-24T17:26:53-05:00"
    },
    "EGV": {
      "Glucose": 149,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:53:17-04:00",
      "DisplayTime": "2016-02-24T17:21:53-05:00"
    },
    "EGV": {
      "Glucose": 149,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:48:17-04:00",
      "DisplayTime": "2016-02-24T17:16:53-05:00"
    },
    "EGV": {
      "Glucose": 143,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:43:17-04:00",
      "DisplayTime": "2016-02-24T17:11:53-05:00"
    },
    "EGV": {
      "Glucose": 130,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:38:17-04:00",
      "DisplayTime": "2016-02-24T17:06:53-05:00"
    },
    "EGV": {
      "Glucose": 116,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:33:17-04:00",
      "DisplayTime": "2016-02-24T17:01:53-05:00"
    },
    "EGV": {
      "Glucose": 110,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 2
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:28:17-04:00",
      "DisplayTime": "2016-02-24T16:56:53-05:00"
    },
    "EGV": {
      "Glucose": 95,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:23:17-04:00",
      "DisplayTime": "2016-02-24T16:51:53-05:00"
    },
    "EGV": {
      "Glucose": 88,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 3
    }
  },
  {
    "Timestamp": {
      "SystemTime": "2015-07-17T20:18:17-04:00",
      "DisplayTime": "2016-02-24T16:46:53-05:00"
    },
    "EGV": {
      "Glucose": 79,
      "DisplayOnly": false,
      "Noise": 1,
      "Trend": 4
    }
  }
]
*/
