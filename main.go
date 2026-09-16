// Command captcha is an OpenPanel plugin: it reads the panel's captcha
// settings from openpanel.config and answers two calls the panel makes over
// a plain subprocess exec (see https://openpanel.com/docs, plugin contract) -
// "widget" for what to render on the login page, and "verify" for whether a
// submitted response token actually passed.
//
// This replaces the old captcha.py Flask module: OpenPanel's containers no
// longer ship Python, so a plugin now has to be a standalone binary the panel
// can exec directly, not code imported into a shared process.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const configFilePath = "/etc/openpanel/openpanel/conf/openpanel.config"

type config map[string]string

// loadConfig reads openpanel.config's flat key=value lines, ignoring section
// headers ("[SECTION]") and comments - key names are unique across the file
// regardless of section, so tracking sections isn't needed here.
func loadConfig(path string) config {
	cfg := config{}
	f, err := os.Open(path)
	if err != nil {
		return cfg
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		cfg[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(val), `"`)
	}
	return cfg
}

func (c config) get(key, def string) string {
	if v, ok := c[key]; ok && v != "" {
		return v
	}
	return def
}

type providerInfo struct {
	FieldName string
	SiteKey   string
}

// resolveProvider returns "" when no (or an unrecognized) provider is configured, which callers treat as "captcha disabled".
func resolveProvider(cfg config) (provider string, info providerInfo) {
	switch cfg.get("captcha_provider", "") {
	case "google":
		return "google", providerInfo{FieldName: "g-recaptcha-response", SiteKey: cfg.get("recaptcha_site_key", "")}
	case "turnstile":
		return "turnstile", providerInfo{FieldName: "cf-turnstile-response", SiteKey: cfg.get("turnstile_site_key", "")}
	case "custom":
		return "custom", providerInfo{FieldName: "custom-captcha", SiteKey: cfg.get("custom_captcha_site_key", "")}
	default:
		return "", providerInfo{}
	}
}

// cmdWidget prints what the login page needs to render the active provider's widget - an empty "provider" means no captcha is configured.
func cmdWidget() {
	provider, info := resolveProvider(loadConfig(configFilePath))
	_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
		"provider":   provider,
		"field_name": info.FieldName,
		"site_key":   info.SiteKey,
	})
}

// cmdVerify checks a submitted response token against the configured provider, printing {"success": bool} and exiting non-zero on failure.
func cmdVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	token := fs.String("token", "", "captcha response token from the login form")
	ip := fs.String("ip", "", "client IP address, forwarded to providers that accept one")
	_ = fs.Parse(args)

	cfg := loadConfig(configFilePath)
	provider, _ := resolveProvider(cfg)

	var success bool
	switch {
	case provider == "":
		success = true // nothing configured, don't block login
	case *token == "":
		success = false
	case provider == "google":
		success = verifySiteverify(cfg, "https://www.google.com/recaptcha/api/siteverify", "recaptcha_secret_key", *token, *ip)
	case provider == "turnstile":
		success = verifySiteverify(cfg, "https://challenges.cloudflare.com/turnstile/v0/siteverify", "turnstile_secret_key", *token, "")
	case provider == "custom":
		success = verifyCustom(cfg, *token)
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]bool{"success": success})
	if !success {
		os.Exit(1)
	}
}

// verifySiteverify posts to a Google/Cloudflare-style siteverify endpoint (shared response shape between the two providers).
func verifySiteverify(cfg config, endpoint, secretKey, token, ip string) bool {
	form := url.Values{"secret": {cfg.get(secretKey, "")}, "response": {token}}
	if ip != "" {
		form.Set("remoteip", ip)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm(endpoint, form)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}
	return result.Success
}

// verifyCustom is a placeholder for forks: replace this with your own verification logic (see README's "Custom" section).
func verifyCustom(cfg config, token string) bool {
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: captcha <widget|verify> [flags]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "widget":
		cmdWidget()
	case "verify":
		cmdVerify(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
