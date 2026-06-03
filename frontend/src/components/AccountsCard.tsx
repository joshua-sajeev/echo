import { useEffect, useRef, useState } from "react";

// ── types ─────────────────────────────────────────────────────────────────────
export interface Account {
  id: number;
  name: string;
  is_archived: boolean;
  created_at: string;
  Balance: number;
}
// ── helpers ───────────────────────────────────────────────────────────────────
// Added division by 100 to safely convert integer Paisa values back to decimals for display
const fmt = (n: number) =>
  "₹" +
  (n / 100).toLocaleString("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

const ACCOUNT_COLORS = [
  "#60a5fa", // blue
  "#34d399", // green
  "#fbbf24", // amber
  "#f472b6", // pink
  "#a78bfa", // purple
];
const getAccountColor = (id: number, isArchived: boolean) => {
  const color = ACCOUNT_COLORS[id % ACCOUNT_COLORS.length];
  return isArchived ? "#6b7280" : color; // muted gray for archived
};
// ── API base ───────────────────────────────────────────────────────────────────
const API_BASE = import.meta.env.VITE_API_URL; 

// ── skeleton ──────────────────────────────────────────────────────────────────
function Skeleton() {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          style={{
            height: 56,
            borderRadius: 10,
            background:
              "linear-gradient(90deg,#161922 25%,#1e2233 50%,#161922 75%)",
            backgroundSize: "200% 100%",
            animation: "shimmer 1.4s infinite",
          }}
        />
      ))}
    </div>
  );
}

// ── rename modal ──────────────────────────────────────────────────────────────
function RenameModal({
  account,
  onClose,
  onDone,
}: {
  account: Account;
  onClose: () => void;
  onDone: () => void;
}) {
  const [name, setName] = useState(account.name);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
    inputRef.current?.select();
  }, []);

  const save = async () => {
    const trimmed = name.trim();
    if (!trimmed || trimmed === account.name) { onClose(); return; }
    setSaving(true);
    setError(null);
    try {
      const r = await fetch(`${API_BASE}/accounts/${account.id}/rename`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        credentials: "include",
        body: JSON.stringify({ name: trimmed }),
      });
      if (!r.ok) throw new Error(`Server returned ${r.status}`);
      onDone();
    } catch (e: any) {
      setError(e.message);
      setSaving(false);
    }
  };

  return (
    <>
      {/* backdrop */}
      <div
        onClick={onClose}
        style={{
          position: "fixed", inset: 0, background: "rgba(0,0,0,0.6)",
          zIndex: 100, backdropFilter: "blur(2px)",
        }}
      />
      {/* sheet */}
      <div
        style={{
          position: "fixed", bottom: 0, left: 0, right: 0,
          background: "#0f1117", borderTop: "0.5px solid #1e2130",
          borderRadius: "16px 16px 0 0", padding: "24px 20px 40px",
          zIndex: 101,
        }}
      >
        {/* drag handle */}
        <div style={{
          width: 36, height: 4, borderRadius: 2, background: "#2a2d3a",
          margin: "0 auto 20px",
        }} />

        <p style={{ color: "#6b7280", fontSize: 11, fontWeight: 600, letterSpacing: "0.08em", textTransform: "uppercase", margin: "0 0 12px" }}>
          Rename account
        </p>

        <input
          ref={inputRef}
          value={name}
          onChange={(e) => setName(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && save()}
          style={{
            width: "100%", boxSizing: "border-box",
            background: "#161922", border: "0.5px solid #2a2d3a",
            borderRadius: 10, padding: "12px 14px",
            color: "#d1d5db", fontSize: 16, outline: "none",
            fontFamily: "inherit",
          }}
          placeholder="Account name"
        />

        {error && (
          <p style={{ color: "#E24B4A", fontSize: 12, margin: "8px 0 0" }}>{error}</p>
        )}

        <div style={{ display: "flex", gap: 8, marginTop: 16 }}>
          <button onClick={onClose} style={sheetBtn("#1a1d27", "#9ca3af")}>
            Cancel
          </button>
          <button
            onClick={save}
            disabled={saving || !name.trim()}
            style={sheetBtn("#1D9E75", "#fff", saving || !name.trim())}
          >
            {saving ? "Saving…" : "Save"}
          </button>
        </div>
      </div>
    </>
  );
}

