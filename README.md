# captcha

Captcha plugin for OpenPanel - discussion: https://github.com/stefanpejcic/OpenPanel/discussions/856

Ships as a standalone Go binary (`captcha`) that OpenPanel execs directly - no Python or other runtime required on the server. The binary is built by CI and committed to this repo, so installing is just cloning it.

Currently supported:
- [Google reCAPTCHA](#google-recaptcha)
- [Cloudflare Turnstile](#cloudflare-turnstile)
- [Custom](#custom)

---

## Google reCAPTCHA

### 1. Install plugin
```bash
cd /etc/openpanel/modules/ && git clone https://github.com/stefanpejcic/captcha
```

### 2. Enable Google reCAPTCHA
```bash
opencli config update captcha_provider google
```

### 3. Add SITE and SECRET keys
```bash
opencli config update recaptcha_site_key SITE_KEY_HERE
```
```bash
opencli config update recaptcha_secret_key SECRET_KEY_HERE
```
### 4. Restart OpenPanel
```bash
docker restart openpanel
```

### 5. Test
Open OpenPanel login page and test.


---

## Cloudflare Turnstile

### 1. Install plugin
```bash
cd /etc/openpanel/modules/ && git clone https://github.com/stefanpejcic/captcha
```

### 2. Enable Cloudflare Turnstile
```bash
opencli config update captcha_provider turnstile
```

### 3. Add SITE and SECRET keys
```bash
opencli config update turnstile_site_key SITE_KEY_HERE
```
```bash
opencli config update turnstile_secret_key SECRET_KEY_HERE
```

### 4. Restart OpenPanel
```bash
docker restart openpanel
```

### 5. Test
Open OpenPanel login page and test.


---

## Custom

'Custom' option should be used by developers as a starting point to integrate a custom CAPTCHA into OpenPanel login page.

Fork the [repo](https://github.com/stefanpejcic/captcha/), edit the `verifyCustom` function in `main.go`, and use your custom repo url in the next step. Push to `main` and the build workflow rebuilds the `captcha` binary for you.

### 1. Install plugin
```bash
cd /etc/openpanel/modules/ && git clone https://YOUR_CUSTOM_FORK
```

### 2. Enable Custom CAPTCHA
```bash
opencli config update captcha_provider custom
```

### 3. Add CUSTOM key
```bash
opencli config update custom_captcha_site_key KEY_HERE
```

### 4. Restart OpenPanel
```bash
docker restart openpanel
```

### 5. Test
Open OpenPanel login page and test.

---

## Plugin contract

OpenPanel execs the plugin binary directly, once per call, and reads its stdout - the same pattern used for every other OpenPanel plugin. There are two subcommands:

**`captcha widget`** - called when rendering the login page. Prints one line of JSON describing what to show:
```json
{"provider": "google", "field_name": "g-recaptcha-response", "site_key": "6Lc..."}
```
`provider` is `""` when no provider is configured - OpenPanel treats that as "don't show a captcha".

**`captcha verify --token=<response_token> [--ip=<client_ip>]`** - called on login submission. Prints `{"success": true|false}` and exits `0` on success, non-zero on failure.
