package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joshu-sajeev/echo/internal/allocation"
	"github.com/joshu-sajeev/echo/internal/repository"
)

type JarHandler struct {
	jarRepo repository.JarRepository
	txRepo  repository.TransactionRepository
}

func NewJarHandler(jarRepo repository.JarRepository, txRepo repository.TransactionRepository) *JarHandler {
	return &JarHandler{jarRepo: jarRepo, txRepo: txRepo}
}

// GET /jars/page — full settings page
func (h *JarHandler) Page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, jarSettingsPage())
}

// GET /jars — returns jar rows partial (used by htmx)
func (h *JarHandler) List(w http.ResponseWriter, r *http.Request) {
	jars, spent, err := h.jarRepo.ListWithSpend(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := h.txRepo.Stats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allocs := allocation.Calculate(stats.MonthlyMasterIncome, jars, spent)
	summary := allocation.Summarise(allocs)

	var sb strings.Builder

	// summary bar at top
	sb.WriteString(fmt.Sprintf(`
	<div class="mb-6 bg-zinc-950 border border-zinc-800 rounded-xl p-4 space-y-2">
		<div class="flex justify-between text-xs text-zinc-500">
			<span>Monthly Master Income</span>
			<span class="text-zinc-300 font-medium">₹%.0f</span>
		</div>
		<div class="flex justify-between text-xs text-zinc-500">
			<span>Total Allocated</span>
			<span class="text-zinc-300">₹%.0f</span>
		</div>
		<div class="flex justify-between text-xs text-zinc-500">
			<span>Total Spent</span>
			<span class="text-red-400">₹%.0f</span>
		</div>
	</div>
	`, stats.MonthlyMasterIncome, summary.TotalAllocated, summary.TotalSpent))

	if len(allocs) == 0 {
		sb.WriteString(`<p class="text-zinc-600 text-xs">No jars yet</p>`)
	}

	for _, a := range allocs {
		sb.WriteString(jarRow(a))
	}

	// add jar form (non-system only)
	sb.WriteString(`
	<div class="mt-6 border-t border-zinc-800 pt-6">
		<p class="text-xs text-zinc-600 uppercase tracking-widest mb-3">Add Custom Jar</p>
		<form hx-post="/jars"
		      hx-target="#jar-list"
		      hx-swap="innerHTML"
		      hx-on::after-request="if(event.detail.successful) this.reset()"
		      class="flex gap-2 items-end flex-wrap">
			<div class="flex-1 min-w-[140px]">
				<label class="block text-xs text-zinc-500 mb-1">Name</label>
				<input name="name" required placeholder="e.g. Travel"
					class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2 text-sm outline-none focus:border-zinc-600"/>
			</div>
			<div class="w-32">
				<label class="block text-xs text-zinc-500 mb-1">Fixed Amount (₹)</label>
				<input name="value" type="number" step="1" min="0" placeholder="0"
					class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2 text-sm outline-none focus:border-zinc-600"/>
			</div>
			<button type="submit"
				class="bg-white text-black px-4 py-2 rounded-lg text-sm font-medium hover:bg-zinc-200 transition">
				Add
			</button>
		</form>
	</div>
	`)

	w.Write([]byte(sb.String()))
}

// PATCH /jars/{id} — update name + allocation_value
func (h *JarHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	var value float64
	fmt.Sscanf(r.FormValue("value"), "%f", &value)

	if err := h.jarRepo.Update(r.Context(), id, name, value); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		log.Println("jar update error:", err)
		return
	}

	w.Header().Set("HX-Trigger", `{"jarUpdated": true}`)
	w.Write([]byte(`<div></div>`))
}

// POST /jars — create a new non-system jar
func (h *JarHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	var value float64
	fmt.Sscanf(r.FormValue("value"), "%f", &value)

	if _, err := h.jarRepo.Create(r.Context(), name, value, 99); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			http.Error(w, fmt.Sprintf(`"%s" already exists`, name), http.StatusConflict)
			return
		}
		http.Error(w, "create failed", http.StatusInternalServerError)
		log.Println("jar create error:", err)
		return
	}

	// re-render the full list
	h.List(w, r)
}

// DELETE /jars/{id} — only non-system jars
func (h *JarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.jarRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "delete failed (system jars cannot be deleted)", http.StatusForbidden)
		return
	}

	w.Header().Set("HX-Trigger", `{"jarUpdated": true}`)
	w.Write([]byte(`<div></div>`))
}

// ── rendering helpers ────────────────────────────────────────────────────────

