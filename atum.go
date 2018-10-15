package libhac

import (
	_ "crypto/sha256"
	_ "crypto/sha512"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (c *HacClient) download(url, path string) error {
	resp, err := c.DoRequest("GET", url, []tls.Certificate{c.DeviceCert}, false, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *HacClient) TestEdgeToken() error {
	id, err := c.GetCNMTID("0100000000010000", 0)
	if err != nil || id == "" {
		return errors.New("edge token is invalid!")
	}

	return nil
}

func (c *HacClient) GetCNMTID(tid string, ver int) (string, error) {
	resp, err := c.DoRequest("HEAD", fmt.Sprintf("https://atum.hac.lp1.d4c.nintendo.net/t/a/%s/%d", tid, ver),
		[]tls.Certificate{c.DeviceCert}, false, true)
	if err != nil {
		return "", err
	}

	cnmtID := resp.Header.Get("X-Nintendo-Content-ID")

	if cnmtID == "" {
		return "", errors.New("title not on cdn")
	}

	return cnmtID, nil
}

func (c *HacClient) DownloadCNMT(cnmtID string, out string) error {
	err := c.download(fmt.Sprintf("https://atum.hac.lp1.d4c.nintendo.net/c/a/%s", cnmtID), out)
	if err != nil {
		return err
	}

	return nil
}

func DecryptNCA(path, out, hactoolPath string) error {
	err := os.MkdirAll(out, 0700)
	if err != nil {
		return err
	}

	err = exec.Command(hactoolPath, "--exefsdir="+out+"/exefs", "--romfsdir="+out+"/romfs",
		"--section0dir="+out+"/section0", "--section1dir="+out+"/section1",
		"--section2dir="+out+"/section2", "--section3dir="+out+"/section3",
		"--header="+out+"/header.bin", path).Run()
	if err != nil {
		return err
	}

	return nil
}

func ParseCNMT(path, headerPath string) (CNMT, error) {
	cnmt, err := os.Open(path)
	if err != nil {
		return CNMT{}, err
	}
	defer cnmt.Close()

	t, err := readHex(cnmt, 0xC, 1, 0)
	if err != nil {
		return CNMT{}, err
	}

	tid, err := readHex(cnmt, 0x0, 8, 0)
	if err != nil {
		return CNMT{}, err
	}

	version, err := readHex(cnmt, 0x8, 4, 0)
	if err != nil {
		return CNMT{}, err
	}

	sysv, err := readHex(cnmt, 0x28, 8, 0)
	if err != nil {
		return CNMT{}, err
	}

	dlsysv, err := readHex(cnmt, 0x18, 8, 0)
	if err != nil {
		return CNMT{}, err
	}

	digest, err := readHex(cnmt, -0x20, 0x20, 2)
	if err != nil {
		return CNMT{}, err
	}

	to, err := readHex(cnmt, 0xE, 1, 0)
	if err != nil {
		return CNMT{}, err
	}

	tableOffset, err := strconv.ParseInt(to, 16, 64)
	if err != nil {
		return CNMT{}, err
	}

	cec, err := readHex(cnmt, 0x10, 1, 0)
	if err != nil {
		return CNMT{}, err
	}

	contentEntryCount, err := strconv.ParseInt(cec, 16, 64)
	if err != nil {
		return CNMT{}, err
	}

	ces := []ContentEntry{}
	var i int64
	for i = 0; i < contentEntryCount; i++ {
		offset := 0x20 + tableOffset + 0x38*i

		hash, err := readHex(cnmt, offset, 32, 0)
		if err != nil {
			return CNMT{}, err
		}

		id, err := readHex(cnmt, offset+0x20, 16, 0)
		if err != nil {
			return CNMT{}, err
		}

		size, err := readHex(cnmt, offset+0x30, 6, 0)
		if err != nil {
			return CNMT{}, err
		}

		ty, err := readHex(cnmt, offset+0x36, 1, 0)
		if err != nil {
			return CNMT{}, err
		}

		ces = append(ces, ContentEntry{
			hash,
			id,
			size,
			getNCAType(ty),
		})
	}

	header, err := os.Open(headerPath)
	if err != nil {
		return CNMT{}, err
	}
	defer header.Close()

	mKeyRev, err := readHex(header, 0x220, 0x1, 0)
	if err != nil {
		return CNMT{}, err
	}

	return CNMT{
		path,
		getCNMTType(t),
		tid,
		version,
		sysv,
		dlsysv,
		digest,
		mKeyRev,
		ces,
	}, nil
}

func (c *HacClient) DownloadContentEntry(ce ContentEntry, out string) error {
	err := c.download(fmt.Sprintf("https://atum.hac.lp1.d4c.nintendo.net/c/c/%s", ce.ID), out)
	if err != nil {
		return err
	}

	return nil
}

func GetRightsID(tid, mKeyRev string) string {
	return fmt.Sprintf("%s%s%s", tid, strings.Repeat("0", 16-len(mKeyRev)),
		mKeyRev)
}

func (c *HacClient) DownloadCetk(rightsID, out string) error {
	err := c.download(fmt.Sprintf("https://atum.hac.lp1.d4c.nintendo.net/r/t/%s", rightsID),
		out)
	if err != nil {
		return err
	}

	return nil
}

func GetTitleKeyFromCetk(path string) (string, error) {
	cetk, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer cetk.Close()

	tk, err := readHex(cetk, 0x180, 16, 0)
	if err != nil {
		return "", err
	}

	return tk, nil
}

func GenerateTicket(in []byte, titleKey, mKeyRev, rightsID, out string) error {
	tk, err := getHexBytes(titleKey)
	if err != nil {
		return err
	}

	mkr, err := getHexBytes(mKeyRev)
	if err != nil {
		return err
	}

	rid, err := getHexBytes(rightsID)
	if err != nil {
		return err
	}

	in[0x180] = tk[0]
	in[0x181] = tk[1]
	in[0x182] = tk[2]
	in[0x183] = tk[3]
	in[0x184] = tk[4]
	in[0x185] = tk[5]
	in[0x186] = tk[6]
	in[0x187] = tk[7]
	in[0x188] = tk[8]
	in[0x189] = tk[9]
	in[0x18A] = tk[10]
	in[0x18B] = tk[11]
	in[0x18C] = tk[12]
	in[0x18D] = tk[13]
	in[0x18E] = tk[14]
	in[0x18F] = tk[15]

	in[0x285] = mkr[0]

	in[0x2A0] = rid[0]
	in[0x2A1] = rid[1]
	in[0x2A2] = rid[2]
	in[0x2A3] = rid[3]
	in[0x2A4] = rid[4]
	in[0x2A5] = rid[5]
	in[0x2A6] = rid[6]
	in[0x2A7] = rid[7]
	in[0x2A8] = rid[8]
	in[0x2A9] = rid[9]
	in[0x2AA] = rid[10]
	in[0x2AB] = rid[11]
	in[0x2AC] = rid[12]
	in[0x2AD] = rid[13]
	in[0x2AE] = rid[14]
	in[0x2AF] = rid[15]

	tik, err := os.Create(out)
	if err != nil {
		return err
	}
	defer tik.Close()

	_, err = tik.Write(in)
	if err != nil {
		return err
	}

	return nil
}

func PackToNSP(path, out string) error {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	n := []string{}
	for _, v := range dir {
		n = append(n, v.Name())
	}

	stringTable := strings.Join(n, "\x00")
	headerSize := 0x10 + (len(dir) * 0x18) + len(stringTable)
	remainder := 0x10 - headerSize%0x10
	headerSize += remainder

	fileSizes := []int64{}
	for _, v := range dir {
		fileSizes = append(fileSizes, v.Size())
	}

	fileOffsets := []int{}

	for i := 0; i < len(dir); i++ {
		fileOffsets = append(fileOffsets, sum64(fileSizes[:i]))
	}

	fileNameLengths := []int{}
	for _, v := range dir {
		fileNameLengths = append(fileNameLengths, len(v.Name())+1)
	}

	stringTableOffsets := []int{}
	for i := 0; i < len(dir); i++ {
		stringTableOffsets = append(stringTableOffsets, sum(fileNameLengths[:i]))
	}

	header := [][]byte{[]byte("PFS0"),
		toBinary32(int32(len(dir))),
		toBinary32(int32(len(stringTable) + remainder)),
		[]byte("\x00\x00\x00\x00"),
	}

	for i := 0; i < len(dir); i++ {
		header = append(header, toBinary64(int64(fileOffsets[i])))
		header = append(header, toBinary64(fileSizes[i]))
		header = append(header, toBinary32(int32(stringTableOffsets[i])))
		header = append(header, []byte("\x00\x00\x00\x00"))
	}

	header = append(header, []byte(stringTable))
	for i := 1; i <= remainder; i++ {
		header = append(header, []byte("\x00"))
	}

	nsp, err := os.Create(out)
	if err != nil {
		return err
	}
	defer nsp.Close()

	for _, v := range header {
		_, err = nsp.Write(v)
		if err != nil {
			return err
		}
	}

	for _, v := range dir {
		f, err := os.Open(fmt.Sprintf("%s/%s", path, v.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(nsp, f)
		if err != nil {
			return err
		}
	}

	return nil
}
