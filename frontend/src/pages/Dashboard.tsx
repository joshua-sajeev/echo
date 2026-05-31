import Navbar from "../components/Navbar";
import AccountsCard from "./AccountsCard";
import { useNavigate } from "react-router-dom";

// ── helpers ───────────────────────────────────────────────────────────────────
const fmt = (n: number) =>
  "₹" + n.toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

const now = new Date();
const MONTH = now.toLocaleString("en-IN", { month: "long" });
const YEAR = now.getFullYear();

// ── mock stats ───────────────────────────────────────────────────────────────
const quickStats = [
  { label: "Income", value: fmt(97000), delta: "+12%", up: true },
  { label: "Expenses", value: fmt(34219), delta: "-8%", up: false },
  { label: "Savings", value: "64.7%", delta: "+3%", up: true },
];

// ── dashboard ─────────────────────────────────────────────────────────────────
export default function Dashboard({ user, setUser }: any) {
  const navigate = useNavigate();

  return (
    <div
      style={{
        background: "#0a0c12",
        minHeight: "100vh",
        fontFamily: "'Syne', sans-serif",
        color: "#e8eaf0",
      }}
    >
      <Navbar setUser={setUser} />

      <div style={{ padding: "20px 16px 80px", maxWidth: 480, margin: "0 auto" }}>
        {/* greeting */}
        <div style={{ marginBottom: 20 }}>
          <p
            style={{
              color: "#4b5563",
              fontSize: 11,
              letterSpacing: "0.1em",
              textTransform: "uppercase",
              margin: "0 0 4px",
            }}
          >
            {MONTH} {YEAR}
          </p>
          <h1 style={{ color: "#e8eaf0", fontSize: 22, fontWeight: 600, margin: 0 }}>
            {user?.name ? `Hey, ${user.name} 👋` : "Hi Joshua 👋"}
          </h1>
        </div>

        {/* ACTION BUTTONS */}
        <div style={{ display: "flex", gap: 10, marginBottom: 14 }}>
          {/* Add Transaction */}
          <button
            onClick={() => navigate("/transactions")}
            style={actionBtn}
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path
                d="M12 5v14M5 12h14"
                stroke="#e8eaf0"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            Add Transaction
          </button>

          {/* Templates */}
          <button
            onClick={() => navigate("/templates")}
            style={actionBtn}
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path
                d="M4 7h16M4 12h10M4 17h16"
                stroke="#e8eaf0"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            Templates
          </button>
        </div>

        {/* stats */}
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "1fr 1fr 1fr",
            gap: 8,
            marginBottom: 12,
          }}
        >
          {quickStats.map((s) => (
            <div key={s.label} style={statCard}>
              <p
                style={{
                  color: "#4b5563",
                  fontSize: 10,
                  letterSpacing: "0.06em",
                  textTransform: "uppercase",
                  margin: "0 0 6px",
                }}
              >
                {s.label}
              </p>
              <p
                style={{
                  fontFamily: "'IBM Plex Mono', monospace",
                  fontSize: 14,
                  fontWeight: 500,
                  color: "#e8eaf0",
                  margin: "0 0 6px",
                }}
              >
                {s.value}
              </p>
              <span
                style={{
                  fontSize: 10,
                  fontWeight: 500,
                  color: s.up ? "#1D9E75" : "#E24B4A",
                  background: s.up ? "#0d2a1f" : "#2a1212",
                  border: `0.5px solid ${s.up ? "#0f6e56" : "#6e1f1f"}`,
                  borderRadius: 4,
                  padding: "1px 5px",
                }}
              >
                {s.delta}
              </span>
            </div>
          ))}
        </div>

        {/* accounts */}
        <AccountsCard />
      </div>
    </div>
  );
}

// ── styles ────────────────────────────────────────────────────────────────────
const statCard: React.CSSProperties = {
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 14,
  padding: "12px 10px",
};

const actionBtn: React.CSSProperties = {
  flex: 1,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  gap: 6,
  padding: "10px 8px",
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 12,
  color: "#e8eaf0",
  fontSize: 12,
  cursor: "pointer",
};
