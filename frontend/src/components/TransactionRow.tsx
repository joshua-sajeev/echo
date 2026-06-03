import React, { useRef } from 'react';
const toTitleCase = (str: string = "") =>
  str
    .toLowerCase()
    .split(" ")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
// Reusable styles
const actionBtn = (color: string): React.CSSProperties => ({
  flex: 1, border: "none", display: "flex", flexDirection: "column",
  alignItems: "center", justifyContent: "center", gap: 2, color,
  fontFamily: "inherit", fontWeight: 600, cursor: "pointer",
});

export function TransactionRow({
  tx,
  accountName,
  jarName,
  isOpen,
  setActiveId,
  setIsDeleteOpen,
  setMenuTarget,
  onEdit,
  fmt,
  formatSmartDate
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
    <div style={{ position: "relative", overflow: "hidden" }} onTouchStart={handleTouchStart} onTouchEnd={handleTouchEnd}>
      {/* ACTIONS */}
      <div style={{ position: "absolute", right: 0, top: 0, bottom: 0, width: ACTION_WIDTH, display: "flex" }}>
        <button onClick={() => { setActiveId(null); onEdit(tx); }} style={actionBtn("#60a5fa")}>
          <span>Edit</span>
        </button>
        <button onClick={() => { setActiveId(null); setMenuTarget(tx); setIsDeleteOpen(true); }} style={actionBtn("#E24B4A")}>
          <span>Delete</span>
        </button>
      </div>

      {/* ROW CONTENT */}
      <div onClick={() => isOpen && setActiveId(null)} style={{ transform: isOpen ? `translateX(-${ACTION_WIDTH}px)` : "translateX(0)", transition: "transform 0.22s ease", display: "flex", justifyContent: "space-between", padding: "12px 0", borderBottom: "1px solid #161922", background: "#0b0c10" }}>
        <div>
          <p style={{ margin: 0, color: "#e8eaf0", fontSize: 14, fontWeight: 500 }}>{tx.name}</p>
          <p style={{ margin: "4px 0 0", color: "#6b7280", fontSize: 11 }}>{toTitleCase(tx.category || "general")} • {formatSmartDate(tx.date)}</p>
        </div>
        <div style={{ textAlign: "right" }}>
          <p style={{ margin: 0, fontSize: 14, fontWeight: 600, color: tx.type === "income" ? "#1D9E75" : tx.type === "transfer" ? "#60a5fa" : "#E24B4A" }}>
            {tx.type === "income" ? "+" : tx.type === "transfer" ? "↔ " : "-"} {fmt(tx.amount)}
          </p>
          <p style={{ margin: "4px 0 0", color: "#6b7280", fontSize: 11 }}>{accountName} {jarName ? ` • ${jarName}` : ""}</p>
        </div>
      </div>
    </div>
  );
}
