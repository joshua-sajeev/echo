import React, { useState, useEffect, useRef } from "react";

// ── Backend Type Blueprints ──────────────────────────────────────────────────
interface Transaction {
  id: number;
  type: "income" | "expense" | "transfer";
  amount: number; 
  name: string;
  date: string;
  from_account_id: number | null;
  to_account_id: number | null;
  category: string | null;
  jar_id: number | null;
  is_master_income: boolean;
  created_at: string;
}

interface Account {
  id: number;
  name: string;
  balance: number;
}

interface Jar {
  id: number;
  name: string;
}

// ── Swipeable Transaction Row Component ──────────────────────────────────────
function TransactionRow({
  tx,
  accountName,
  jarLabel,
  meta,
  formatSmartDate,
  fmt,
  onEdit,
  onDelete,
  activeId,
  setActiveId,
}: {
  tx: Transaction;
  accountName: string;
  jarLabel: string;
  meta: { textClass: string; prefix: string };
  formatSmartDate: (d: string) => string;
  fmt: (n: number) => string;
  onEdit: (tx: Transaction) => void;
  onDelete: (tx: Transaction) => void;
  activeId: number | null;
  setActiveId: (id: number | null) => void;
}) {
  const isOpen = activeId === tx.id;
  const touchStartX = useRef<number>(0);
  const touchStartY = useRef<number>(0);
  const didSwipe = useRef(false);

  const ACTION_W = 140; // Total layout width allocation for absolute action buttons

  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX;
    touchStartY.current = e.touches[0].clientY;
    didSwipe.current = false;
  };

  const handleTouchEnd = (e: React.TouchEvent) => {
    const dx = e.changedTouches[0].clientX - touchStartX.current;
    const dy = e.changedTouches[0].clientY - touchStartY.current;
    
    // Lock track logic to bypass vertical jitter scrolling
    if (Math.abs(dx) < Math.abs(dy) * 1.5) return; 

    if (dx < -30) {
      setActiveId(tx.id);
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
      style={{ position: "relative", overflow: "hidden", borderRadius: "6px" }}
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
    >
      {/* Background Action Row — Clear backdrop canvas, styling affects text and icons only */}
      <div
        style={{
          position: "absolute",
          right: 0,
          top: 0,
          bottom: 0,
          display: "flex",
          alignItems: "stretch",
          width: `${ACTION_W}px`,
          zIndex: 0,
        }}
      >
        {/* Edit Action */}
        <button
          onClick={() => {
            setActiveId(null);
            onEdit(tx);
          }}
          style={{
            flex: 1,
            border: "none",
            background: "transparent",
            color: "#60a5fa",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer",
            gap: "2px",
          }}
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2.5">
            <path strokeLinecap="round" strokeLinejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
          <span style={{ fontSize: "10px", fontWeight: 600, letterSpacing: "0.02em" }}>Edit</span>
        </button>

        {/* Delete Action */}
        <button
          onClick={() => {
            setActiveId(null);
            onDelete(tx);
          }}
          style={{
            flex: 1,
            border: "none",
            background: "transparent",
            color: "#f87171",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer",
            gap: "2px",
          }}
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2.5">
            <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          <span style={{ fontSize: "10px", fontWeight: 600, letterSpacing: "0.02em" }}>Delete</span>
        </button>
      </div>

      {/* Main Sliding Content Layer */}
      <div
        onClick={handleRowClick}
        className="flex items-center justify-between py-3 px-2 bg-[#0b0c10] border-b border-[#161922] transition-transform duration-200"
        style={{
          position: "relative",
          zIndex: 1,
          transform: isOpen ? `translateX(-${ACTION_W}px)` : "translateX(0px)",
          transitionTimingFunction: "cubic-bezier(0.25, 1, 0.5, 1)",
        }}
      >
        <div className="min-w-0 flex items-center gap-2">
          <div className="min-w-0">
            <p className="text-sm font-medium text-zinc-100 truncate">{tx.name}</p>
            <p className="text-xs text-zinc-500 mt-0.5 truncate capitalize">
              {tx.category ? tx.category : "General"} · {formatSmartDate(tx.date)}
            </p>
          </div>
          
          {!isOpen && (
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="#2a2d3a" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" className="opacity-40">
              <polyline points="15 18 9 12 15 6" />
            </svg>
          )}
        </div>
        
        <div className="text-right shrink-0 ml-4">
          <p className={`text-sm font-semibold ${meta.textClass}`}>
            {meta.prefix}{fmt(tx.amount)}
          </p>
          <p className="text-xs text-zinc-600 mt-0.5 truncate uppercase font-medium tracking-wide">
            {accountName}{jarLabel}
          </p>
        </div>
      </div>
    </div>
  );
}

