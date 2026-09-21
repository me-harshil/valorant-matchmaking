import { useQuery } from "@tanstack/react-query";
import { fetchPlayers, fetchAllMatches, fetchMatches } from "@/lib/api";

export function usePlayers() {
    return useQuery({
        queryKey: ["players"],
        queryFn: fetchPlayers,
        refetchInterval: 3000,
    });
}

export function usePendingMatches() {
    return useQuery({
        queryKey: ["matches", "pending"],
        queryFn: () => fetchMatches("pending"),
        refetchInterval: 2000,
    });
}

export function useAllMatches() {
    return useQuery({
        queryKey: ["matches", "all"],
        queryFn: fetchAllMatches,
        refetchInterval: 3000,
    });
}