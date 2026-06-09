interface DeleteModalProps {
  open: boolean;
  title?: string;
  itemName?: string;
  loading?: boolean;
  onClose: () => void;
  onConfirm: () => void;
}

export default function DeleteModal({
  open,
  title = "Delete Item",
  itemName,
  loading = false,
  onClose,
  onConfirm,
}: DeleteModalProps) {
  if (!open) return null;

  return (
    <>
      <div
        className="fixed inset-0 bg-black/60 z-50"
        onClick={onClose}
      />

      <div className="fixed bottom-0 left-0 right-0 bg-[#0f1117] border-t border-[#1e2130] rounded-t-2xl p-6 pb-10 z-50">
        <div className="w-9 h-1 bg-[#2a2d3a] rounded-full mx-auto mb-5" />

        <div className="flex flex-col items-center justify-center gap-1.5 mb-3 text-red-400">
          <svg
            className="w-6 h-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth="2.5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
            />
          </svg>

          <p className="text-sm font-bold uppercase tracking-wider">
            {title}
          </p>
        </div>

        <p className="text-xs text-zinc-400 text-center mb-6">
          Are you sure you want to delete{" "}
          <span className="text-zinc-200 font-medium">
            "{itemName}"
          </span>
          ?
        </p>

        <div className="flex gap-3">
          <button
            onClick={onClose}
            disabled={loading}
            className="flex-1 py-3 bg-[#1a1d27] text-zinc-400 font-semibold text-sm rounded-lg"
          >
            Cancel
          </button>

          <button
            onClick={onConfirm}
            disabled={loading}
            className="flex-1 py-3 bg-red-600 text-white font-semibold text-sm rounded-lg"
          >
            {loading ? "Deleting..." : "Delete"}
          </button>
        </div>
      </div>
    </>
  );
}
