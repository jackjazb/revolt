import { initialState, type Action, type Block, type State } from "./types";
import { randomName } from "./utils";

export enum MessageType {
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

const WS_TIMEOUT = 5000;

export enum ClientStatus {
    Default = "default",
    Connecting = "connecting",
    Connected = "connected",
    Failed = "failed"
}

export class Client {
    private socket?: WebSocket;
    state: State = initialState;
    status: ClientStatus = ClientStatus.Default;

    constructor(private onStateUpdate: (state: State) => void) { };

    async connect(uri: string, playerName?: string) {
        this.status = ClientStatus.Connecting;

        // Initial connection should include the player's alias in the URL search params.
        if (!playerName) {
            playerName = randomName();
        }

        const url = new URL(uri);
        url.searchParams.set('name', playerName);
        const socket = new WebSocket(url);

        return new Promise<void>(resolve => {
            const timeout = setTimeout(() => {
                socket.close();
            }, WS_TIMEOUT);

            socket.onopen = () => {
                this.status = ClientStatus.Connected;

                clearTimeout(timeout);
                socket.onmessage = e => this.handleMessage(e);
                this.socket = socket;
                resolve();
            };

            socket.onclose = () => {
                this.status = ClientStatus.Failed;

                clearTimeout(timeout);
                resolve();
            };
        });
    }

    leave() {
        this.socket?.close();
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
        const message = JSON.parse(event.data) as State;
        this.state = message;
        this.onStateUpdate(this.state);
        console.log('received state update:', JSON.stringify(this.state, undefined, 2));
    }
}