import { useState, useEffect } from "react";

const API_BASE = import.meta.env.VITE_API_URL;

interface Goal {
  id: number;
  name: string;
  allocation_percentage: number;
  saved_amount: number;
  target_amount: number;
}

interface AllocationModalProps {
  goals: Goal[];
  allocationType: "manual" | "auto";
  onClose: () => void;
  onSuccess: () => void;
}

interface AllocationPreview {
  goalId: number;
  name: string;
  percentage: number;
  amount: number;
}

export default function AllocationModal({
  goals,
  allocationType,
  onClose,
  onSuccess,
}: AllocationModalProps) {
  const [amount, setAmount] = useState("");
  const [selectedGoalId, setSelectedGoalId] = useState<number | null>(
    allocationType === "manual" ? goals[0]?.id || null : null
  );
  const [preview, setPreview] = useState<AllocationPreview[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // "custom" for entering amount, "leisure_leftover" for last month's leisure leftover
  const [autoMode, setAutoMode] = useState<"custom" | "leisure_leftover">("custom");
  const [leisureLeftover, setLeisureLeftover] = useState<number | null>(null);
  const [alreadyAllocated, setAlreadyAllocated] = useState(false);

  // Fetch leisure leftover
  useEffect(() => {
    if (allocationType === "auto" && autoMode === "leisure_leftover") {
      fetch(`${API_BASE}/allocation/leisure-leftover`, { credentials: "include" })
        .then((res) => {
          if (!res.ok) throw new Error("Failed to fetch leftover");
          return res.json();
        })
        .then((data) => {
          setLeisureLeftover(data.amount);
          setAlreadyAllocated(!!data.already_allocated);
          setAmount((data.amount / 100).toString());
        })
        .catch((err) => {
          setError(err.message || "Failed to fetch last month's leisure leftover");
        });
    }
  }, [autoMode, allocationType]);

  // Calculate preview for auto allocation
  useEffect(() => {
    if (allocationType === "auto" && amount) {
      const totalAmount = Math.round(parseFloat(amount) * 100);
      const previews: AllocationPreview[] = [];
      let allocated = 0;

      goals.forEach((goal, idx) => {
        let alloc: number;
        if (idx === goals.length - 1) {
          // Last goal gets remainder
          alloc = totalAmount - allocated;
        } else {
          alloc = Math.floor((totalAmount * goal.allocation_percentage) / 100);
          allocated += alloc;
        }

        previews.push({
          goalId: goal.id,
          name: goal.name,
          percentage: goal.allocation_percentage,
          amount: alloc,
        });
      });

      setPreview(previews);
    }
  }, [amount, allocationType, goals]);

  const canSubmit =
    allocationType === "auto"
      ? amount && parseFloat(amount) > 0 && preview.length > 0 && !(autoMode === "leisure_leftover" && alreadyAllocated)
      : amount && parseFloat(amount) > 0 && selectedGoalId;

  const handleSubmit = async () => {
    setSubmitting(true);
    setError(null);

    try {
      const totalAmount = Math.round(parseFloat(amount) * 100);

      if (allocationType === "auto") {
        const payload: { type: string; amount?: number } = {
          type: autoMode === "custom" ? "automatic_splitting" : "leisure_leftover",
        };
        if (autoMode === "custom") {
          payload.amount = totalAmount;
        }

        const res = await fetch(`${API_BASE}/allocation/distribute`, {
          method: "POST",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });

        if (!res.ok) {
          const data = await res.json();
          throw new Error(data.error?.message || `Server error: ${res.status}`);
        }
      } else {
        const res = await fetch(`${API_BASE}/allocation/run`, {
          method: "POST",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            goal_id: selectedGoalId,
            amount: totalAmount,
          }),
        });

        if (!res.ok) {
          const data = await res.json();
          throw new Error(data.error?.message || `Server error: ${res.status}`);
        }
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
            margin: "0 0 8px",
          }}
        >
          {allocationType === "auto" ? "Auto Allocate" : "Manual Allocation"}
        </h2>

        <p
          style={{
            fontSize: 13,
            color: "#6b7280",
            margin: "0 0 20px",
            lineHeight: 1.5,
          }}
        >
          {allocationType === "auto"
            ? "Automatically distribute funds to goals based on their allocation percentages. Carryover from previous months will be proportionally distributed."
            : "Allocate funds to a specific goal. Useful for targeted savings toward one objective."}
        </p>

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

        {allocationType === "auto" ? (
          /* ═══ AUTO ALLOCATION ═══ */
          <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
            {/* Auto Mode Tabs */}
            <div style={{ display: "flex", gap: 8, background: "#161922", padding: 4, borderRadius: 10, border: "0.5px solid #1e2130" }}>
              <button
                type="button"
                onClick={() => {
                  setAutoMode("custom");
                  setAmount("");
                  setPreview([]);
                }}
                style={{
                  flex: 1,
                  padding: "8px 12px",
                  background: autoMode === "custom" ? "#1D9E75" : "transparent",
                  color: autoMode === "custom" ? "#fff" : "#9ca3af",
                  border: "none",
                  borderRadius: 8,
                  cursor: "pointer",
                  fontSize: 12,
                  fontWeight: 600,
                  transition: "all 0.2s",
                }}
              >
                Custom Amount
              </button>
              <button
                type="button"
                onClick={() => {
                  setAutoMode("leisure_leftover");
                  setAmount("");
                  setPreview([]);
                }}
                style={{
                  flex: 1,
                  padding: "8px 12px",
                  background: autoMode === "leisure_leftover" ? "#1D9E75" : "transparent",
                  color: autoMode === "leisure_leftover" ? "#fff" : "#9ca3af",
                  border: "none",
                  borderRadius: 8,
                  cursor: "pointer",
                  fontSize: 12,
                  fontWeight: 600,
                  transition: "all 0.2s",
                }}
              >
                Leisure Leftover
              </button>
            </div>

            {/* Input / Info based on mode */}
            {autoMode === "custom" ? (
              <div>
                <label style={labelStyle}>Amount to Allocate (₹) *</label>
                <div style={{ position: "relative" }}>
                  <span
                    style={{
                      position: "absolute",
                      left: 14,
                      top: "50%",
                      transform: "translateY(-50%)",
                      color: "#6b7280",
                      fontSize: 14,
                      pointerEvents: "none",
                      fontFamily: "'IBM Plex Mono', monospace",
                    }}
                  >
                    ₹
                  </span>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    placeholder="0.00"
                    style={{
                      ...inputStyle,
                      paddingLeft: 30,
                      fontFamily: "'IBM Plex Mono', monospace",
                    }}
                  />
                </div>
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
                <div
                  style={{
                    background: "#161922",
                    border: alreadyAllocated ? "0.5px solid #2a2d3a" : "0.5px solid #1e2130",
                    borderRadius: 10,
                    padding: "14px 16px",
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    opacity: alreadyAllocated ? 0.6 : 1,
                  }}
                >
                  <div>
                    <p style={{ fontSize: 11, color: "#6b7280", textTransform: "uppercase", fontWeight: 600, margin: 0 }}>
                      Last Month Leisure Leftover
                    </p>
                    <p style={{ fontSize: 12, color: "#9ca3af", margin: "4px 0 0" }}>
                      {alreadyAllocated ? "Leisure leftover has already been allocated." : "Automatically calculated from last month's leftover"}
                    </p>
                  </div>
                  <strong
                    style={{
                      color: alreadyAllocated ? "#6b7280" : "#1D9E75",
                      fontSize: 16,
                      fontFamily: "'IBM Plex Mono', monospace",
                      textDecoration: alreadyAllocated ? "line-through" : "none",
                    }}
                  >
                    {leisureLeftover !== null ? `₹${(leisureLeftover / 100).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : "Loading..."}
                  </strong>
                </div>

                {alreadyAllocated && (
                  <div
                    style={{
                      background: "rgba(217, 119, 6, 0.1)",
                      border: "0.5px solid rgba(217, 119, 6, 0.4)",
                      borderRadius: 10,
                      padding: "10px 12px",
                      color: "#f59e0b",
                      fontSize: 12,
                      display: "flex",
                      alignItems: "center",
                      gap: 8,
                    }}
                  >
                    <span>⚠️</span>
                    <span>This allocation has already been executed for this month.</span>
                  </div>
                )}
              </div>
            )}

            {/* Preview */}
            {preview.length > 0 && (
              <div>
                <p style={labelStyle}>Distribution Preview</p>
                <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                  {preview.map((p) => (
                    <div
                      key={p.goalId}
                      style={{
                        background: "#161922",
                        border: "0.5px solid #1e2130",
                        borderRadius: 10,
                        padding: "10px 12px",
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                      }}
                    >
                      <div>
                        <p
                          style={{
                            fontSize: 12,
                            fontWeight: 600,
                            color: "#d1d5db",
                            margin: 0,
                          }}
                        >
                          {p.name}
                        </p>
                        <p style={{ fontSize: 11, color: "#6b7280", margin: "2px 0 0" }}>
                          {p.percentage}% allocation
                        </p>
                      </div>
                      <p
                        style={{
                          fontSize: 13,
                          fontWeight: 700,
                          color: "#3b82f6",
                          margin: 0,
                          fontFamily: "'IBM Plex Mono', monospace",
                        }}
                      >
                        ₹{(p.amount / 100).toLocaleString("en-IN", {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })}
                      </p>
                    </div>
                  ))}
                </div>

                {/* Total */}
                <div
                  style={{
                    marginTop: 12,
                    padding: "10px 12px",
                    background: "#1a2a2a",
                    border: "0.5px solid #1D9E75",
                    borderRadius: 10,
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                  }}
                >
                  <p
                    style={{
                      fontSize: 12,
                      fontWeight: 600,
                      color: "#1D9E75",
                      margin: 0,
                    }}
                  >
                    Total to Allocate
                  </p>
                  <p
                    style={{
                      fontSize: 13,
                      fontWeight: 700,
                      color: "#1D9E75",
                      margin: 0,
                      fontFamily: "'IBM Plex Mono', monospace",
                    }}
                  >
                    ₹{(preview.reduce((sum, p) => sum + p.amount, 0) / 100).toLocaleString("en-IN", {
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2,
                    })}
                  </p>
                </div>
              </div>
            )}

            {/* Submit */}
            <button
              onClick={handleSubmit}
              disabled={!canSubmit || submitting}
              style={{
                ...submitBtn,
                opacity: canSubmit ? 1 : 0.5,
                cursor: canSubmit ? "pointer" : "not-allowed",
              }}
            >
              {submitting ? "Allocating…" : "Allocate Funds"}
            </button>
          </div>
        ) : (
          /* ═══ MANUAL ALLOCATION ═══ */
          <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
            {/* Goal Selection */}
            <div>
              <label style={labelStyle}>Select Goal *</label>
              <select
                value={selectedGoalId || ""}
                onChange={(e) => setSelectedGoalId(parseInt(e.target.value))}
                style={{
                  ...inputStyle,
                  appearance: "none",
                  paddingRight: 36,
                }}
              >
                <option value="">Choose a goal</option>
                {goals.map((goal) => (
                  <option key={goal.id} value={goal.id}>
                    {goal.name} ({goal.allocation_percentage}% allocation)
                  </option>
                ))}
              </select>
              <svg
                style={{
                  position: "absolute",
                  right: 12,
                  top: "50%",
                  transform: "translateY(-50%)",
                  pointerEvents: "none",
                }}
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#4b5563"
                strokeWidth="2.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <polyline points="6 9 12 15 18 9" />
              </svg>
            </div>

            {/* Amount Input */}
            <div>
              <label style={labelStyle}>Amount to Allocate (₹) *</label>
              <div style={{ position: "relative" }}>
                <span
                  style={{
                    position: "absolute",
                    left: 14,
                    top: "50%",
                    transform: "translateY(-50%)",
                    color: "#6b7280",
                    fontSize: 14,
                    pointerEvents: "none",
                    fontFamily: "'IBM Plex Mono', monospace",
                  }}
                >
                  ₹
                </span>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  placeholder="0.00"
                  style={{
                    ...inputStyle,
                    paddingLeft: 30,
                    fontFamily: "'IBM Plex Mono', monospace",
                  }}
                />
              </div>
            </div>

            {/* Selected Goal Info */}
            {selectedGoalId && (
              <div
                style={{
                  background: "#161922",
                  border: "0.5px solid #3b82f6",
                  borderRadius: 10,
                  padding: "12px",
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                }}
              >
                <div>
                  <p
                    style={{
                      fontSize: 12,
                      fontWeight: 600,
                      color: "#d1d5db",
                      margin: 0,
                    }}
                  >
                    {goals.find((g) => g.id === selectedGoalId)?.name}
                  </p>
                  <p style={{ fontSize: 11, color: "#6b7280", margin: "2px 0 0" }}>
                    Allocation:{" "}
                    {goals.find((g) => g.id === selectedGoalId)?.allocation_percentage}%
                  </p>
                </div>
                {amount && (
                  <p
                    style={{
                      fontSize: 13,
                      fontWeight: 700,
                      color: "#3b82f6",
                      margin: 0,
                      fontFamily: "'IBM Plex Mono', monospace",
                    }}
                  >
                    +₹{parseFloat(amount).toLocaleString("en-IN", {
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2,
                    })}
                  </p>
                )}
              </div>
            )}

            {/* Submit */}
            <button
              onClick={handleSubmit}
              disabled={!canSubmit || submitting}
              style={{
                ...submitBtn,
                opacity: canSubmit ? 1 : 0.5,
                cursor: canSubmit ? "pointer" : "not-allowed",
              }}
            >
              {submitting ? "Allocating…" : "Allocate to Goal"}
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
            fontFamily: "inherit",
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
  position: "relative",
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
