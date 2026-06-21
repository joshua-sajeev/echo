import { useEffect, useState } from "react";

// ── types ─────────────────────────────────────────────────────────────────────
export interface Account {
  id: number;
  name: string;
  Balance: number;
}

export interface Jar {
  id: number;
  name: string;
  balance?: number;
}

export type TxType = "expense" | "income" | "transfer";

export interface TransactionFormValues {
  type: TxType;
  date: string;
  name: string;
  amount: string;
  category: string;
  accountId: string;
  fromId: string;
  toId: string;
  jarId: string;
  isMasterIncome: boolean;
}

export interface TransactionFormProps {
  /** Pre-fill values for edit mode. Omit (or pass undefined) for create mode. */
  initialValues?: Partial<TransactionFormValues>;
  /** Called with the built payload on submit. Return a promise; the form handles loading/error state. */
  onSubmit: (payload: Record<string, unknown>) => Promise<void>;
  /** Label shown on the submit button (default: "Save Transaction") */
  submitLabel?: string;
}

// ── API base ──────────────────────────────────────────────────────────────────
const API_BASE = import.meta.env.VITE_API_URL; 

// ── helpers ───────────────────────────────────────────────────────────────────
const fmt = (n: number) =>
  "₹" +
  (n / 100).toLocaleString("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

const todayStr = () => new Date().toISOString().split("T")[0];

const defaults: TransactionFormValues = {
  type: "expense",
  date: todayStr(),
  name: "",
  amount: "",
  category: "",
  accountId: "",
  fromId: "",
  toId: "",
  jarId: "",
  isMasterIncome: false,
};
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
  options: {
    id: string | number;
    label: string;
    sub?: string;
  }[];
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
            {o.label}
            {o.sub ? ` — ${o.sub}` : ""}
          </option>
          ))}
        </select>
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
export default function TransactionForm({
  initialValues,
  onSubmit,
  submitLabel = "Save Transaction",
}: TransactionFormProps) {
  const init = { ...defaults, ...initialValues };

  const [type, setType] = useState<TxType>(init.type);
  const [date, setDate] = useState(init.date);
  const [name, setName] = useState(init.name);
  const [amount, setAmount] = useState(init.amount);
  const [accountId, setAccountId] = useState(init.accountId);
  const [fromId, setFromId] = useState(init.fromId);
  const [toId, setToId] = useState(init.toId);
  const [jarId, setJarId] = useState(init.jarId);
  const [isMasterIncome, setIsMasterIncome] = useState(init.isMasterIncome);

  const [accounts, setAccounts] = useState<Account[]>([]);
  const [jars, setJars] = useState<Jar[]>([]);
  const [loadingAccounts, setLoadingAccounts] = useState(true);
  const [loadingJars, setLoadingJars] = useState(true);

  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const [category, setCategory] = useState(init.category);
useEffect(() => {
  if (!initialValues) return;

  const values = {
    ...defaults,
    ...initialValues,
  };

  setType(values.type);
  setDate(values.date);
  setName(values.name);
  setAmount(values.amount);
  setCategory(values.category);
  setAccountId(values.accountId);
  setFromId(values.fromId);
  setToId(values.toId);
  setJarId(values.jarId);
  setIsMasterIncome(values.isMasterIncome);
}, [initialValues]);
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

  // reset dependent fields when type changes —
  // but only on user-driven type changes, not on initial render with initialValues
  const [hasTypeChanged, setHasTypeChanged] = useState(false);
  useEffect(() => {
    if (!hasTypeChanged) return;
    setAccountId("");
    setFromId("");
    setToId("");
    setJarId("");
    setIsMasterIncome(false);
  }, [type]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleTypeChange = (t: TxType) => {
    setHasTypeChanged(true);
    setType(t);
  };

  useEffect(() => {
    if (isMasterIncome) setJarId("");
  }, [isMasterIncome]);

  const CATEGORIES = [
    "Food",
    "Transport",
    "Shopping",
    "Donations",
    "Entertainment",
    "Health",
    "Income",
    "Investment",
    "Housing",
    "Transfers",
  ];
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

  const numericAmount = parseFloat(amount) || 0;
  const amountInPaisa = Math.round(numericAmount * 100);


  const payload: Record<string, unknown> = {
    type,
    date: new Date(date).toISOString(),
    name: name.trim(),
    amount: amountInPaisa,
    category: category || null,
  };

  // default ALL to null (important for clean updates)
  payload.from_account_id = null;
  payload.to_account_id = null;
  payload.jar_id = null;
  payload.is_master_income = null;

  if (type === "expense") {
    if (!accountId) throw new Error("Account required");

    payload.from_account_id = Number(accountId);
    if (jarId) payload.jar_id = Number(jarId);
  }

  if (type === "income") {
    if (!accountId) throw new Error("Account required");

    payload.to_account_id = Number(accountId);
    payload.is_master_income = isMasterIncome;

    if (!isMasterIncome && jarId) {
      payload.jar_id = Number(jarId);
    }
  }

  if (type === "transfer") {
    if (!fromId || !toId) throw new Error("Both accounts required");
    if (fromId === toId) throw new Error("From and To cannot be same");

    payload.from_account_id = Number(fromId);
    payload.to_account_id = Number(toId);
  }

  try {
    await onSubmit(payload);
  } catch (e: any) {
    setSubmitError(e.message);
  } finally {
    setSubmitting(false);
  }
};

  const canSubmit =
    !submitting &&
    !!date &&
    !!name.trim() &&
    !!amount &&
    parseFloat(amount) > 0 &&
    (type === "expense" ? !!accountId :
     type === "income"  ? !!accountId :
     /* transfer */       !!fromId && !!toId && fromId !== toId);

  return (
    <>
      <style>{`
        * { box-sizing: border-box; }
        select option { background: #0f1117; color: #d1d5db; }
        input[type="date"]::-webkit-calendar-picker-indicator { filter: invert(0.4); }
      `}</style>

      {/* type tabs */}
      <div style={tabBar}>
        {(["expense", "income", "transfer"] as TxType[]).map((t) => (
          <button
            key={t}
            onClick={() => handleTypeChange(t)}
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
       {/* category */}
        <SelectField
          label="Category"
          name="category"
          value={category}
          onChange={setCategory}
          options={CATEGORIES.map(c => ({
            id: c,
            label: c,
          }))}
        />
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
              label="Jar"
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
                label="Jar"
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
          {submitting ? "Saving…" : submitLabel}
        </button>

      </div>
    </>
  );
}

// ── styles ────────────────────────────────────────────────────────────────────
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
