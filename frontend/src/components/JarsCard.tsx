import React  from "react";

interface Jar {
  id: number;
  name: string;
  allocation_type: string;
  value: number;
  allocated_amount: number;
  balance: number;
  spent_this_month: number;
}

interface Props {
  jars: Jar[];
}

const fmt = (n: number) =>
  "₹" +
  (n / 100).toLocaleString("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

const fmtCompact = (n: number) => {
  const abs = Math.abs(n / 100);
  const sign = n < 0 ? "-" : "+";
  if (abs >= 100000) return `${sign}₹${(abs / 100000).toFixed(1)}L`;
  if (abs >= 1000) return `${sign}₹${(abs / 1000).toFixed(1)}k`;
  return `${sign}₹${abs.toFixed(0)}`;
};

const JAR_COLORS = [
  "#60a5fa",
  "#34d399",
  "#f59e0b",
  "#f472b6",
  "#a78bfa",
];

export default function JarsCard({ jars }: Props) {

  if (!jars || jars.length === 0) {
    return (
      <div style={card}>
        <span style={sectionLabel}>Jars</span>
        <p style={{ color: "#4b5563", fontSize: 13, textAlign: "center", padding: "12px 0" }}>
          No jars found
        </p>
      </div>
    );
  }

  const totalBalance = jars.reduce((sum, jar) => sum + (jar.balance ?? 0), 0);
  const totalAllocated = jars.reduce((sum, jar) => sum + (jar.allocated_amount ?? 0), 0);

  return (
    <div style={{ ...card, padding: "16px" }}>
      {/* header */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "baseline", marginBottom: 4 }}>
        <span style={sectionLabel}>Jars</span>
        <span style={{ fontSize: 18, fontWeight: 700, color: "#f3f4f6" }}>
          {fmt(totalBalance)}
        </span>
      </div>
      <p style={{ color: "#4b5563", fontSize: 11, margin: "0 0 14px", textAlign: "right" }}>
        This month +{fmt(totalAllocated)}
      </p>

      <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
        {jars.map((jar, idx) => {
          const color = JAR_COLORS[idx % JAR_COLORS.length];
          const balance = jar.balance ?? 0;
          const allocated = jar.allocated_amount ?? 0;
          const spent = jar.spent_this_month ?? 0;

          // balance = carryover + this month's allocation - spent
          // It's the single "how much do I have available" number
          const leftThisMonth = allocated - spent;
          const carryover = balance - leftThisMonth;

          // % of this month's allocation spent
          const spentPct = allocated > 0
            ? Math.min(Math.round((spent / allocated) * 100), 100)
            : 0;

          const progressPct = allocated > 0
            ? Math.min(Math.max((spent / allocated) * 100, 0), 100)
            : 0;

          const barColor = progressPct >= 90 ? "#E24B4A"
            : progressPct >= 70 ? "#f59e0b"
            : color;

          const isNegative = balance < 0;

          return (
            <div key={jar.id} style={{ ...jarRow, padding: "12px 14px" }}>
              {/* jar name */}
              <p style={{ margin: "0 0 10px", fontSize: 11, fontWeight: 600, color, letterSpacing: "0.06em", textTransform: "uppercase" }}>
                {jar.name}
              </p>

              {/* 2×2 grid */}
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", marginBottom: 10 }}>
                {/* TOP LEFT: total available (carryover + this month - spent) = balance */}
                <div>
                  <p style={{
                    margin: 0,
                    fontSize: 18,
                    fontWeight: 700,
                    fontFamily: "'IBM Plex Mono', monospace",
                    color: isNegative ? "#E24B4A" : "#f3f4f6",
                    lineHeight: 1.1,
                  }}>
                    {fmt(balance)}
                  </p>
                  <p style={{ margin: "3px 0 0", fontSize: 10, color: "#4b5563" }}>
                    available
                  </p>
                </div>

                {/* TOP RIGHT: % of this month's allocation spent */}
                <div style={{ textAlign: "right" }}>
                  <p style={{
                    margin: 0,
                    fontSize: 18,
                    fontWeight: 700,
                    fontFamily: "'IBM Plex Mono', monospace",
                    color: spentPct >= 90 ? "#E24B4A" : spentPct >= 70 ? "#f59e0b" : "#6b7280",
                    lineHeight: 1.1,
                  }}>
                    {allocated > 0 ? `${spentPct}%` : "—"}
                  </p>
                  <p style={{ margin: "3px 0 0", fontSize: 10, color: "#4b5563" }}>
                    of {fmt(allocated)} spent
                  </p>
                </div>
              </div>

              {/* progress bar */}
              <div style={{ height: 3, background: "#1a1d27", borderRadius: 999, overflow: "hidden", marginBottom: 8 }}>
                <div style={{
                  width: `${progressPct}%`,
                  height: "100%",
                  background: barColor,
                  borderRadius: 999,
                  transition: "width 0.4s ease, background 0.3s ease",
                }} />
              </div>

              {/* BOTTOM row */}
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr" }}>
                {/* BOTTOM LEFT: spent this month */}
                <p style={{ margin: 0, fontSize: 11, color: "#4b5563" }}>
                  {fmt(spent)} spent
                </p>

                {/* BOTTOM RIGHT: carryover from previous months */}
                <p style={{ margin: 0, fontSize: 11, textAlign: "right" }}>
                  {carryover === 0 ? (
                    <span style={{ color: "#374151" }}>no carryover</span>
                  ) : (
                    <span style={{ color: carryover > 0 ? "#1D9E75" : "#E24B4A" }}>
                      {fmtCompact(carryover)} carryover
                    </span>
                  )}
                </p>
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
