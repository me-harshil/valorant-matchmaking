import { useAllMatches } from "@/hooks/useDashboardData";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";
import { Badge } from "./ui/badge";
import { ScrollArea } from "./ui/scroll-area";

const STATUS_VARIANT: Record<string, string> = {
  pending: "bg-amber-500/15 text-amber-300 border-amber-500/30",
  in_progress: "bg-sky-500/15 text-sky-300 border-sky-500/30",
  completed: "bg-emerald-500/15 text-emerald-300 border-emerald-500/30",
  cancelled: "bg-red-500/15 text-red-300 border-red-500/30",
};

export function RecentMatches() {
  const { data: matches, isLoading, isError } = useAllMatches();

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader>
        <CardTitle className="text-sm font-semibold text-zinc-200">
          Match History
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && <div className="text-sm text-zinc-400">Loading…</div>}
        {isError && (
          <div className="text-sm text-red-400">Failed to load matches.</div>
        )}
        {matches && matches.length === 0 && (
          <div className="text-sm text-zinc-500">No matches yet.</div>
        )}
        {matches && matches.length > 0 && (
          <ScrollArea className="h-[380px] pr-3">
            <div className="space-y-2">
              {matches.map((m) => (
                <div
                  key={m.match_id}
                  className="flex items-center justify-between rounded-md border border-zinc-800 bg-zinc-900 px-3 py-2"
                >
                  <span className="text-sm font-medium text-zinc-200">
                    {m.map}
                  </span>
                  <Badge
                    variant="outline"
                    className={
                      STATUS_VARIANT[m.status] ?? "bg-zinc-700 text-zinc-200"
                    }
                  >
                    {m.status}
                  </Badge>
                  <span className="font-mono text-xs text-zinc-500">
                    {m.match_id.slice(0, 8)}
                  </span>
                </div>
              ))}
            </div>
          </ScrollArea>
        )}
      </CardContent>
    </Card>
  );
}
