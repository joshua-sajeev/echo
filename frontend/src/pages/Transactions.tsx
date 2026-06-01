import { useEffect, useState } from "react";

// ── types ─────────────────────────────────────────────────────────────────────
interface Account {
  id: number;
  name: string;
  Balance: number;
}

interface Jar {
  id: number;
  name: string;
  balance?: number;
}

type TxType = "expense" | "income" | "transfer";

// ── API base ──────────────────────────────────────────────────────────────────
const API_BASE = "/api/v1";

// ── helpers ───────────────────────────────────────────────────────────────────
const fmt = (n: number) =>
  "₹" + n.toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

const today = () => new Date().toISOString().split("T")[0];

// ── sub-components ────────────────────────────────────────────────────────────

function SelectField({
  label,
  name,
  value,
  onChange,
  options,
  loading,
  placeholder = "Select…",
  required,
}: {
  label: string;
  name: string;
  value: string;
  onChange: (v: string) => void;
  options: { id: number; label: string; sub?: string }[];
  loading?: boolean;
  placeholder?: string;
  required?: boolean;
}) {
  return (
    <div>
      <label style={labelStyle}>{label}{required && " *"}</label>
      <div style={{ position: "relative" }}>
        <select
          name={name}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={required}
          disabled={loading}
          style={{
            ...inputStyle,
            appearance: "none",
            paddingRight: 36,
            color: value ? "#d1d5db" : "#4b5563",
            opacity: loading ? 0.5 : 1,
          }}
        >
          <option value="" disabled>{loading ? "Loading…" : placeholder}</option>
          {options.map((o) => (
            <option key={o.id} value={String(o.id)}>
              {o.label}{o.sub ? ` — ${o.sub}` : ""}
            </option>
          ))}
        </select>
        {/* chevron */}
        <svg
          style={{ position: "absolute", right: 12, top: "50%", transform: "translateY(-50%)", pointerEvents: "none" }}
          width="14" height="14" viewBox="0 0 24 24" fill="none"
          stroke="#4b5563" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"
        >
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </div>
    </div>
  );
}

