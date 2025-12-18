
export const ActionType = {
    HELLO_THERE: 0,
    ADD_ITEM_TAKER: 1,
    REMOVE_ITEM_TAKER: 2,
    ADD_ITEM: 3,
    REMOVE_ITEM: 4,
    EDIT_ITEM: 5,
    ADD_TAKER_ID: 6,
    REMOVE_TAKER_ID: 7,
    BYE_BYE: 8
} as const;

export type ActionType = typeof ActionType[keyof typeof ActionType];

export interface IAction {
    actionId: number;
    actionType: ActionType;
    splitId: string;
    itemId: string;
    itemName: string;
    takerId: string;
    price: number;
}

interface ActionResponse {
    success: boolean;
    message: string;
    action: IAction;
}

// responsible for sending, receiving and managing the actions via websocket
class ActionDispatcher {
    private socket: WebSocket | null = null;
    private messageListeners: ((msg: ActionResponse) => void)[] = [];
    private statusListeners: ((isConnected: boolean) => void)[] = [];

    public connect(url: string): void {
        if (this.socket) return;

        this.socket = new WebSocket(url);

        this.socket.onopen = () => {
            console.log('ActionDispatcher: Connected to WebSocket');
            this.notifyStatus(true);
        };

        this.socket.onmessage = (event) => {
            try {
                const parsedData: ActionResponse = JSON.parse(event.data);
                this.notifyMessage(parsedData)
            } catch (error) {
                console.error('ActionDispatcher: Error parsing message:', error);
            }
        };

        this.socket.onclose = () => {
            console.log('ActionDispatcher: Disconnected');
            this.notifyStatus(false);
            this.socket = null;
        };

        this.socket.onerror = (error) => {
            console.error('ActionDispatcher: WebSocket error:', error);
            this.socket?.close();
        };
    }

    public disconnect(): void {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
            this.notifyStatus(false);
        }
    }

    public sendAction(action: IAction): void {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            const payload = JSON.stringify(action);
            this.socket.send(payload);
        } else {
            console.warn('ActionDispatcher: Cannot send action, socket not connected');
        }
    }

    public onMessage(callback: (msg: ActionResponse) => void): void {
        this.messageListeners.push(callback);
    }

    public offMessage(callback: (msg: ActionResponse) => void): void {
        this.messageListeners = this.messageListeners.filter((cb) => cb !== callback);
    }

    public onStatusChange(callback: (isConnected: boolean) => void): void {
        this.statusListeners.push(callback);
    }

    public offStatusChange(callback: (isConnected: boolean) => void): void {
        this.statusListeners = this.statusListeners.filter((cb) => cb !== callback);
    }

    private notifyMessage(msg: ActionResponse): void {
        this.messageListeners.forEach((cb) => cb(msg));
    }

    private notifyStatus(isConnected: boolean): void {
        this.statusListeners.forEach((cb) => cb(isConnected));
    }
}

// responsible for exposing the APIs that construct the actions and dispatches them via ActionDispatcher
export class ActionManager {
    private actionDispatcher: ActionDispatcher;

    constructor() {
        this.actionDispatcher = new ActionDispatcher();
        this.actionDispatcher.onMessage(this.OnActionResponseRecived);
    }

    public connect(url: string) {
        this.actionDispatcher.connect(url);
    }

    public disconnect() {
        this.actionDispatcher.disconnect();
    }

    public onMessage(callback: (msg: ActionResponse) => void) {
        this.actionDispatcher.onMessage(callback);
    }

    public offMessage(callback: (msg: ActionResponse) => void) {
        this.actionDispatcher.offMessage(callback);
    }

    private createAction(type: ActionType, splitId: string, itemId: string = "", itemName: string = "", takerId: string = "", price: number = 0): IAction {
        return {
            actionId: 0, // Generating a unique ID for the action
            actionType: type,
            splitId,
            itemId,
            itemName,
            takerId,
            price
        };
    }

    public helloThere(splitId: string) {
        const action = this.createAction(ActionType.HELLO_THERE, splitId);
        this.actionDispatcher.sendAction(action);
    }

    public addItemTaker(splitId: string, itemId: string, takerId: string) {
        const action = this.createAction(ActionType.ADD_ITEM_TAKER, splitId, itemId, "", takerId);
        this.actionDispatcher.sendAction(action);
    }

    public removeItemTaker(splitId: string, itemId: string, takerId: string) {
        const action = this.createAction(ActionType.REMOVE_ITEM_TAKER, splitId, itemId, "", takerId);
        this.actionDispatcher.sendAction(action);
    }

    public addItem(splitId: string, itemName: string, price: number) {
        const action = this.createAction(ActionType.ADD_ITEM, splitId, "", itemName, "", price);
        this.actionDispatcher.sendAction(action);
    }

    public removeItem(splitId: string, itemId: string) {
        const action = this.createAction(ActionType.REMOVE_ITEM, splitId, itemId);
        this.actionDispatcher.sendAction(action);
    }

    public editItem(splitId: string, itemId: string, itemName: string, price: number) {
        const action = this.createAction(ActionType.EDIT_ITEM, splitId, itemId, itemName, "", price);
        this.actionDispatcher.sendAction(action);
    }

    public addTakerId(splitId: string, itemId: string, takerId: string) {
        const action = this.createAction(ActionType.ADD_TAKER_ID, splitId, "", "", takerId);
        this.actionDispatcher.sendAction(action);
    }

    public removeTakerId(splitId: string, takerId: string) {
        const action = this.createAction(ActionType.REMOVE_TAKER_ID, splitId, "", "", takerId);
        this.actionDispatcher.sendAction(action);
    }

    public byeBye(splitId: string) {
        const action = this.createAction(ActionType.BYE_BYE, splitId);
        this.actionDispatcher.sendAction(action);
    }

    private OnActionResponseRecived(msg: ActionResponse) {

    }
}
