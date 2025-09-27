package georesolver

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// fileMD5 computes the MD5 hash of the given file path.
func fileMD5(path string) (string, error) {
	err := ensureDir(path)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write([]byte(path))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// bytesToUint128 converts a 16-byte array to a uint128 structure.
func bytesToUint128(b []byte) uint128{
	if len(b) != 16 {
		return uint128{}
	}
	hi := uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	lo := uint64(b[8])<<56 | uint64(b[9])<<48 | uint64(b[10])<<40 | uint64(b[11])<<32 |
		uint64(b[12])<<24 | uint64(b[13])<<16 | uint64(b[14])<<8 | uint64(b[15])
	return uint128{Hi: hi, Lo: lo}
}

// uint128ToBytes converts a uint128 structure to a 16-byte array.
func uint128ToBytes(u uint128) []byte {
	b := make([]byte, 16)

	b[0] = byte(u.Hi >> 56)
	b[1] = byte(u.Hi >> 48)
	b[2] = byte(u.Hi >> 40)
	b[3] = byte(u.Hi >> 32)
	b[4] = byte(u.Hi >> 24)
	b[5] = byte(u.Hi >> 16)
	b[6] = byte(u.Hi >> 8)
	b[7] = byte(u.Hi)
	b[8] = byte(u.Lo >> 56)
	b[9] = byte(u.Lo >> 48)
	b[10] = byte(u.Lo >> 40)
	b[11] = byte(u.Lo >> 32)
	b[12] = byte(u.Lo >> 24)
	b[13] = byte(u.Lo >> 16)
	b[14] = byte(u.Lo >> 8)
	b[15] = byte(u.Lo)
	return b
}

// encodeRangeValue encodes a IP Range structure into two byte arrays.
func encodeRangeValue(r Range) ([]byte, []byte) {
	end := uint128ToBytes(r.End)
	iso := []byte(r.IOS2)


	return uint128ToBytes(r.Start), append(end, iso...)
}

// decodeRangeValue decodes byte arrays into a IP Range structure.
func decodeRangeValue(key []byte, value []byte) (Range, error){
	if len(value) < 16 {
		return nil, ErrInvalidIP
	}

	start := bytesToUint128(key)
	end := bytesToUint128(value[:16])
	iso := string(value[16:])

	return Range{Start: start, End: end, ISO2: iso}, nil
}

// ensureDir checks if the path exists and returns an error if it does not.
func ensureDir(path string) error{
	_, err := os.Stat(path)
	return err
}