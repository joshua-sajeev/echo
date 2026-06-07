import { useMemo } from "react";

const CATEGORIES = [
  "Food", "Transport", "Shopping", "Donations", "Entertainment",
  "Health", "Income", "Investment", "Housing", "Transfers",
];

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
  jarId, setJarId, accounts, jars, onNext
}: any) {
  
  const accent = useMemo(() => {
    switch (type) {
      case "expense": return "#ef4444";
      case "income": return "#22c55e";
      default: return "#3b82f6";
    }
  }, [type]);

  const canContinue = !!(category && accountId && jarId);

  return (
    <div style={containerStyle}>
      <SelectionGroup label="Category" items={CATEGORIES.map(c => ({ id: c, name: c }))} selectedId={category} onSelect={setCategory} accent={accent} />
      <SelectionGroup label="Account" items={accounts} selectedId={accountId} onSelect={setAccountId} accent={accent} />
      <SelectionGroup label="Jar" items={jars} selectedId={jarId} onSelect={setJarId} accent={accent} />

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
