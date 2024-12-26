import { ActionType, Card, TurnState, type State } from "./types";

export function randomName(): string {
    const adjectives = ['Heated', 'Evil', 'Holy', 'Furious', 'Raging', 'Purple', 'Red', 'Dreaded'];
    const nouns = ['Buffalo', 'Otter', 'Giraffe', 'Carrot', 'Loaf', 'Tiger', 'Eagle', 'Snake'];
    const adjective = adjectives[Math.floor(Math.random() * (adjectives.length - 1))];
    const noun = nouns[Math.floor(Math.random() * (nouns.length - 1))];
    return `${adjective}${noun}`;
}

/**
 * Returns a player from `State` by their ID. 
 */
export function getPlayerById(state: State, id: string): string | undefined {
    if (id === state.self.id) {
        return state.self.name;
    }
    return state.peers.filter(p => p.id === id).at(0)?.name;
}

/**
 * Returns the name of the current game leader.
 */
export function getLeader(state: State): string {
    if (state.self.leading) {
        return state.self.name;
    }
    const leaders = state.peers.filter(p => p.leading);
    if (leaders.length === 1) {
        return leaders[0].name;
    }
    // Shouldn't be able to get here.
    return "";
}

/**
 * Returns true if `self` is allowed to perform `action` with their current cards.
 */
export function isAllowedAction(state: State, action: ActionType): boolean {
    if (!state.self.allowedActions) {
        return false;
    }
    return state.self.allowedActions.includes(action);
}


/**
 * Defines which cards block a given action.
 */
export const BLOCKED_BY: Partial<Record<ActionType, Card[]>> = {
    [ActionType.ForeignAid]: [Card.Duke],
    [ActionType.Assassinate]: [Card.Contessa],
    [ActionType.Steal]: [Card.Captain, Card.Ambassador],
};

/**
 * Returns a list of cards that block the current action.
 */
export function getCurrentActionBlockers(state: State): Card[] {
    if (!state.pendingAction) {
        return [];
    }
    return BLOCKED_BY[state.pendingAction.type] ?? [];
}

/**
 * Returns true only if `state.turnState` is equal to one of the provided states.
 */
export function stateIn(state: State, ...states: TurnState[]): boolean {
    for (const turnState of states) {
        if (state.turnState === turnState) {
            return true;
        }
    }
    return false;
}

export function formatActionType(type: ActionType): string {
    switch (type) {
        case ActionType.Income:
            return "Income";
        case ActionType.ForeignAid:
            return "Foreign Aid";
        case ActionType.Tax:
            return "Tax";
        case ActionType.Assassinate:
            return "Assassinate";
        case ActionType.Revolt:
            return "Revolt";
        case ActionType.Exchange:
            return "Exchange";
        case ActionType.Steal:
            return "Steal";
    }
    return "";
}

export function formatCurrency(value: number): string {
    return `${value}₡`;
}
