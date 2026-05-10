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

// GET /templates — the full templates page
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

	var jarOpts string
	for _, j := range jars {
		jarOpts += fmt.Sprintf(`<option value="%d">%s</option>`, j.ID, j.Name)
	}

	var accOpts string
	for _, a := range accounts {
		accOpts += fmt.Sprintf(`<option value="%d">%s</option>`, a.ID, a.Name)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, templatesPage(tmpls, jarOpts, accOpts))
}

// POST /templates — create a new template
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
		Name:   name,
		Type:   txType,
		Amount: amount,
	}

	if jarIDStr := r.FormValue("jar_id"); jarIDStr != "" {
		var jid int64
		if _, err := fmt.Sscanf(jarIDStr, "%d", &jid); err == nil {
			t.JarID = &jid
		}
	}

	if _, err := h.repo.Create(r.Context(), t); err != nil {
		http.Error(w, "create failed: "+err.Error(), http.StatusInternalServerError)
		log.Println("template create error:", err)
		return
	}

	// re-render the list partial
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

// GET /templates/list — HTMX partial: just the cards
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
		if t.Type == "income" {
			typeColor = "text-emerald-400"
			typeBg = "bg-emerald-500/10 border-emerald-500/20"
			typeLabel = "Income"
		}

		amtStr := ""
		if t.Amount > 0 {
			amtStr = fmt.Sprintf("₹%.0f", t.Amount)
		} else {
			amtStr = "any amount"
		}

		jar := t.JarName
		if jar == "" {
			jar = "—"
		}

		// clicking the card goes to /transactions/new?template=ID
		sb.WriteString(fmt.Sprintf(`
<div class="relative group bg-zinc-950 border border-zinc-800 hover:border-zinc-600 rounded-xl p-4 cursor-pointer transition-all active:scale-95"
     onclick="useTemplate(%d, '%s', '%s', %.2f, '%s')">

  <div class="flex items-start justify-between gap-2 mb-3">
    <p class="font-medium text-sm text-zinc-100 leading-tight">%s</p>
    <span class="text-[10px] border rounded px-1.5 py-0.5 shrink-0 %s %s">%s</span>
  </div>

  <div class="space-y-1">
    <p class="text-lg font-semibold %s">%s</p>
    <p class="text-xs text-zinc-600">%s</p>
  </div>

  <!-- delete button — top-right on hover -->
  <button
    onclick="event.stopPropagation(); deleteTemplate(%d)"
    class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 text-zinc-700 hover:text-red-400 transition text-sm p-1 rounded">
    ✕
  </button>
</div>
`, t.ID, t.Name, t.Type, t.Amount, t.JarName,
			t.Name,
			typeBg, typeColor, typeLabel,
			typeColor, amtStr,
			jar,
			t.ID))
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
    Tap a template to open a pre-filled transaction form. Amount can be overridden before saving.
  </p>

  <!-- template cards grid -->
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

      <!-- type toggle -->
      <div class="flex rounded-lg bg-zinc-900 border border-zinc-800 p-1">
        <button type="button" id="tt-expense" onclick="setTmplType('expense')"
          class="tt-btn flex-1 py-2 text-sm rounded-md font-medium transition bg-zinc-800 text-white">
          Expense
        </button>
        <button type="button" id="tt-income" onclick="setTmplType('income')"
          class="tt-btn flex-1 py-2 text-sm rounded-md font-medium transition text-zinc-500">
          Income
        </button>
      </div>
      <input type="hidden" name="type" id="tmpl-type" value="expense"/>

      <div>
        <label class="block text-xs text-zinc-500 mb-1">Name *</label>
        <input name="name" required placeholder="e.g. Zomato, Petrol, SIP"
          class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2.5 text-sm outline-none focus:border-zinc-600 placeholder-zinc-700"/>
      </div>

      <div>
        <label class="block text-xs text-zinc-500 mb-1">Amount (₹) — leave 0 to fill each time</label>
        <input name="amount" type="number" step="1" min="0" placeholder="0"
          class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2.5 text-sm outline-none focus:border-zinc-600 placeholder-zinc-700"/>
      </div>

      <div>
        <label class="block text-xs text-zinc-500 mb-1">Category (Jar)</label>
        <select name="jar_id"
          class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2.5 text-sm outline-none focus:border-zinc-600 text-zinc-300">
          <option value="">Select jar</option>
          ` + jarOpts + `
        </select>
      </div>

      <button type="submit"
        class="w-full bg-white text-black py-2.5 rounded-lg text-sm font-medium hover:bg-zinc-200 transition">
        Save Template
      </button>
    </form>
  </div>

