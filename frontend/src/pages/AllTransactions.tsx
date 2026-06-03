import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { TransactionRow } from "../components/TransactionRow";

const API_BASE = import.meta.env.VITE_API_URL; 

/* ── Helpers (Shared) ── */
const fmt = (amount: number) =>
  "₹" + (amount / 100).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

const formatSmartDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString("en-IN", { day: "numeric", month: "short", year: "numeric" });
};

const actionBtn = (color: string): React.CSSProperties => ({
  flex: 1, border: "none", display: "flex", flexDirection: "column",
  alignItems: "center", justifyContent: "center", gap: 2, color,
  fontFamily: "inherit", fontWeight: 600, cursor: "pointer",
});

export default function AllTransactionsPage() {
  const navigate = useNavigate();
  const [transactions, setTransactions] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  
  // Swipe State
  const [activeId, setActiveId] = useState<number | null>(null);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [menuTarget, setMenuTarget] = useState<any>(null);

  // Filters
  const [type, setType] = useState("all");
  const [search, setSearch] = useState("");
  const [showFilters, setShowFilters] = useState(false);
  const [page, setPage] = useState(1);

  const fetchData = async () => {
    setLoading(true);
    const params = new URLSearchParams({ page: String(page), limit: "50" });
    if (type !== "all") params.set("type", type);
    if (search) params.set("search", search);

    const res = await fetch(`${API_BASE}/transactions?${params}`, { credentials: "include" });
    const data = await res.json();
    setTransactions(data);
    setLoading(false);
  };

  useEffect(() => { fetchData(); }, [type, search, page]);

  return (
    <div className="min-h-screen bg-[#0b0c10] text-zinc-200">
      {/* NAV */}
      <div className="border-b border-[#161922] px-4 py-4 flex items-center gap-3">
        <button onClick={() => navigate(-1)} className="text-xs text-zinc-500">← Back</button>
        <h1 className="text-sm font-semibold flex-1">All Transactions</h1>
      </div>

      <div className="p-4 space-y-4">
        {/* Filters UI (same as your provided structure) */}
        <div className="flex gap-2 items-center">
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search transactions..."
            className="flex-1 bg-[#0f1117] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
          />
        </div>

        {/* LIST */}
        <div className="bg-[#0f1117] border border-[#161922] rounded-xl overflow-hidden">
          {loading ? (
            <div className="p-4 text-xs text-zinc-500">Loading...</div>
          ) : transactions.length === 0 ? (
            <div className="p-4 text-xs text-zinc-500">No transactions found</div>
          ) : (
            transactions.map((tx) => (
              <TransactionRow
                key={tx.id}
                tx={tx}
                isOpen={activeId === tx.id}
                setActiveId={setActiveId}
                setIsDeleteOpen={setIsDeleteOpen}
                setMenuTarget={setMenuTarget}
                onEdit={(t: any) => navigate(`/transactions/${t.id}/edit`)}
                fmt={fmt}
                formatSmartDate={formatSmartDate} 
              />
            ))
          )}
        </div>
      </div>

      {/* Delete Modal (Omitted for brevity, use same as RecentTransactions component) */}
    </div>
  );
}

/* Include your TransactionRow component here (it is already perfectly themed) */
