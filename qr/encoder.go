package qr

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	MIN_Version = 1
	MAX_Version = 40
)

type Ecc int // error correction code level
type Ecm int // encoding mode like numeric,alphanumeric,byte,kanji

const (
	low      Ecc = iota // 7% of codewords can be restored
	medium              // 15% of codewords can be restored
	quartile            // 25% of codewords can be restored
	high                // 30% of codewords can be restored
)

const (
	Numeric      Ecm = iota //0, numeric mode
	Alphanumeric            //1, alphanumeric mode
	Byte                    //2, byte mode
	Kanji                   //3, kanji mode
)

type EncoderData struct {
	Data         []byte
	EncodingMode *Ecm
	Ecl          Ecc
	Datalength   int
	Version      *int
}

func (dt *EncoderData) Encoder() string {
	//for byte mode, mode indicator is 0100
	var modeIndicator string
	if dt.EncodingMode != nil {
		modeIndicator = getModeIndicator(*dt.EncodingMode)
	} else {
		bytemode := Byte // if any Encoding mode is not given consider Byte as default Encoding mode
		dt.EncodingMode = &bytemode
		modeIndicator = getModeIndicator(bytemode)
	}

	//if Version number is not given detect the Version number based on the datalength and ecl
	if dt.Version == nil {
		Version, err := detectWhichVersionToUse(dt.Datalength, dt.Ecl)
		if err != nil {
			log.Fatalf("Failed to detect Version: %v\n", err)
		}
		dt.Version = &Version
	}

	//number of bits required for the sizes of the character count indicator for given Version and Encoding mode
	noOfBits := charCountIndicator(*dt.Version, *dt.EncodingMode)
	binaryInputTextLength := fmt.Sprintf("%0*b", noOfBits, dt.Datalength)
	var encodedText strings.Builder
	encodedText.WriteString(modeIndicator)
	encodedText.WriteString(binaryInputTextLength)
	for i := 0; i < dt.Datalength; i++ {
		encodedText.WriteString(fmt.Sprintf("%08b", dt.Data[i]))
	}
	noOfDataBits := getDataBytesStored(*dt.Version, dt.Ecl) * 8
	if (encodedText.Len()) > noOfDataBits {
		panic("input text is too long!")
	}
	//for 7-Q Version total of 88 bytes means 88*8=704 bits are required
	//adding 4 bit terminator (0000) to the end

	paddingLength := noOfDataBits - (encodedText.Len())
	if paddingLength >= 4 {
		encodedText.WriteString("0000")
	} else {
		encodedText.WriteString(strings.Repeat("0", paddingLength))
	}
	//adding 0 as pad bytes for making encodedText multiple of 8
	encodedText.WriteString(strings.Repeat("0", ((encodedText.Len()) % 8)))

	//adding 11101100 00010001 untill encodedText reach upto 704 bits
	for encodedText.Len() != noOfDataBits {
		encodedText.WriteString("11101100")
		if encodedText.Len() == noOfDataBits {
			break
		}
		encodedText.WriteString("00010001")
	}

	// fmt.Println(inputTextLength)
	fmt.Printf("%d\n", (encodedText.Len()))
	fmt.Println(*dt.Version)
	// fmt.Println(inputText)
	// fmt.Println(errorCorrectionMode)
	return encodedText.String()

}

func getModeIndicator(mode Ecm) string {
	modeMap := map[Ecm]string{
		Numeric:      "0001",
		Alphanumeric: "0010",
		Byte:         "0100",
		Kanji:        "1000",
	}

	if value, exists := modeMap[mode]; exists {
		return value
	}
	panic("mode not found")
}

// get the number of data bits that can be stored in a QR code of given Version number, after
// all functions modules are excluded.
func getNumRawDataModules(ver int) int {
	if ver < MIN_Version || ver > MAX_Version {
		panic("Version number out of range")
	}
	size := ver*4 + 17
	result := size * size     //number of module in qr code
	result -= 8 * 8 * 3       //subtract the three finder pattern with separators
	result -= (size - 16) * 2 //subrtract the timing pattern
	result -= (15 * 2) + 1    //subtract dark module and format information area

	if ver >= 2 {
		alp := 2 + (ver / 7)                 //total number of alignment pattern location
		result -= (alp - 1) * (alp - 1) * 25 //subtract alignment patterns not overlapping timing pattern
		result -= (alp - 2) * 2 * 20         // subtract alignment pattern which overlap timing pattern
		if ver >= 7 {
			result -= 6 * 3 * 2 //subtract Version information area
		}
	}
	return result
}

// get the number of data bytes can be stored in given qr Version excluding the error correction code
func getDataBytesStored(ver int, ecl Ecc) int {
	getEccCodeWordsPerBlock := GetEccCodeWordsPerBlock()
	getNumOfErrorCorrectionBlock := GetNumOfErrCorrectionBlocks()
	numRawDataModules := getNumRawDataModules(ver)
	return (numRawDataModules / 8) - int(getEccCodeWordsPerBlock[ecl][ver])*int(getNumOfErrorCorrectionBlock[ecl][ver])
}

// Automatically detect which Version to use based on input text length and error correction level
func detectWhichVersionToUse(textLength int, ecl Ecc) (int, error) {
	for i := 1; i <= 40; i++ {
		if getDataBytesStored(i, ecl) >= textLength+3 {
			return i, nil
		}
	}
	return -1, errors.New("input text is too Long for this EC level")
}

// no of bits required for storing the size of data based on Version and encoding mode
func charCountIndicator(Version int, encodingMode Ecm) int {
	ecmByte := [][]int8{
		//0,1,2,3
		{10, 9, 8, 8},    //for Version from 1 to 9
		{12, 11, 16, 10}, //for Version from 10 to 26
		{14, 13, 16, 12}, //for Version from 27 to 40
	}
	// Ensure encodingMode is within bounds (0 to 3)
	if encodingMode < 0 || encodingMode > 3 {
		panic("encodingMode out of range")
	}
	if Version >= 1 && Version <= 9 {
		return int(ecmByte[0][encodingMode])
	} else if Version >= 10 && Version <= 26 {
		return int(ecmByte[1][encodingMode])
	}
	return int(ecmByte[2][encodingMode])
}
