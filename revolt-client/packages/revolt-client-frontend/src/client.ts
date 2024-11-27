export enum MessageType {
    CreateGame = 'create_game',
    JoinGame = 'join_game',
}
export interface Message {
    type: MessageType,
    payload?: Record<string, any>;
}

export class Client {

    private socket: WebSocket | undefined;

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

    private sendMessage(message: Message) {
        if (!this.socket) {
            return;
        }
        this.socket.send(JSON.stringify(message));
    }

    private handleMessage(event: MessageEvent) {
        console.log(event.data);
        if (!event.data) {
            return;
        }
        const message = JSON.parse(event.data);
        if (message.type === MessageType.CreateGame) {
            console.log(message.payload.gameId);
        }
    }
}