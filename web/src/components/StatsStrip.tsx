import { usePlayers, usePendingMatches } from "@/hooks/useDashboardData";
import { Card, CardContent } from "./ui/card";

function StatCard({ label, value }: { label: string; value: string | number }) {
  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardContent className="px-5 py-4">
        <div className="text-xs font-medium uppercase tracking-wide text-zinc-500">
          {label}
        </div>
        <div className="mt-1 text-2xl font-semibold text-zinc-100">{value}</div>
      </CardContent>
    </Card>
  );
}

export function StatsStrip() {
  const { data: players } = usePlayers();
  const { data: pendingMatches } = usePendingMatches();

  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
      <StatCard label="Total Players" value={players?.length ?? "—"} />
      <StatCard label="Live Matches" value={pendingMatches?.length ?? "—"} />
      <StatCard
        label="Players In-Match"
        value={pendingMatches ? pendingMatches.length * 10 : "—"}
      />
      <StatCard label="Status" value="Simulating" />
    </div>
  );
}
