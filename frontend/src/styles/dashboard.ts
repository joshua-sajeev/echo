import React from "react";

export const statCard: React.CSSProperties = {
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 14,
  padding: "14px 12px",
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-between"
};

export const actionBtn: React.CSSProperties = {
  flex: 1,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  gap: 8,
  padding: "12px 14px",
  background: "#0f1117",
  border: "0.5px solid #1e2130",
  borderRadius: 12,
  color: "#e8eaf0",
  fontSize: 13,
  fontWeight: 500,
  cursor: "pointer",
  transition: "background 0.2s ease",
};
