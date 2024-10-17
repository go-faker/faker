package faker

import (
	"fmt"
	"reflect"
	"time"

	"github.com/cj/faker/pkg/interfaces"
)

var (
	browserTypes          = []string{"chrome", "firefox", "internetExplorer", "msedge", "opera", "safari"}
	windowsPlatformTokens = []string{"Windows NT 6.2", "Windows NT 6.1", "Windows NT 6.0", "Windows NT 5.2",
		"Windows NT 5.1", "Windows NT 5.01", "Windows NT 5.0", "Windows NT 4.0", "Windows 98; Win 9x 4.90",
		"Windows 98", "Windows 95", "Windows CE",
	}
	linuxProcessor = []string{"i686", "x86_64"}
	macProcessor   = []string{"Intel", "PPC", "U; Intel", "U; PPC"}
	// langCodes contains language codes extracted from https://gist.github.com/JamieMason/3748498.
	langCodes = []string{"af-ZA", "ar-AE", "ar-BH", "ar-DZ", "ar-EG", "ar-IQ", "ar-JO", "ar-KW", "ar-LB", "ar-LY",
		"ar-MA", "ar-OM", "ar-QA", "ar-SA", "ar-SY", "ar-TN", "ar-YE", "be-BY", "bg-BG", "ca-ES", "cs-CZ", "cy-GB",
		"da-DK", "de-AT", "de-DE", "de-CH", "de-LI", "de-LU", "dv-MV", "el-GR", "en-AU", "en-BZ", "en-CA", "en-GB",
		"en-IE", "en-JM", "en-NZ", "en-PH", "en-TT", "en-US", "en-ZA", "en-ZW", "es-AR", "es-BO", "es-CL", "es-CO",
		"es-CR", "es-DO", "es-EC", "es-ES", "es-GT", "es-HN", "es-MX", "es-NI", "es-PA", "es-PE", "es-PR", "es-PY",
		"es-SV", "es-UY", "es-VE", "et-EE", "eu-ES", "fa-IR", "fi-FI", "fo-FO", "fr-BE", "fr-CA", "fr-FR", "fr-CH",
		"fr-LU", "fr-MC", "gl-ES", "gu-IN", "he-IL", "hi-IN", "hr-BA", "hr-HR", "hu-HU", "hy-AM", "id-ID", "is-IS",
		"it-CH", "it-IT", "ja-JP", "ka-GE", "kk-KZ", "kn-IN", "ko-KR", "ky-KG", "lt-LT", "lv-LV", "mi-NZ", "mk-MK",
		"mn-MN", "mr-IN", "ms-BN", "ms-MY", "mt-MT", "nb-NO", "nl-BE", "nl-NL", "nn-NO", "ns-ZA", "pa-IN", "pl-PL",
		"pt-BR", "pt-PT", "ro-RO", "ru-RU", "sa-IN", "se-FI", "se-NO", "se-SE", "sk-SK", "sl-SI", "sq-AL", "sv-FI",
		"sv-SE", "sw-KE", "ta-IN", "te-IN", "th-TH", "tn-ZA", "tr-TR", "tt-RU", "uk-UA", "ur-PK", "vi-VN", "xh-ZA",
		"zh-CN", "zh-HK", "zh-MO", "zh-SG", "zh-TW", "zu-ZA",
	}
)

func GetUserAgent() UserAgenter {
	ua := &UserAgent{}
	return ua
}

type UserAgenter interface {
	UserAgent(v reflect.Value) (interface{}, error)
}

// UserAgent implements logic loosely based on https://fakerphp.org/formatters/user-agent/.
type UserAgent struct{}

func (ua UserAgent) UserAgent(reflect.Value) (interface{}, error) {
	browserType := randomElementFromSliceString(browserTypes)
	switch browserType {
	case "chrome":
		return ua.chrome(), nil
	case "firefox":
		return ua.firefox(), nil
	case "internetExplorer":
		return ua.internetExplorer(), nil
	case "msedge":
		return ua.msEdge(), nil
	case "opera":
		return ua.opera(), nil
	case "safari":
		return ua.safari(), nil
	}
	return nil, fmt.Errorf(`unknown browser type: "%s"`, browserType)
}

// chrome generates a user agent string that resembles Google Chrome.
// Example: "Mozilla/5.0 (X11; Linux i686) AppleWebKit/532.2 (KHTML, like Gecko) Chrome/69.0.3460.334 Safari/532.2"
func (ua UserAgent) chrome() string {
	chromeVer := fmt.Sprintf("%d.0.%d.%d",
		ua.intBetween(36, 127), ua.intBetween(0, 6000), ua.intBetween(0, 400),
	)
	platformToken := randomElementFromSliceString([]string{
		ua.linuxPlatformToken(),
		ua.macPlatformToken("_"),
		ua.windowsPlatformToken(),
	})
	safariType := randomElementFromSliceString([]string{"Safari", "Mobile Safari"})
	webKitVer := ua.webKitVersion(false)
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/%s (KHTML, like Gecko) Chrome/%s %s/%s",
		platformToken, webKitVer, chromeVer, safariType, webKitVer,
	)
}

