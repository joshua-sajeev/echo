import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useDashboard } from "../hooks/useDashboard";

interface TransactionTemplate {
  id: string;
  templateName: string;
  type: "expense" | "income" | "transfer";
  name: string;
  amount: string;
  category: string;
  accountId: string;
  fromId: string;
  toId: string;
  jarId: string;
  isMasterIncome: boolean;
}

const CATEGORIES = [
  "Food", "Transport", "Shopping", "Donations", "Entertainment",
  "Health", "Income", "Investment", "Housing", "Transfers",
];

const DEFAULT_TEMPLATES: TransactionTemplate[] = [
  {
    id: "default-salary",
    templateName: "Salary Deposit",
    type: "income",
    name: "Monthly Salary",
    amount: "85000",
    category: "Income",
    accountId: "",
    fromId: "",
    toId: "",
    jarId: "",
    isMasterIncome: true,
  },
  {
    id: "default-rent",
    templateName: "House Rent",
    type: "expense",
    name: "Rent Payment",
    amount: "18000",
    category: "Housing",
    accountId: "",
    fromId: "",
    toId: "",
    jarId: "",
    isMasterIncome: false,
  },
  {
    id: "default-coffee",
    templateName: "Daily Coffee",
    type: "expense",
    name: "Coffee",
    amount: "150",
    category: "Food",
    accountId: "",
    fromId: "",
    toId: "",
    jarId: "",
    isMasterIncome: false,
  }
];

