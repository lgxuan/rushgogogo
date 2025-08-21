package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"regexp"
	"rushgogogo/pkgs/filter"
	"rushgogogo/pkgs/utils"

	"github.com/elazarl/goproxy"
)

func loadData() {
	config, err := utils.GetConfigurationFromYaml("config.yaml")
	if err != nil {
		fmt.Println("config file read error: ", err)
		// 如果配置文件读取失败，使用默认过滤器
		fmt.Println("using default filters")
		return
	}

	if config.Filters != nil && len(config.Filters) > 0 {
		filters = []filter.InformationFilter{}
		for _, fil := range config.Filters {
			// 安全地编译正则表达式
			regex, err := regexp.Compile(fil.FilterRegex)
			if err != nil {
				fmt.Printf("invalid regex for filter %s: %v\n", fil.FilterName, err)
				continue
			}

			filters = append(filters, filter.NewInformationFilter(
				fil.FilterName,
				fil.FilterType,
				regex,
				fil.FilterResource,
				fil.FilterLevel,
				fil.FilterEnabled,
			))
		}
		fmt.Printf("loaded %d filters from config\n", len(filters))
	} else {
		fmt.Println("no filters found in config, using default filters")
	}

	// 安全地更新证书
	cert := config.Cert.Cert
	certKey := config.Cert.Key
	if cert != "" || certKey != "" {
		_caCert = []byte(cert)
		_caKey = []byte(certKey)
	}
}

var (
	filters = []filter.InformationFilter{
		filter.NewInformationFilter(
			"cloud-key",
			"reg",
			regexp.MustCompile(`(((access)([\-_])(key)([\-_])(id|secret))|(LTAI[a-z0-9]{12,20}))`),
			"body",
			"low",
			true,
		),
		filter.NewInformationFilter(
			"windows file/dir Path",
			"reg",
			regexp.MustCompile(`[^\w]([a-zA-Z]:\\\\?(?:[^<>:/\\|?*]+\\\\?)*)([^<>:/\\|?*]+(?:\.[^<>:/\\|?*]+)?)`),
			"body",
			"low",
			true,
		),
		filter.NewInformationFilter(
			"Password Field",
			"reg",
			regexp.MustCompile(`(((|\\)(|'|")(|[\.\w]{1,10})([p](ass|wd|asswd|assword))(|[\.\w]{1,10})(|\\)(|'|")( |)(:|[=]{1,3}|![=]{1,2}|[\)]{0,1}\.val\()( |)(|\\)('|")([^'"]+?)(|\\)('|")(|,|\)))|((|\\)('|")([^'"]+?)(|\\)('|")(|\\)(|'|")( |)(:|[=]{1,3}|![=]{1,2})( |)(|[\.\w]{1,10})([p](ass|wd|asswd|assword))(|[\.\w]{1,10})(|\\)(|'|")))`),
			"body",
			"low",
			true,
		),
		filter.NewInformationFilter(
			"username field",
			"reg",
			regexp.MustCompile(`(((|\\)(|'|")(|[\.\w]{1,10})(([u](ser|name|sername))|(account)|((((create|update)((d|r)|(by|on|at)))|(creator))))(|[\.\w]{1,10})(|\\)(|'|")( |)(:|=|!=|[\)]{0,1}\.val\()( |)(|\\)('|")([^'"]+?)(|\\)('|")(|,|\)))|((|\\)('|")([^'"]+?)(|\\)('|")(|\\)(|'|")( |)(:|[=]{1,3}|![=]{1,2})( |)(|[\.\w]{1,10})(([u](ser|name|sername))|(account)|((((create|update)((d|r)|(by|on|at)))|(creator))))(|[\.\w]{1,10})(|\\)(|'|")))`),
			"body",
			"low",
			true,
		),
		filter.NewInformationFilter(
			"email field",
			"reg",
			regexp.MustCompile(`(([a-z0-9]+[_|\.])*[a-z0-9]+@([a-z0-9]+[-|_|\.])*[a-z0-9]+\.([a-z]{2,5}(?:\.[a-z]{2,5})*))`),
			"body",
			"low",
			true,
		),

		filter.NewInformationFilter(
			"chinese phone number",
			"reg",
			regexp.MustCompile(`[^\w]((?:(?:\+|0{0,2})86)?1(?:(?:3[\d])|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\d])|(?:9[189]))\d{8})[^\w]`),
			"body",
			"low",
			true,
		),

		filter.NewInformationFilter(
			"sensitive data",
			"reg",
			regexp.MustCompile(`(((\[)?('|")?([\.\w]{0,10})(key|secret|token|config|auth|access|admin|ticket)([\.\w]{0,10})('|")?(\])?( |)(:|=|!=|[\)]{0,1}\.val\()( |)('|")([^'"]+?)(\('|")(|,|\)))|((|\\)('|")([^'"]+?)(|\\)('|")(|\\)(|'|")( |)(:|[=]{1,3}|![=]{1,2})( |)(|[\.\w]{1,10})(key|secret|token|config|auth|access|admin|ticket)(|[\.\w]{1,10})(|\\)(|'|")))`),
			"body",
			"low",
			true,
		),
	}
)

func ListenAddress(address string) {
	loadData()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	cert, _ := parseCA(_caCert, _caKey)

	customCaMitm := &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(cert)}
	var customAlwaysMitm goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return customCaMitm, host
	}

	proxy.OnRequest().HandleConnect(customAlwaysMitm)
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		// 安全检查：确保响应不为空
		if r == nil {
			return nil
		}

		newResponse := filter.FilterResponse(r, filters)
		return newResponse
	})
	http.ListenAndServe(address, proxy)

	// log.Fatal()
}

func parseCA(caCert, caKey []byte) (*tls.Certificate, error) {
	parsedCert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, err
	}
	if parsedCert.Leaf, err = x509.ParseCertificate(parsedCert.Certificate[0]); err != nil {
		return nil, err
	}
	return &parsedCert, nil
}
