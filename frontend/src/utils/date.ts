export const getCurrentMonthYear = () => {
  const now = new Date();

  return {
    month: now.toLocaleString("en-IN", { month: "long" }),
    year: now.getFullYear(),
  };
};
