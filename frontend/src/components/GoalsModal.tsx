import { useState, useEffect } from "react";

const API_BASE = import.meta.env.VITE_API_URL;

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

interface GoalItem {
  id: number | null; // null = new goal
  name: string;
  target_amount: number;
  allocation_percentage: number;
  deadline: string | null;
}

interface GoalsModalProps {
  existingGoals: Goal[];
  onClose: () => void;
  onSuccess: () => void;
}

export default function GoalsModal({
  existingGoals,
  onClose,
  onSuccess,
}: GoalsModalProps) {
  const isSimpleMode = existingGoals.length === 0;

  // Simple mode: just new goal
  const [simpleName, setSimpleName] = useState("");
  const [simpleTarget, setSimpleTarget] = useState("");
  const [simpleDeadline, setSimpleDeadline] = useState("");
  const [simplePercentage ] = useState("100");

  // Rebalance mode: existing + new goal
  const [allGoals, setAllGoals] = useState<GoalItem[]>([]);
  const [newGoalName, setNewGoalName] = useState("");
  const [newGoalTarget, setNewGoalTarget] = useState("");
  const [newGoalDeadline, setNewGoalDeadline] = useState("");
  const [newGoalPercentage, setNewGoalPercentage] = useState("0");

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Initialize rebalance mode goals
  useEffect(() => {
    if (!isSimpleMode) {
      setAllGoals(
        existingGoals.map((g) => ({
          id: g.id,
          name: g.name,
          target_amount: g.target_amount,
          allocation_percentage: g.allocation_percentage,
          deadline: g.deadline,
        }))
      );
    }
  }, [isSimpleMode, existingGoals]);

  // Calculate totals
  const simpleTotal = parseInt(simplePercentage) || 0;
  const rebalanceTotal = allGoals.reduce((sum, g) => sum + g.allocation_percentage, 0) + (parseInt(newGoalPercentage) || 0);

  const canSubmitSimple =
    simpleName.trim() &&
    simpleTarget &&
    parseInt(simpleTarget) > 0 &&
    simpleTotal === 100;

  const canSubmitRebalance =
    newGoalName.trim() &&
    newGoalTarget &&
    parseInt(newGoalTarget) > 0 &&
    rebalanceTotal === 100 &&
    allGoals.every((g) => g.name.trim() && g.target_amount > 0);

  // Auto-balance: distribute equally among all goals
  const handleAutoBalance = () => {
    const totalGoals = allGoals.length + 1; // existing + new
    const percentPerGoal = Math.floor(100 / totalGoals);
    const remainder = 100 % totalGoals;

    const updated = allGoals.map((g, idx) => ({
      ...g,
      allocation_percentage: percentPerGoal + (idx === 0 ? remainder : 0),
    }));

    setAllGoals(updated);
    setNewGoalPercentage(String(percentPerGoal));
  };

  // Handle simple creation
  const handleSimpleCreate = async () => {
    setSubmitting(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE}/goals`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: simpleName.trim(),
          target_amount: Math.round(parseFloat(simpleTarget) * 100),
          allocation_percentage: parseInt(simplePercentage),
          deadline: simpleDeadline || null,
        }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Server error: ${res.status}`);
      }

      onSuccess();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setSubmitting(false);
    }
  };

  // Handle rebalance creation
  const handleRebalanceCreate = async () => {
    setSubmitting(true);
    setError(null);

    try {
      const goalsPayload = [
        ...allGoals.map((g) => ({
          id: g.id,
          name: g.name.trim(),
          target_amount: Math.round(g.target_amount * 100),
          allocation_percentage: g.allocation_percentage,
          deadline: g.deadline || null,
        })),
        {
          id: null,
          name: newGoalName.trim(),
          target_amount: Math.round(parseFloat(newGoalTarget) * 100),
          allocation_percentage: parseInt(newGoalPercentage),
          deadline: newGoalDeadline || null,
        },
      ];

      const res = await fetch(`${API_BASE}/goals/rebalance`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ goals: goalsPayload }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Server error: ${res.status}`);
      }

      onSuccess();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      {/* Backdrop */}
      <div
        onClick={onClose}
        style={{
          position: "fixed",
          inset: 0,
          background: "rgba(0,0,0,0.6)",
          zIndex: 100,
          backdropFilter: "blur(2px)",
        }}
      />

      {/* Modal */}
      <div
        style={{
          position: "fixed",
          bottom: 0,
          left: 0,
          right: 0,
          background: "#0f1117",
          borderTop: "0.5px solid #1e2130",
          borderRadius: "20px 20px 0 0",
          maxHeight: "90vh",
          overflow: "auto",
          zIndex: 101,
          padding: "24px 20px 40px",
        }}
      >
        {/* Drag handle */}
        <div
          style={{
            width: 36,
            height: 4,
            borderRadius: 2,
            background: "#2a2d3a",
            margin: "0 auto 20px",
          }}
        />

        <h2
          style={{
            fontSize: 18,
            fontWeight: 700,
            color: "#f3f4f6",
            margin: "0 0 20px",
          }}
        >
          {isSimpleMode ? "Create Your First Goal" : "Add New Goal & Rebalance"}
        </h2>

        {/* Error */}
        {error && (
          <div
            style={{
              padding: "10px 12px",
              background: "#3d1515",
              border: "0.5px solid #E24B4A",
              borderRadius: 10,
              color: "#E24B4A",
              fontSize: 12,
              marginBottom: 16,
            }}
          >
            {error}
          </div>
        )}

        {/* ═══ SIMPLE MODE ═══ */}
        {isSimpleMode ? (
          <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
            {/* Name */}
            <div>
              <label style={labelStyle}>Goal Name *</label>
              <input
                type="text"
                value={simpleName}
                onChange={(e) => setSimpleName(e.target.value)}
                placeholder="e.g. Emergency Fund"
                style={inputStyle}
              />
            </div>

            {/* Target Amount */}
            <div>
              <label style={labelStyle}>Target Amount (₹) *</label>
              <input
                type="number"
                step="0.01"
                min="0"
                value={simpleTarget}
                onChange={(e) => setSimpleTarget(e.target.value)}
                placeholder="100000"
                style={inputStyle}
              />
            </div>

            {/* Deadline */}
            <div>
              <label style={labelStyle}>Deadline (Optional)</label>
              <input
                type="date"
                value={simpleDeadline}
                onChange={(e) => setSimpleDeadline(e.target.value)}
                style={inputStyle}
              />
            </div>

            {/* Percentage (always 100 in simple mode) */}
            <div>
              <label style={labelStyle}>Allocation %</label>
              <div
                style={{
                  padding: "12px 14px",
                  background: "#161922",
                  border: "0.5px solid #1e2130",
                  borderRadius: 10,
                  color: "#d1d5db",
                  fontSize: 14,
                }}
              >
                100% (First goal gets all allocation)
              </div>
            </div>

            {/* Submit */}
            <button
              onClick={handleSimpleCreate}
              disabled={!canSubmitSimple || submitting}
              style={{
                ...submitBtn,
                opacity: canSubmitSimple ? 1 : 0.5,
                cursor: canSubmitSimple ? "pointer" : "not-allowed",
              }}
            >
              {submitting ? "Creating…" : "Create Goal"}
            </button>
          </div>
        ) : (
          /* ═══ REBALANCE MODE ═══ */
          <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
            {/* Existing Goals */}
            <div>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
                <label style={labelStyle}>Existing Goals</label>
                <button
                  onClick={handleAutoBalance}
                  style={{
                    fontSize: 11,
                    padding: "4px 8px",
                    background: "#1D9E75",
                    border: "none",
                    borderRadius: 6,
                    color: "#fff",
                    fontWeight: 600,
                    cursor: "pointer",
                  }}
                >
                  Auto-Balance
                </button>
              </div>

              <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
                {allGoals.map((goal, idx) => (
                  <div
                    key={goal.id}
                    style={{
                      background: "#161922",
                      border: "0.5px solid #1e2130",
                      borderRadius: 10,
                      padding: 12,
                    }}
                  >
                    {/* Goal Name + Target */}
                    <div style={{ marginBottom: 10 }}>
                      <p
                        style={{
                          fontSize: 12,
                          fontWeight: 600,
                          color: "#d1d5db",
                          margin: "0 0 6px",
                        }}
                      >
                        {goal.name}
                      </p>
                      <p
                        style={{
                          fontSize: 11,
                          color: "#6b7280",
                          margin: 0,
                        }}
                      >
                        Target: ₹{(goal.target_amount / 100).toLocaleString("en-IN")}
                      </p>
                    </div>

                    {/* Allocation Slider */}
                    <div>
                      <div
                        style={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                          marginBottom: 6,
                        }}
                      >
                        <label style={{ fontSize: 11, color: "#9ca3af" }}>
                          Allocation %
                        </label>
                        <input
                          type="number"
                          min="0"
                          max="100"
                          value={goal.allocation_percentage}
                          onChange={(e) => {
                            const newVal = parseInt(e.target.value) || 0;
                            setAllGoals((prev) =>
                              prev.map((g, i) =>
                                i === idx
                                  ? { ...g, allocation_percentage: newVal }
                                  : g
                              )
                            );
                          }}
                          style={{
                            width: 50,
                            padding: "4px 6px",
                            background: "#0f1117",
                            border: "0.5px solid #1e2130",
                            borderRadius: 6,
                            color: "#d1d5db",
                            fontSize: 12,
                            textAlign: "center",
                          }}
                        />
                      </div>

                      {/* Slider */}
                      <input
                        type="range"
                        min="0"
                        max="100"
                        value={goal.allocation_percentage}
                        onChange={(e) => {
                          const newVal = parseInt(e.target.value) || 0;
                          setAllGoals((prev) =>
                            prev.map((g, i) =>
                              i === idx
                                ? { ...g, allocation_percentage: newVal }
                                : g
                            )
                          );
                        }}
                        style={{
                          width: "100%",
                          height: 4,
                          borderRadius: 999,
                          background: "#0f1117",
                          outline: "none",
                          cursor: "pointer",
                        }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* New Goal */}
            <div>
              <label style={labelStyle}>New Goal *</label>
              <div
                style={{
                  background: "#161922",
                  border: "0.5px solid #1D9E75",
                  borderRadius: 10,
                  padding: 12,
                  display: "flex",
                  flexDirection: "column",
                  gap: 10,
                }}
              >
                <input
                  type="text"
                  value={newGoalName}
                  onChange={(e) => setNewGoalName(e.target.value)}
                  placeholder="Goal name"
                  style={{
                    ...inputStyle,
                    background: "#0f1117",
                    fontSize: 13,
                  }}
                />

                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={newGoalTarget}
                  onChange={(e) => setNewGoalTarget(e.target.value)}
                  placeholder="Target amount"
                  style={{
                    ...inputStyle,
                    background: "#0f1117",
                    fontSize: 13,
                  }}
                />

                <input
                  type="date"
                  value={newGoalDeadline}
                  onChange={(e) => setNewGoalDeadline(e.target.value)}
                  style={{
                    ...inputStyle,
                    background: "#0f1117",
                    fontSize: 13,
                  }}
                />

                {/* Allocation Slider */}
                <div>
                  <div
                    style={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center",
                      marginBottom: 6,
                    }}
                  >
                    <label style={{ fontSize: 11, color: "#9ca3af" }}>
                      Allocation %
                    </label>
                    <input
                      type="number"
                      min="0"
                      max="100"
                      value={newGoalPercentage}
                      onChange={(e) =>
                        setNewGoalPercentage(e.target.value)
                      }
                      style={{
                        width: 50,
                        padding: "4px 6px",
                        background: "#161922",
                        border: "0.5px solid #1e2130",
                        borderRadius: 6,
                        color: "#d1d5db",
                        fontSize: 12,
                        textAlign: "center",
                      }}
                    />
                  </div>

                  <input
                    type="range"
                    min="0"
                    max="100"
                    value={newGoalPercentage}
                    onChange={(e) =>
                      setNewGoalPercentage(e.target.value)
                    }
                    style={{
                      width: "100%",
                      height: 4,
                      borderRadius: 999,
                      background: "#0f1117",
                      outline: "none",
                      cursor: "pointer",
                    }}
                  />
                </div>
              </div>
            </div>

            {/* Total % Indicator */}
            <div
              style={{
                padding: "12px",
                background: rebalanceTotal === 100 ? "#1a3a2a" : "#3a2a1a",
                border: `0.5px solid ${rebalanceTotal === 100 ? "#1D9E75" : "#f59e0b"}`,
                borderRadius: 10,
                textAlign: "center",
              }}
            >
              <p
                style={{
                  fontSize: 12,
                  fontWeight: 600,
                  color:
                    rebalanceTotal === 100 ? "#1D9E75" : "#f59e0b",
                  margin: 0,
                }}
              >
                Total Allocation: {rebalanceTotal}%
              </p>
              {rebalanceTotal !== 100 && (
                <p
                  style={{
                    fontSize: 11,
                    color: "#f59e0b",
                    margin: "4px 0 0",
                  }}
                >
                  Must equal 100% to continue
                </p>
              )}
            </div>

            {/* Submit */}
            <button
              onClick={handleRebalanceCreate}
              disabled={!canSubmitRebalance || submitting}
              style={{
                ...submitBtn,
                opacity: canSubmitRebalance ? 1 : 0.5,
                cursor: canSubmitRebalance ? "pointer" : "not-allowed",
              }}
            >
              {submitting ? "Creating…" : "Create & Rebalance"}
            </button>
          </div>
        )}

        {/* Cancel */}
        <button
          onClick={onClose}
          disabled={submitting}
          style={{
            width: "100%",
            marginTop: 12,
            padding: "12px 0",
            background: "#1a1d27",
            border: "0.5px solid #2a2d3a",
            borderRadius: 10,
            color: "#9ca3af",
            fontSize: 14,
            fontWeight: 600,
            cursor: "pointer",
          }}
        >
          Cancel
        </button>
      </div>
    </>
  );
}

const labelStyle: React.CSSProperties = {
  display: "block",
  fontSize: 11,
  fontWeight: 600,
  letterSpacing: "0.07em",
  textTransform: "uppercase",
  color: "#6b7280",
  marginBottom: 8,
};

const inputStyle: React.CSSProperties = {
  width: "100%",
  background: "#161922",
  border: "0.5px solid #1e2130",
  borderRadius: 10,
  padding: "12px 14px",
  color: "#d1d5db",
  fontSize: 14,
  outline: "none",
  fontFamily: "inherit",
};

const submitBtn: React.CSSProperties = {
  width: "100%",
  padding: "13px 0",
  borderRadius: 10,
  border: "none",
  background: "#1D9E75",
  color: "#fff",
  fontSize: 14,
  fontWeight: 600,
  cursor: "pointer",
  fontFamily: "inherit",
};
