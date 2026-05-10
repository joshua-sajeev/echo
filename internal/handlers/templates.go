package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joshu-sajeev/echo/internal/models"
	"github.com/joshu-sajeev/echo/internal/repository"
)

type TxTemplateHandler struct {
	repo    repository.TxTemplateRepository
	jarRepo repository.JarRepository
	accRepo repository.AccountRepository
}

func NewTxTemplateHandler(
	repo repository.TxTemplateRepository,
	jarRepo repository.JarRepository,
	accRepo repository.AccountRepository,
) *TxTemplateHandler {
	return &TxTemplateHandler{repo: repo, jarRepo: jarRepo, accRepo: accRepo}
}

// GET /templates — full page
func (h *TxTemplateHandler) Page(w http.ResponseWriter, r *http.Request) {
	tmpls, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jars, err := h.jarRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accounts, err := h.accRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var jarOpts, accOpts string
	for _, j := range jars {
		jarOpts += fmt.Sprintf(`<option value="%d">%s</option>`, j.ID, j.Name)
	}
	for _, a := range accounts {
		accOpts += fmt.Sprintf(`<option value="%d">%s</option>`, a.ID, a.Name)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, templatesPage(tmpls, jarOpts, accOpts))
}

// POST /templates
func (h *TxTemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	txType := r.FormValue("type")
	if name == "" || txType == "" {
		http.Error(w, "name and type required", http.StatusBadRequest)
		return
	}

	var amount float64
	fmt.Sscanf(r.FormValue("amount"), "%f", &amount)

	t := models.TxTemplate{
		Name:           name,
		Type:           txType,
		Amount:         amount,
		IsMasterIncome: r.FormValue("master_income") == "true",
	}

	if jarIDStr := r.FormValue("jar_id"); jarIDStr != "" && jarIDStr != "0" {
		var jid int64
		if _, err := fmt.Sscanf(jarIDStr, "%d", &jid); err == nil {
			t.JarID = &jid
		}
	}
	if fromStr := r.FormValue("from_account_id"); fromStr != "" {
		var id int64
		if _, err := fmt.Sscanf(fromStr, "%d", &id); err == nil {
			t.FromAccountID = &id
		}
	}
	if toStr := r.FormValue("to_account_id"); toStr != "" {
		var id int64
		if _, err := fmt.Sscanf(toStr, "%d", &id); err == nil {
			t.ToAccountID = &id
		}
	}

	if _, err := h.repo.Create(r.Context(), t); err != nil {
		http.Error(w, "create failed: "+err.Error(), http.StatusInternalServerError)
		log.Println("template create error:", err)
		return
	}

	h.listPartial(w, r)
}

// DELETE /templates/{id}
func (h *TxTemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	h.listPartial(w, r)
}

func (h *TxTemplateHandler) listPartial(w http.ResponseWriter, r *http.Request) {
	tmpls, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, templateCards(tmpls))
}

func (h *TxTemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	h.listPartial(w, r)
}

// ── rendering ────────────────────────────────────────────────────────────────

