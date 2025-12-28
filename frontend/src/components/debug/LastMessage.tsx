import React, { useContext } from 'react';
import { SocketContextProvider } from '../split-grid/action-manager/SocketContext';
import './LastMessage.css';

const LastMessage: React.FC = () => {
    const socketContext = useContext(SocketContextProvider);

    if (!socketContext) {
        return null;
    }

    const { lastMessage } = socketContext;

    return (
        <div className="last-message-container">
            <div className="last-message-header">
                <span className="last-message-title">Last Socket Message Debugger</span>
            </div>
            <div className="last-message-content">
                <pre>
                    {lastMessage ? JSON.stringify(JSON.parse(lastMessage.data), null, 4) : 'No messages received yet'}
                </pre>
            </div>
        </div>
    );
};

export default LastMessage;
