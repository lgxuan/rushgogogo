package utils

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

func ConvertToUTF8(input []byte) (string, error) {

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(input)
	if err != nil {
		return "", fmt.Errorf("charset detection error: %v", err)
	}

	var enc encoding.Encoding
	switch strings.ToLower(result.Charset) {
	case "utf-8":
		return string(input), nil
	case "gbk", "gb18030":
		enc = simplifiedchinese.GBK
	case "gb2312":
		enc = simplifiedchinese.HZGB2312
	case "big5":
		enc = traditionalchinese.Big5
	case "shift_jis":
		enc = japanese.ShiftJIS
	case "euc-kr":
		enc = korean.EUCKR
	case "windows-1252":
		enc = charmap.Windows1252
	case "iso-8859-1":
		enc = charmap.ISO8859_1
	default:
		return "", fmt.Errorf("unsupported charset: %s", result.Charset)
	}
	reader := transform.NewReader(bytes.NewReader(input), enc.NewDecoder())
	output, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("charset conversion error: %v", err)
	}

	return string(output), nil
}
