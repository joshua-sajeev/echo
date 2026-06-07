import { useEffect, useState, useCallback } from "react";
import { getDashboard } from "../api/dashboard";

export function useDashboard() {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    setLoading(true);

    try {
      const dashboard = await getDashboard();


      setData(dashboard);
    } catch (err) {
      console.error("[useDashboard] failed to fetch dashboard:", err);
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
    refresh,
  };
}