func templateCards(tmpls []models.TxTemplate) string {
	if len(tmpls) == 0 {
		return `<p class="text-zinc-600 text-xs col-span-2 py-4 text-center">No templates yet — add one below</p>`
	}

	var sb strings.Builder
	for _, t := range tmpls {
		typeColor := "text-red-400"
		typeBg := "bg-red-500/10 border-red-500/20"
		typeLabel := "Expense"
		switch t.Type {
		case "income":
			typeColor = "text-emerald-400"
			typeBg = "bg-emerald-500/10 border-emerald-500/20"
			typeLabel = "Income"
		case "transfer":
			typeColor = "text-sky-400"
			typeBg = "bg-sky-500/10 border-sky-500/20"
			typeLabel = "Transfer"
		}

		amtStr := "any amount"
		if t.Amount > 0 {
			amtStr = fmt.Sprintf("₹%.0f", t.Amount)
		}

		var details []string
		switch t.Type {
		case "expense":
			if t.FromAccountName != "" {
				details = append(details, t.FromAccountName)
			}
			if t.JarName != "" {
				details = append(details, t.JarName)
			}
		case "income":
			if t.ToAccountName != "" {
				details = append(details, t.ToAccountName)
			}
			if t.IsMasterIncome {
				details = append(details, "master income")
			} else if t.JarName != "" {
				details = append(details, t.JarName)
			}
		case "transfer":
			if t.FromAccountName != "" && t.ToAccountName != "" {
				details = append(details, t.FromAccountName+" → "+t.ToAccountName)
			}
		}
		subtitle := strings.Join(details, " · ")
		if subtitle == "" {
			subtitle = "—"
		}

		fromAcc := strings.ReplaceAll(t.FromAccountName, "'", "\\'")
		toAcc := strings.ReplaceAll(t.ToAccountName, "'", "\\'")
		jarName := strings.ReplaceAll(t.JarName, "'", "\\'")
		tName := strings.ReplaceAll(t.Name, "'", "\\'")

		fromID := "0"
		if t.FromAccountID != nil {
			fromID = fmt.Sprintf("%d", *t.FromAccountID)
		}
		toID := "0"
		if t.ToAccountID != nil {
			toID = fmt.Sprintf("%d", *t.ToAccountID)
		}
		jarID := "0"
		if t.JarID != nil {
			jarID = fmt.Sprintf("%d", *t.JarID)
		}
		isMaster := "false"
		if t.IsMasterIncome {
			isMaster = "true"
		}

		sb.WriteString(fmt.Sprintf(`
<div class="relative group bg-zinc-950 border border-zinc-800 hover:border-zinc-600 rounded-xl p-4 cursor-pointer transition-all active:scale-95"
     onclick="useTemplate(%d,'%s','%s',%.2f,'%s',%s,%s,%s,%s,'%s','%s')">
  <div class="flex items-start justify-between gap-2 mb-3">
    <p class="font-medium text-sm text-zinc-100 leading-tight">%s</p>
    <span class="text-[10px] border rounded px-1.5 py-0.5 shrink-0 %s %s">%s</span>
  </div>
  <div class="space-y-1">
    <p class="text-lg font-semibold %s">%s</p>
    <p class="text-xs text-zinc-600 truncate">%s</p>
  </div>
  <button
    onclick="event.stopPropagation(); deleteTemplate(%d)"
    class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 text-zinc-700 hover:text-red-400 transition text-sm p-1 rounded">
    ✕
  </button>
</div>
`,
			t.ID, tName, t.Type, t.Amount, jarName,
			fromID, toID, jarID, isMaster, fromAcc, toAcc,
			t.Name,
			typeBg, typeColor, typeLabel,
			typeColor, amtStr,
			subtitle,
			t.ID,
		))
	}
	return sb.String()
}

