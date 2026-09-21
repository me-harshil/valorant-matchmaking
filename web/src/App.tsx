import { Leaderboard } from "@/components/Leaderboard";
import { PendingMatches } from "@/components/PendingMatches";
import { RecentMatches } from "@/components/RecentMatches";
import { StatsStrip } from "@/components/StatsStrip";

function App() {
  return (
    <div className="min-h-screen bg-zinc-900 px-6 py-8 text-zinc-100">
      <div className="mx-auto max-w-6xl space-y-6">
        <header>
          <h1 className="text-xl font-bold tracking-tight">
            Valorant Matchmaking — Live Dashboard
          </h1>
          <p className="mt-1 text-sm text-zinc-500">
            Real-time view of the matchmaking, match, and rating services.
          </p>
        </header>

        <StatsStrip />

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <div className="lg:col-span-2">
            <Leaderboard />
          </div>
          <div className="space-y-6">
            <PendingMatches />
            <RecentMatches />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
