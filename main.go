package main
import(
	"fmt"
	"os"
	"strings"
	"github.com/abhishekbrt/qr-code-generator/qr"
)


func main() {
	text:=os.Args[1:]
	errorCorrectionMode:=text[0]
	inputText:=strings.Join(text[1:]," ")
	// length:=len(inputText)

	str:=qr.Encoder(inputText,errorCorrectionMode)
	fmt.Println(str)
	// fmt.Println(length)
	// fmt.Println(inputText)
	// fmt.Println(errorCorrectionMode)
	
}