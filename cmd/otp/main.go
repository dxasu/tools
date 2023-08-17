package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"bytes"
	"io/ioutil"

	"bay.core/lancet/errdo"
	_ "github.com/dxasu/tools/lancet/version"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// otpauth://totp/luck/luck@sz.com?secret=qlt6vmy6svfx4bt4rpmisaiyol6hihca&issuer=luck
func main() {
	if len(os.Args) < 2 {
		println(`input:
otp "otpauth://totp/xxxxxx"
otp [-p] luck luck@sz.com`)
		return
	}
	if len(os.Args) >= 3 {
		var needPng bool
		issuer, account := os.Args[1], os.Args[2]
		if len(os.Args) >= 4 {
			needPng = os.Args[1] == "-p"
			issuer, account = os.Args[1], os.Args[2]
		}
		generate(issuer, account, needPng)
		return
	}

	otpUrl := os.Args[1]
	w, _ := otp.NewKeyFromURL(otpUrl)
	sec := w.Secret()
	code, _ := totp.GenerateCode(sec, time.Now().UTC())
	// valid := totp.Validate(code, sec)
	fmt.Println(code)
}

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to otp.png....")
	ioutil.WriteFile("otp.png", data, 0644)
}

// func promptForPasscode() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Passcode: ")
// 	text, err := reader.ReadString('\n')
// 	errdo.ExitIf(err)
// 	return text
// }

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom parameters
// func GeneratePassCode(utf8string string) string {
// 	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
// 	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
// 		Period:    30,
// 		Skew:      1,
// 		Digits:    otp.DigitsSix,
// 		Algorithm: otp.AlgorithmSHA512,
// 	})
// 	errdo.ExitIf(err)
// 	return passcode
// }

func generate(issuer, account string, needPng bool) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	errdo.ExitIf(err)
	if needPng {
		var buf bytes.Buffer
		img, err := key.Image(200, 200)
		errdo.ExitIf(err)
		png.Encode(&buf, img)
		display(key, buf.Bytes())
	}
	println(key.String())
}
