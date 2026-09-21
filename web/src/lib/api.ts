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
    return res.data;
}

export async function fetchMatches(status: string = "pending"): Promise<MatchSummary[]> {
    const res = await api.get<MatchSummary[]>("/matches", { params: { status } });
    return res.data;
}

export async function fetchParticipants(matchID: string): Promise<Participant[]> {
    const res = await api.get<Participant[]>(`/matches/${matchID}/participants`);
    return res.data;
}