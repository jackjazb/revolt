import type { State } from "./types";

export function randomName(): string {
    const adjectives = ['Cool', 'Angry', 'Purple', 'Red', 'Awful', 'Cool'];
    const nouns = ['Buffalo', 'Bean', 'Giraffe', 'Carrot', 'Loaf'];
    const adjective = adjectives[Math.floor(Math.random() * (adjectives.length - 1))];
    const noun = nouns[Math.floor(Math.random() * (nouns.length - 1))];
    return `${adjective}${noun}`;
}

export function getPlayerById(state: State, id: string): string | undefined {
    if (id === state.self.id) {
        return state.self.name;
    }
    return state.peers.filter(p => p.id === id).at(0)?.name;
}