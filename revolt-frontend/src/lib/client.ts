import { type Action, type Block, type CreateGameResponse, type State } from "./types";

export enum MessageType {
    CreateGame = 'create_game',
    JoinGame = 'join_game',
    ChangeName = 'change_name',
    StartGame = 'start_game',
    AttemptAction = 'attempt_action',
    AttemptBlock = 'attempt_block',
    Challenge = 'challenge',
    ResolveDeath = 'resolve_death',
    CommitTurn = 'commit_turn',
    EndTurn = 'end_turn'
}

export interface Message {
    type: MessageType,
    payload?: Record<string, any>;
}

export class Client {

    private socket?: WebSocket;
    private onStateUpdate: (state: State) => void;
    state?: State;

    constructor(onStateUpdate: (state: State) => void) {
        this.onStateUpdate = onStateUpdate;
    }

    async connect(uri: string) {
        const socket = new WebSocket(uri);
        return new Promise<void>(resolve => {
            socket.onopen = () => {
                socket.onmessage = e => this.handleMessage(e);
                this.socket = socket;
                resolve();
            };
        });
    }

    joinGame(gameId: string, playerName: string) {
        this.sendMessage({
            type: MessageType.JoinGame,
            payload: {
                gameId,
                playerName
            }
        });
    }

    changeName(playerName: string) {
        this.sendMessage({
            type: MessageType.ChangeName,
            payload: {
                playerName
            }
        });
    }

    startGame() {
        this.sendMessage({
            type: MessageType.StartGame,
        });
    }


    /**
     * The initial action in a turn. Represents an attempt to perform an action.
     */
    attemptAction(action: Action) {
        this.sendMessage({
            type: MessageType.AttemptAction,
            payload: {
                action
            }
        });
    }

    /**
     * Represents an attempt to block an action.
     */
    attemptBlock(block: Block) {
        this.sendMessage({
            type: MessageType.AttemptBlock,
            payload: {
                block
            }
        });
    }

    /** 
     * Represents a challenge of a block or action.
    */
    challenge() {
        this.sendMessage({
            type: MessageType.Challenge,
        });
    }

    /**
     * Moves the game on from a pending death state by killing `card` 
     * on the next player to die.
     */
    resolveDeath(card: number) {
        this.sendMessage({
            type: MessageType.ResolveDeath,
            payload: {
                card
            }
        });
    }

    commitTurn() {
        this.sendMessage({
            type: MessageType.CommitTurn,
        });
    }

    endTurn() {
        this.sendMessage({
            type: MessageType.EndTurn
        });
    }

    private sendMessage(message: Message) {
        if (!this.socket) {
            return;
        }
        console.log('sending message', JSON.stringify(message, undefined, 2));
        this.socket.send(JSON.stringify(message));
    }

    private handleMessage(event: MessageEvent) {
        if (!event.data) {
            return;
        }
        const message = JSON.parse(event.data) as State | CreateGameResponse;
        if ('id' in message) {
            console.log('connected - given ID', message.id);
            return;
        }
        this.state = message;
        this.onStateUpdate(this.state);
        console.log('received state update:', JSON.stringify(this.state, undefined, 2));
    }
}