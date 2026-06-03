import { useDashboard } from "../hooks/useDashboard";
import { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
/* ───────────────────────── Helpers ───────────────────────── */
const fmt = (amount: number) =>
  "₹" +
    (amount / 100).toLocaleString("en-IN", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

const API_BASE = "/api/v1";

const actionBtn = (color: string): React.CSSProperties => ({
  flex: 1,
  border: "none",
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  gap: 2,
  color,
  fontFamily: "inherit",
  fontWeight: 600,
  cursor: "pointer",
});

const formatSmartDate = (dateStr: string) => {
  const date = new Date(dateStr);

  const today = new Date();
  const yesterday = new Date();
  yesterday.setDate(today.getDate() - 1);

  if (date.toDateString() === today.toDateString()) return "Today";
  if (date.toDateString() === yesterday.toDateString()) return "Yesterday";

  return date.toLocaleDateString("en-IN", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
};
/* ───────────────────────── Main Component ───────────────────────── */

export default function RecentTransactions() {
const { data, loading, refresh } = useDashboard();

const navigate = useNavigate();
  const [activeId, setActiveId] = useState<number | null>(null);

  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [menuTarget, setMenuTarget] = useState<any>(null);

  const transactions = (data?.transactions ?? []).slice(0, 10);
  const accounts = data?.accounts ?? [];
  const jars = data?.jars ?? [];

  const accountMap = Object.fromEntries(
    accounts.map((a: any) => [a.id, a.name])
  );

  const jarMap = Object.fromEntries(
    jars.map((j: any) => [j.id, j.name])
  );

  const handleDeleteSubmit = async () => {
    if (!menuTarget) return;

    try {
      const res = await fetch(`${API_BASE}/transactions/${menuTarget.id}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) throw new Error("Delete failed");

      setIsDeleteOpen(false);
      setMenuTarget(null);

      await refresh();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="bg-[#0b0c10] border border-[#161922] rounded-xl p-4 relative">

      {/* tap outside to close swipe */}
      {activeId !== null && (
        <div
          onClick={() => setActiveId(null)}
          style={{ position: "fixed", inset: 0 }}
        />
      )}

      <div className="flex items-center justify-between mb-2">
        <h2 className="text-sm font-semibold text-zinc-100">Recent Transactions</h2>
        <div className="flex items-center gap-2">
          <button className="text-[11px] text-zinc-500 bg-transparent border border-[#1e2130] rounded-md px-2 py-0.5 font-medium hover:text-zinc-300 transition-colors">
            View all
          </button>
          <button onClick={refresh} className="bg-transparent text-zinc-500 hover:text-zinc-300 transition-colors p-1" title="Refresh" >
            <svg className={`w-4 h-4 ${loading ? "animate-spin text-zinc-400" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        </div>
      </div>

      {!loading && transactions.length > 0 && (
        <p className="text-[10px] text-zinc-700 text-right mb-2">
          ← swipe row to manage
        </p>
      )}
      {loading ? (
        <p className="text-xs text-zinc-500">Loading...</p>
      ) : transactions.length === 0 ? (
          <p className="text-xs text-zinc-500">No transactions found.</p>
        ) : (
            <div>
              {transactions.map((tx: any) => {
                const accountName =
                  accountMap[tx.from_account_id] ||
                    accountMap[tx.to_account_id] ||
                    "Unknown";

                const jarName = tx.jar_id ? jarMap[tx.jar_id] : null;

                return (
                  <TransactionRow
                    key={tx.id}
                    tx={tx}
                    accountName={accountName}
                    jarName={jarName}
                    isOpen={activeId === tx.id}
                    setActiveId={setActiveId}
                    setIsDeleteOpen={setIsDeleteOpen}
                    setMenuTarget={setMenuTarget}
                    onEdit={(tx: any) =>
                      navigate(`/transactions/${tx.id}/edit`)
                    }
                  />
                );
              })}
            </div>
          )}

      {/* ── DELETE MODAL ────────────────────────────────────── */}
      {isDeleteOpen && (
        <>
          <div
            className="fixed inset-0 bg-black/60 z-50"
            onClick={() => setIsDeleteOpen(false)}
          />

          <div className="fixed bottom-0 left-0 right-0 bg-[#0f1117] border-t border-[#1e2130] rounded-t-2xl p-6 pb-10 z-50">
            <div className="w-9 h-1 bg-[#2a2d3a] rounded-full mx-auto mb-5" />

            <div className="flex flex-col items-center justify-center gap-1.5 mb-3 text-red-400">
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2.5">
                <path strokeLinecap="round" strokeLinejoin="round"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
              <p className="text-sm font-bold uppercase tracking-wider">
                Delete Transaction
              </p>
            </div>

            <p className="text-xs text-zinc-400 text-center mb-6">
              Are you sure you want to delete{" "}
              <span className="text-zinc-200 font-medium">
                "{menuTarget?.name}"
              </span>
              ?
            </p>

            <div className="flex gap-3">
              <button
                onClick={() => setIsDeleteOpen(false)}
                className="flex-1 py-3 bg-[#1a1d27] text-zinc-400 font-semibold text-sm rounded-lg"
              >
                Cancel
              </button>

              <button
                onClick={handleDeleteSubmit}
                className="flex-1 py-3 bg-red-600 text-white font-semibold text-sm rounded-lg"
              >
                Delete
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

/* ───────────────────────── Swipe Row ───────────────────────── */

function TransactionRow({
  tx,
  accountName,
  jarName,
  isOpen,
  setActiveId,
  setIsDeleteOpen,
  setMenuTarget,
  onEdit,
}: any) {
  const touchStartX = useRef(0);
  const touchStartY = useRef(0);

  const ACTION_WIDTH = 140;

  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX;
    touchStartY.current = e.touches[0].clientY;
  };

  const handleTouchEnd = (e: React.TouchEvent) => {
    const dx = e.changedTouches[0].clientX - touchStartX.current;
    const dy = e.changedTouches[0].clientY - touchStartY.current;

    if (Math.abs(dy) > Math.abs(dx)) return;

    if (dx < -40) setActiveId(tx.id);
      else if (dx > 40) setActiveId(null);
  };

  return (
    <div
      style={{ position: "relative", overflow: "hidden" }}
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
    >
      {/* ACTIONS */}
      <div
        style={{
          position: "absolute",
          right: 0,
          top: 0,
          bottom: 0,
          width: ACTION_WIDTH,
          display: "flex",
        }}
      >
        {/* EDIT */}
        <button
          onClick={() => {
            setActiveId(null);
            onEdit(tx);
          }}
          style={actionBtn("#60a5fa")}
        >
          <svg
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zm17.71-10.04a1.003 1.003 0 0 0 0-1.42l-2.5-2.5a1.003 1.003 0 0 0-1.42 0l-1.96 1.96 3.75 3.75 2.13-1.79z"/>
          </svg>
          <span style={{ fontSize: 11, marginTop: 3 }}>Edit</span>
        </button>

        {/* DELETE */}
        <button
          onClick={() => {
            setActiveId(null);
            setMenuTarget(tx);
            setIsDeleteOpen(true);
          }}
          style={actionBtn("#E24B4A")}
        >
          <svg
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <polyline points="3 6 5 6 21 6" />
            <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
            <path d="M10 11v6" />
            <path d="M14 11v6" />
            <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
          </svg>
          <span style={{ fontSize: 11, marginTop: 3 }}>Delete</span>
        </button>
      </div>

      {/* ROW */}
      <div
        onClick={() => isOpen && setActiveId(null)}
        style={{
          transform: isOpen ? `translateX(-${ACTION_WIDTH}px)` : "translateX(0)",
          transition: "transform 0.22s ease",
          display: "flex",
          justifyContent: "space-between",
          padding: "12px 0",
          borderBottom: "1px solid #161922",
          background: "#0b0c10",
        }}
      >
        <div>
          <p style={{ margin: 0, color: "#e8eaf0", fontSize: 14, fontWeight: 500, }} >
            {tx.name}
          </p>
          <p style={{ margin: "4px 0 0", color: "#6b7280", fontSize: 11, }} >
            {(tx.category || "general")} •{" "}
            {formatSmartDate(tx.date)}
          </p>
        </div>

        <div
          style={{
            textAlign: "right",
          }}
        >
          <p
            style={{
              margin: 0,
              fontSize: 14,
              fontWeight: 600,
              color:
              tx.type === "income"
                ? "#1D9E75"
                : tx.type === "transfer"
                  ? "#60a5fa"
                  : "#E24B4A",
            }}
          >
            {tx.type === "income"
              ? "+"
              : tx.type === "transfer"
                ? "↔ "
                : "-"}
            {fmt(tx.amount)}
          </p>

          <p
            style={{
              margin: "4px 0 0",
              color: "#6b7280",
              fontSize: 11,
            }}
          >
            {accountName}
            {jarName ? ` • ${jarName}` : ""}
          </p>
        </div>
      </div>
    </div>
  );
}
