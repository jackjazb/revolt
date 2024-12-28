/*
Generally, types defined in this file match the JSON messages sent by the server.
*/


export type IconType = "dice" | "coin" | "lie" | "skull" | "crown";
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

/**
 * Possible types of action
 */
export enum ActionType {
    Empty = "",
    Income = "income",
    ForeignAid = "foreign_aid",
    Tax = "tax",
    Assassinate = "assassinate",
    Revolt = "revolt",
    Exchange = "exchange",
    Steal = "steal",
}

/**
 * An attempt to perform an action by a player.
 */
export interface Action {
    type: ActionType,
    target?: string;
}

/**
 * Possible states for a turn.
 */
export enum TurnState {
    Default = "default",
    ActionPending = "action_pending",
    BlockPending = "block_pending",
    ExchangePending = "exchange_pending",
    PlayerLostChallenge = "player_lost_challenge",
    LeaderLostChallenge = "leader_lost_challenge",
    PlayerKilled = "player_killed",
    PlayerWon = "player_won",
    Finished = "finished",
}

/**
 * The status of a game
 */
export enum GameStatus {
    Default = "default",
    Lobby = "lobby",
    InProgress = "in_progress",
    CompleteGameStatus = "complete",
}

export interface CardState {
    card: Card,
    alive: boolean;
}

export interface Block {
    card: Card;
    initiator?: string;
}

export interface Challenge {
    initiator: string;
}

/**
 * A player - can be the player or a peer.
 */
export interface Peer {
    id: string;
    name: string;
    cards: CardState[];
    credits: number;
    leading: boolean;
    /**
     * Allowed actions - should only appear on `self`.
     */
    allowedActions?: ActionType[];
}

export interface CreateGameResponse {
    id: string;
}

/**
 * A state update received from the server.
 */
export interface State {
    timestamp: string,
    gameId: string,
    ownerId: string;
    winner: string;
    self: Peer;
    peers: Peer[];
    status: GameStatus,
    nextDeath: string,
    turnState: TurnState,
    pendingAction: Action;
    pendingBlock: Block;
    pendingChallenge: Challenge;
}

export const initialState: State = {
    timestamp: "",
    gameId: "",
    ownerId: "",
    winner: "",
    nextDeath: "",
    self: {
        id: "",
        name: "",
        cards: [],
        credits: 0,
        leading: false
    },
    peers: [],
    status: GameStatus.Default,
    turnState: TurnState.Default,
    pendingAction: {
        type: ActionType.Empty,
        target: ""
    },
    pendingBlock: {
        card: "" as Card,
        initiator: ""
    },
    pendingChallenge: {
        initiator: ""
    }
};