import { createContext, useRef, useContext } from "react";
import type { SocketProvider, ActionResponse, IAction, BillData } from "../../../common/interfaces";
import { processAction } from "./ActionUtils";
import useSocket from "react-use-websocket";
import type { Options } from "react-use-websocket";
import { ReadyState } from "react-use-websocket";

export const SocketContextProvider = createContext<SocketProvider | null>(null);

interface SocketContextProps {
    children: React.ReactNode
    setBillData: React.Dispatch<React.SetStateAction<BillData>>;
}

export const SocketContextComponent = ({ children, setBillData }: SocketContextProps) => {
    const socketStatus = useRef<string>('Disconnected');
    const messageListeners = useRef<((msg: ActionResponse) => void)[]>([]);
    const statusListeners = useRef<((status: string) => void)[]>([]);
    const socketOptions: Options = {
        share: true,
        reconnectInterval: 500,
        reconnectAttempts: 5,
        shouldReconnect: () => true, // must be updated to handle reconnection logic
        onMessage: (event) => {
            try {
                const actionResponse: ActionResponse = JSON.parse(event.data);
                console.log("Action received:", actionResponse);
                if (actionResponse && actionResponse.success) {
                    setBillData((oldBillData) => processAction(actionResponse.action, oldBillData));
                    messageListeners.current.forEach(listener => listener(actionResponse));
                }
            } catch (err) {
                console.error("Error parsing WebSocket message:", err);
            }
        },
        onOpen: () => {
            console.log('WebSocket connected');
            socketStatus.current = "Session connected";
            statusListeners.current.forEach(listener => listener(socketStatus.current));
        },
        onClose: () => {
            console.log('WebSocket disconnected');
            socketStatus.current = "Disconnected";
            statusListeners.current.forEach(listener => listener(socketStatus.current));
        },
        onError: (error) => {
            console.error('WebSocket error:', error);
            socketStatus.current = "Socket error";
            statusListeners.current.forEach(listener => listener(socketStatus.current));
        }
    }

    const { sendMessage, readyState } = useSocket(import.meta.env.VITE_WEBSOCKET_URL, socketOptions)

    const subscribeToMessages = (callBack: (msg: ActionResponse) => void) => {
        messageListeners.current.push(callBack);
    }

    const subscribeToStatus = (callBack: (status: string) => void) => {
        statusListeners.current.push(callBack);
    }

    const publishAction = (action: IAction) => {
        sendMessage(JSON.stringify(action));
    }

    return (
        <SocketContextProvider.Provider value={{ subscribeToMessages, subscribeToStatus, publishAction, readyState }}>
            {children}
        </SocketContextProvider.Provider>
    );
};