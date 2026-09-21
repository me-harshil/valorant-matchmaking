import { usePlayers } from "@/hooks/useDashboardData";
import { conservativeRating, tierFromRating } from "@/lib/api";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "./ui/table";
import { Badge } from "./ui/badge";
import { ScrollArea } from "./ui/scroll-area";

const TIER_BADGE_CLASS: Record<string, string> = {
  Iron: "bg-zinc-700 text-zinc-100 border-zinc-600",
  Bronze: "bg-amber-900 text-amber-100 border-amber-700",
  Silver: "bg-slate-400 text-slate-900 border-slate-300",
  Gold: "bg-yellow-400 text-yellow-950 border-yellow-300",
  Platinum: "bg-teal-400 text-teal-950 border-teal-300",
  Diamond: "bg-sky-400 text-sky-950 border-sky-300",
  Ascendant: "bg-emerald-400 text-emerald-950 border-emerald-300",
  Immortal: "bg-fuchsia-500 text-fuchsia-50 border-fuchsia-400",
  Radiant:
    "bg-gradient-to-r from-yellow-200 via-white to-yellow-200 text-yellow-950 border-yellow-300",
};

export function Leaderboard() {
  const { data: players, isLoading, isError } = usePlayers();

  const ranked = players
    ? [...players]
        .map((p) => ({ ...p, rating: conservativeRating(p.mu, p.sigma) }))
        .sort((a, b) => b.rating - a.rating)
    : [];

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader>
        <CardTitle className="text-sm font-semibold text-zinc-200">
          Leaderboard
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && (
          <div className="text-sm text-zinc-400">Loading players…</div>
        )}
        {isError && (
          <div className="text-sm text-red-400">Failed to load players.</div>
        )}
        {players && (
          <ScrollArea className="h-[520px] pr-3">
            <Table>
              <TableHeader>
                <TableRow className="border-zinc-800 hover:bg-transparent">
                  <TableHead className="w-10 text-zinc-500">#</TableHead>
                  <TableHead className="text-zinc-500">Player</TableHead>
                  <TableHead className="text-zinc-500">Tier</TableHead>
                  <TableHead className="text-right text-zinc-500">
                    Rating
                  </TableHead>
                  <TableHead className="text-right text-zinc-500">
                    μ / σ
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {ranked.map((p, i) => {
                  const tier = tierFromRating(p.rating);
                  const badgeClass =
                    TIER_BADGE_CLASS[tier.name] ?? "bg-zinc-700 text-zinc-100";
                  return (
                    <TableRow key={p.player_id} className="border-zinc-900">
                      <TableCell className="text-zinc-500">{i + 1}</TableCell>
                      <TableCell className="font-medium text-zinc-100">
                        {p.username}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant="outline"
                          className={`font-semibold ${badgeClass}`}
                        >
                          {tier.name}
                          {tier.sub > 0 ? ` ${tier.sub}` : ""}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right tabular-nums text-zinc-200">
                        {p.rating.toFixed(1)}
                      </TableCell>
                      <TableCell className="text-right tabular-nums text-zinc-500">
                        {p.mu.toFixed(1)} / {p.sigma.toFixed(2)}
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </ScrollArea>
        )}
      </CardContent>
    </Card>
  );
}
