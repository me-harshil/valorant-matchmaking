import axios from "axios";

const api = axios.create({
    baseURL: "http://localhost:8080",
});

export interface Player {
    player_id: string;
    username: string;
    mu: number;
    sigma: number;
}

export interface MatchSummary {
    match_id: string;
    status: string;
    map: string;
}

export interface Participant {
    participant_id: string;
    player_id: string;
    team: string;
}

export async function fetchPlayers(): Promise<Player[]> {
    const res = await api.get<Player[]>("/players");
    return res.data ?? [];
}

export async function fetchMatches(status: string = "pending"): Promise<MatchSummary[]> {
    const res = await api.get<MatchSummary[]>("/matches", { params: { status } });
    return res.data ?? [];
}

export async function fetchAllMatches(): Promise<MatchSummary[]> {
    const [pending, completed] = await Promise.all([
        fetchMatches("pending"),
        fetchMatches("completed"),
    ]);
    return [...pending, ...completed];
}

export async function fetchParticipants(matchID: string): Promise<Participant[]> {
    const res = await api.get<Participant[]>(`/matches/${matchID}/participants`);
    return res.data ?? [];
}

const TIER_BOUNDARIES: { min: number; name: string; sub: number }[] = [
    { min: 110, name: "Radiant", sub: 0 },
    { min: 105, name: "Immortal", sub: 3 },
    { min: 102, name: "Immortal", sub: 2 },
    { min: 99, name: "Immortal", sub: 1 },
    { min: 96, name: "Ascendant", sub: 3 },
    { min: 93, name: "Ascendant", sub: 2 },
    { min: 90, name: "Ascendant", sub: 1 },
    { min: 85, name: "Diamond", sub: 3 },
    { min: 80, name: "Diamond", sub: 2 },
    { min: 75, name: "Diamond", sub: 1 },
    { min: 70, name: "Platinum", sub: 3 },
    { min: 65, name: "Platinum", sub: 2 },
    { min: 60, name: "Platinum", sub: 1 },
    { min: 55, name: "Gold", sub: 3 },
    { min: 50, name: "Gold", sub: 2 },
    { min: 45, name: "Gold", sub: 1 },
    { min: 40, name: "Silver", sub: 3 },
    { min: 35, name: "Silver", sub: 2 },
    { min: 30, name: "Silver", sub: 1 },
    { min: 25, name: "Bronze", sub: 3 },
    { min: 20, name: "Bronze", sub: 2 },
    { min: 15, name: "Bronze", sub: 1 },
    { min: 10, name: "Iron", sub: 3 },
    { min: 5, name: "Iron", sub: 2 },
    { min: -100, name: "Iron", sub: 1 },
];

export function conservativeRating(mu: number, sigma: number): number {
    return mu - 3 * sigma;
}

export function tierFromRating(rating: number): { name: string; sub: number } {
    const match = TIER_BOUNDARIES.find((b) => rating >= b.min);
    return match ?? { name: "Iron", sub: 1 };
}