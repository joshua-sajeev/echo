import { useState } from "react";
import GoalsModal from "./GoalsModal";
import AllocationModal from "./AllocationModal";

interface Goal {
  id: number;
  name: string;
  target_amount: number;
  saved_amount: number;
  allocation_percentage: number;
  is_archived: boolean;
  deadline: string | null;
  created_at: string;
  updated_at: string;
}

interface GoalsCardProps {
  goals: Goal[];
  onRefresh: () => void;
}

const fmt = (n: number) =>
  "₹" +
  (n / 100).toLocaleString("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

const GOAL_COLORS = [
  "#60a5fa", // blue
  "#34d399", // green
  "#fbbf24", // amber
  "#f472b6", // pink
  "#a78bfa", // purple
];

const getGoalColor = (id: number) => GOAL_COLORS[id % GOAL_COLORS.length];

export default function GoalsCard({ goals, onRefresh }: GoalsCardProps) {
  const [showGoalsModal, setShowGoalsModal] = useState(false);
  const [showAllocationModal, setShowAllocationModal] = useState(false);
  const [allocationType, setAllocationType] = useState<"manual" | "auto" | null>(null);

  const activeGoals = goals.filter((g) => !g.is_archived);

  const totalAllocated = activeGoals.reduce(
    (sum, g) => sum + g.allocation_percentage,
    0
  );

  const totalSaved = activeGoals.reduce((sum, g) => sum + g.saved_amount, 0);

  const hasAllocatedGoals = activeGoals.some((g) => g.allocation_percentage > 0);

  // Auto allocate handler
  const handleAutoAllocate = async () => {
    setAllocationType("auto");
    setShowAllocationModal(true);
  };

  // Manual allocate handler
  const handleManualAllocate = () => {
    setAllocationType("manual");
    setShowAllocationModal(true);
  };

  if (!activeGoals || activeGoals.length === 0) {
    return (
      <div style={card}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 14 }}>
          <span style={sectionLabel}>Goals</span>
          <button
            style={iconBtn}
            onClick={() => setShowGoalsModal(true)}
            title="Create goal"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M12 5v14" />
              <path d="M5 12h14" />
            </svg>
          </button>
        </div>
        <p style={{ color: "#4b5563", fontSize: 13, textAlign: "center", padding: "12px 0" }}>
          No active goals found
        </p>

        {showGoalsModal && (
          <GoalsModal
            existingGoals={[]}
            onClose={() => setShowGoalsModal(false)}
            onSuccess={() => {
              setShowGoalsModal(false);
              onRefresh();
            }}
          />
        )}
      </div>
    );
  }

  return (
    <>
      <div style={card}>
        {/* Header */}
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 14 }}>
          <span style={sectionLabel}>Goals</span>
          <button
            style={iconBtn}
            onClick={() => setShowGoalsModal(true)}
            title="Add goal"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M12 5v14" />
              <path d="M5 12h14" />
            </svg>
          </button>
        </div>

        {/* Summary */}
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10, marginBottom: 14 }}>
          <div style={{ ...summaryBox, borderColor: "#60a5fa" }}>
            <p style={{ color: "#60a5fa", fontSize: 11, margin: "0 0 6px" }}>Total Saved</p>
            <p style={{ fontSize: 16, fontWeight: 700, color: "#f3f4f6", margin: 0, fontFamily: "'IBM Plex Mono', monospace" }}>
              {fmt(totalSaved)}
            </p>
          </div>
          <div style={{ ...summaryBox, borderColor: totalAllocated === 100 ? "#1D9E75" : "#f59e0b" }}>
            <p style={{ color: totalAllocated === 100 ? "#1D9E75" : "#f59e0b", fontSize: 11, margin: "0 0 6px" }}>Allocation</p>
            <p style={{ fontSize: 16, fontWeight: 700, color: totalAllocated === 100 ? "#1D9E75" : "#f59e0b", margin: 0, fontFamily: "'IBM Plex Mono', monospace" }}>
              {totalAllocated}%
            </p>
          </div>
        </div>

        {/* Goals List */}
        <div style={{ display: "flex", flexDirection: "column", gap: 10, marginBottom: 14 }}>
          {activeGoals.map((goal) => {
            const color = getGoalColor(goal.id);
            const progress = goal.target_amount > 0 ? (goal.saved_amount / goal.target_amount) * 100 : 0;

            return (
              <div key={goal.id} style={{ ...goalRow, borderColor: color + "40" }}>
                {/* Top: Name + Allocation */}
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "baseline", marginBottom: 10 }}>
                  <p style={{ color, fontSize: 12, fontWeight: 600, margin: 0, letterSpacing: "0.05em", textTransform: "uppercase" }}>
                    {goal.name}
                  </p>
                  <p style={{ color: "#6b7280", fontSize: 11, margin: 0 }}>
                    {goal.allocation_percentage}%
                  </p>
                </div>

                {/* Progress Bar */}
                <div style={{ height: 3, background: "#1a1d27", borderRadius: 999, overflow: "hidden", marginBottom: 8 }}>
                  <div
                    style={{
                      width: `${Math.min(progress, 100)}%`,
                      height: "100%",
                      background: color,
                      borderRadius: 999,
                      transition: "width 0.3s ease",
                    }}
                  />
                </div>

                {/* Bottom: Saved / Target */}
                <div style={{ display: "flex", justifyContent: "space-between" }}>
                  <p style={{ color: "#6b7280", fontSize: 11, margin: 0 }}>
                    {fmt(goal.saved_amount)} / {fmt(goal.target_amount)}
                  </p>
                  <p style={{ color: "#6b7280", fontSize: 11, margin: 0 }}>
                    {progress.toFixed(0)}%
                  </p>
                </div>
              </div>
            );
          })}
        </div>

        {/* Allocation Buttons - Only show if goals have allocation % */}
        {hasAllocatedGoals && (
          <div style={{ display: "flex", gap: 10, paddingTop: 14, borderTop: "0.5px solid #1e2130" }}>
            <button
              onClick={handleAutoAllocate}
              style={{
                ...allocBtn,
                background: "#1D9E75",
                borderColor: "#1D9E75",
              }}
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8" />
                <polyline points="16 6 12 2 8 6" />
                <line x1="12" y1="2" x2="12" y2="15" />
              </svg>
              <span>Auto Allocate</span>
            </button>

            <button
              onClick={handleManualAllocate}
              style={{
                ...allocBtn,
                background: "#3b82f6",
                borderColor: "#3b82f6",
              }}
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14" />
                <path d="M5 12h14" />
              </svg>
              <span>Manual</span>
            </button>
          </div>
        )}
      </div>

      {/* Goals Modal - for creating/editing goals */}
      {showGoalsModal && (
        <GoalsModal
          existingGoals={activeGoals}
          onClose={() => setShowGoalsModal(false)}
          onSuccess={() => {
            setShowGoalsModal(false);
            onRefresh();
          }}
        />
      )}

      {/* Allocation Modal - for auto and manual allocation */}
      {showAllocationModal && allocationType && (
        <AllocationModal
          goals={activeGoals.filter((g) => g.allocation_percentage > 0)}
          allocationType={allocationType}
          onClose={() => {
            setShowAllocationModal(false);
            setAllocationType(null);
          }}
          onSuccess={() => {
            setShowAllocationModal(false);
            setAllocationType(null);
            onRefresh();
          }}
        />
      )}
    </>
  );
}

const card: React.CSSProperties = {
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 14,
  padding: 16,
  marginBottom: 12,
};

const sectionLabel: React.CSSProperties = {
  fontSize: 11,
  fontWeight: 600,
  letterSpacing: "0.08em",
  textTransform: "uppercase",
  color: "#6b7280",
};

const iconBtn: React.CSSProperties = {
  width: 30,
  height: 30,
  borderRadius: 8,
  border: "0.5px solid #2a2d3a",
  background: "#1a1d27",
  color: "#9ca3af",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  cursor: "pointer",
};

const summaryBox: React.CSSProperties = {
  background: "#161922",
  border: "0.5px solid",
  borderRadius: 10,
  padding: "10px 12px",
};

const goalRow: React.CSSProperties = {
  background: "#161922",
  border: "0.5px solid",
  borderRadius: 10,
  padding: "10px 12px",
};

const allocBtn: React.CSSProperties = {
  flex: 1,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  gap: 8,
  padding: "12px 14px",
  border: "0.5px solid",
  borderRadius: 10,
  color: "#fff",
  fontSize: 13,
  fontWeight: 600,
  cursor: "pointer",
  fontFamily: "inherit",
  transition: "opacity 0.2s ease",
};