// firefox generates a user agent string that resembles Firefox.
// Example: "Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20130219 Firefox/58.0"
func (ua UserAgent) firefox() string {
	firefoxVer := fmt.Sprintf("%d.0", ua.intBetween(35, 129))
	platformToken := randomElementFromSliceString([]string{
		ua.linuxPlatformToken(),
		ua.macPlatformToken("."),
		ua.windowsPlatformToken(),
	})
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%s) Gecko/%s Firefox/%s",
		platformToken, firefoxVer, ua.geckoVersion(), firefoxVer,
	)
}

// internetExplorer generates a user agent string that resembles Internet Explorer
// Example: "Mozilla/5.0 (compatible; MSIE 9.0; Windows 98; Win 9x 4.90; Trident/5.0)"
func (ua UserAgent) internetExplorer() string {
	return fmt.Sprintf("Mozilla/5.0 (compatible; MSIE %d.0; %s; Trident/%d.0)",
		ua.intBetween(5, 11), ua.windowsPlatformToken(), ua.intBetween(3, 7),
	)
}

// msEdge generates a user agent string that resembles MS Edge
// Example: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/533.0 (KHTML, like Gecko) Chrome/86.0.4685.32 Safari/533.0 Edg/86.0.4685.32"
func (ua UserAgent) msEdge() string {
	webKitVer := ua.webKitVersion(false)
	chromiumVer := fmt.Sprintf("%d.0.%d.%d", ua.intBetween(79, 99), ua.intBetween(4000, 4844), ua.intBetween(10, 99))
	platformToken := randomElementFromSliceString([]string{
		ua.linuxPlatformToken(),
		ua.windowsPlatformToken(),
		ua.macPlatformToken("_"),
	})
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/%s (KHTML, like Gecko) Chrome/%s Safari/%s Edg/%s",
		platformToken, webKitVer, chromiumVer, webKitVer, chromiumVer,
	)
}

// opera generates a user agent string that resembles Opera.
// Example: "Opera/8.40 (X11; Linux i686; zh-TW) Presto/2.12.181 Version/10.57"
func (ua UserAgent) opera() string {
	platformToken := randomElementFromSliceString([]string{
		ua.linuxPlatformToken(),
		ua.windowsPlatformToken(),
	})
	operaVer := fmt.Sprintf("%d.%d", ua.intBetween(8, 9), ua.intBetween(10, 99))
	prestoVer := fmt.Sprintf("2.%d.%d", ua.intBetween(8, 12), ua.intBetween(160, 355))
	ver := fmt.Sprintf("%d.%d", ua.intBetween(10, 12), ua.intBetween(0, 70))
	return fmt.Sprintf("Opera/%s (%s; %s) Presto/%s Version/%s",
		operaVer, platformToken, randomElementFromSliceString(langCodes), prestoVer, ver,
	)
}

// safari generates a user agent string that resembles Safari.
// Example: "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_7_3; bg-BG) AppleWebKit/548.43.16 (KHTML, like Gecko) Version/7.5 Safari/548.43.16"
func (ua UserAgent) safari() string {
	webKitVer := ua.webKitVersion(true)
	safariVer := fmt.Sprintf("%d.%d", ua.intBetween(4, 17), ua.intBetween(0, 5))
	if ua.intBetween(0, 100) < 50 {
		safariVer = fmt.Sprintf("%s.%d", safariVer, ua.intBetween(1, 20))
	}
	platform := randomElementFromSliceString([]string{
		fmt.Sprintf("Windows; U; %s; %s", ua.windowsPlatformToken(), randomElementFromSliceString(langCodes)),
		fmt.Sprintf("%s; %s", ua.macPlatformToken("_"), randomElementFromSliceString(langCodes)),
	})

	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/%s (KHTML, like Gecko) Version/%s Safari/%s",
		platform, webKitVer, safariVer, webKitVer,
	)
}

func (ua UserAgent) linuxPlatformToken() string {
	return fmt.Sprintf("X11; Linux %s", randomElementFromSliceString(linuxProcessor))
}

func (ua UserAgent) macPlatformToken(sep string) string {
	return fmt.Sprintf("Macintosh; %s Mac OS X %d%s%d%s%d",
		randomElementFromSliceString(macProcessor), ua.intBetween(10, 17), sep, ua.intBetween(5, 8), sep, ua.intBetween(0, 9),
	)
}

func (ua UserAgent) windowsPlatformToken() string {
	return randomElementFromSliceString(windowsPlatformTokens)
}

// webKitVersion generates a random string that resembles an AppleWebKit version string.
func (ua UserAgent) webKitVersion(withPatchVer bool) string {
	majorVer := ua.intBetween(531, 618)
	if withPatchVer {
		return fmt.Sprintf("%d.%d.%d", majorVer, ua.intBetween(1, 50), ua.intBetween(1, 20))
	}
	return fmt.Sprintf("%d.%d", majorVer, ua.intBetween(0, 10))
}

// geckoVersion generates a random string that resembles a Gecko version string.
func (ua UserAgent) geckoVersion() string {
	const YYYYMMDD = "20060102"
	minDate := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Now()
	randDuration := time.Duration(rand.Int63n(maxDate.Unix()-minDate.Unix())) * time.Second
	return minDate.Add(randDuration).Format(YYYYMMDD)
}

// intBetween generates a random integer in the range [min, max].
func (ua UserAgent) intBetween(min, max int) int {
	boundary := interfaces.RandomIntegerBoundary{Start: min, End: max + 1}
	return randomIntegerWithBoundary(boundary)
}
