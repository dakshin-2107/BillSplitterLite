import { createContext, useRef } from "react";
import type { SocketProvider, ActionResponse, IAction, BillData } from "../../../common/interfaces";
import { processAction } from "./ActionUtils";
import useSocket from "react-use-websocket";
import type { Options } from "react-use-websocket";
import type { Tally, TallyResponse, ApiResponse } from "../../../common/interfaces";
import { ActionType } from "../../../common/interfaces";
import { URLProvider } from "../../../common/urlProvider";

export const SocketContextProvider = createContext<SocketProvider | null>(null);

interface SocketContextProps {
    children: React.ReactNode
    setBillData: React.Dispatch<React.SetStateAction<BillData | null>>;
    setTallyData: React.Dispatch<React.SetStateAction<Tally | null>>;
}

export const SocketContextComponent = ({ children, setBillData, setTallyData }: SocketContextProps) => {
    const messageListeners = useRef<((msg: ActionResponse) => void)[]>([]);
    const tallyListeners = useRef<((msg: Tally) => void)[]>([]);
    const socketOptions: Options = {
        share: true,
        reconnectInterval: 500,
        reconnectAttempts: 5,
        shouldReconnect: () => true, // must be updated to handle reconnection logic
        onMessage: (event) => {
            try {
                console.log("Message received:", event.data);
                const actionResponse: ActionResponse = JSON.parse(event.data);
                if (actionResponse && actionResponse.success && actionResponse.action) {
                    console.log("Action received:", actionResponse);
                    setBillData((oldBillData) => processAction(actionResponse.action, oldBillData));
                    messageListeners.current.forEach(listener => listener(actionResponse));
                    return;
                }

                const tallyResponse: TallyResponse = JSON.parse(event.data);
                if (tallyResponse && tallyResponse.success && tallyResponse.tally) {
                    console.log("Tally received:", tallyResponse);
                    setTallyData(tallyResponse.tally);
                    tallyListeners.current.forEach(listener => listener(tallyResponse.tally));
                    return;
                }

                const billResponse: ApiResponse = JSON.parse(event.data);
                if (billResponse && billResponse.success) {
                    if (billResponse.bill) {
                        console.log("Bill received:", billResponse);
                        setBillData(billResponse.bill);
                    }
                    else {
                        console.log("Bill deleted:", billResponse);
                        setBillData(null);
                    }
                }
            } catch (err) {
                console.error("Error parsing WebSocket message:", err);
            }
        },
        onOpen: () => {
            console.log('WebSocket connected');
            publishAction({
                actionType: ActionType.SYNC_BILL_STATE,
                itemId: "0"
            });
        },
        onClose: () => {
            console.log('WebSocket disconnected');
        },
        onError: (error) => {
            console.error('WebSocket error:', error);
        }
    }

    const socketUrl = URLProvider.getSocketUrl();
    const { sendMessage, lastMessage, readyState } = useSocket(socketUrl, socketOptions)

    const subscribeToMessages = (callBack: (msg: ActionResponse) => void) => {
        messageListeners.current.push(callBack);
    }

    const publishAction = (action: IAction) => {
        sendMessage(JSON.stringify(action));
    }

    return (
        <SocketContextProvider.Provider value={{ subscribeToMessages, publishAction, readyState, lastMessage }}>
            {children}
        </SocketContextProvider.Provider>
    );
};