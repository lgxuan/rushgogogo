package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// installCertCmd represents the install-cert command
var installCertCmd = &cobra.Command{
	Use:   "install-cert",
	Short: "Install CA certificate to system trust store",
	Long: `Install the CA certificate to the system's trust store.
This command supports Windows, macOS, and Linux systems.
On Windows, it uses certutil to add to trusted root authorities.
On macOS, it uses security to add to system keychain.
On Linux, it copies to ca-certificates directory and updates.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := installCertificate(); err != nil {
			fmt.Printf("Error installing certificate: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("CA certificate installed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(installCertCmd)
}

func installCertificate() error {
	// 创建临时证书文件
	tmpFile, err := os.CreateTemp("", "ca-cert-*.crt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入证书内容
	if err := os.WriteFile(tmpFile.Name(), getCACert(), 0644); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// 根据操作系统安装证书
	switch runtime.GOOS {
	case "windows":
		return installOnWindows(tmpFile.Name())
	case "darwin":
		return installOnMacOS(tmpFile.Name())
	default:
		return installOnLinux(tmpFile.Name())
	}
}

func installOnWindows(certPath string) error {
	cmd := exec.Command("certutil", "-addstore", "-f", "Root", certPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installOnMacOS(certPath string) error {
	cmd := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain", certPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installOnLinux(certPath string) error {
	// 复制到ca-certificates目录
	caBundle := "/usr/local/share/ca-certificates/goproxy-demo.crt"
	copyCmd := exec.Command("sudo", "cp", certPath, caBundle)
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy certificate: %w", err)
	}

	// 更新ca-certificates
	updateCmd := exec.Command("sudo", "update-ca-certificates")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	return updateCmd.Run()
}

func getCACert() []byte {
	return []byte(`-----BEGIN CERTIFICATE-----
MIID1zCCAr+gAwIBAgIUFoZhiSvErDyRQwUMbsjs3/+GjgwwDQYJKoZIhvcNAQEL
BQAwezELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcM
DVNhbiBGcmFuY2lzY28xGTAXBgNVBAoMEENsb3VkRmxhcmUsIEluYy4xJDAiBgNV
BAMMG2Nkbi5jbG91ZGZsYXJlLW9wdGltaXplLmNvbTAeFw0yNTA3MzExNjIwMzNa
Fw0zMDA3MzAxNjIwMzNaMHsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9y
bmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKDBBDbG91ZEZsYXJl
LCBJbmMuMSQwIgYDVQQDDBtjZG4uY2xvdWRmbGFyZS1vcHRpbWl6ZS5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC1i54o/NSKrNGZypRV494dF+Pz
PvW6yPLUgoZIbTK9Dhl6QWpa23FwzKvtZx6YuvWAtp8MThMwKuiIzyD0gxc9qqxC
5VMEKrYxzv9POkMv2C7ejirFiHZGCvXxMyBK88c2HAen0hs6j0zS1Ch5P+d9ySZz
o6SjKSy2ciUesXwFgdNN5Fu8fWcgT2+vK8Z9lXo+wcMboBjf1uItiADPCFaV5Ql5
f3D7ShcMF/QSQDVIfeEf5aatBwlhd3HGGbLBPxXQPHX6fiCMsrWpa4l4H76OVlND
91QN50W5U1S+YPlp0ALsdTwFdrU8W4qUvRC+uYkE0AnYJeFYJjM+sCHRPhxNAgMB
AAGjUzBRMB0GA1UdDgQWBBQZAkfiyJvcCJw+aFX0urhMBiJSDDAfBgNVHSMEGDAW
gBQZAkfiyJvcCJw+aFX0urhMBiJSDDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQB7q0SwKkvoNPxrjXn1UOCXvkkBcpDmnvezlbWcORApVjmLAQY1
3uU5YHFB50NcXiuQTAGJa70FCirL/8uZlHnD+fhLC1f2iyBy1gogCMh6ASRW8I69
kkhyFtcqHtJFkgfOb4WVgpxy+VWFdFXrRuD77346XqD9XG51yVPnIQpHawgdXmps
wHQ58nNM27X6Fz0klEyn77BrZy1EBGbtWGwGjF25x863giaPK/R7Z2HjkM0Ka+Gj
BQLfsGdB97ecAkLLIVwojLEEHSi6fgtYytTfav6FlyG/BMpRigmnldcL51Hotufn
c5zvXq/DSOMJrSgD6YOx9neCmkWIfiAo0pmy
-----END CERTIFICATE-----`)
}
