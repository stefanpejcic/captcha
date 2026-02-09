# captcha

Captcha plugin for OpenPanel - discussion: https://github.com/stefanpejcic/OpenPanel/discussions/856

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

Fork the [repo](https://github.com/stefanpejcic/captcha/), edit `_verify_custom` function and use your custom repo url in the next step.

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




