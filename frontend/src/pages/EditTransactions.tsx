import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import Navbar from "../components/Navbar";
import type { TransactionFormValues } from "../components/TransactionForm";
import TransactionForm from "../components/TransactionForm";
const API_BASE = import.meta.env.VITE_API_URL; 

export default function EditTransactions({
  setUser,
}: {
  setUser: any;
}) {
  const { id } = useParams();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(true);

  const [initialValues, setInitialValues] =
    useState<Partial<TransactionFormValues>>();

  useEffect(() => {
    async function loadTransaction() {
      try {
        const res = await fetch(
          `${API_BASE}/transactions/${id}`,
          {
            credentials: "include",
          }
        );

        if (!res.ok) {
          throw new Error("Failed to load transaction");
        }

        const tx = await res.json();


        setInitialValues({
          type: tx.type,
          date: tx.date.split("T")[0],
          name: tx.name,
          amount: String(tx.amount / 100),
          category: tx.category ?? "",

          accountId:
            tx.type === "expense"
              ? String(tx.from_account_id ?? "")
              : tx.type === "income"
              ? String(tx.to_account_id ?? "")
              : "",

          fromId: String(tx.from_account_id ?? ""),
          toId: String(tx.to_account_id ?? ""),
          jarId: String(tx.jar_id ?? ""),
          isMasterIncome: tx.is_master_income ?? false,
        });
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    }

    loadTransaction();
  }, [id]);

  async function handleUpdate(
    payload: Record<string, unknown>
  ) {
    const res = await fetch(
      `${API_BASE}/transactions/${id}`,
      {
        method: "PUT",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      }
    );

    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || "Update failed");
    }

    navigate("/");
  }

  if (loading) {
    return (
      <div
        style={{
          minHeight: "100vh",
          background: "#0a0c12",
          color: "#d1d5db",
          display: "grid",
          placeItems: "center",
        }}
      >
        Loading...
      </div>
    );
  }

  return (
    <div
      style={{
        background: "#0a0c12",
        minHeight: "100vh",
        color: "#d1d5db",
      }}
    >
      <Navbar setUser={setUser} />

      <div
        style={{
          maxWidth: 480,
          margin: "0 auto",
          padding: "24px 16px 80px",
        }}
      >
        <h1
          style={{
            fontSize: 22,
            fontWeight: 600,
            marginBottom: 24,
            color: "#e8eaf0",
          }}
        >
          Edit Transaction
        </h1>

        <TransactionForm
          initialValues={initialValues}
          onSubmit={handleUpdate}
          submitLabel="Update Transaction"
        />
      </div>
    </div>
  );
}
