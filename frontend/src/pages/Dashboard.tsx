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

const now = new Date();

const currentMonthTransactions = data.transactions.filter((t: any) => {
  const txDate = new Date(t.date);

  return (
    txDate.getMonth() === now.getMonth() &&
    txDate.getFullYear() === now.getFullYear()
  );
});

const income = currentMonthTransactions
  .filter((t: any) => t.type === "income")
  .reduce((sum: number, t: any) => sum + t.amount, 0);

const expenses = currentMonthTransactions
  .filter((t: any) => t.type === "expense")
  .reduce((sum: number, t: any) => sum + t.amount, 0);

const netCashFlow = income - expenses;

const savingsRate =
  income > 0
    ? (((income - expenses) / income) * 100).toFixed(1)
    : "0";

const quickStats = [
  {
    label: "Income",
    value: formatCurrency(income),
    color: "#1D9E75",
  },
  {
    label: "Expenses",
    value: formatCurrency(expenses),
    color: "#E24B4A",
  },
  {
    label: "Net",
    value: formatCurrency(netCashFlow),
    color: netCashFlow >= 0 ? "#4F7CFF" : "#E24B4A",
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
  <div
    key={s.label}
    style={{
      ...statCard,
      border: `1px solid ${s.color}75`,
    }}
  >
    <p
      style={{
        color: s.color,
        fontSize: 12,
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
    fontSize: "clamp(10px, 2.8vw, 14px)",
    fontWeight: 600,
    color: s.color,
    margin: 0,
    overflow: "hidden",
    textOverflow: "ellipsis",
    whiteSpace: "nowrap",
    minWidth: 0,
  }}
  title={s.value}
>
  {s.value}
</p>
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
