import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { TransactionRow } from "../components/TransactionRow";

import { useDashboard } from "../hooks/useDashboard";
const API_BASE = import.meta.env.VITE_API_URL; 

const fmt = (amount: number) =>
  "₹" +
    (amount / 100).toLocaleString("en-IN", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

const formatSmartDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString("en-IN", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
};

export default function AllTransactionsPage() {
  const navigate = useNavigate();

  const { data, loading, refresh } = useDashboard();

  const [allTransactions, setAllTransactions] = useState<any[]>([]);
  const accounts = data?.accounts ?? [];
  const jars = data?.jars ?? [];

  const accountMap = Object.fromEntries(
    accounts.map((a: any) => [a.id, a.name])
  );

  const jarMap = Object.fromEntries(
    jars.map((j: any) => [j.id, j.name])
  );
  useEffect(() => {
    if (data?.transactions) {
      setAllTransactions(data.transactions);
    }
  }, [data]);
  const [transactions, setTransactions] = useState<any[]>([]);

  const [activeId, setActiveId] = useState<number | null>(null);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [menuTarget, setMenuTarget] = useState<any>(null);
  const [showFilters, setShowFilters] = useState(false);

  const [type, setType] = useState<"all" | "expense" | "income" | "transfer">(
    "all"
  );
  const [search, setSearch] = useState("");

  const [filters, setFilters] = useState({
    accountId: "",
    jarId: "",
    minAmount: "",
    maxAmount: "",
    fromDate: "",
    toDate: "",
    category: "",
  });

  // ---------------- FILTER ENGINE ----------------
  const applyFilters = () => {
    let data = [...allTransactions];

    // TYPE
    if (type !== "all") {
      data = data.filter((t) => t.type === type);
    }

    // SEARCH
    if (search.trim()) {
      const s = search.toLowerCase();
      data = data.filter((t) => t.name?.toLowerCase().includes(s));
    }

    // CATEGORY
    if (filters.category) {
      data = data.filter(
        (t) => t.category?.toLowerCase() === filters.category.toLowerCase()
      );
    }

    // ACCOUNT
    if (filters.accountId) {
      data = data.filter(
        (t) =>
          String(t.from_account_id) === filters.accountId ||
            String(t.to_account_id) === filters.accountId
      );
    }

    // JAR
    if (filters.jarId) {
      data = data.filter((t) => String(t.jar_id) === filters.jarId);
    }

    // AMOUNT RANGE (PAISE)
    if (filters.minAmount) {
      data = data.filter(
        (t) => t.amount >= Number(filters.minAmount) * 100
      );
    }

    if (filters.maxAmount) {
      data = data.filter(
        (t) => t.amount <= Number(filters.maxAmount) * 100
      );
    }

    // DATE RANGE
    if (filters.fromDate) {
      const from = new Date(filters.fromDate);
      data = data.filter((t) => new Date(t.date) >= from);
    }

    if (filters.toDate) {
      const to = new Date(filters.toDate);
      data = data.filter((t) => new Date(t.date) <= to);
    }

    setTransactions(data);
  };

  useEffect(() => {
    applyFilters();
  }, [allTransactions, type, search, filters]);

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
    <div className="min-h-screen bg-[#0b0c10] text-zinc-200">
      {/* HEADER */}
      <div className="border-b border-[#161922] px-4 py-4 flex items-center gap-3">
        <button onClick={() => navigate(-1)} className="text-xs text-zinc-500">
          ← Back
        </button>
        <h1 className="text-sm font-semibold flex-1">All Transactions</h1>
      </div>

      <div className="p-4 space-y-4">
        {/* TYPE FILTER */}
        <div className="flex gap-2 pb-1">
          {(["all", "expense", "income", "transfer"] as const).map((t) => (
            <button
              key={t}
              onClick={() => {
                setType(t);
              }}
              className={`px-4 py-1.5 rounded-full text-xs capitalize ${
type === t
? "bg-zinc-200 text-black font-semibold"
: "bg-[#0f1117] border border-[#1e2130] text-zinc-400"
}`}
            >
              {t}
            </button>
          ))}
        </div>

        {/* SEARCH */}
        <div className="flex gap-2">
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search name..."
            className="flex-1 bg-[#0f1117] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
          />
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="px-4 rounded-lg text-xs border bg-[#0f1117] border-[#1e2130]"
          >
            Filters
          </button>
        </div>

        {/* ADVANCED FILTERS */}
        {showFilters && (
          <div className="bg-[#0f1117] border border-[#1e2130] rounded-xl p-4 space-y-4">

            {/* GRID */}
            <div className="grid grid-cols-2 gap-4">

              {/* ACCOUNT */}
              <div className="space-y-1 col-span-2">
                <label className="text-[10px] text-zinc-500 uppercase">
                  Account
                </label>

                <select
                  value={filters.accountId}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, accountId: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                >
              <option value="">All accounts</option>
              {accounts.map((account: any) => (
                <option key={account.id} value={account.id}>
                  {account.name}
                </option>
              ))}
                </select>
              </div>

              {/* JARS */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase">
                  Jars
                </label>
                <select
                  value={filters.jarId}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, jarId: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                >
<option value="">All jars</option>
{jars.map((jar: any) => (
  <option key={jar.id} value={jar.id}>
    {jar.name}
  </option>
))}
                </select>
              </div>
              {/* CATEGORY */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase">
                  Category
                </label>

                <select
                  value={filters.category}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, category: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                >
                  <option value="">All categories</option>
                  <option value="salary">Salary</option>
                  <option value="food">Food</option>
                  <option value="entertainment">Entertainment</option>
                  <option value="shopping">Shopping</option>
                  <option value="transport">Transport</option>
                  <option value="bills">Bills</option>
                  <option value="health">Health</option>
                  <option value="others">Others</option>
                </select>
              </div>
              {/* MIN AMOUNT */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase">
                  Min Amount (₹)
                </label>
                <input
                  type="number"
                  value={filters.minAmount}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, minAmount: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                />
              </div>

              {/* MAX AMOUNT */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase">
                  Max Amount (₹)
                </label>
                <input
                  type="number"
                  value={filters.maxAmount}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, maxAmount: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                />
              </div>

              {/* FROM DATE */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase">
                  From Date
                </label>
                <input
                  type="date"
                  value={filters.fromDate}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, fromDate: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                />
              </div>

              {/* TO DATE */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase">
                  To Date
                </label>
                <input
                  type="date"
                  value={filters.toDate}
                  onChange={(e) =>
                    setFilters((p) => ({ ...p, toDate: e.target.value }))
                  }
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs"
                />
              </div>
            </div>

            {/* FOOTER ACTION */}
            <div className="flex justify-end">
              <button
                onClick={() =>
                  setFilters({
                    accountId: "",
                    jarId: "",
                    minAmount: "",
                    maxAmount: "",
                    fromDate: "",
                    toDate: "",
                    category: "",
                  })
                }
                className="text-[11px] text-zinc-500 hover:text-zinc-300 underline"
              >
                Clear all filters
              </button>
            </div>
          </div>
        )}

        {/* LIST */}
        <div className="bg-[#0f1117] border border-[#161922] rounded-xl overflow-hidden">
          {loading ? (
            <div className="p-4 text-xs text-zinc-500">Loading...</div>
          ) : transactions.length === 0 ? (
              <div className="p-4 text-xs text-zinc-500">
                No transactions found
              </div>
            ) : (
                transactions.map((tx) => {
                  const accountName =
                    accountMap[tx.from_account_id] ||
                      accountMap[tx.to_account_id] ||
                      "Unknown";

                  const jarName = tx.jar_id
                    ? jarMap[tx.jar_id]
                    : null;

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
                      onEdit={(t: any) =>
                        navigate(`/transactions/${t.id}/edit`)
                      }
                      fmt={fmt}
                      formatSmartDate={formatSmartDate}
                    />
                  );
                })
              )}

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
      </div>
    </div>
  );
}
