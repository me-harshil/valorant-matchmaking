import { useQuery } from "@tanstack/react-query";
import { fetchPlayers } from "../lib/api";

export function usePlayers() {
    return useQuery({
        queryKey: ["players"],
        queryFn: fetchPlayers,
        refetchInterval: 3000, // poll every 3s, matches the simulator's queueing cadence
    });
}