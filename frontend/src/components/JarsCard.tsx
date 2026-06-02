import React from "react";

interface Jar {
  id: number;
  name: string;
  allocation_type: string;
  value: number;
  allocated_amount: number;
}

interface Transaction {
  id: number;
  jar_id: number | null;
  amount: number;
  type: string;
}

interface Props {
  jars: Jar[];
  transactions: Transaction[];
}

const fmt = (n: number) =>
  "₹" +
  (n / 100).toLocaleString("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

const JAR_COLORS = [
  "#60a5fa",
  "#34d399",
  "#f59e0b",
  "#f472b6",
  "#a78bfa",
];

export default function JarsCard({
  jars,
  transactions,
}: Props) {
  const usage: Record<number, number> = {};

  transactions.forEach((t) => {
    if (
      t.type?.toLowerCase() === "expense" &&
      t.jar_id
    ) {
      usage[t.jar_id] =
        (usage[t.jar_id] || 0) + Number(t.amount || 0);
    }
  });

  const total = jars.reduce(
    (sum, jar) => sum + Number(jar.allocated_amount || 0),
    0
  );

  return (
    <div style={{ ...card, padding: "16px" }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "baseline",
          marginBottom: 16,
        }}
      >
        <span style={sectionLabel}>Jars</span>

        <span
          style={{
            fontSize: 18,
            fontWeight: 700,
            color: "#f3f4f6",
          }}
        >
          {fmt(total)}
        </span>
      </div>

      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: 10,
        }}
      >
        {jars.map((jar, idx) => {
          const color =
            JAR_COLORS[idx % JAR_COLORS.length];

          const used = usage[jar.id] || 0;

          const progressPercent =
            jar.allocated_amount > 0
              ? (used / jar.allocated_amount) * 100
              : 0;

          return (
            <div
              key={jar.id}
              style={{
                ...jarRow,
                padding: "10px 12px",
              }}
            >
              <div
                style={{
                  display: "flex",
                  justifyContent:
                    "space-between",
                  alignItems: "center",
                  marginBottom: 6,
                }}
              >
                <span
                  style={{
                    color,
                    fontSize: 14,
                    fontWeight: 600,
                  }}
                >
                  {jar.name}
                </span>

                <span
                  style={{
                    color: "#6b7280",
                    fontSize: 12,
                    fontWeight: 600,
                  }}
                >
                  {Math.min(
                    progressPercent,
                    100
                  ).toFixed(0)}
                  %
                </span>
              </div>

              <div
                style={{
                  display: "flex",
                  justifyContent:
                    "space-between",
                  alignItems: "center",
                  marginBottom: 6,
                }}
              >
                <span
                  style={{
                    color: "#9ca3af",
                    fontSize: 13,
                    fontWeight: 500,
                  }}
                >
                  {fmt(used)} /{" "}
                  {fmt(jar.allocated_amount)}
                </span>
              </div>

              <div
                style={{
                  height: 5,
                  background: "#1a1d27",
                  borderRadius: 999,
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    width: `${Math.min(
                      progressPercent,
                      100
                    )}%`,
                    height: "100%",
                    background: color,
                  }}
                />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

const card: React.CSSProperties = {
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 14,
};

const jarRow: React.CSSProperties = {
  background: "#161922",
  border: "0.5px solid #1e2130",
  borderRadius: 12,
};

const sectionLabel: React.CSSProperties = {
  fontSize: 12,
  fontWeight: 600,
  letterSpacing: "0.05em",
  textTransform: "uppercase",
  color: "#6b7280",
};
