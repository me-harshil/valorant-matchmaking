import { usePendingMatches } from "@/hooks/useDashboardData";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";
import { Badge } from "./ui/badge";

export function PendingMatches() {
  const { data: matches, isLoading, isError } = usePendingMatches();

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-sm font-semibold text-zinc-200">
          Live Matches
        </CardTitle>
        <span className="relative flex h-2.5 w-2.5">
          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
          <span className="relative inline-flex h-2.5 w-2.5 rounded-full bg-emerald-500" />
        </span>
      </CardHeader>
      <CardContent className="space-y-2">
        {isLoading && <div className="text-sm text-zinc-400">Loading…</div>}
        {isError && (
          <div className="text-sm text-red-400">Failed to load matches.</div>
        )}
        {matches && matches.length === 0 && (
          <div className="text-sm text-zinc-500">
            No matches in progress right now.
          </div>
        )}
        {matches?.map((m) => (
          <div
            key={m.match_id}
            className="flex items-center justify-between rounded-md border border-zinc-800 bg-zinc-900 px-3 py-2"
          >
            <span className="text-sm font-medium text-zinc-200">{m.map}</span>
            <Badge variant="secondary" className="font-mono text-xs">
              {m.match_id.slice(0, 8)}
            </Badge>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