// ── Main Component ────────────────────────────────────────────────────────────
export default function RecentTransactions() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [jars, setJars] = useState<Jar[]>([]);
  
  const [loading, setLoading] = useState(false);
  const [loadingAccounts, setLoadingAccounts] = useState(true);
  const [loadingJars, setLoadingJars] = useState(true);
  
  // Track open interactive swipe rows
  const [activeRowId, setActiveRowId] = useState<number | null>(null);
  const [menuTarget, setMenuTarget] = useState<Transaction | null>(null);
  
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editForm, setEditForm] = useState({ name: "", amount: 0, date: "" });
  const [editError, setEditError] = useState("");
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

const API_BASE = "/api/v1";

  const fetchTransactions = async () => {
    setLoading(true);
    setActiveRowId(null);
    try {
      const response = await fetch(`${API_BASE}/transactions/`, {
        method: "GET",
        credentials: "include", 
        headers: { "Content-Type": "application/json" }
      });
      
      if (!response.ok) throw new Error(`Status error: ${response.status}`);
      const data = await response.json();
      
      if (Array.isArray(data)) {
        setTransactions(data);
      } else if (data && typeof data === 'object' && Array.isArray((data as any).data)) {
        setTransactions((data as any).data);
      } else {
        setTransactions([]);
      }
    } catch (err) {
      console.error("Fetch Error:", err);
    } finally {
      setLoading(false);
    }
  }; 

  useEffect(() => {
    fetchTransactions();

    fetch(`${API_BASE}/accounts/balances`, { credentials: "include", headers: { Accept: "application/json" } })
      .then((r) => r.json())
      .then((d: Account[]) => setAccounts(d))
      .catch((err) => console.error("Accounts Fetch Error:", err))
      .finally(() => setLoadingAccounts(false));

    fetch(`${API_BASE}/jars`, { credentials: "include", headers: { Accept: "application/json" } })
      .then((r) => r.json())
      .then((d: Jar[]) => setJars(d))
      .catch((err) => console.error("Jars Fetch Error:", err))
      .finally(() => setLoadingJars(false));
  }, []);

  const handleEditTrigger = (tx: Transaction) => {
    setMenuTarget(tx);
    setEditForm({
      name: tx.name,
      amount: Math.abs(tx.amount / 100),
      date: new Date(tx.date).toISOString().split('T')[0]
    });
    setEditError("");
    setIsEditOpen(true);
  };

  const handleDeleteTrigger = (tx: Transaction) => {
    setMenuTarget(tx);
    setIsDeleteOpen(true);
  };

  const handleSaveEdit = async () => {
    if (!menuTarget || !editForm.name || !editForm.amount) {
      setEditError("Name and amount are required.");
      return;
    }
    try {
      const updatedPayload = {
        ...menuTarget,
        name: editForm.name,
        amount: editForm.amount * 100, // Safe integer conversion back to backend Paisa structure
        date: new Date(editForm.date).toISOString()
      };

      const response = await fetch(`${API_BASE}/transactions/${menuTarget.id}`, {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updatedPayload)
      });

      if (!response.ok) throw new Error("Update request failed");
      
      setTransactions(prev => prev.map(t => t.id === menuTarget.id ? updatedPayload : t));
      setIsEditOpen(false);
    } catch (err) {
      setEditError("Could not update record on target node.");
    }
  };

  const handleDeleteSubmit = async () => {
    if (!menuTarget) return;
    try {
      const response = await fetch(`${API_BASE}/transactions/${menuTarget.id}`, {
        method: "DELETE",
        credentials: "include"
      });

      if (!response.ok) throw new Error("Delete action rejected");

      setTransactions(prev => prev.filter(t => t.id !== menuTarget.id));
      setIsDeleteOpen(false);
    } catch (err) {
      console.error("Delete Error:", err);
    }
  };

  const getAccountName = (id: number | null) => {
    if (id === null) return "";
    const matched = accounts.find((a) => a.id === id);
    return matched ? matched.name : `Acc ${id}`;
  };

  const getJarName = (id: number | null) => {
    if (id === null) return "";
    const matched = jars.find((j) => j.id === id);
    return matched ? matched.name : `Jar ${id}`;
  };

  const formatSmartDate = (dateStr: string) => {
    try {
      const txDate = new Date(dateStr);
      const today = new Date();
      const yesterday = new Date();
      yesterday.setDate(today.getDate() - 1);

      if (txDate.toDateString() === today.toDateString()) {
        return "Today";
      } else if (txDate.toDateString() === yesterday.toDateString()) {
        return "Yesterday";
      } else {
        return txDate.toLocaleDateString("en-GB", {
          day: "2-digit", month: "short", year: "numeric",
        });
      }
    } catch {
      return dateStr;
    }
  };

  const fmt = (n: number) =>
    "₹" + Math.abs(n / 100).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

  const getTransactionMeta = (type: "income" | "expense" | "transfer") => {
    switch (type) {
      case "income": return { textClass: "text-emerald-400", prefix: "+ " };
      case "expense": return { textClass: "text-rose-400", prefix: "- " };
      case "transfer": return { textClass: "text-blue-400", prefix: "± " };
      default: return { textClass: "text-zinc-100", prefix: "" };
    }
  };

  return (
    <div className="flex flex-col font-sans select-none relative">
      
      {/* Tap-out Overlay logic to instantly close active swipe handles */}
      {activeRowId !== null && (
        <div 
          onClick={() => setActiveRowId(null)}
          className="fixed inset-0 z-10 bg-transparent"
        />
      )}

      {/* ── CARD HOLDER CONTAINER ────────────────────────────────────────── */}
      <div className="bg-[#0b0c10] border border-[#161922] rounded-xl p-4 relative z-0">
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-sm font-semibold text-zinc-100">Recent Transactions</h2>
          <div className="flex items-center gap-2">
            <button className="text-[11px] text-zinc-500 bg-transparent border border-[#1e2130] rounded-md px-2 py-0.5 font-medium hover:text-zinc-300 transition-colors">
              View all
            </button>
            <button onClick={fetchTransactions} className="bg-transparent text-zinc-500 hover:text-zinc-300 transition-colors p-1" title="Refresh">
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

        {/* Ledger Loop Row Entries */}
        <div className="flex flex-col min-h-[80px] divide-y divide-[#161922]">
          {loading ? (
            <div className="flex items-center justify-center py-6 text-xs text-zinc-500">Loading system records...</div>
          ) : transactions.length === 0 ? (
            <div className="flex items-center justify-center py-6 text-xs text-zinc-500">No system transactions processed.</div>
          ) : (
            transactions.map((tx) => {
              const meta = getTransactionMeta(tx.type);
              const targetAccountId = tx.type === "income" ? tx.to_account_id : tx.from_account_id;
              const accountName = getAccountName(targetAccountId);
              const jarName = getJarName(tx.jar_id);
              const jarLabel = jarName ? ` · ${jarName}` : "";

              return (
                <TransactionRow
                  key={tx.id}
                  tx={tx}
                  accountName={accountName}
                  jarLabel={jarLabel}
                  meta={meta}
                  formatSmartDate={formatSmartDate}
                  fmt={fmt}
                  onEdit={handleEditTrigger}
                  onDelete={handleDeleteTrigger}
                  activeId={activeRowId}
                  setActiveId={setActiveRowId}
                />
              );
            })
          )}
        </div>
      </div>

      {/* ── EDIT BOTTOM SHEET / MODAL ────────────────────────────────────── */}
      {isEditOpen && (
        <>
          <div className="fixed inset-0 bg-black/60 backdrop-blur-[2px] z-50" onClick={() => setIsEditOpen(false)} />
          <div className="fixed bottom-0 left-0 right-0 bg-[#0f1117] border-t border-[#1e2130] rounded-t-2xl p-6 pb-10 z-50">
            <div className="w-9 h-1 bg-[#2a2d3a] rounded-full mx-auto mb-5" />
            
            <div className="flex flex-col items-center justify-center gap-1.5 mb-5 text-[#60a5fa]">
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2.5">
                <path strokeLinecap="round" strokeLinejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              <p className="text-xs font-bold uppercase tracking-wider">Edit Transaction</p>
            </div>
            
            <div className="space-y-4">
              <input
                type="text"
                value={editForm.name}
                onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                className="w-full bg-[#161922] border border-[#2a2d3a] rounded-lg p-3 text-sm text-zinc-300 outline-none placeholder-zinc-600"
                placeholder="Transaction text name"
              />
              <input
                type="number"
                value={editForm.amount || ""}
                onChange={(e) => setEditForm({ ...editForm, amount: parseFloat(e.target.value) })}
                className="w-full bg-[#161922] border border-[#2a2d3a] rounded-lg p-3 text-sm text-zinc-300 outline-none placeholder-zinc-600"
                placeholder="Amount (₹)"
              />
              <input
                type="date"
                value={editForm.date}
                onChange={(e) => setEditForm({ ...editForm, date: e.target.value })}
                className="w-full bg-[#161922] border border-[#2a2d3a] rounded-lg p-3 text-sm text-zinc-300 outline-none text-left"
              />
            </div>

            {editError && <p className="text-red-400 text-xs mt-2 text-center">{editError}</p>}

            <div className="flex gap-3 mt-6">
              <button onClick={() => setIsEditOpen(false)} className="flex-1 py-3 bg-[#1a1d27] text-zinc-400 font-semibold text-sm rounded-lg border-none cursor-pointer">Cancel</button>
              <button onClick={handleSaveEdit} className="flex-1 py-3 bg-[#1D9E75] text-white font-semibold text-sm rounded-lg border-none cursor-pointer">Save Ledger</button>
            </div>
          </div>
        </>
      )}

      {/* ── DELETE MODAL CONFIRMATION ────────────────────────────────────── */}
      {isDeleteOpen && (
        <>
          <div className="fixed inset-0 bg-black/60 z-50" onClick={() => setIsDeleteOpen(false)} />
          <div className="fixed bottom-0 left-0 right-0 bg-[#0f1117] border-t border-[#1e2130] rounded-t-2xl p-6 pb-10 z-50">
            <div className="w-9 h-1 bg-[#2a2d3a] rounded-full mx-auto mb-5" />

            <div className="flex flex-col items-center justify-center gap-1.5 mb-3 text-red-400">
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2.5">
                <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              <p className="text-sm font-bold uppercase tracking-wider">Purge Record</p>
            </div>

            <p className="text-xs text-zinc-400 text-center max-w-xs mx-auto mb-6">
              Are you sure you want to delete <span className="text-zinc-200 font-medium">"{menuTarget?.name}"</span>? This will re-calculate downstream metrics.
            </p>

            <div className="flex gap-3">
              <button onClick={() => setIsDeleteOpen(false)} className="flex-1 py-3 bg-[#1a1d27] text-zinc-400 font-semibold text-sm rounded-lg border-none cursor-pointer">Keep Entry</button>
              <button onClick={handleDeleteSubmit} className="flex-1 py-3 bg-red-600 text-white font-semibold text-sm rounded-lg border-none cursor-pointer">Confirm Delete</button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
