import { createContext, useRef } from "react";
import type { SocketProvider, ActionResponse, IAction, SplitData } from "../../../common/interfaces";
import { processAction } from "./ActionUtils";
import useSocket from "react-use-websocket";
import type { Options } from "react-use-websocket";
import type { Tally, TallyResponse, BillTallyResponse, ApiResponse } from "../../../common/interfaces";
import { ActionType } from "../../../common/interfaces";
import { URLProvider } from "../../../common/urlProvider";

export const SocketContextProvider = createContext<SocketProvider | null>(null);

interface SocketContextProps {
    children: React.ReactNode
    setSplitData: React.Dispatch<React.SetStateAction<SplitData | null>>;
    setTallyData: React.Dispatch<React.SetStateAction<Tally | null>>;
}

export const SocketContextComponent = ({ children, setSplitData, setTallyData }: SocketContextProps) => {
    const messageListeners = useRef<((msg: ActionResponse) => void)[]>([]);
    const tallyListeners = useRef<((msg: Tally) => void)[]>([]);
    const socketOptions: Options = {
        share: true,
        reconnectInterval: 500,
        reconnectAttempts: 5,
        shouldReconnect: () => true,
        onMessage: (event) => {
            try {
                console.log("Message received:", event.data);

                // action sync
                const actionResponse: ActionResponse = JSON.parse(event.data);
                if (actionResponse && actionResponse.success && actionResponse.action) {
                    console.log("Action received:", actionResponse);
                    setSplitData((oldSplitData) => processAction(actionResponse.action, oldSplitData));
                    messageListeners.current.forEach(listener => listener(actionResponse));
                    return;
                }

                // full tally sync
                const tallyResponse: TallyResponse = JSON.parse(event.data);
                if (tallyResponse && tallyResponse.success && tallyResponse.tally) {
                    console.log("Full tally received:", tallyResponse);
                    setTallyData(tallyResponse.tally);
                    tallyListeners.current.forEach(listener => listener(tallyResponse.tally));
                    return;
                }

                // partial tally sync — merge dirty bill shares into cached tally
                const billTallyResponse: BillTallyResponse = JSON.parse(event.data);
                if (billTallyResponse && billTallyResponse.success && billTallyResponse.billShares) {
                    console.log("Partial tally received:", billTallyResponse);
                    setTallyData(prev => {
                        if (!prev) return prev;
                        const merged: Tally = {
                            ...prev,
                            billShares: { ...prev.billShares, ...billTallyResponse.billShares }
                        };
                        tallyListeners.current.forEach(listener => listener(merged));
                        return merged;
                    });
                    return;
                }

                // split sync
                const billResponse: ApiResponse = JSON.parse(event.data);
                if (billResponse && billResponse.success) {
                    if (billResponse.split) {
                        console.log("Bill received:", billResponse);
                        setSplitData(billResponse.split);
                    }
                    else {
                        // ignore else case for now since tally is not developed 
                        //console.log("Bill deleted:", billResponse);
                        //setSplitData(null);
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
                itemId: 0
            });
        },
        onClose: (event) => {
            console.warn('[WebSocket] Connection closed', {
                timestamp: new Date().toISOString(),
                code: event.code,
                reason: event.reason || '(no reason provided)',
                wasClean: event.wasClean,
            });
        },
        onError: (event) => {
            const ws = event.target as WebSocket;
            const readyStateMap: Record<number, string> = {
                0: 'CONNECTING',
                1: 'OPEN',
                2: 'CLOSING',
                3: 'CLOSED',
            };
            console.error('[WebSocket] Error event fired', {
                timestamp: new Date().toISOString(),
                // WebSocket state at the time of the error
                readyState: readyStateMap[ws?.readyState] ?? ws?.readyState,
                url: ws?.url,
                protocol: ws?.protocol || '(none)',
                // The error event itself carries no message — the close event
                // that follows will have the code and reason.
                note: 'A close event with a code/reason should follow immediately.',
            });
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