package main

import (
	"github.com/spf13/cobra"
)

var liveCmd = &cobra.Command{
	Use: "live",
	Run: func(c *cobra.Command, args []string) {
		initProxies()
	},
}

func init() {
	addAndBindPFlagsBoolP(liveCmd, configViper, []boolFlagParamSet{
		boolFlagParamSet{"explicit-only", "", false, "Boot Explicit Proxies only", "LiveServer.StartExplicitOnly"},
		boolFlagParamSet{"iptables", "", false, "Enable automatic iptables configuration", "LiveServer.Iptables.EnableAutoConfig"},
	})

	addAndBindPFlagsIntSliceP(liveCmd, configViper, []intSliceFlagParamSet{
		intSliceFlagParamSet{"tcp-dports", "", []int{22}, "TCP Proxy dports, as `port1,port2,...`", "LiveServer.TCP.DestPorts"},
	})

	addAndBindPFlagsStringP(liveCmd, configViper, []stringFlagParamSet{
		stringFlagParamSet{"all-proxy", "", "", "Upstream proxy address for any protocol", "LiveServer.UpstreamProxyURL"},
		stringFlagParamSet{"explicit-listen", "", ":3132", "Explicit Proxy listen address for HTTP/HTTPS, as `[host]:port` Note: This proxy doesn't use authentication info of the `http_proxy` and `https_proxy` environment variables", "LiveServer.Explicit.ListenAddress"},
		stringFlagParamSet{"explicit-with-auth-listen", "", ":3133", "Explicit Proxy with auth listen address for HTTP/HTTPS, as `[host]:port` Note: This proxy uses authentication info of the `http_proxy` and `https_proxy` environment variables", "LiveServer.ExplicitWithAuth.ListenAddress"},
		stringFlagParamSet{"explicit-proxy", "", "", "Upstream proxy address for any protocol, used by explicit proxies", "LiveServer.Explicit.UpstreamProxyURL"},
		stringFlagParamSet{"http-listen", "", ":3129", "HTTP Proxy listen address, as `[host]:port`", "LiveServer.HTTP.ListenAddress"},
		stringFlagParamSet{"http-proxy", "", "", "Upstream proxy address for HTTP", "LiveServer.HTTP.UpstreamProxyURL"},
		stringFlagParamSet{"https-listen", "", ":3130", "HTTPS Proxy listen address, as `[host]:port`", "LiveServer.HTTPS.ListenAddress"},
		stringFlagParamSet{"https-proxy", "", "", "Upstream proxy address for HTTPS", "LiveServer.HTTPS.UpstreamProxyURL"},
		stringFlagParamSet{"tcp-listen", "", ":3128", "TCP Proxy listen address, as `[host]:port`", "LiveServer.TCP.ListenAddress"},
		stringFlagParamSet{"tcp-proxy", "", "", "Upstream proxy address for TCP", "LiveServer.TCP.UpstreamProxyURL"},
	})

	addAndBindPFlagsStringSliceP(liveCmd, configViper, []stringSliceFlagParamSet{
		stringSliceFlagParamSet{"no-proxy", "", []string{}, "Comma-separated list of domain extensions, CIDR, and/or IP addresses which proxy should not be used for", "LiveServer.NoProxy"},
	})

	// configViper.bindEnv(envSetsToBind[i]...) are called in cobra.OnInitialize after it is decided if this program should use env
	// Use http_proxy as default upstream proxy for all connection method
	envSetsToBind = append(envSetsToBind, []string{"LiveServer.UpstreamProxyURL", "http_proxy"})
	envSetsToBind = append(envSetsToBind, []string{"LiveServer.NoProxy", "no_proxy"})

	rootCmd.AddCommand(liveCmd)
}
