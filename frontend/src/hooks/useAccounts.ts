import { useEffect, useState } from "react";
import { getAccounts } from "../api/accounts";
import { Account } from "../types/account";

export function useAccounts() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getAccounts()
      .then(setAccounts)
      .finally(() => setLoading(false));
  }, []);

  return { accounts, loading };
}
