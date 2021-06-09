package ip

import (
	"context"
	"strings"
)

var IPGetters = []IPGetter{
	{
		ServiceURL: "https://httpbin.org/ip",
		GetIP: func(url string, ctx context.Context) (string, error) {
			s, err := HTTPJSONGetField("origin")(url, ctx)
			if err != nil {
				return "", err
			}
			return strings.Split(s, ", ")[0], nil
		},
	},
	{
		ServiceURL: "https://api.myip.com",
		GetIP:      HTTPJSONGetField("ip"),
	},
	{
		ServiceURL: "https://ip4.seeip.org/json",
		GetIP:      HTTPJSONGetField("ip"),
	},
	{
		ServiceURL: "https://ipv4bot.whatismyipaddress.com/",
		GetIP:      HTTPGetString,
	},
	{
		ServiceURL: "https://api.ipify.org?format=json",
		GetIP:      HTTPJSONGetField("ip"),
	},
}