</main>

<!-- USE TEMPLATE MODAL — amount override before navigating -->
<div id="use-modal"
  class="hidden fixed inset-0 z-50 flex items-end justify-center bg-black/60 backdrop-blur-sm">
  <div class="bg-zinc-900 border border-zinc-800 rounded-t-3xl p-6 w-full max-w-lg space-y-4 shadow-2xl">
    <div class="w-10 h-1 bg-zinc-700 rounded-full mx-auto mb-2"></div>

    <div>
      <p id="modal-label" class="font-semibold text-base"></p>
      <p id="modal-jar" class="text-xs text-zinc-500 mt-0.5"></p>
    </div>

    <div>
      <label class="block text-xs text-zinc-500 mb-1.5">Amount (₹)</label>
      <input id="modal-amount" type="number" step="1" min="0"
        class="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-3 text-lg font-semibold outline-none focus:border-zinc-600"
        onkeydown="if(event.key==='Enter') confirmUse()"/>
    </div>

    <div class="grid grid-cols-2 gap-3 pt-1">
      <button onclick="closeModal()"
        class="py-3 rounded-xl border border-zinc-800 text-zinc-400 text-sm hover:bg-zinc-800 transition">
        Cancel
      </button>
      <button onclick="confirmUse()"
        class="py-3 rounded-xl bg-white text-black text-sm font-semibold hover:bg-zinc-200 transition">
        Use →
      </button>
    </div>
  </div>
</div>

<script>
let _tmpl = {};

function setTmplType(type) {
  document.getElementById('tmpl-type').value = type;
  document.querySelectorAll('.tt-btn').forEach(b => {
    b.classList.remove('bg-zinc-800','text-white');
    b.classList.add('text-zinc-500');
  });
  const btn = document.getElementById('tt-' + type);
  btn.classList.add('bg-zinc-800','text-white');
  btn.classList.remove('text-zinc-500');
}

function resetTmplForm() {
  document.getElementById('tmpl-form').reset();
  setTmplType('expense');
}

// Called when a template card is tapped
function useTemplate(id, name, type, amount, jar) {
  _tmpl = { id, name, type, amount, jar };
  document.getElementById('modal-label').textContent = name;
  document.getElementById('modal-jar').textContent   = jar && jar !== '—' ? '📂 ' + jar : '';
  document.getElementById('modal-amount').value      = amount > 0 ? amount : '';
  document.getElementById('modal-amount').placeholder = amount > 0 ? amount : 'Enter amount';
  document.getElementById('use-modal').classList.remove('hidden');
  setTimeout(() => document.getElementById('modal-amount').focus(), 100);
}

function closeModal() {
  document.getElementById('use-modal').classList.add('hidden');
}

function confirmUse() {
  const amt = document.getElementById('modal-amount').value || _tmpl.amount || 0;
  const params = new URLSearchParams({
    prefill_name:   _tmpl.name,
    prefill_type:   _tmpl.type,
    prefill_amount: amt,
    prefill_jar:    _tmpl.jar,
  });
  window.location.href = '/transactions/new?' + params.toString();
}

// Close modal on backdrop tap
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
