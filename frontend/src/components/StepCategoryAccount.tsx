import { useMemo } from "react";
import type { TxType, Account, Jar } from "./TransactionForm";
const CATEGORIES = [
  "Food", "Transport", "Shopping", "Donations", "Entertainment",
  "Health", "Income", "Investment", "Housing", "Transfers",
];
interface StepCategoryAccountProps {
  type: TxType;
  category: string;
  setCategory: (c: string) => void;
  accountId: string;
  setAccountId: (id: string) => void;
  fromAccountId?: string;
  setFromAccountId?: (id: string) => void;
  toAccountId?: string;
  setToAccountId?: (id: string) => void;
  jarId: string;
  setJarId: (id: string) => void;
  isMasterIncome: boolean; // Add this
  setIsMasterIncome: (val: boolean) => void; // Add this
  accounts: Account[];
  jars: Jar[];
  onNext: () => void;
}
const SelectionGroup = ({
  label, items, selectedId, onSelect, accent,
}: {
  label: string; items: { id: string | number; name: string }[]; 
  selectedId: string; onSelect: (id: string) => void; accent: string;
}) => (
  <section style={{ marginBottom: 24 }}>
    <h3 style={heading}>{label}</h3>
    <div style={chipWrap}>
      {items.map((item) => {
        const isSelected = selectedId === String(item.id);
        return (
          <button
            key={item.id}
            type="button"
            onClick={() => onSelect(String(item.id))}
            style={{
              ...chip,
              backgroundColor: isSelected ? accent : "#18181b",
              border: `1px solid ${isSelected ? accent : "#2d2d31"}`,
              color: isSelected ? "#ffffff" : "#a1a1aa",
            }}
          >
            {item.name}
          </button>
        );
      })}
    </div>
  </section>
);

export default function StepCategoryAccount({
  type, category, setCategory, accountId, setAccountId, 
  fromAccountId, setFromAccountId, toAccountId, setToAccountId,
  jarId, setJarId, isMasterIncome, setIsMasterIncome, accounts, jars, onNext, 
}: StepCategoryAccountProps){
  
  const accent = useMemo(() => {
    switch (type) {
      case "expense": return "#ef4444";
      case "income": return "#22c55e";
      default: return "#3b82f6";
    }
  }, [type]);

  const canContinue = 
    type === "transfer"
      ? category !== "" &&
        fromAccountId !== "" &&
        toAccountId !== "" &&
        fromAccountId !== toAccountId
      : category !== "" && 
        accountId !== "" && 
        (isMasterIncome || jarId !== "");

return (
    <div style={containerStyle}>
      <SelectionGroup 
        label="Category" 
        items={CATEGORIES.map(c => ({ id: c, name: c }))} 
        selectedId={category} 
        onSelect={setCategory} 
        accent={accent} 
      />
      
      {type === "transfer" ? (
        <>
          <SelectionGroup 
            label="From Account" 
            items={accounts} 
            selectedId={fromAccountId || ""} 
            onSelect={setFromAccountId || (() => {})} 
            accent={accent} 
          />
          <SelectionGroup 
            label="To Account" 
            items={accounts.filter(a => String(a.id) !== fromAccountId)} 
            selectedId={toAccountId || ""} 
            onSelect={setToAccountId || (() => {})} 
            accent={accent} 
          />
        </>
      ) : (
        <SelectionGroup 
          label="Account" 
          items={accounts} 
          selectedId={accountId} 
          onSelect={setAccountId} 
          accent={accent} 
        />
      )}

      {type === "income" && (
        <section 
          style={{ 
            marginBottom: 24, 
            display: "flex", 
            alignItems: "center", 
            justifyContent: "space-between", 
            padding: "16px", 
            backgroundColor: "#111827", 
            borderRadius: 14, 
            border: "1px solid #1f2937" 
          }}
        >
          <div>
            <div style={{ fontSize: 14, fontWeight: 600, color: "#fff" }}>Master Income</div>
            <div style={{ fontSize: 12, color: "#9ca3af", marginTop: 2 }}>Auto-distribute to jars</div>
          </div>
          <button
            type="button"
            onClick={() => 
            {
                const next = !isMasterIncome;
                setIsMasterIncome(next);

                if (next) {
                  setJarId("");
                }
              }}
            style={{
              width: 44,
              height: 24,
              borderRadius: 12,
              border: "none",
              backgroundColor: isMasterIncome ? accent : "#374151",
              position: "relative",
              cursor: "pointer",
              transition: "background-color 0.2s ease",
              padding: 0
            }}
          >
            <div style={{
              width: 20,
              height: 20,
              borderRadius: "50%",
              backgroundColor: "white",
              position: "absolute",
              top: 2,
              left: isMasterIncome ? 22 : 2,
              transition: "left 0.2s ease"
            }} />
          </button>
        </section>
      )}

      {/* Only show jars if type is not transfer AND (Master Income is disabled OR type is not income) */}
      {type !== "transfer" && !(type === "income" && isMasterIncome) && (
        <SelectionGroup 
          label="Jar" 
          items={jars} 
          selectedId={jarId} 
          onSelect={setJarId} 
          accent={accent} 
        />
      )}

      <button
        type="button"
        onClick={onNext}
        disabled={!canContinue}
        style={{
          ...buttonStyle,
          backgroundColor: canContinue ? accent : "#2d2d31",
          cursor: canContinue ? "pointer" : "not-allowed",
          opacity: canContinue ? 1 : 0.6,
        }}
      >
        Continue
      </button>
    </div>
  );
}

const containerStyle: React.CSSProperties = {
  backgroundColor: "#09090b",
  padding: 24,
  borderRadius: 20,
  border: "1px solid #1f1f23",
};

const heading: React.CSSProperties = {
  fontSize: 11,
  fontWeight: 700,
  color: "#71717a",
  textTransform: "uppercase",
  letterSpacing: "0.08em",
  marginBottom: 12,
};

const chipWrap: React.CSSProperties = {
  display: "flex",
  flexWrap: "wrap",
  gap: 8,
};

const chip: React.CSSProperties = {
  padding: "8px 14px",
  borderRadius: 10,
  fontSize: 13,
  fontWeight: 500,
  border: "none",
  cursor: "pointer",
  transition: "all 0.2s ease",
};

const buttonStyle: React.CSSProperties = {
  width: "100%",
  padding: 14,
  borderRadius: 12,
  border: "none",
  fontWeight: 600,
  fontSize: 14,
  color: "#fff",
  marginTop: 16,
};
