package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"
	"time"

	"bytes"
	"io/fs"

	"github.com/atotto/clipboard"
	"github.com/dxasu/pure/rain"
	_ "github.com/dxasu/pure/version"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// otpauth://totp/luck/luck@sz.com?secret=qlt6vmy6svfx4bt4rpmisaiyol6hihca&issuer=luck
func main() {
	if len(os.Args) < 2 || rain.NeedHelp() {
		fmt.Println(`Usage:
otp [filename] otpauth://totp/xxxxxx
otp -[grpcn] 
Flags:
	-g generate link. otp -g luck luck@sz.com
	-G -g and output png
	-r read code from file
	-c copy to clipboard
	-n print nothing but error
`)
		return
	}

	params := make([]string, len(os.Args))
	copy(params, os.Args)
	needCopy := false
	notPrint := false
RESTART:
	var cmd string
	var printText strings.Builder
	if strings.HasPrefix(params[1], "otpauth") || len(params) >= 3 && strings.HasPrefix(params[2], "otpauth") {
		startIdx := 1
		otpName := ""
		if len(params) >= 3 && strings.HasPrefix(params[2], "otpauth") {
			otpName = params[1]
			startIdx = 2
		}
		otpUrl := strings.Join(params[startIdx:], " ")
		w, err := otp.NewKeyFromURL(otpUrl)
		rain.ExitIf(err)
		key := w.Secret()
		code, err := totp.GenerateCode(key, time.Now().UTC())
		rain.ExitIf(err)
		printText.WriteString(code)
		// valid := totp.Validate(code, key)
		if otpName != "" {
			err = os.WriteFile(otpName, []byte(otpUrl), fs.ModePerm)
			rain.ExitIf(err)
		}
	} else {
		cmd = params[1]
		if cmd[0] != '-' {
			rain.ExitIf(fmt.Errorf("invalid params: %s", cmd))
		}
		needCopy, notPrint = strings.ContainsRune(cmd, 'c'), strings.ContainsRune(cmd, 'n')
		if len(params) >= 4 && strings.ContainsAny(cmd, "Gg") {
			needPng := strings.ContainsRune(cmd, 'G')
			issuer, account := params[2], params[3]
			key := generate(issuer, account, needPng)
			printText.WriteString(key)
		} else if strings.ContainsRune(cmd, 'r') && len(params) >= 3 {
			fileName := params[2]
			bys, err := os.ReadFile(fileName)
			rain.ExitIf(err)
			params = []string{os.Args[0], string(bys)}
			goto RESTART
		} else {
			rain.ExitIf(fmt.Errorf("unimplement or invalid params: %s", cmd))
		}
	}

	if printText.Len() != 0 {
		result := printText.String()
		if needCopy {
			clipboard.WriteAll(result)
		}

		if notPrint {
			return
		}
		fmt.Println(result)
	}
}

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to otp.png....")
	os.WriteFile("otp.png", data, 0644)
}

// func promptForPasscode() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter Passcode: ")
// 	text, err := reader.ReadString('\n')
// 	rain.ExitIf(err)
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
// 	rain.ExitIf(err)
// 	return passcode
// }

func generate(issuer, account string, needPng bool) string {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	rain.ExitIf(err)
	if needPng {
		var buf bytes.Buffer
		img, err := key.Image(200, 200)
		rain.ExitIf(err)
		png.Encode(&buf, img)
		display(key, buf.Bytes())
	}
	return key.String()
}
