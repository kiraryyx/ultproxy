package main

// Config is struct to store entire configuration
type Config struct {
	LiveServer struct {
		NoProxy           []string
		StartExplicitOnly bool
		UpstreamProxyURL  string

		Iptables struct {
			EnableAutoConfig bool
		}

		Explicit struct {
			ListenAddress    string
			UpstreamProxyURL string
		}

		ExplicitWithAuth struct {
			ListenAddress    string
			UpstreamProxyURL string
		}

		HTTP struct {
			ListenAddress    string
			UpstreamProxyURL string
		}

		HTTPS struct {
			ListenAddress    string
			UpstreamProxyURL string
		}

		TCP struct {
			DestPorts        []int
			ListenAddress    string
			UpstreamProxyURL string
		}
	}

	Logging struct {
		Level string
	}
}