function CreateAccountModal({
  onClose,
  onDone,
}: {
  onClose: () => void;
  onDone: () => void;
}) {
  const [name, setName] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const create = async () => {
    const trimmed = name.trim();
    if (!trimmed) return;

    setSaving(true);
    setError(null);

    try {
      const r = await fetch(`${API_BASE}/accounts`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ name: trimmed }),
      });

      if (!r.ok) throw new Error(`Server returned ${r.status}`);

      onDone();
    } catch (e: any) {
      setError(e.message);
      setSaving(false);
    }
  };

  return (
    <>
      <div
        onClick={onClose}
        style={{
          position: "fixed",
          inset: 0,
          background: "rgba(0,0,0,0.6)",
          zIndex: 100,
        }}
      />

      <div
        style={{
          position: "fixed",
          bottom: 0,
          left: 0,
          right: 0,
          background: "#0f1117",
          borderTop: "0.5px solid #1e2130",
          borderRadius: "16px 16px 0 0",
          padding: "24px 20px 40px",
          zIndex: 101,
        }}
      >
        <p style={{ color: "#6b7280", fontSize: 11, marginBottom: 10 }}>
          Create account
        </p>

        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Account name"
          style={{
            width: "100%",
            padding: "12px",
            borderRadius: 10,
            background: "#161922",
            border: "0.5px solid #2a2d3a",
            color: "#d1d5db",
            outline: "none",
          }}
        />

        {error && (
          <p style={{ color: "#E24B4A", fontSize: 12 }}>{error}</p>
        )}

        <div style={{ display: "flex", gap: 8, marginTop: 16 }}>
          <button onClick={onClose} style={sheetBtn("#1a1d27", "#9ca3af")}>
            Cancel
          </button>

          <button
            onClick={create}
            disabled={!name.trim() || saving}
            style={sheetBtn("#1D9E75", "#fff", !name.trim() || saving)}
          >
            {saving ? "Creating..." : "Create"}
          </button>
        </div>
      </div>
    </>
  );
}
const sheetBtn = (bg: string, color: string, disabled = false): React.CSSProperties => ({
  flex: 1, padding: "13px 0", borderRadius: 10, border: "none",
  background: disabled ? "#1a1d27" : bg,
  color: disabled ? "#4b5563" : color,
  fontSize: 15, fontWeight: 600, cursor: disabled ? "not-allowed" : "pointer",
  fontFamily: "inherit",
});

