package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/BurntSushi/toml"
	"github.com/valyala/fasthttp"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
)

// config is must point
func LoadConfig(fileName string, config interface{}) error {
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if _, err := toml.DecodeFile(fileName, config); err != nil {
		return err
	}
	return err
}

func WaitGroupTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func GetOTP() string {
	return ""
}

func CreateFastClient() *fasthttp.Client {
	return &fasthttp.Client{MaxConnsPerHost: 50000, MaxIdleConnDuration: 1 * time.Second}
}

func CreateDefaultFastClient() *fasthttp.Client {
	return &fasthttp.Client{}
}

func HandleError(err interface{}) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[E] %v %s:%d", err, fn, line)
	}
}

func GetSlug(input string) string {
	if len(input) == 0 {
		return "user-name"
	}
	trans := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	out, _, _ := transform.String(trans, input)
	out = strings.ToLower(out)
	out = strings.Trim(out, ` `)
	var re = regexp.MustCompile(`[ ]{1,}`)
	out = re.ReplaceAllString(out, `-`)
	out = strings.Replace(out, "đ", "d", -1)
	return out
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func GenHMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(message)
	messageMAC := mac.Sum(nil)
	return messageMAC
}
