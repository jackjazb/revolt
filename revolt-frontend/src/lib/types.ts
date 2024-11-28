

/**
 * Possible card types
 */
export enum Card {
    Duke = "duke",
    Assassin = "assassin",
    Ambassador = "ambassador",
    Captain = "captain",
    Contessa = "contessa",
}

export enum MessageType {
    CreateGame = 'create_game',
    JoinGame = 'join_game',
    StartGame = 'start_game',
}

export interface StateUpdate {
    gameId: string,
    ownerId: string,
    clientId: string,
    leader: number,
    number: number,
    status: string,
    clients: [],
    self: {
        cards: {
            card: Card,
            alive: boolean;
        }[],
        credits: number;
    },
    turnState: number,
    pendingAction: {
        type: number,
        targetPlayer: number;
    },
    pendingBlock: {
        card: string,
        initiator: number;
    },
    pendingChallenge: {
        initiator: number;
    };
}

/**
 * The status of a game
 */
export enum GameStatus {
    Lobby = "lobby",
    InProgress = "in_progress",
    CompleteGameStatus = "complete",
}

export interface Message {
    type: MessageType,
    payload?: Record<string, any>;
}

// All the information about a collected connected peer a client is allowed to have.
export interface Peer {
    name: string;
    number: number;
}

export interface State {
    gameId: string,
    ownerId: string,
    status: GameStatus,
    game: Record<string, unknown>;
    clients: Record<string, Peer>;
}