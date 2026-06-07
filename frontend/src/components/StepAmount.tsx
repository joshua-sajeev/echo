import type { TxType } from "./TransactionForm";

interface StepAmountProps {
  amount: string;
  setAmount: (value: string) => void;

  type: TxType;
  setType: (value: TxType) => void;

  onNext: () => void;
}

const keypad = [
  "7",
  "8",
  "9",
  "4",
  "5",
  "6",
  "1",
  "2",
  "3",
  ".",
  "0",
  "⌫",
];

export default function StepAmount({
  amount,
  setAmount,
  type,
  setType,
  onNext,
}: StepAmountProps) {
  const accent =
    type === "expense"
      ? "#ef4444"
      : type === "income"
      ? "#22c55e"
      : "#3b82f6";

  const handleKeyPress = (key: string) => {
    if (key === "⌫") {
      setAmount(
        amount.length <= 1
          ? "0"
          : amount.slice(0, -1)
      );
      return;
    }

    if (key === ".") {
      if (amount.includes(".")) return;

      setAmount(`${amount}.`);
      return;
    }

    if (amount === "0") {
      setAmount(key);
      return;
    }

    setAmount(amount + key);
  };

  const canContinue =
    amount !== "0" &&
    amount !== "" &&
    amount !== ".";

  return (
    <>
      <div
        style={{
          textAlign: "center",
          marginTop: 32,
          marginBottom: 32,
        }}
      >
        <div
          style={{
            color: accent,
            fontSize: 14,
            fontWeight: 600,
            marginBottom: 16,
            textTransform: "capitalize",
          }}
        >
          {type}
        </div>

        <div
          style={{
            fontSize: 56,
            fontWeight: 700,
            color: "#fff",
            letterSpacing: "-2px",
          }}
        >
          ₹{amount}
        </div>
      </div>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "1fr 1fr 1fr",
          gap: 8,
          marginBottom: 32,
        }}
      >
        {(
          [
            "expense",
            "income",
            "transfer",
          ] as TxType[]
        ).map((t) => {
          const color =
            t === "expense"
              ? "#ef4444"
              : t === "income"
              ? "#22c55e"
              : "#3b82f6";

          return (
            <button
              key={t}
              onClick={() => setType(t)}
              style={{
                height: 48,
                borderRadius: 999,
                border:
                  type === t
                    ? `1px solid ${color}`
                    : "1px solid #1f2937",
                background:
                  type === t
                    ? `${color}20`
                    : "#111827",
                color:
                  type === t
                    ? color
                    : "#9ca3af",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              {t.charAt(0).toUpperCase() +
                t.slice(1)}
            </button>
          );
        })}
      </div>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(3, 1fr)",
          gap: 12,
        }}
      >
        {keypad.map((key) => (
          <button
            key={key}
            onClick={() => handleKeyPress(key)}
            style={{
              height: 72,
              borderRadius: 18,
              border: "1px solid #1f2937",
              background: "#111827",
              color: "#fff",
              fontSize: 28,
              fontWeight: 600,
              cursor: "pointer",
            }}
          >
            {key}
          </button>
        ))}
      </div>

      <button
        onClick={onNext}
        disabled={!canContinue}
        style={{
          width: "100%",
          marginTop: 24,
          height: 56,
          border: "none",
          borderRadius: 14,
          background: accent,
          color: "#fff",
          fontWeight: 700,
          fontSize: 16,
          cursor: canContinue
            ? "pointer"
            : "not-allowed",
          opacity: canContinue ? 1 : 0.5,
        }}
      >
        Next
      </button>
    </>
  );
}