func templatesPage(tmpls []models.TxTemplate, jarOpts, accOpts string) string {
	cards := templateCards(tmpls)

	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Templates – Echo Finance</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body class="bg-zinc-950 text-white min-h-screen">

<nav class="border-b border-zinc-800 bg-zinc-950 sticky top-0 z-10">
  <div class="max-w-lg mx-auto px-6 py-4 flex items-center gap-4">
    <a href="/" class="text-zinc-500 hover:text-white transition text-sm shrink-0">← Back</a>
    <h1 class="text-base font-semibold flex-1">Transaction Templates</h1>
  </div>
</nav>

<main class="max-w-lg mx-auto px-6 py-6 space-y-6">

  <p class="text-xs text-zinc-600">
    Tap a template to instantly record a transaction. If no amount is set, you'll be asked to enter one.
  </p>

  <div id="template-list" class="grid grid-cols-2 gap-3">
    ` + cards + `
  </div>

  <!-- add template form -->
  <div class="border-t border-zinc-800 pt-6 space-y-4">
    <p class="text-xs text-zinc-600 uppercase tracking-widest">New Template</p>

    <form id="tmpl-form" class="space-y-3"
          hx-post="/templates"
          hx-target="#template-list"
          hx-swap="innerHTML"
          hx-on::after-request="if(event.detail.successful) resetTmplForm()">

      <div class="flex rounded-lg bg-zinc-900 border border-zinc-800 p-1">
        <button type="button" id="tt-expense" onclick="setTmplType('expense')"
          class="tt-btn flex-1 py-2 text-sm rounded-md font-medium transition bg-zinc-800 text-white">
          Expense
        </button>
        <button type="button" id="tt-income" onclick="setTmplType('income')"
          class="tt-btn flex-1 py-2 text-sm rounded-md font-medium transition text-zinc-500">
          Income
        </button>
        <button type="button" id="tt-transfer" onclick="setTmplType('transfer')"
          class="tt-btn flex-1 py-2 text-sm rounded-md font-medium transition text-zinc-500">
          Transfer
        </button>
      </div>
      <input type="hidden" name="type" id="tmpl-type" value="expense"/>

      <div>
        <label class="block text-xs text-zinc-500 mb-1">Name *</label>
        <input name="name" required placeholder="e.g. Zomato, Petrol, SIP"
          class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2.5 text-sm outline-none focus:border-zinc-600 placeholder-zinc-700"/>
      </div>

      <div>
        <label class="block text-xs text-zinc-500 mb-1">Amount (₹) — leave 0 to enter each time</label>
        <input name="amount" type="number" step="1" min="0" placeholder="0"
          class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2.5 text-sm outline-none focus:border-zinc-600 placeholder-zinc-700"/>
      </div>

      <div id="tmpl-dynamic-fields" class="space-y-3"></div>

      <button type="submit"
        class="w-full bg-white text-black py-2.5 rounded-lg text-sm font-medium hover:bg-zinc-200 transition">
        Save Template
      </button>
    </form>
  </div>

</main>

<!-- MODAL: shown when amount=0 (ask amount+date) or amount>0 (ask nothing, post directly) -->
<div id="use-modal"
  class="hidden fixed inset-0 z-50 flex items-end justify-center bg-black/60 backdrop-blur-sm">
  <div class="bg-zinc-900 border border-zinc-800 rounded-t-3xl p-6 w-full max-w-lg shadow-2xl"
       onclick="event.stopPropagation()">
    <div class="w-10 h-1 bg-zinc-700 rounded-full mx-auto mb-5"></div>

    <!-- template name + subtitle -->
    <p id="modal-label" class="font-semibold text-base mb-0.5"></p>
    <p id="modal-sub"   class="text-xs text-zinc-500 mb-5"></p>

    <!-- amount field — only shown when template amount = 0 -->
    <div id="modal-amount-row" class="mb-4">
      <label class="block text-xs text-zinc-500 mb-1.5">Amount (₹) *</label>
      <input id="modal-amount" type="number" step="1" min="1"
        placeholder="Enter amount"
        class="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-3 text-lg font-semibold outline-none focus:border-zinc-600"
        onkeydown="if(event.key==='Enter') submitModal()"/>
    </div>

    <!-- date field — only shown when template amount = 0 -->
    <div id="modal-date-row" class="mb-5">
      <label class="block text-xs text-zinc-500 mb-1.5">Date *</label>
      <input id="modal-date" type="date"
        class="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-3 text-sm outline-none focus:border-zinc-600"/>
    </div>

    <div id="modal-error" class="hidden text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2 mb-4"></div>

    <div class="grid grid-cols-2 gap-3">
      <button onclick="closeModal()"
        class="py-3 rounded-xl border border-zinc-800 text-zinc-400 text-sm hover:bg-zinc-800 transition">
        Cancel
      </button>
      <button id="modal-submit-btn" onclick="submitModal()"
        class="py-3 rounded-xl bg-white text-black text-sm font-semibold hover:bg-zinc-200 transition">
        Save →
      </button>
    </div>
  </div>
</div>

<!-- result toast -->
<div id="result-toast"
  class="hidden fixed bottom-6 left-1/2 -translate-x-1/2 z-50 px-5 py-3 rounded-xl shadow-2xl text-sm font-medium transition-all">
</div>

<script>
const _jarOptsHTML = ` + "`" + `<option value="">Select jar</option>` + "`" + ` + ` + "`" + jarOpts + "`" + `;
const _accOptsHTML = ` + "`" + `<option value="">Select account</option>` + "`" + ` + ` + "`" + accOpts + "`" + `;

// ── template type form ────────────────────────────────────────
let _tmplType = 'expense';

function setTmplType(type) {
  _tmplType = type;
  document.getElementById('tmpl-type').value = type;
  document.querySelectorAll('.tt-btn').forEach(b => {
    b.classList.remove('bg-zinc-800','text-white');
    b.classList.add('text-zinc-500');
  });
  document.getElementById('tt-' + type).classList.add('bg-zinc-800','text-white');
  document.getElementById('tt-' + type).classList.remove('text-zinc-500');
  renderTmplDynamicFields(type, false);
}

function sel(name, label, optsHTML, required) {
  return ` + "`" + `<div>
    <label class="block text-xs text-zinc-500 mb-1">${label}</label>
    <select name="${name}" ${required ? 'required' : ''}
      class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2.5 text-sm outline-none focus:border-zinc-600 text-zinc-300">
      ${optsHTML}
    </select>
  </div>` + "`" + `;
}

function renderTmplDynamicFields(type, isMaster) {
  let html = '';
  if (type === 'expense') {
    html += sel('from_account_id', 'From Account *', _accOptsHTML, true);
    html += sel('jar_id', 'Category (Jar) *', _jarOptsHTML, true);
  } else if (type === 'income') {
    html += sel('to_account_id', 'To Account *', _accOptsHTML, true);
    html += ` + "`" + `<div><label class="flex items-center gap-3 cursor-pointer">
      <input type="checkbox" id="tmpl-master" name="master_income" value="true"
        onchange="onTmplMasterChange(this)" class="w-4 h-4 accent-emerald-400"/>
      <span class="text-sm text-zinc-300">Master Income <span class="text-xs text-zinc-500">(auto-split to jars)</span></span>
    </label></div>` + "`" + `;
    if (!isMaster) {
      html += sel('jar_id', 'Category (Jar)', _jarOptsHTML, false);
    }
  } else if (type === 'transfer') {
    html += sel('from_account_id', 'From Account *', _accOptsHTML, true);
    html += sel('to_account_id',   'To Account *',   _accOptsHTML, true);
  }
  document.getElementById('tmpl-dynamic-fields').innerHTML = html;
}

function onTmplMasterChange(cb) {
  renderTmplDynamicFields('income', cb.checked);
  const newCb = document.getElementById('tmpl-master');
  if (newCb) newCb.checked = cb.checked;
}

function resetTmplForm() {
  document.getElementById('tmpl-form').reset();
  setTmplType('expense');
}

setTmplType('expense');

// ── use template ──────────────────────────────────────────────
let _tmpl = {};

function useTemplate(id, name, type, amount, jarName, fromAccID, toAccID, jarID, isMaster, fromAccName, toAccName) {
  _tmpl = { id, name, type, amount, jarName, fromAccID, toAccID, jarID,
            isMaster: isMaster === 'true', fromAccName, toAccName };

  // build subtitle line
  let sub = '';
  if (type === 'expense')
    sub = [fromAccName ? 'from ' + fromAccName : '', jarName].filter(Boolean).join(' · ');
  else if (type === 'income')
    sub = [toAccName ? 'to ' + toAccName : '', isMaster === 'true' ? 'master income' : jarName].filter(Boolean).join(' · ');
  else if (type === 'transfer')
    sub = fromAccName + ' → ' + toAccName;

  document.getElementById('modal-label').textContent = name;
  document.getElementById('modal-sub').textContent   = sub;
  document.getElementById('modal-error').classList.add('hidden');

  const today = new Date().toISOString().split('T')[0];

  if (amount > 0) {
    // Has fixed amount — only ask for date, hide amount field
    document.getElementById('modal-amount-row').classList.add('hidden');
    document.getElementById('modal-date-row').classList.remove('hidden');
    document.getElementById('modal-date').value = today;
    document.getElementById('modal-submit-btn').textContent = 'Save →';
    document.getElementById('use-modal').classList.remove('hidden');
    setTimeout(() => document.getElementById('modal-date').focus(), 100);
  } else {
    // No fixed amount — ask for both amount and date
    document.getElementById('modal-amount-row').classList.remove('hidden');
    document.getElementById('modal-date-row').classList.remove('hidden');
    document.getElementById('modal-amount').value = '';
    document.getElementById('modal-date').value   = today;
    document.getElementById('modal-submit-btn').textContent = 'Save →';
    document.getElementById('use-modal').classList.remove('hidden');
    setTimeout(() => document.getElementById('modal-amount').focus(), 100);
  }
}

function closeModal() {
  document.getElementById('use-modal').classList.add('hidden');
}

// Post the transaction directly from the modal — no form page needed
async function submitModal() {
  const errEl  = document.getElementById('modal-error');
  errEl.classList.add('hidden');

  const amount = _tmpl.amount > 0
    ? _tmpl.amount
    : parseFloat(document.getElementById('modal-amount').value);

  const date = document.getElementById('modal-date').value;

  if (!amount || amount <= 0) {
    errEl.textContent = 'Please enter an amount.';
    errEl.classList.remove('hidden');
    document.getElementById('modal-amount').focus();
    return;
  }
  if (!date) {
    errEl.textContent = 'Please select a date.';
    errEl.classList.remove('hidden');
    return;
  }

  // Build form body matching what POST /transactions expects
  const body = new URLSearchParams({
    type:   _tmpl.type,
    name:   _tmpl.name,
    amount: amount,
    date:   date,
  });

  if (_tmpl.type === 'expense') {
    if (!_tmpl.fromAccID || _tmpl.fromAccID === '0') {
      errEl.textContent = 'Template has no account set — edit the template first.';
      errEl.classList.remove('hidden');
      return;
    }
    body.set('account', _tmpl.fromAccID);
    if (_tmpl.jarID && _tmpl.jarID !== '0') body.set('jar_id', _tmpl.jarID);

  } else if (_tmpl.type === 'income') {
    if (!_tmpl.toAccID || _tmpl.toAccID === '0') {
      errEl.textContent = 'Template has no account set — edit the template first.';
      errEl.classList.remove('hidden');
      return;
    }
    body.set('account', _tmpl.toAccID);
    if (_tmpl.isMaster) {
      body.set('master_income', 'true');
    } else if (_tmpl.jarID && _tmpl.jarID !== '0') {
      body.set('jar_id', _tmpl.jarID);
    }

  } else if (_tmpl.type === 'transfer') {
    if (!_tmpl.fromAccID || _tmpl.fromAccID === '0' || !_tmpl.toAccID || _tmpl.toAccID === '0') {
      errEl.textContent = 'Template has no accounts set — edit the template first.';
      errEl.classList.remove('hidden');
      return;
    }
    body.set('from', _tmpl.fromAccID);
    body.set('to',   _tmpl.toAccID);
  }

  const btn = document.getElementById('modal-submit-btn');
  btn.disabled = true;
  btn.textContent = 'Saving…';

  try {
    const res = await fetch('/transactions', {
      method:  'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body:    body,
    });

    if (res.ok) {
      closeModal();
      showToast('✓ ' + _tmpl.name + ' saved', 'success');
    } else {
      const t = await res.text();
      errEl.textContent = t || 'Failed to save transaction.';
      errEl.classList.remove('hidden');
      btn.disabled = false;
      btn.textContent = 'Save →';
    }
  } catch (e) {
    errEl.textContent = 'Network error — please try again.';
    errEl.classList.remove('hidden');
    btn.disabled = false;
    btn.textContent = 'Save →';
  }
}

function showToast(msg, type) {
  const t = document.getElementById('result-toast');
  t.textContent = msg;
  t.className = 'fixed bottom-6 left-1/2 -translate-x-1/2 z-50 px-5 py-3 rounded-xl shadow-2xl text-sm font-medium ' +
    (type === 'success'
      ? 'bg-emerald-500/20 border border-emerald-500/30 text-emerald-300'
      : 'bg-red-500/20 border border-red-500/30 text-red-300');
  t.classList.remove('hidden');
  setTimeout(() => t.classList.add('hidden'), 3000);
}

// close modal on backdrop tap
document.getElementById('use-modal').addEventListener('click', function(e) {
  if (e.target === this) closeModal();
});

function deleteTemplate(id) {
  if (!confirm('Delete this template?')) return;
  fetch('/templates/' + id, { method: 'DELETE' })
    .then(r => r.text())
    .then(html => { document.getElementById('template-list').innerHTML = html; });
}
</script>

</body>
</html>`
}