// ── main component ────────────────────────────────────────────────────────────
export default function NewTransaction() {
  const [type, setType] = useState<TxType>("expense");
  const [date, setDate] = useState(today());
  const [name, setName] = useState("");
  const [amount, setAmount] = useState("");

  // dynamic fields
  const [accountId, setAccountId] = useState("");     // expense / income single account
  const [fromId, setFromId] = useState("");            // transfer
  const [toId, setToId] = useState("");                // transfer
  const [jarId, setJarId] = useState("");
  const [isMasterIncome, setIsMasterIncome] = useState(false);

  // remote data
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [jars, setJars] = useState<Jar[]>([]);
  const [loadingAccounts, setLoadingAccounts] = useState(true);
  const [loadingJars, setLoadingJars] = useState(true);

  // submission
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  // load accounts + jars once
  useEffect(() => {
    fetch(`${API_BASE}/accounts/balances`, { credentials: "include", headers: { Accept: "application/json" } })
      .then((r) => r.json())
      .then((d: Account[]) => setAccounts(d))
      .catch(() => {})
      .finally(() => setLoadingAccounts(false));

    fetch(`${API_BASE}/jars`, { credentials: "include", headers: { Accept: "application/json" } })
      .then((r) => r.json())
      .then((d: Jar[]) => setJars(d))
      .catch(() => {})
      .finally(() => setLoadingJars(false));
  }, []);

  // reset dependent fields when type changes
  useEffect(() => {
    setAccountId("");
    setFromId("");
    setToId("");
    setJarId("");
    setIsMasterIncome(false);
  }, [type]);

  // when master income toggled, clear jar
  useEffect(() => {
    if (isMasterIncome) setJarId("");
  }, [isMasterIncome]);

  const accountOptions = accounts.map((a) => ({
    id: a.id,
    label: a.name,
    sub: fmt(a.Balance),
  }));

  const jarOptions = jars.map((j) => ({
    id: j.id,
    label: j.name,
    sub: j.balance !== undefined ? fmt(j.balance) : undefined,
  }));

const handleSubmit = async () => {
  setSubmitting(true);
  setSubmitError(null);

  // 1. Convert standard decimal string input to integer Paisa
  const numericAmount = parseFloat(amount) || 0;
  const amountInPaisa = Math.round(numericAmount * 100); 

  // 2. Build the payload matching the Go struct tags exactly
  const base: Record<string, unknown> = {
    type,
    date: new Date(date).toISOString(), // Ensures valid time.Time parsing on backend
    name: name.trim(),
    amount: amountInPaisa, // Sent as an integer over the wire
  };

  if (type === "expense") {
    // Note: Map frontend accountId to destination field expected by Go
    // Since expense subtracts money, you might use account_id or from_account_id depending on your implementation.
    // If your Go struct uses from_account_id/to_account_id for everything:
    base.from_account_id = parseInt(accountId); 
    if (jarId) base.jar_id = parseInt(jarId);
  } else if (type === "income") {
    base.to_account_id = parseInt(accountId);
    base.is_master_income = isMasterIncome; // FIX: Fixed snake_case mapping to match Go struct tag
    if (!isMasterIncome && jarId) base.jar_id = parseInt(jarId);
  } else if (type === "transfer") {
    base.from_account_id = parseInt(fromId);
    base.to_account_id = parseInt(toId);
  }

  try {
    const r = await fetch(`${API_BASE}/transactions`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json", Accept: "application/json" },
      body: JSON.stringify(base),
    });
    if (!r.ok) {
      const text = await r.text();
      throw new Error(`Server returned ${r.status}: ${text.slice(0, 200)}`);
    }
    window.location.href = "/";
  } catch (e: any) {
    setSubmitError(e.message);
    setSubmitting(false);
  }
};

  const canSubmit =
    !submitting &&
    date &&
    name.trim() &&
    amount &&
    parseFloat(amount) > 0 &&
    (type === "expense" ? !!accountId :
     type === "income"  ? !!accountId :
     /* transfer */       !!fromId && !!toId && fromId !== toId);

  return (
    <>
      <style>{`
        * { box-sizing: border-box; }
        body { margin: 0; background: #080a0f; }
        select option { background: #0f1117; color: #d1d5db; }
        input[type="date"]::-webkit-calendar-picker-indicator { filter: invert(0.4); }
      `}</style>

      {/* nav */}
      <nav style={{
        borderBottom: "0.5px solid #1e2130",
        background: "#080a0f",
        position: "sticky", top: 0, zIndex: 10,
      }}>
        <div style={navInner}>
          <a href="/" style={{ color: "#6b7280", fontSize: 14, textDecoration: "none" }}>← Back</a>
          <span style={{ fontSize: 15, fontWeight: 600, color: "#e5e7eb" }}>New Transaction</span>
          <div style={{ width: 48 }} />
        </div>
      </nav>

      <main style={{ maxWidth: 480, margin: "0 auto", padding: "24px 16px 80px" }}>

        {/* type tabs */}
        <div style={tabBar}>
          {(["expense", "income", "transfer"] as TxType[]).map((t) => (
            <button
              key={t}
              onClick={() => setType(t)}
              style={{
                ...tabBtn,
                background: type === t ? "#1e2130" : "transparent",
                color: type === t ? "#e5e7eb" : "#4b5563",
                fontWeight: type === t ? 600 : 400,
              }}
            >
              {t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>

        <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>

          {/* date */}
          <div>
            <label style={labelStyle}>Date *</label>
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              style={inputStyle}
              required
            />
          </div>

          {/* name */}
          <div>
            <label style={labelStyle}>Name *</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Zomato Order"
              style={inputStyle}
              required
            />
          </div>

          {/* amount */}
          <div>
            <label style={labelStyle}>Amount *</label>
            <div style={{ position: "relative" }}>
              <span style={{
                position: "absolute", left: 14, top: "50%", transform: "translateY(-50%)",
                color: "#6b7280", fontSize: 16, pointerEvents: "none",
                fontFamily: "'IBM Plex Mono', monospace",
              }}>₹</span>
              <input
                type="number"
                step="0.01"
                min="0"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="0.00"
                style={{ ...inputStyle, paddingLeft: 30, fontFamily: "'IBM Plex Mono', monospace" }}
                required
              />
            </div>
          </div>

          {/* ── dynamic fields by type ── */}

          {type === "expense" && (
            <>
              <SelectField
                label="Account"
                name="account"
                value={accountId}
                onChange={setAccountId}
                options={accountOptions}
                loading={loadingAccounts}
                placeholder="Select account"
                required
              />
              <SelectField
                label="Jar (optional)"
                name="jar_id"
                value={jarId}
                onChange={setJarId}
                options={jarOptions}
                loading={loadingJars}
                placeholder="No jar"
              />
            </>
          )}

          {type === "income" && (
            <>
              <SelectField
                label="Into account"
                name="account"
                value={accountId}
                onChange={setAccountId}
                options={accountOptions}
                loading={loadingAccounts}
                placeholder="Select account"
                required
              />

              {/* master income toggle */}
              <div
                onClick={() => setIsMasterIncome((v) => !v)}
                style={{
                  display: "flex", alignItems: "center", justifyContent: "space-between",
                  padding: "12px 14px",
                  background: "#0f1117", border: "0.5px solid #1e2130", borderRadius: 10,
                  cursor: "pointer", userSelect: "none",
                }}
              >
                <div>
                  <p style={{ margin: 0, fontSize: 14, color: "#d1d5db", fontWeight: 500 }}>Master income</p>
                  <p style={{ margin: "2px 0 0", fontSize: 12, color: "#4b5563" }}>
                    Distributes across all jars by allocation
                  </p>
                </div>
                {/* toggle pill */}
                <div style={{
                  width: 40, height: 24, borderRadius: 12, flexShrink: 0,
                  background: isMasterIncome ? "#1D9E75" : "#1e2130",
                  position: "relative", transition: "background 0.2s",
                }}>
                  <div style={{
                    position: "absolute", top: 3,
                    left: isMasterIncome ? 19 : 3,
                    width: 18, height: 18, borderRadius: "50%",
                    background: "#fff", transition: "left 0.2s",
                  }} />
                </div>
              </div>

              {!isMasterIncome && (
                <SelectField
                  label="Jar (optional)"
                  name="jar_id"
                  value={jarId}
                  onChange={setJarId}
                  options={jarOptions}
                  loading={loadingJars}
                  placeholder="No jar"
                />
              )}
            </>
          )}

          {type === "transfer" && (
            <>
              <SelectField
                label="From account"
                name="from"
                value={fromId}
                onChange={setFromId}
                options={accountOptions.filter((a) => String(a.id) !== toId)}
                loading={loadingAccounts}
                placeholder="Select account"
                required
              />
              <SelectField
                label="To account"
                name="to"
                value={toId}
                onChange={setToId}
                options={accountOptions.filter((a) => String(a.id) !== fromId)}
                loading={loadingAccounts}
                placeholder="Select account"
                required
              />
            </>
          )}

          {/* error */}
          {submitError && (
            <div style={{
              padding: "10px 14px", borderRadius: 10,
              background: "#1a0e0e", border: "0.5px solid #3d1515",
              color: "#E24B4A", fontSize: 13,
            }}>
              {submitError}
            </div>
          )}

          {/* submit */}
          <button
            onClick={handleSubmit}
            disabled={!canSubmit}
            style={{
              width: "100%", padding: "15px 0", borderRadius: 12, border: "none",
              background: canSubmit ? "#e5e7eb" : "#1a1d27",
              color: canSubmit ? "#080a0f" : "#374151",
              fontSize: 15, fontWeight: 700, cursor: canSubmit ? "pointer" : "not-allowed",
              transition: "background 0.15s, color 0.15s",
              marginTop: 4,
            }}
          >
            {submitting ? "Saving…" : "Save Transaction"}
          </button>

        </div>
      </main>
    </>
  );
}

// ── styles ────────────────────────────────────────────────────────────────────
const navInner: React.CSSProperties = {
  maxWidth: 480, margin: "0 auto", padding: "14px 16px",
  display: "flex", alignItems: "center", justifyContent: "space-between",
};

const tabBar: React.CSSProperties = {
  display: "flex", gap: 4,
  background: "#0f1117", border: "0.5px solid #1e2130",
  borderRadius: 12, padding: 4, marginBottom: 20,
};

const tabBtn: React.CSSProperties = {
  flex: 1, padding: "8px 0", borderRadius: 9, border: "none",
  fontSize: 14, cursor: "pointer", transition: "all 0.15s", fontFamily: "inherit",
};

const labelStyle: React.CSSProperties = {
  display: "block", fontSize: 11, fontWeight: 600,
  letterSpacing: "0.07em", textTransform: "uppercase",
  color: "#4b5563", marginBottom: 8,
};

const inputStyle: React.CSSProperties = {
  width: "100%", background: "#0f1117",
  border: "0.5px solid #1e2130", borderRadius: 10,
  padding: "12px 14px", color: "#d1d5db", fontSize: 15,
  outline: "none", fontFamily: "inherit",
};
