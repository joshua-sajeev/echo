package handlers

import (
	"net/http"

	"github.com/joshu-sajeev/echo/internal/auth"
)

// LoginHandler serves the PIN login page and processes submissions.
type LoginHandler struct{}

func NewLoginHandler() *LoginHandler { return &LoginHandler{} }

func (h *LoginHandler) Page(w http.ResponseWriter, r *http.Request) {
	// already logged in → home
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(loginPage()))
}

func (h *LoginHandler) Submit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	pin := r.FormValue("pin")
	if !auth.VerifyPIN(pin) {
		// Return error fragment for HTMX swap
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<p id="login-error" class="text-red-400 text-sm text-center mt-2 animate-pulse">Incorrect PIN</p>`))
		return
	}

	auth.SetSession(w)
	// Tell HTMX to do a full-page redirect (not a swap)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSession(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Echo Finance</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body class="bg-zinc-950 text-white min-h-screen flex items-center justify-center px-6">

<div class="w-full max-w-xs space-y-8">

  <!-- Logo / title -->
  <div class="text-center space-y-1">
    <div class="w-14 h-14 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center mx-auto mb-4">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
      </svg>
    </div>
    <h1 class="text-xl font-semibold">Echo Finance</h1>
    <p class="text-zinc-500 text-sm">Enter your PIN to continue</p>
  </div>

  <!-- PIN dots display -->
  <div class="flex justify-center gap-4" id="pin-dots">
    <div class="dot w-4 h-4 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
    <div class="dot w-4 h-4 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
    <div class="dot w-4 h-4 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
    <div class="dot w-4 h-4 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
    <div class="dot w-4 h-4 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
    <div class="dot w-4 h-4 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
  </div>

  <!-- Hidden input + error slot -->
  <input type="hidden" id="pin-value"/>
  <div id="login-error" class="min-h-[20px]"></div>

  <!-- Numpad -->
  <div class="grid grid-cols-3 gap-3">
    <button class="numpad-btn" onclick="pressDigit('1')">1</button>
    <button class="numpad-btn" onclick="pressDigit('2')">2</button>
    <button class="numpad-btn" onclick="pressDigit('3')">3</button>
    <button class="numpad-btn" onclick="pressDigit('4')">4</button>
    <button class="numpad-btn" onclick="pressDigit('5')">5</button>
    <button class="numpad-btn" onclick="pressDigit('6')">6</button>
    <button class="numpad-btn" onclick="pressDigit('7')">7</button>
    <button class="numpad-btn" onclick="pressDigit('8')">8</button>
    <button class="numpad-btn" onclick="pressDigit('9')">9</button>
    <!-- Fingerprint / biometric button -->
    <button id="bio-btn" onclick="tryBiometric()"
      class="numpad-btn flex items-center justify-center text-zinc-400 hover:text-emerald-400"
      title="Use fingerprint">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M12 11c0-1.657-1.343-3-3-3S6 9.343 6 11v1a6 6 0 0012 0v-1c0-1.657-1.343-3-3-3s-3 1.343-3 3"/>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M12 21a9 9 0 110-18 9 9 0 010 18z"/>
      </svg>
    </button>
    <button class="numpad-btn text-2xl font-light" onclick="pressDigit('0')">0</button>
    <button class="numpad-btn text-zinc-400" onclick="backspace()">⌫</button>
  </div>

</div>

<style>
  .numpad-btn {
    @apply w-full aspect-square rounded-2xl bg-zinc-900 border border-zinc-800
           text-xl font-medium text-white
           hover:bg-zinc-800 active:scale-95
           transition-all duration-100 select-none;
  }
</style>

<script>
const PIN_LEN = 6;
let _pin = '';

function updateDots() {
  const dots = document.querySelectorAll('.dot');
  dots.forEach((d, i) => {
    if (i < _pin.length) {
      d.classList.add('bg-emerald-400', 'border-emerald-400');
    } else {
      d.classList.remove('bg-emerald-400', 'border-emerald-400');
    }
  });
}

function shake() {
  const dots = document.getElementById('pin-dots');
  dots.classList.add('translate-x-2');
  setTimeout(() => dots.classList.add('-translate-x-2'), 80);
  setTimeout(() => dots.classList.remove('translate-x-2', '-translate-x-2'), 160);
}

function clearError() {
  document.getElementById('login-error').innerHTML = '';
}

function pressDigit(d) {
  clearError();
  if (_pin.length >= PIN_LEN) return;
  _pin += d;
  updateDots();
  if (_pin.length === PIN_LEN) submitPIN();
}

function backspace() {
  clearError();
  _pin = _pin.slice(0, -1);
  updateDots();
}

function submitPIN() {
  htmx.ajax('POST', '/login', {
    target: '#login-error',
    swap: 'innerHTML',
    values: { pin: _pin },
  }).then(() => {
    // on failure the dots shake and reset
    const err = document.getElementById('login-error').textContent.trim();
    if (err) {
      shake();
      _pin = '';
      updateDots();
    }
  });
}

// ── Biometric (WebAuthn / platform authenticator) ──────────────
// On mobile Safari/Chrome this triggers Face ID or fingerprint.
// We use a credential stored in localStorage as the "user handle".
// The actual secret is still the PIN hash on the server;
// biometric just auto-submits the stored PIN via a passkey assertion.
//
// Simpler approach: store an encrypted PIN blob in localStorage,
// decrypt with biometric → submit. We use the even simpler approach:
// store the PIN itself locally, submit on biometric success.
// This is fine for a single-user personal app.

async function tryBiometric() {
  if (!window.PublicKeyCredential) {
    document.getElementById('login-error').innerHTML =
      '<p class="text-zinc-500 text-xs text-center">Biometric not supported on this browser</p>';
    return;
  }

  // Check if platform authenticator is available
  const available = await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
  if (!available) {
    document.getElementById('login-error').innerHTML =
      '<p class="text-zinc-500 text-xs text-center">No fingerprint/Face ID found</p>';
    return;
  }

  // If we have a stored credential, use it
  const storedCredId = localStorage.getItem('echo_cred_id');
  const storedPin    = localStorage.getItem('echo_bio_pin');

  if (!storedCredId || !storedPin) {
    document.getElementById('login-error').innerHTML =
      '<p class="text-zinc-500 text-xs text-center mt-1">No biometric set up.<br/>Enter PIN once then press 🖐 to register.</p>';
    return;
  }

  try {
    // Trigger biometric prompt via a get() assertion
    const credId = Uint8Array.from(atob(storedCredId), c => c.charCodeAt(0));
    await navigator.credentials.get({
      publicKey: {
        challenge:        crypto.getRandomValues(new Uint8Array(32)),
        rpId:             location.hostname,
        userVerification: 'required',
        allowCredentials: [{ type: 'public-key', id: credId }],
        timeout:          60000,
      }
    });
    // Biometric passed → submit stored PIN
    _pin = storedPin;
    updateDots();
    submitPIN();
  } catch (e) {
    document.getElementById('login-error').innerHTML =
      '<p class="text-zinc-500 text-xs text-center">Biometric cancelled</p>';
  }
}

// Register biometric after a successful PIN login
// Call this from the console or add a "Set up fingerprint" button:
// registerBiometric('your-6-digit-pin')
async function registerBiometric(pin) {
  if (!window.PublicKeyCredential) return alert('Not supported');

  const uid = crypto.getRandomValues(new Uint8Array(16));
  const cred = await navigator.credentials.create({
    publicKey: {
      challenge:              crypto.getRandomValues(new Uint8Array(32)),
      rp:                     { name: 'Echo Finance', id: location.hostname },
      user:                   { id: uid, name: 'owner', displayName: 'Owner' },
      pubKeyCredParams:       [{ type: 'public-key', alg: -7 }],
      authenticatorSelection: { authenticatorAttachment: 'platform', userVerification: 'required' },
      timeout:                60000,
    }
  });

  // Store credential id (base64) and pin for future assertions
  localStorage.setItem('echo_cred_id', btoa(String.fromCharCode(...new Uint8Array(cred.rawId))));
  localStorage.setItem('echo_bio_pin', pin);
  alert('Fingerprint registered! You can now use the fingerprint button.');
}

// Keyboard support
document.addEventListener('keydown', e => {
  if (e.key >= '0' && e.key <= '9') pressDigit(e.key);
  if (e.key === 'Backspace') backspace();
});

// Hide bio button if not available
PublicKeyCredential?.isUserVerifyingPlatformAuthenticatorAvailable?.()
  .then(ok => { if (!ok) document.getElementById('bio-btn').classList.add('invisible'); });
</script>

</body>
</html>`
}
