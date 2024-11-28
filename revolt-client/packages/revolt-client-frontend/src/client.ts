export enum MessageType {
    CreateGame = 'create_game',
    JoinGame = 'join_game',
    StartGame = 'start_game',
}

export interface StateUpdate {
    clientId: string;
    state: State;
}
export enum GameStatus {

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
    game: Record<string, any>;
    clients: Record<string, Peer>;
}
export class Client {

    private socket?: WebSocket;
    private clientId?: string;
    private state?: State;

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

    createGame() {
        this.sendMessage({
            type: MessageType.CreateGame
        });
    }

    joinGame(gameId: string) {
        this.sendMessage({
            type: MessageType.JoinGame,
            payload: {
                gameId
            }
        });
    }

    startGame() {
        this.sendMessage({
            type: MessageType.StartGame,
        });
    }

    private sendMessage(message: Message) {
        if (!this.socket) {
            return;
        }
        this.socket.send(JSON.stringify(message));
    }

    private handleMessage(event: MessageEvent) {
        if (!event.data) {
            return;
        }
        const message = JSON.parse(event.data) as StateUpdate;
        this.clientId = message.clientId;
        this.state = message.state;

        console.log('my id:', this.clientId);
        console.log(JSON.stringify(this.state, undefined, 2));

    }
}