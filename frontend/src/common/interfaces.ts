import { ReadyState } from "react-use-websocket";

// split grid stuff
export interface Item {
    id: string;
    name: string;
    price: number;
    takers: Record<string, number>;
}

export interface BillData {
    splitId: string;
    date: string;
    location: string;
    total: number;
    items: Record<string, Item>;
    participants: Record<string, string>;
}

// bill form stuff 
export interface FormData {
    place: string;
    dateTime: string;
    names: string[];
    image: File | null;
}

export interface ApiResponse {
    success: boolean;
    message: string;
    bill?: BillData;
}

// action manager stuff
// export const ActionType = {
//     HELLO_THERE: 0,
//     ADD_ITEM_TAKER: 1,
//     REMOVE_ITEM_TAKER: 2,
//     ADD_ITEM: 3,
//     REMOVE_ITEM: 4,
//     EDIT_ITEM: 5,
//     ADD_TAKER_ID: 6,
//     REMOVE_TAKER_ID: 7,
//     BYE_BYE: 8
// } as const;

export const ActionType = {
    // initiate the session
    HELLO_THERE: 0,

    // just to the ping the backend
    PING: 1,

    // sync the current bill state with the client
    SYNC_BILL_STATE: 2,

    // basic taker CRUD operations
    ADD_TAKER_FOR_ITEM: 3,
    DELETE_TAKER_FOR_ITEM: 4,

    // basic item CRUD operations
    ADD_NEW_ITEM: 5,
    DELETE_ITEM: 6,
    EDIT_ITEM: 7,

    // create/remove the taker only from the item
    ADD_NEW_TAKER: 8,
    DELETE_TAKER: 9,
    EDIT_TAKER: 10,

    // increment/decrement the number of shares for a taker in an item
    INCREMENT_TAKER_ID: 11,
    DECREMENT_TAKER_ID: 12,

    // close the session and delete all the data associated with it
    BYE_BYE: 13
} as const;

export type ActionType = typeof ActionType[keyof typeof ActionType];

export interface IAction {
    actionId?: number;
    actionType: ActionType;
    splitId?: string;
    itemId: string;
    itemName?: string;
    takerId: string;
    price?: number;
    total?: number;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    action: IAction;
}

export interface SocketProvider {
    subscribeToMessages: (callback: (msg: ActionResponse) => void) => void;
    subscribeToStatus: (callback: (status: string) => void) => void;
    publishAction: (action: IAction) => void;
    readyState: ReadyState;
}