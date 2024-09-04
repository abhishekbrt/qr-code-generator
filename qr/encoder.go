package qr

import (
	"fmt"
	"strings"
)

// for demo purpose we are using error correction mode as Q and version 7 which is 45*45 matrix
func Encoder(inputText, errorCorrectionMode string) string {
	//for byte mode, mode indicator is 0100
	modeIndicator := "0100"
	inputTextLength := len(inputText)
	// for version 1-9, 8 bits are required to store the input character length in byte mode
	binaryInputTextLength := fmt.Sprintf("%08b", inputTextLength)
	var encodedText strings.Builder
	encodedText.WriteString(modeIndicator)
	encodedText.WriteString(binaryInputTextLength)
	for i := 0; i < inputTextLength; i++ {
		encodedText.WriteString(fmt.Sprintf("%08b", inputText[i]))
	}
	if (encodedText.Len()) > 704 {
		return "input text is too long!"
	}
	//for 7-Q version total of 88 bytes means 88*8=704 bits are required
	//adding 4 bit terminator (0000) to the end

	paddingLength := 704 - (encodedText.Len())
	if paddingLength >= 4 {
		encodedText.WriteString("0000")
	} else {
		encodedText.WriteString(strings.Repeat("0", paddingLength))
	}
	//adding 0 as pad bytes for making encodedText multiple of 8
	encodedText.WriteString(strings.Repeat("0", ((encodedText.Len()) % 8)))

	//adding 11101100 00010001 untill encodedText reach upto 704 bits
	for encodedText.Len() != 704 {
		encodedText.WriteString("11101100")
		if encodedText.Len() == 704 {
			break
		}
		encodedText.WriteString("00010001")
	}

	fmt.Println(inputTextLength)
	fmt.Printf("%d\n", (encodedText.Len()))
	fmt.Println(inputText)
	fmt.Println(errorCorrectionMode)
	return encodedText.String()

}
