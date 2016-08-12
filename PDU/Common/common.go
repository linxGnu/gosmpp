package Common

func DecodeUnsigned(signed byte) int16 {
	if signed >= 0 {
		return int16(signed)
	}

	return int16(256 + int16(signed))
}

func DecodeUnsignedFromInt16(signed int16) int {
	if signed >= 0 {
		return int(signed)
	}

	return int(65536 + int32(signed))
}

func EncodeUnsigned(positive int16) byte {
	if positive < 128 {
		return byte(positive)
	}

	return byte(-(256 - positive))
}

func EncodeUnsignedFromInt(positive int) int16 {
	if positive < 32768 {
		return int16(positive)
	}

	return int16(-(65536 - positive))
}
