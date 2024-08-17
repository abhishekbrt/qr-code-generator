package qr

import (
	"fmt"
	"strings"
)

// for demo purpose we are using error correction mode as Q and version 7 which 45*45 matrix
func Encoder(inputText, errorCorrectionMode string) string {
	//for byte mode, mode indicator is 0100
	modeIndicator := "0100"
	inputTextLength := len(inputText)
	// for version 1-9, 8 bits are required to store the input character length in byte mode
	binaryInputTextLength := fmt.Sprintf("%08b", inputTextLength)

	encodedText := modeIndicator + binaryInputTextLength
	for i := 0; i < inputTextLength; i++ {
		encodedText += fmt.Sprintf("%08b", inputText[i])
	}
	if len(encodedText) > 704 {
		return "input text is too long!"
	}
	//for 7-Q version total of 88 bytes means 88*8=704 bits are required
	//adding 4 bit terminator (0000) to the end

	paddingLength := 704 - len(encodedText)
	if paddingLength >= 4 {
		encodedText += "0000"
	} else {
		encodedText += strings.Repeat("0", paddingLength)
	}
	//adding 0 as pad bytes for making encodedText multiple of 8
	encodedText += strings.Repeat("0", (len(encodedText) % 8))

	//adding 11101100 00010001 untill encodedText reach upto 704 bits
	paddingLen := (704 - len(encodedText)) / 16
	encodedText += strings.Repeat("1110110000010001", paddingLen)

	fmt.Println(inputTextLength)
	fmt.Printf("%d\n", len(encodedText))
	fmt.Println(inputText)
	fmt.Println(errorCorrectionMode)
	return encodedText

}
