package libatum

import (
	"encoding/binary"
	"encoding/hex"
	"os"
)

func readHex(file *os.File, offset int64, size int64, whence int) (string, error) {
	_, err := file.Seek(offset, whence)
	if err != nil {
		return "", err
	}

	s := make([]byte, size)
	_, err = file.Read(s)
	if err != nil {
		return "", err
	}

	d := hex.EncodeToString(s)

	return d, nil
}

func getHexBytes(in string) ([]byte, error) {
	d, err := hex.DecodeString(in)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func toBinary32(in int32) []byte {
	out := make([]byte, binary.Size(in))
	binary.LittleEndian.PutUint32(out, uint32(in))

	return out
}

func toBinary64(in int64) []byte {
	out := make([]byte, binary.Size(in))
	binary.LittleEndian.PutUint64(out, uint64(in))

	return out
}

func getCNMTType(val string) string {
	switch val {
	case "80":
		return "Application"
	case "81":
		return "Patch"
	case "82":
		return "AddOnContent"
	case "83":
		return "Delta"
	}

	return ""
}

func getNCAType(val string) string {
	switch val {
	case "00":
		return "Meta"
	case "01":
		return "Program"
	case "02":
		return "Data"
	case "03":
		return "Control"
	case "04":
		return "HtmlDocument"
	case "05":
		return "LegalInformation"
	case "06":
		return "DeltaFragment"
	}

	return ""
}

func sum(array []int) int {
	n := 0
	for _, v := range array {
		n += v
	}

	return n
}

func sum64(array []int64) int {
	n := 0
	for _, v := range array {
		n += int(v)
	}

	return n
}
