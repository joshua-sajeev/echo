import { useEffect, useState, useCallback } from "react";
import { getDashboard } from "../api/dashboard";

export function useDashboard() {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const dashboard = await getDashboard();
      setData(dashboard);
    } catch (err) {
      console.error("[useDashboard] failed to fetch dashboard:", err);
      setError("Failed to load dashboard");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return {
    data,
    loading,
    error,
    refresh,
  };
}
