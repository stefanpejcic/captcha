# captcha.py

# import flask app
from flask import Flask, render_template, render_template_string, request

# import what is needed for this plugin
import os
import requests

# For translations
from flask_babel import Babel, _ # https://python-babel.github.io/flask-babel/

# Import config
from modules.core.config import read_config_file

# =========================
# Configuration
# =========================

# read openpanel config file 
config_file_path = '/etc/openpanel/openpanel/conf/openpanel.config'
config = read_config_file(config_file_path)

# 'google', 'turnstile', 'custom'
CAPTCHA_PROVIDER = config.get('captcha_provider', 'google')

# Google reCAPTCHA
GOOGLE_SITE_KEY = config.get('recaptcha_site_key', 'your_google_site_key')
GOOGLE_SECRET_KEY = config.get('recaptcha_secret_key', 'your_google_secret_key')
GOOGLE_FIELD_NAME = 'g-recaptcha-response'

# Cloudflare Turnstile
TURNSTILE_SITE_KEY = config.get('turnstile_site_key', 'your_turnstile_site_key')
TURNSTILE_SECRET_KEY = config.get('turnstile_secret_key', 'your_turnstile_secret_key')
TURNSTILE_FIELD_NAME = 'cf-turnstile-response'

# Custom CAPTCHA (just a placeholder for now, example)
CUSTOM_SITE_KEY = config.get('custom_captcha_site_key', 'custom_placeholder')
CUSTOM_FIELD_NAME = 'custom-captcha'


FIELD_NAME = {
    'google': GOOGLE_FIELD_NAME,
    'turnstile': TURNSTILE_FIELD_NAME,
    'custom': CUSTOM_FIELD_NAME
}.get(CAPTCHA_PROVIDER, GOOGLE_FIELD_NAME)

SITE_KEY = {
    'google': GOOGLE_SITE_KEY,
    'turnstile': TURNSTILE_SITE_KEY,
    'custom': CUSTOM_SITE_KEY
}.get(CAPTCHA_PROVIDER, GOOGLE_SITE_KEY)


# =========================
# Verification functions
# =========================

def verify(response_token: str, remote_ip: str = None) -> bool:
    """
    Verify CAPTCHA based on the selected provider.
    """
    if not response_token:
        return False

    if CAPTCHA_PROVIDER == 'google':
        return _verify_google(response_token, remote_ip)
    elif CAPTCHA_PROVIDER == 'turnstile':
        return _verify_turnstile(response_token)
    elif CAPTCHA_PROVIDER == 'custom':
        return _verify_custom(response_token)
    else:
        return _verify_google(response_token, remote_ip)


# -------------------------
# Google reCAPTCHA
# -------------------------
def _verify_google(response_token: str, remote_ip: str = None) -> bool:
    url = 'https://www.google.com/recaptcha/api/siteverify'
    payload = {
        'secret': GOOGLE_SECRET_KEY,
        'response': response_token
    }
    if remote_ip:
        payload['remoteip'] = remote_ip
    try:
        r = requests.post(url, data=payload, timeout=5)
        return r.json().get('success', False)
    except requests.RequestException:
        return False


# -------------------------
# Cloudflare Turnstile
# -------------------------
def _verify_turnstile(response_token: str) -> bool:
    url = 'https://challenges.cloudflare.com/turnstile/v0/siteverify'
    data = {'secret': TURNSTILE_SECRET_KEY, 'response': response_token}
    try:
        r = requests.post(url, data=data, timeout=5)
        return r.json().get('success', False)
    except requests.RequestException:
        return False


# -------------------------
# Custom CAPTCHA (placeholder)
# -------------------------
def _verify_custom(response_token: str) -> bool:
    """
    Example: check against a value stored in session.
    """
    # For now, always pass
    return True