// ── swipeable account row ─────────────────────────────────────────────────────
function AccountRow({
  account,
  isArchived,
  onRename,
  onArchive,
  onUnarchive,
  activeId,
  setActiveId,
}: {
  account: Account;
  isArchived: boolean;
  onRename: (a: Account) => void;
  onArchive: (a: Account) => void;
  onUnarchive: (a: Account) => void;
  activeId: number | null;
  setActiveId: (id: number | null) => void;
}) {
  const isOpen = activeId === account.id;

  // touch tracking
  const touchStartX = useRef<number>(0);
  const touchStartY = useRef<number>(0);
  const didSwipe = useRef(false);

  const ACTION_W = isArchived ? 90 : 170; // unarchive only vs rename+archive

  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX;
    touchStartY.current = e.touches[0].clientY;
    didSwipe.current = false;
  };

  const handleTouchEnd = (e: React.TouchEvent) => {
    const dx = e.changedTouches[0].clientX - touchStartX.current;
    const dy = e.changedTouches[0].clientY - touchStartY.current;
    if (Math.abs(dx) < Math.abs(dy) * 1.5) return; // mostly vertical → ignore
    if (dx < -30) {
      setActiveId(account.id);
      didSwipe.current = true;
    } else if (dx > 30) {
      setActiveId(null);
      didSwipe.current = true;
    }
  };

  const handleRowClick = () => {
    if (didSwipe.current) return;
    if (isOpen) setActiveId(null);
  };

  return (
    <div
      style={{ position: "relative", borderRadius: 10, overflow: "hidden" }}
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
    >
      {/* action buttons behind the row */}
      <div
        style={{
          position: "absolute", right: 0, top: 0, bottom: 0,
          display: "flex", alignItems: "stretch",
          width: ACTION_W,
        }}
      >
        {!isArchived && (
          <>
            {/* Rename */}
            <button
              onClick={() => { setActiveId(null); onRename(account); }}
              style={actionBtn("#60a5fa")}
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
              </svg>
              <span style={{ fontSize: 11, marginTop: 3 }}>Rename</span>
            </button>

            {/* Archive */}
            <button
              onClick={() => { setActiveId(null); onArchive(account); }}
              style={actionBtn("#E24B4A")}
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="4" width="18" height="5" rx="1" />
                <path d="M5 9v11a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V9" />
                <path d="M10 13h4" />
              </svg>
              <span style={{ fontSize: 11, marginTop: 3 }}>Archive</span>
            </button>
          </>
        )}

        {isArchived && (
          /* Unarchive */
          <button
            onClick={() => { setActiveId(null); onUnarchive(account); }}
            style={actionBtn("#1D9E75")}
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <rect x="3" y="4" width="18" height="5" rx="1" />
              <path d="M5 9v11a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V9" />
              <path d="M9 13l3-3 3 3" />
              <path d="M12 10v6" />
            </svg>
            <span style={{ fontSize: 11, marginTop: 3 }}>Restore</span>
          </button>
        )}
      </div>

      {/* main row — slides left to reveal actions */}
      <div
        onClick={handleRowClick}
        style={{
          ...row,
          position: "relative",
          transform: isOpen ? `translateX(-${ACTION_W}px)` : "translateX(0)",
          transition: "transform 0.22s cubic-bezier(0.25,1,0.5,1)",
          zIndex: 1,
          cursor: "default",
          // subtle hint chevron visible when closed
          borderRight: isOpen ? "none" : "0.5px solid #1e2130",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <p
            style={{
              color: getAccountColor(account.id, isArchived),
              fontSize: 15,
              fontWeight: 600,
              margin: 0,
            }}
          >
            {account.name}
          </p>
          {/* swipe hint — tiny chevrons, disappear once user has swiped */}
          {!isOpen && (
            <svg
              width="12"
              height="12"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#2a2d3a"
              strokeWidth="2.5"
              strokeLinecap="round"
              strokeLinejoin="round"
              style={{ flexShrink: 0 }}
            >
              <polyline points="15 18 9 12 15 6" />
              <polyline points="21 18 15 12 21 6" />
            </svg>
          )}
        </div>

        <span
          style={{
            fontFamily: "'IBM Plex Mono', monospace",
            fontSize: 15,
            color: account.Balance < 0 ? "#E24B4A" : "#1D9E75",
          }}
        >
          {fmt(account.Balance)}
        </span>
      </div>
    </div>
  );
}

const actionBtn = (color: string): React.CSSProperties => ({
  flex: 1,
  border: "none",
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  cursor: "pointer",
  fontFamily: "inherit",
  fontWeight: 600,
  gap: 2,
  color,
});

// ── main component ────────────────────────────────────────────────────────────
export default function AccountsCard() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showArchived, setShowArchived] = useState(false);
  const [creating, setCreating] = useState(false);
  // which row has its actions revealed
  const [activeRowId, setActiveRowId] = useState<number | null>(null);

  // rename modal
  const [renaming, setRenaming] = useState<Account | null>(null);

  // per-row loading state (archive/unarchive)
  const [busyId, setBusyId] = useState<number | null>(null);

  const load = () => {
    setLoading(true);
    setError(null);
    setActiveRowId(null);

    const endpoint = showArchived ? "/accounts/archived" : "/accounts/balances";

    fetch(`${API_BASE}${endpoint}`, {
      headers: { Accept: "application/json" },
      credentials: "include",
    })
      .then(async (r) => {
        const text = await r.text();
        if (!r.ok) throw new Error(`Server returned ${r.status}`);
        try { return JSON.parse(text) as Account[]; }
        catch { throw new Error(`Not JSON: ${text.slice(0, 120)}`); }
      })
      .then((data) => { 
        // FIX: Fallback to an empty array if data is explicitly null or undefined
        setAccounts(data || []); 
        setLoading(false); 
      })
      .catch((e) => { 
        setError(e.message); 
        setLoading(false); 
      });
  };

  useEffect(() => { load(); }, []);
  useEffect(() => { load(); }, [showArchived]);

  const handleArchive = async (a: Account) => {
    setBusyId(a.id);
    try {
      const r = await fetch(`${API_BASE}/accounts/${a.id}/archive`, {
        method: "PATCH", credentials: "include",
      });
      if (!r.ok) throw new Error(`Server returned ${r.status}`);
      load();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setBusyId(null);
    }
  };

  const handleUnarchive = async (a: Account) => {
    setBusyId(a.id);
    try {
      const r = await fetch(`${API_BASE}/accounts/${a.id}/unarchive`, {
        method: "PATCH", credentials: "include",
      });
      if (!r.ok) throw new Error(`Server returned ${r.status}`);
      load();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setBusyId(null);
    }
  };

  return (
    <>
      <style>{`
        @keyframes shimmer { to { background-position: -200% 0; } }
        @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
      `}</style>

      {/* tap outside to close any open row */}
      {activeRowId !== null && (
        <div
          onClick={() => setActiveRowId(null)}
          style={{ position: "fixed", inset: 0, zIndex: 0 }}
        />
      )}

      <div style={card}>
        {/* header */}
        <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: 14 }}>
          <span style={sectionLabel}>Accounts</span>

          <div style={{ display: "flex", gap: 8 }}>
            {/* refresh */}
            <button style={iconBtn} onClick={load} title="Refresh">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"
                style={loading ? { animation: "spin 0.8s linear infinite" } : {}}>
                <path d="M21 12a9 9 0 1 1-3-6.7" />
                <polyline points="21 3 21 9 15 9" />
              </svg>
            </button>

            {/* add */}
            <button style={iconBtn} title="Add account" onClick={() => setCreating(true)} >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M12 5v14" /><path d="M5 12h14" />
              </svg>
            </button>

            {/* archive toggle */}
            <button
              style={{ ...iconBtn, borderColor: showArchived ? "#3b82f6" : "#2a2d3a", color: showArchived ? "#3b82f6" : "#9ca3af" }}
              onClick={() => setShowArchived((v) => !v)}
              title="Toggle archived accounts"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="4" width="18" height="5" rx="1" />
                <path d="M5 9v11a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V9" />
                <path d="M10 13h4" />
              </svg>
            </button>
          </div>
        </div>

        {/* hint text */}
        {!loading && !error && accounts && accounts.length > 0 && (
          <p style={{ color: "#374151", fontSize: 11, margin: "0 0 10px", textAlign: "right" }}>
            ← swipe row to edit
          </p>
        )}

        {/* states */}
        {loading && <Skeleton />}

        {!loading && error && (
          <div style={{ color: "#E24B4A", fontSize: 12 }}>{error}</div>
        )}

        {!loading && !error && (!accounts || accounts.length === 0) && (
          <p style={{ color: "#4b5563", fontSize: 13, textAlign: "center", padding: "12px 0" }}>
            No {showArchived ? "archived" : "active"} accounts found
          </p>
        )}

        {!loading && !error && accounts && accounts.length > 0 && (
          <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            {accounts.map((a) =>
              busyId === a.id ? (
                // row is being mutated — show a muted placeholder
                <div key={a.id} style={{ ...row, opacity: 0.4 }}>
                  <p style={{ color: "#d1d5db", fontSize: 15, fontWeight: 500, margin: 0 }}>{a.name}</p>
                  <span style={{ color: "#6b7280", fontSize: 13 }}>…</span>
                </div>
              ) : (
                <AccountRow
                  key={a.id}
                  account={a}
                  isArchived={showArchived}
                  onRename={setRenaming}
                  onArchive={handleArchive}
                  onUnarchive={handleUnarchive}
                  activeId={activeRowId}
                  setActiveId={setActiveRowId}
                />
              )
            )}
          </div>
        )}
      </div>

      {/* rename modal */}
      {renaming && (
        <RenameModal
          account={renaming}
          onClose={() => setRenaming(null)}
          onDone={() => { setRenaming(null); load(); }}
        />
      )}
      {creating && (
        <CreateAccountModal
          onClose={() => setCreating(false)}
          onDone={() => {
            setCreating(false);
            load();
          }}
        />
      )}
    </>
  );
}

// ── styles ────────────────────────────────────────────────────────────────────
const card: React.CSSProperties = {
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 14,
  padding: 16,
  marginBottom: 12,
  position: "relative",
};

const row: React.CSSProperties = {
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
  padding: "10px 12px",
  background: "#161922",
  border: "0.5px solid #1e2130",
  borderRadius: 10,
};

const iconBtn: React.CSSProperties = {
  width: 30, height: 30,
  borderRadius: 8, border: "0.5px solid #2a2d3a",
  background: "#1a1d27", color: "#9ca3af",
  display: "flex", alignItems: "center", justifyContent: "center",
  cursor: "pointer",
};

const sectionLabel: React.CSSProperties = {
  fontSize: 11, fontWeight: 600,
  letterSpacing: "0.08em", textTransform: "uppercase",
  color: "#6b7280",
};
