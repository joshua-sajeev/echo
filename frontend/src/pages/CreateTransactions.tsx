import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import TransactionForm, {
  type Account,
  type Jar,
  type TxType,
} from "../components/TransactionForm";

import StepAmount from "../components/StepAmount";
import StepCategoryAccount from "../components/StepCategoryAccount";

const API_BASE = import.meta.env.VITE_API_URL;

export default function CreateTransactions() {
  const navigate = useNavigate();

  const [step, setStep] = useState(1);

  // Step 1
  const [amount, setAmount] = useState("0");
  const [type, setType] = useState<TxType>("expense");

  // Step 2
  const [category, setCategory] = useState("");
  const [accountId, setAccountId] = useState("");
  const [jarId, setJarId] = useState("");

  // API data
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [jars, setJars] = useState<Jar[]>([]);

  useEffect(() => {
    fetch(`${API_BASE}/accounts/balances`, {
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    })
      .then((r) => r.json())
      .then((d: Account[]) => setAccounts(d))
      .catch(console.error);

    fetch(`${API_BASE}/jars`, {
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    })
      .then((r) => r.json())
      .then((d: Jar[]) => setJars(d))
      .catch(console.error);
  }, []);

  const handleSubmit = async (
    payload: Record<string, unknown>,
  ) => {
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

      throw new Error(
        `Server returned ${r.status}: ${text.slice(
          0,
          200,
        )}`,
      );
    }

    navigate("/");
  };

  return (
    <>
      <style>{`
        * {
          box-sizing: border-box;
        }

        body {
          margin: 0;
          background: #080a0f;
          color: white;
        }

        button {
          font-family: inherit;
        }
      `}</style>

      <nav
        style={{
          borderBottom:
            "0.5px solid #1e2130",
          background: "#080a0f",
          position: "sticky",
          top: 0,
          zIndex: 10,
        }}
      >
        <div style={navInner}>
          <button
            onClick={() => {
              if (step === 1) {
                navigate("/");
              } else {
                setStep((s) => s - 1);
              }
            }}
            style={backButton}
          >
            ← Back
          </button>

          <span
            style={{
              fontSize: 15,
              fontWeight: 600,
              color: "#e5e7eb",
            }}
          >
            New Transaction
          </span>

          <div style={{ width: 48 }} />
        </div>
      </nav>

      <main
        style={{
          maxWidth: 480,
          margin: "0 auto",
          padding: 16,
          minHeight:
            "calc(100vh - 60px)",
        }}
      >
        <StepIndicator step={step} />

        {step === 1 && (
          <StepAmount
            amount={amount}
            setAmount={setAmount}
            type={type}
            setType={setType}
            onNext={() => setStep(2)}
          />
        )}

        {step === 2 && (
          <StepCategoryAccount
            type={type}
            category={category}
            setCategory={setCategory}
            accountId={accountId}
            setAccountId={setAccountId}
            jarId={jarId}
            setJarId={setJarId}
            accounts={accounts}
            jars={jars}
            onNext={() => setStep(3)}
          />
        )}

        {step === 3 && (
          <TransactionForm
            initialValues={{
              amount,
              type,
              category,
              accountId,
              jarId,
            }}
            submitLabel="Save Transaction"
            onSubmit={handleSubmit}
          />
        )}
      </main>
    </>
  );
}

function StepIndicator({
  step,
}: {
  step: number;
}) {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        gap: 10,
        marginTop: 8,
        marginBottom: 24,
      }}
    >
      {[1, 2, 3].map((s) => (
        <div
          key={s}
          style={{
            width: 10,
            height: 10,
            borderRadius: "50%",
            background:
              s <= step
                ? "#ffffff"
                : "#374151",
          }}
        />
      ))}
    </div>
  );
}

const navInner: React.CSSProperties = {
  maxWidth: 480,
  margin: "0 auto",
  padding: "14px 16px",
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
};

const backButton: React.CSSProperties = {
  background: "transparent",
  border: "none",
  color: "#9ca3af",
  cursor: "pointer",
  fontSize: 14,
};