export default function TemplatesPage() {
  const navigate = useNavigate();
  const { data } = useDashboard();
  const accounts = data?.accounts ?? [];
  const jars = data?.jars ?? [];

  const [templates, setTemplates] = useState<TransactionTemplate[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<TransactionTemplate | null>(null);

  // Form State
  const [templateName, setTemplateName] = useState("");
  const [type, setType] = useState<"expense" | "income" | "transfer">("expense");
  const [name, setName] = useState("");
  const [amount, setAmount] = useState("");
  const [category, setCategory] = useState("");
  const [accountId, setAccountId] = useState("");
  const [fromId, setFromId] = useState("");
  const [toId, setToId] = useState("");
  const [jarId, setJarId] = useState("");
  const [isMasterIncome, setIsMasterIncome] = useState(false);

  // Load from localStorage
  useEffect(() => {
    const stored = localStorage.getItem("echo_templates");
    if (stored) {
      try {
        setTemplates(JSON.parse(stored));
      } catch (e) {
        setTemplates(DEFAULT_TEMPLATES);
      }
    } else {
      setTemplates(DEFAULT_TEMPLATES);
      localStorage.setItem("echo_templates", JSON.stringify(DEFAULT_TEMPLATES));
    }
  }, []);

  const saveTemplates = (newTemplates: TransactionTemplate[]) => {
    setTemplates(newTemplates);
    localStorage.setItem("echo_templates", JSON.stringify(newTemplates));
  };

  const handleOpenAdd = () => {
    setEditingTemplate(null);
    setTemplateName("");
    setType("expense");
    setName("");
    setAmount("");
    setCategory("");
    setAccountId(accounts[0]?.id ? String(accounts[0].id) : "");
    setFromId(accounts[0]?.id ? String(accounts[0].id) : "");
    setToId(accounts[1]?.id ? String(accounts[1].id) : "");
    setJarId("");
    setIsMasterIncome(false);
    setIsOpen(true);
  };

  const handleOpenEdit = (t: TransactionTemplate) => {
    setEditingTemplate(t);
    setTemplateName(t.templateName);
    setType(t.type);
    setName(t.name);
    setAmount(t.amount);
    setCategory(t.category);
    setAccountId(t.accountId);
    setFromId(t.fromId);
    setToId(t.toId);
    setJarId(t.jarId);
    setIsMasterIncome(t.isMasterIncome);
    setIsOpen(true);
  };

  const handleDelete = (id: string) => {
    const updated = templates.filter((t) => t.id !== id);
    saveTemplates(updated);
  };

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    if (!templateName.trim() || !name.trim()) return;

    const newTemplate: TransactionTemplate = {
      id: editingTemplate ? editingTemplate.id : Math.random().toString(36).substr(2, 9),
      templateName: templateName.trim(),
      type,
      name: name.trim(),
      amount,
      category,
      accountId: type !== "transfer" ? accountId : "",
      fromId: type === "transfer" ? fromId : "",
      toId: type === "transfer" ? toId : "",
      jarId: (type === "expense" || (type === "income" && !isMasterIncome)) ? jarId : "",
      isMasterIncome: type === "income" ? isMasterIncome : false,
    };

    if (editingTemplate) {
      saveTemplates(templates.map((t) => (t.id === editingTemplate.id ? newTemplate : t)));
    } else {
      saveTemplates([...templates, newTemplate]);
    }
    setIsOpen(false);
  };

  const handleApply = (t: TransactionTemplate) => {
    // Fill in default accounts/jars if not set in template
    const applied = {
      ...t,
      accountId: t.accountId || (accounts[0]?.id ? String(accounts[0].id) : ""),
      fromId: t.fromId || (accounts[0]?.id ? String(accounts[0].id) : ""),
      toId: t.toId || (accounts[1]?.id ? String(accounts[1].id) : ""),
    };
    navigate("/transactions/new", { state: { template: applied } });
  };

  const getAccent = (t: string) => {
    switch (t) {
      case "expense": return "#ef4444";
      case "income": return "#22c55e";
      default: return "#3b82f6";
    }
  };

  return (
    <div className="min-h-screen bg-[#0b0c10] text-zinc-200" style={{ fontFamily: "'Syne', sans-serif" }}>
      {/* HEADER */}
      <div className="border-b border-[#161922] bg-[#0b0c10] sticky top-0 z-10">
        <div style={{ maxWidth: 480, margin: "0 auto" }} className="px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button onClick={() => navigate(-1)} className="text-xs text-zinc-500 bg-transparent border-none cursor-pointer">
              ← Back
            </button>
            <h1 className="text-sm font-semibold m-0">Transaction Templates</h1>
          </div>
          <button
            onClick={handleOpenAdd}
            className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-zinc-200 text-black hover:bg-white transition-colors cursor-pointer border-none"
          >
            + Add New
          </button>
        </div>
      </div>

      <div style={{ maxWidth: 480, margin: "0 auto" }} className="p-4 space-y-4">
        {templates.length === 0 ? (
          <div className="text-center py-12 text-zinc-500 text-xs">
            No templates configured yet. Click "+ Add New" to create one.
          </div>
        ) : (
          <div className="space-y-3">
            {templates.map((t) => {
              const accentColor = getAccent(t.type);
              return (
                <div
                  key={t.id}
                  className="bg-[#0f1117] border border-[#1e2130] rounded-xl p-4 flex flex-col justify-between gap-4 transition-all hover:border-zinc-700"
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <h2 className="text-sm font-bold text-zinc-100 m-0 flex items-center gap-2">
                        {t.templateName}
                        <span
                          className="text-[9px] px-2 py-0.5 rounded-full font-semibold uppercase"
                          style={{
                            backgroundColor: `${accentColor}15`,
                            color: accentColor,
                            border: `1px solid ${accentColor}30`,
                          }}
                        >
                          {t.type}
                        </span>
                      </h2>
                      <p className="text-xs text-zinc-500 mt-1 mb-0">
                        {t.name} {t.category ? `• ${t.category}` : ""}
                      </p>
                    </div>
                    {t.amount && (
                      <span className="font-mono text-sm font-bold" style={{ color: accentColor }}>
                        {t.type === "income" ? "+" : t.type === "transfer" ? "↔ " : "-"} ₹{t.amount}
                      </span>
                    )}
                  </div>

                  <div className="flex items-center justify-between border-t border-[#161922] pt-3">
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleOpenEdit(t)}
                        className="text-[11px] text-zinc-400 hover:text-zinc-200 bg-transparent border-none cursor-pointer"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(t.id)}
                        className="text-[11px] text-[#E24B4A] hover:text-red-400 bg-transparent border-none cursor-pointer"
                      >
                        Delete
                      </button>
                    </div>
                    <button
                      onClick={() => handleApply(t)}
                      className="px-3 py-1 rounded-lg text-xs font-semibold cursor-pointer border-none text-white transition-all active:scale-95"
                      style={{ backgroundColor: accentColor }}
                    >
                      Use Template
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* FORM MODAL */}
      {isOpen && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center p-4 z-50 animate-fade-in">
          <div className="bg-[#0f1117] border border-[#1e2130] rounded-2xl w-full max-w-sm overflow-hidden flex flex-col">
            <div className="px-4 py-3 border-b border-[#161922] flex justify-between items-center">
              <h3 className="text-xs font-bold uppercase tracking-wider text-zinc-400 m-0">
                {editingTemplate ? "Edit Template" : "New Template"}
              </h3>
              <button
                onClick={() => setIsOpen(false)}
                className="text-zinc-500 hover:text-zinc-300 bg-transparent border-none cursor-pointer text-sm"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleSave} className="p-4 space-y-4 overflow-y-auto max-h-[75vh]">
              {/* Template Name */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Template Label</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Monthly Rent"
                  value={templateName}
                  onChange={(e) => setTemplateName(e.target.value)}
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                />
              </div>

              {/* Type Selection */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Transaction Type</label>
                <div className="flex gap-2 bg-[#0b0c10] border border-[#1e2130] p-1 rounded-lg">
                  {(["expense", "income", "transfer"] as const).map((t) => (
                    <button
                      key={t}
                      type="button"
                      onClick={() => setType(t)}
                      className={`flex-1 py-1.5 rounded-md text-xs font-semibold capitalize border-none cursor-pointer transition-all ${
                        type === t ? "bg-zinc-200 text-black" : "bg-transparent text-zinc-500"
                      }`}
                    >
                      {t}
                    </button>
                  ))}
                </div>
              </div>

              {/* Name / Desc */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Name / Desc *</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Rent Payment"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                />
              </div>

              {/* Category */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Category</label>
                <select
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                >
                  <option value="">Select category...</option>
                  {CATEGORIES.map((c) => (
                    <option key={c} value={c}>{c}</option>
                  ))}
                </select>
              </div>

              {/* Amount */}
              <div className="space-y-1">
                <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Default Amount (₹)</label>
                <input
                  type="number"
                  placeholder="0.00"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500 font-mono"
                />
              </div>

              {/* Dynamic Fields */}
              {type !== "transfer" ? (
                <>
                  <div className="space-y-1">
                    <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Account</label>
                    <select
                      value={accountId}
                      onChange={(e) => setAccountId(e.target.value)}
                      className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                    >
                      <option value="">Choose account...</option>
                      {accounts.map((a: any) => (
                        <option key={a.id} value={a.id}>{a.name}</option>
                      ))}
                    </select>
                  </div>

                  {type === "income" && (
                    <div
                      onClick={() => setIsMasterIncome(!isMasterIncome)}
                      className="flex items-center justify-between p-2.5 bg-[#0b0c10] border border-[#1e2130] rounded-lg cursor-pointer select-none"
                    >
                      <div>
                        <span className="text-xs font-semibold text-zinc-300">Master income</span>
                      </div>
                      <div
                        className="w-8 h-5 rounded-full relative transition-all"
                        style={{ backgroundColor: isMasterIncome ? "#1D9E75" : "#1e2130" }}
                      >
                        <div
                          className="w-3.5 h-3.5 rounded-full bg-white absolute top-0.5 transition-all"
                          style={{ left: isMasterIncome ? "14px" : "2px" }}
                        />
                      </div>
                    </div>
                  )}

                  {!(type === "income" && isMasterIncome) && (
                    <div className="space-y-1">
                      <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">Jar</label>
                      <select
                        value={jarId}
                        onChange={(e) => setJarId(e.target.value)}
                        className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                      >
                        <option value="">No jar</option>
                        {jars.map((j: any) => (
                          <option key={j.id} value={j.id}>{j.name}</option>
                        ))}
                      </select>
                    </div>
                  )}
                </>
              ) : (
                <>
                  <div className="space-y-1">
                    <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">From Account</label>
                    <select
                      value={fromId}
                      onChange={(e) => setFromId(e.target.value)}
                      className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                    >
                      <option value="">Choose source...</option>
                      {accounts.map((a: any) => (
                        <option key={a.id} value={a.id}>{a.name}</option>
                      ))}
                    </select>
                  </div>
                  <div className="space-y-1">
                    <label className="text-[10px] text-zinc-500 uppercase font-bold tracking-wide">To Account</label>
                    <select
                      value={toId}
                      onChange={(e) => setToId(e.target.value)}
                      className="w-full bg-[#0b0c10] border border-[#1e2130] rounded-lg px-3 py-2 text-xs text-zinc-200 outline-none focus:border-zinc-500"
                    >
                      <option value="">Choose destination...</option>
                      {accounts.filter((a: any) => String(a.id) !== fromId).map((a: any) => (
                        <option key={a.id} value={a.id}>{a.name}</option>
                      ))}
                    </select>
                  </div>
                </>
              )}

              <button
                type="submit"
                className="w-full py-2.5 rounded-lg text-xs font-bold text-black bg-zinc-200 hover:bg-white transition-colors cursor-pointer border-none"
              >
                {editingTemplate ? "Save Changes" : "Create Template"}
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
