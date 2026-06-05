import TransactionForm from "../components/TransactionForm";
import { useNavigate } from "react-router-dom";

const API_BASE = import.meta.env.VITE_API_URL;  

export default function NewTransaction() {

const navigate = useNavigate();
const handleSubmit = async (payload: Record<string, unknown>) => {
  const r = await fetch(`${API_BASE}/transactions`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!r.ok) {
    const text = await r.text();
    throw new Error(`Server returned ${r.status}: ${text.slice(0, 200)}`);
  }

  navigate("/");
};

  return (
    <>
      <style>{`* { box-sizing: border-box; } body { margin: 0; background: #080a0f; }`}</style>

      <nav style={{
        borderBottom: "0.5px solid #1e2130",
        background: "#080a0f",
        position: "sticky", top: 0, zIndex: 10,
      }}>
        <div style={navInner}>
          <a href="/" style={{ color: "#6b7280", fontSize: 14, textDecoration: "none" }}>← Back</a>
          <span style={{ fontSize: 15, fontWeight: 600, color: "#e5e7eb" }}>New Transaction</span>
          <div style={{ width: 48 }} />
        </div>
      </nav>

      <main style={{ maxWidth: 480, margin: "0 auto", padding: "24px 16px 80px" }}>
        <TransactionForm
          onSubmit={handleSubmit}
          submitLabel="Save Transaction"
        />
      </main>
    </>
  );
}

const navInner: React.CSSProperties = {
  maxWidth: 480, margin: "0 auto", padding: "14px 16px",
  display: "flex", alignItems: "center", justifyContent: "space-between",
};
