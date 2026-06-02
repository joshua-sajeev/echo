import Navbar from "../components/Navbar";
import AccountsCard from "../components/AccountsCard";
import JarsCard from "../components/JarsCard";
import RecentTransactions from "../components/RecentTransactions";

import { useDashboard } from "../hooks/useDashboard";
import { useNavigate } from "react-router-dom";

import { formatCurrency } from "../utils/currency";
import { getCurrentMonthYear } from "../utils/date";

import { statCard, actionBtn } from "../styles/dashboard";

const { month, year } = getCurrentMonthYear();

export default function Dashboard({ user, setUser }: any) {
  const navigate = useNavigate();

  const { data, loading } = useDashboard();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!data) {
    return <div>Failed to load dashboard</div>;
  }

  const income = data.transactions
    .filter((t: any) => t.type === "income")
    .reduce((sum: number, t: any) => sum + t.amount, 0);

  const expenses = data.transactions
    .filter((t: any) => t.type === "expense")
    .reduce((sum: number, t: any) => sum + t.amount, 0);

  const savings =
    income > 0
      ? (((income - expenses) / income) * 100).toFixed(1)
      : "0";

  const quickStats = [
    {
      label: "Income",
      value: formatCurrency(income),
      delta: "",
      up: true,
    },
    {
      label: "Expenses",
      value: formatCurrency(expenses),
      delta: "",
      up: false,
    },
    {
      label: "Savings",
      value: `${savings}%`,
      delta: "",
      up: true,
    },
  ];

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

      <div
        style={{
          padding: "24px 16px 80px",
          maxWidth: 480,
          margin: "0 auto",
          display: "flex",
          flexDirection: "column",
          gap: "20px",
        }}
      >
        <div>
          <p
            style={{
              color: "#4b5563",
              fontSize: 11,
              letterSpacing: "0.1em",
              textTransform: "uppercase",
              margin: "0 0 6px 0",
            }}
          >
            {month} {year}
          </p>

          <h1
            style={{
              color: "#e8eaf0",
              fontSize: 22,
              fontWeight: 600,
              margin: 0,
            }}
          >
            {user?.name ? `Hey, ${user.name} 👋` : "Hi 👋"}
          </h1>
        </div>

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "1fr 1fr 1fr",
            gap: 10,
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
                  margin: "0 0 8px 0",
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
                  margin: "0 0 10px 0",
                }}
              >
                {s.value}
              </p>

              {s.delta && (
                <div style={{ display: "flex" }}>
                  <span
                    style={{
                      fontSize: 10,
                      fontWeight: 500,
                      color: s.up ? "#1D9E75" : "#E24B4A",
                      background: s.up ? "#0d2a1f" : "#2a1212",
                      border: `0.5px solid ${
                        s.up ? "#0f6e56" : "#6e1f1f"
                      }`,
                      borderRadius: 4,
                      padding: "2px 6px",
                      display: "inline-flex",
                      alignItems: "center",
                    }}
                  >
                    {s.delta}
                  </span>
                </div>
              )}
            </div>
          ))}
        </div>

        <div style={{ display: "flex", gap: 12 }}>
          <button
            onClick={() => navigate("/transactions")}
            style={actionBtn}
          >
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              style={{ flexShrink: 0 }}
            >
              <path
                d="M12 5v14M5 12h14"
                stroke="#e8eaf0"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            <span style={{ transform: "translateY(0.5px)" }}>
              Add Transaction
            </span>
          </button>

          <button
            onClick={() => navigate("/templates")}
            style={actionBtn}
          >
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              style={{ flexShrink: 0 }}
            >
              <path
                d="M4 7h16M4 12h10M4 17h16"
                stroke="#e8eaf0"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
            <span style={{ transform: "translateY(0.5px)" }}>
              Templates
            </span>
          </button>
        </div>

        <AccountsCard accounts={data.accounts} />

        <JarsCard
          jars={data.jars}
          transactions={data.transactions}
        />

        <RecentTransactions
          transactions={data.transactions}
        />
      </div>
    </div>
  );
}
