package handlers

import (
	"net/http"
)

type BiometricHandler struct{}

func NewBiometricHandler() *BiometricHandler { return &BiometricHandler{} }

func (h *BiometricHandler) SetupPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(biometricSetupPage()))
}

func biometricSetupPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no"/>
  <title>Set up Fingerprint – Echo Finance</title>
  <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-zinc-950 text-white min-h-screen flex flex-col">

<nav class="border-b border-zinc-800 bg-zinc-950">
  <div class="max-w-lg mx-auto px-6 py-4 flex items-center gap-4">
    <a href="/" class="text-zinc-500 hover:text-white transition text-sm">← Back</a>
    <h1 class="text-base font-semibold">Fingerprint / Face ID</h1>
  </div>
</nav>

<main class="flex-1 flex flex-col items-center justify-center px-6 max-w-xs mx-auto w-full space-y-8">

  <!-- icon -->
  <div id="icon-wrap"
    class="w-24 h-24 rounded-3xl bg-zinc-900 border border-zinc-800 flex items-center justify-center transition-all duration-300">
    <svg id="bio-icon" xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-zinc-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
        d="M12 11c0-1.657-1.343-3-3-3S6 9.343 6 11v1a6 6 0 0012 0v-1c0-1.657-1.343-3-3-3s-3 1.343-3 3"/>
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
        d="M12 21a9 9 0 110-18 9 9 0 010 18z"/>
    </svg>
  </div>

  <div class="text-center space-y-2">
    <h2 class="text-lg font-semibold" id="setup-title">Set up biometric login</h2>
    <p class="text-zinc-500 text-sm" id="setup-subtitle">
      Use fingerprint or Face ID instead of your PIN next time.
    </p>
  </div>

  <!-- PIN entry (needed once to register) -->
  <div id="pin-section" class="w-full space-y-4">
    <p class="text-xs text-zinc-600 text-center uppercase tracking-wider">Confirm your PIN first</p>

    <!-- dots -->
    <div id="pin-dots" class="flex justify-center gap-3">
      <div class="dot w-3.5 h-3.5 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
      <div class="dot w-3.5 h-3.5 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
      <div class="dot w-3.5 h-3.5 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
      <div class="dot w-3.5 h-3.5 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
      <div class="dot w-3.5 h-3.5 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
      <div class="dot w-3.5 h-3.5 rounded-full border-2 border-zinc-700 transition-all duration-150"></div>
    </div>

    <!-- numpad -->
    <div class="grid grid-cols-3 gap-2.5">
      <button class="numpad-btn" onclick="pressDigit('1')">1</button>
      <button class="numpad-btn" onclick="pressDigit('2')">2</button>
      <button class="numpad-btn" onclick="pressDigit('3')">3</button>
      <button class="numpad-btn" onclick="pressDigit('4')">4</button>
      <button class="numpad-btn" onclick="pressDigit('5')">5</button>
      <button class="numpad-btn" onclick="pressDigit('6')">6</button>
      <button class="numpad-btn" onclick="pressDigit('7')">7</button>
      <button class="numpad-btn" onclick="pressDigit('8')">8</button>
      <button class="numpad-btn" onclick="pressDigit('9')">9</button>
      <div></div>
      <button class="numpad-btn" onclick="pressDigit('0')">0</button>
      <button class="numpad-btn text-zinc-400" onclick="backspace()">⌫</button>
    </div>

    <div id="pin-error" class="min-h-[20px] text-center"></div>
  </div>

  <!-- shown after success -->
  <div id="success-section" class="hidden text-center space-y-4">
    <p class="text-emerald-400 font-medium">Fingerprint registered!</p>
    <p class="text-zinc-500 text-sm">Use the fingerprint button on the login screen next time.</p>
    <a href="/" class="block w-full py-3 rounded-xl bg-white text-black text-sm font-medium text-center hover:bg-zinc-200 transition">
      Go to Dashboard
    </a>
  </div>

  <!-- shown if already registered -->
  <div id="already-section" class="hidden w-full space-y-3">
    <div class="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-4 text-center">
      <p class="text-emerald-400 text-sm font-medium">Biometric already set up ✓</p>
      <p class="text-zinc-500 text-xs mt-1">Active on this device</p>
    </div>
    <button onclick="resetBiometric()"
      class="w-full py-3 rounded-xl border border-red-900 text-red-400 text-sm hover:bg-red-950 transition">
      Remove &amp; re-register
    </button>
    <a href="/"
      class="block w-full py-3 rounded-xl bg-zinc-800 text-zinc-300 text-sm font-medium text-center hover:bg-zinc-700 transition">
      Back to Dashboard
    </a>
  </div>

  <!-- shown if WebAuthn not available -->
  <div id="unsupported-section" class="hidden text-center space-y-3">
    <div class="bg-zinc-900 border border-zinc-800 rounded-xl p-4">
      <p class="text-zinc-400 text-sm">Biometric login isn't available on this browser/device.</p>
      <p class="text-zinc-600 text-xs mt-1">Try Chrome or Safari on a phone with fingerprint or Face ID.</p>
    </div>
    <a href="/" class="block text-zinc-500 text-sm hover:text-zinc-300 transition">← Back</a>
  </div>

