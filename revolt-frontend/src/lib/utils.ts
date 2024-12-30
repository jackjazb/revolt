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
 * Returns a list of cards that block the current pending action.
 */
export function getCurrentActionBlockers(state: State): Card[] {
    const blockedBy: Partial<Record<ActionType, Card[]>> = {
        [ActionType.ForeignAid]: [Card.Duke],
        [ActionType.Assassinate]: [Card.Contessa],
        [ActionType.Steal]: [Card.Captain, Card.Ambassador],
    };
    if (!state.pendingAction) {
        return [];
    }
    if (state.pendingAction.target && state.pendingAction.target != state.self.id) {
        return [];

    }
    return blockedBy[state.pendingAction.type] ?? [];
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

export function formatCard(card: Card): string {
    return `${card}`.charAt(0).toUpperCase() + `${card}`.slice(1);
}

export function formatCurrentAction(state: State): string {
    const leader = getLeader(state);
    if (state.pendingAction.target) {
        const target = state.pendingAction.target === state.self.id ? 'you' : getPlayerById(state, state.pendingAction.target);
        switch (state.pendingAction.type) {
            case ActionType.Assassinate:
                return `${leader} has attempted to assassinate ${target}.`;
            case ActionType.Revolt:
                return `${leader} has revolted against ${target}.`;
            case ActionType.Steal:
                return `${leader} has attempted to steal from ${target}.`;
        }
    }
    return `${leader} has attempted ${formatActionType(state.pendingAction.type)}.`;
}

export function formatCurrentBlock(state: State): string {
    if (!state.pendingBlock.initiator) {
        return "";
    }
    const blocker = getPlayerById(state, state.pendingBlock.initiator);
    const card = formatCard(state.pendingBlock.card);
    return `${blocker} has blocked your action with their ${card}.`;
}