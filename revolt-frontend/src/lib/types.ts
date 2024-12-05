/*
Generally, types defined in this file match the JSON messages sent by the server.
*/

import type { Client } from "./Client";

export type IconType = 'dice' | 'coin';
/**
 * Possible card types
 */
export enum Card {
    Empty = "",
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
    target?: number;
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
    Finished = "finished",
}

/**
 * The status of a game
 */
export enum GameStatus {
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
    initiator: number;
}

export interface Challenge {
    initiator: number;
}

export interface Peer {
    id: string;
    name: string;
    number: string;
}
/**
 * A state update received from the server.
 */
export interface State {
    gameId: string,
    ownerId: string,
    clientId: string,
    clientName: string;
    leader: number,
    isLeader: boolean;
    number: number,
    status: GameStatus,
    clients: Peer[],
    turnState: TurnState,
    pendingAction: Action;
    pendingBlock: Block;
    pendingChallenge: Challenge;
    self: {
        cards: CardState[],
        credits: number;
    },
}

/**
 * Type for props of components that need to render game stuff.
 */
export interface GameComponentContext {
    gameState: State;
    client: Client;
}