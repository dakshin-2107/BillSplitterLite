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


// HELLO_THERE = iota

// 	// just to the ping the backend
// 	PING

// 	// sync the current bill state with the client
// 	SYNC_BILL_STATE

// 	// basic taker CRUD operations
// 	ADD_TAKER_FOR_ITEM // operation must be idempotent
// 	DELETE_TAKER_FOR_ITEM

// 	// increment/decrement the number of shares for a taker in an item
// 	INCREMENT_TAKER_ID
// 	DECREMENT_TAKER_ID

// 	// basic item CRUD operations
// 	ADD_NEW_ITEM
// 	DELETE_ITEM
// 	EDIT_ITEM

// 	// create/remove the taker only from the item
// 	ADD_NEW_TAKER
// 	DELETE_TAKER
// 	EDIT_TAKER

// 	// edit bill information
// 	EDIT_BILL_INFORMATION

// 	// close the session and delete all the data associated with it
// 	BYE_BYE

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

    // increment/decrement the number of shares for a taker in an item
    INCREMENT_TAKER_ID: 5,
    DECREMENT_TAKER_ID: 6,

    // basic item CRUD operations
    ADD_NEW_ITEM: 7,
    DELETE_ITEM: 8,
    EDIT_ITEM: 9,

    // create/remove the taker only from the item
    ADD_NEW_TAKER: 10,
    DELETE_TAKER: 11,
    EDIT_TAKER: 12,

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
    takerId?: string;
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
    publishAction: (action: IAction) => void;
    readyState: ReadyState;
}