</main>

<style>
  .numpad-btn {
    width: 100%;
    aspect-ratio: 1;
    border-radius: 1rem;
    background: rgb(24 24 27);
    border: 1px solid rgb(39 39 42);
    font-size: 1.25rem;
    font-weight: 500;
    color: white;
    transition: all 0.1s;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
  }
  .numpad-btn:active { transform: scale(0.93); background: rgb(39 39 42); }
</style>

<script>
const PIN_LEN = 6;
let _pin = '';

// ── check state on load ───────────────────────────────────────
(async function init() {
  // 1. WebAuthn supported?
  if (!window.PublicKeyCredential) {
    show('unsupported-section'); return;
  }
  const available = await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable().catch(() => false);
  if (!available) {
    show('unsupported-section'); return;
  }

  // 2. Already registered on this device?
  if (localStorage.getItem('echo_cred_id')) {
    show('already-section'); return;
  }

  // 3. Show PIN entry
  show('pin-section');
})();

function show(id) {
  ['pin-section','success-section','already-section','unsupported-section'].forEach(s => {
    document.getElementById(s).classList.add('hidden');
  });
  document.getElementById(id).classList.remove('hidden');
}

// ── numpad ───────────────────────────────────────────────────
function updateDots() {
  document.querySelectorAll('.dot').forEach((d, i) => {
    d.classList.toggle('bg-emerald-400', i < _pin.length);
    d.classList.toggle('border-emerald-400', i < _pin.length);
    d.classList.toggle('border-zinc-700', i >= _pin.length);
  });
}

function shake() {
  const dots = document.getElementById('pin-dots');
  dots.style.transform = 'translateX(8px)';
  setTimeout(() => dots.style.transform = 'translateX(-8px)', 80);
  setTimeout(() => dots.style.transform = '', 160);
}

function pressDigit(d) {
  if (_pin.length >= PIN_LEN) return;
  _pin += d;
  updateDots();
  if (_pin.length === PIN_LEN) verifythenRegister();
}

function backspace() {
  _pin = _pin.slice(0, -1);
  updateDots();
  document.getElementById('pin-error').innerHTML = '';
}

// ── verify PIN with server, then trigger WebAuthn registration ─
async function verifythenRegister() {
  const errEl = document.getElementById('pin-error');
  errEl.innerHTML = '';

  // Verify PIN with server (reuse the same /login endpoint)
  const res = await fetch('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({ pin: _pin }),
  });

  if (!res.ok) {
    shake();
    _pin = '';
    updateDots();
    errEl.innerHTML = '<span class="text-red-400 text-sm">Incorrect PIN</span>';
    return;
  }

  // PIN correct — now register the biometric
  const pinToStore = _pin; // capture before async
  _pin = '';
  updateDots();

  try {
    errEl.innerHTML = '<span class="text-zinc-500 text-sm">Touch sensor / Face ID…</span>';

    const uid = crypto.getRandomValues(new Uint8Array(16));
    const cred = await navigator.credentials.create({
      publicKey: {
        challenge:        crypto.getRandomValues(new Uint8Array(32)),
        rp:               { name: 'Echo Finance', id: location.hostname },
        user:             { id: uid, name: 'owner', displayName: 'Owner' },
        pubKeyCredParams: [{ type: 'public-key', alg: -7 }, { type: 'public-key', alg: -257 }],
        authenticatorSelection: {
          authenticatorAttachment: 'platform',
          userVerification: 'required',
          residentKey: 'preferred',
        },
        timeout: 60000,
      }
    });

    // Store credential id + PIN locally (single-user personal app)
    localStorage.setItem('echo_cred_id', btoa(String.fromCharCode(...new Uint8Array(cred.rawId))));
    localStorage.setItem('echo_bio_pin', pinToStore);

    // Update icon to success
    document.getElementById('icon-wrap').classList.add('border-emerald-500/30', 'bg-emerald-500/10');
    document.getElementById('bio-icon').classList.replace('text-zinc-400', 'text-emerald-400');

    show('success-section');

  } catch (e) {
    errEl.innerHTML = '<span class="text-zinc-500 text-sm">Cancelled — try again</span>';
  }
}

function resetBiometric() {
  localStorage.removeItem('echo_cred_id');
  localStorage.removeItem('echo_bio_pin');
  show('pin-section');
}

// keyboard support (desktop testing)
document.addEventListener('keydown', e => {
  if (e.key >= '0' && e.key <= '9') pressDigit(e.key);
  if (e.key === 'Backspace') backspace();
});
</script>

</body>
</html>`
}
