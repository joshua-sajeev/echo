import { useDashboard } from "../hooks/useDashboard";
import {  useState } from "react";
import { useNavigate } from "react-router-dom";

import DeleteModal from "../components/DeleteModal";
import { TransactionRow } from "../components/TransactionRow";
/* ───────────────────────── Helpers ───────────────────────── */
const fmt = (amount: number) =>
  "₹" +
    (amount / 100).toLocaleString("en-IN", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

const API_BASE = import.meta.env.VITE_API_URL; 


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
          <button 

            onClick={() => navigate("/transactions")}
            className="text-[11px] text-zinc-500 bg-transparent border border-[#1e2130] rounded-md px-2 py-0.5 font-medium hover:text-zinc-300 transition-colors">
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
                    fmt={fmt}
                    formatSmartDate={formatSmartDate}
                  />
                );
              })}
            </div>
          )}

      <DeleteModal
        open={isDeleteOpen}
        title="Delete Transaction"
        itemName={menuTarget?.name}
        onClose={() => setIsDeleteOpen(false)}
        onConfirm={handleDeleteSubmit}
      />
    </div>
  );
}

/* ───────────────────────── Swipe Row ───────────────────────── */