func jarRow(a allocation.JarAllocation) string {
	pct := 0
	if a.Allocated > 0 {
		pct = int(a.Spent / a.Allocated * 100)
		if pct > 100 {
			pct = 100
		}
	}

	barColor := "bg-emerald-400"
	if pct >= 100 {
		barColor = "bg-red-400"
	} else if pct >= 75 {
		barColor = "bg-yellow-400"
	}

	allocLabel := fmt.Sprintf("₹%.0f / ₹%.0f", a.Spent, a.Allocated)
	if a.IsCap {
		allocLabel = fmt.Sprintf("₹%.0f / ₹%.0f cap", a.Spent, a.Allocated)
	}

	remainLabel := ""
	if a.Remaining < 0 {
		remainLabel = fmt.Sprintf(`<span class="text-red-400 text-xs">₹%.0f over</span>`, -a.Remaining)
	} else {
		remainLabel = fmt.Sprintf(`<span class="text-zinc-600 text-xs">₹%.0f left</span>`, a.Remaining)
	}

	systemBadge := ""
	deleteBtn := ""
	if a.IsSystem {
		systemBadge = `<span class="text-[10px] text-zinc-600 border border-zinc-800 rounded px-1.5 py-0.5 ml-1">system</span>`
	} else {
		deleteBtn = fmt.Sprintf(`
		<button
			hx-delete="/jars/%d"
			hx-target="#jar-list"
			hx-swap="innerHTML"
			hx-confirm="Delete jar '%s'? Existing transactions keep their jar association."
			class="text-zinc-700 hover:text-red-400 transition text-xs px-2 py-1 rounded border border-zinc-800 hover:border-red-800">
			Delete
		</button>`, a.ID, a.Name)
	}

	// value label depends on jar type
	valuePlaceholder := "Fixed amount (₹)"
	valueNote := ""
	switch a.Name {
	case "Charity":
		valuePlaceholder = "Percentage of income"
		valueNote = `<p class="text-[10px] text-zinc-600 mt-0.5">% of total master income, rounded down to ₹100</p>`
	case "SIP":
		valueNote = `<p class="text-[10px] text-zinc-600 mt-0.5">Fixed amount deducted from pool each month</p>`
	case "Chitty":
		valuePlaceholder = "Spending cap (₹)"
		valueNote = `<p class="text-[10px] text-zinc-600 mt-0.5">Max spending cap — not deducted from allocation pool</p>`
	case "Necessities":
		valueNote = `<p class="text-[10px] text-zinc-600 mt-0.5">Auto-calculated: pool − leisure, rounded up to ₹500</p>`
	case "Leisure":
		valueNote = `<p class="text-[10px] text-zinc-600 mt-0.5">Auto-calculated: remainder after all other jars (≤ 10% of income)</p>`
	}

	// necessities and leisure have no editable value
	valueInput := fmt.Sprintf(`
		<input name="value" type="number" step="1" min="0"
			value="%.0f"
			placeholder="%s"
			class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2 text-sm outline-none focus:border-zinc-600"/>
		%s`, a.AllocationValue, valuePlaceholder, valueNote)

	if a.Name == "Necessities" || a.Name == "Leisure" {
		valueInput = fmt.Sprintf(`
		<div class="bg-zinc-950 border border-zinc-800/50 rounded-lg px-3 py-2 text-sm text-zinc-600 italic">
			₹%.0f (auto-calculated)
		</div>
		%s`, a.Allocated, valueNote)
	}

	return fmt.Sprintf(`
	<div class="bg-zinc-950 border border-zinc-800 rounded-xl p-4 space-y-3" id="jar-%d">

		<!-- progress bar -->
		<div>
			<div class="flex justify-between items-center mb-1.5">
				<div class="flex items-center gap-1">
					<span class="text-sm font-medium text-zinc-200">%s</span>
					%s
				</div>
				<div class="flex items-center gap-2">
					%s
					%s
				</div>
			</div>
			<div class="flex justify-between text-xs text-zinc-500 mb-1">
				<span>%s</span>
				<span>%d%%</span>
			</div>
			<div class="w-full bg-zinc-800 h-1.5 rounded-full">
				<div class="%s h-1.5 rounded-full transition-all" style="width:%d%%"></div>
			</div>
		</div>

		<!-- inline edit form -->
		<form hx-patch="/jars/%d"
		      hx-target="#jar-list"
		      hx-swap="innerHTML"
		      class="grid grid-cols-2 gap-2 pt-1 border-t border-zinc-800/50">
			<div>
				<label class="block text-[10px] text-zinc-600 mb-1 uppercase tracking-wider">Name</label>
				<input name="name" value="%s" required
					class="w-full bg-zinc-950 border border-zinc-800 rounded-lg px-3 py-2 text-sm outline-none focus:border-zinc-600"/>
			</div>
			<div>
				<label class="block text-[10px] text-zinc-600 mb-1 uppercase tracking-wider">%s</label>
				%s
			</div>
			<div class="col-span-2 flex justify-end">
				<button type="submit"
					class="text-xs bg-zinc-800 hover:bg-zinc-700 text-zinc-300 px-3 py-1.5 rounded-lg transition">
					Save
				</button>
			</div>
		</form>

	</div>
	`,
		a.ID,
		a.Name, systemBadge,
		remainLabel, deleteBtn,
		allocLabel, pct,
		barColor, pct,
		a.ID,
		a.Name,
		valuePlaceholder,
		valueInput,
	)
}

func jarSettingsPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Jars – Echo Finance</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body class="bg-zinc-950 text-white min-h-screen">

<nav class="border-b border-zinc-800 bg-zinc-950 sticky top-0 z-10">
  <div class="max-w-lg mx-auto px-6 py-4 flex items-center gap-4">
    <a href="/" class="text-zinc-500 hover:text-white transition text-sm shrink-0">← Back</a>
    <h1 class="text-base font-semibold flex-1">Allocation Jars</h1>
  </div>
</nav>

<main class="max-w-lg mx-auto px-6 py-6">
  <p class="text-xs text-zinc-600 mb-6">
    Allocation is calculated from this month's master income.
    System jars use fixed rules — only their values are editable.
  </p>

  <div id="jar-list"
    hx-get="/jars"
    hx-trigger="load, jarUpdated from:body"
    hx-swap="innerHTML"
    class="space-y-3">
  </div>
</main>
</body>
</html>`
}
