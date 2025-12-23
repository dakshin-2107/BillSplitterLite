import React, { useContext } from 'react';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { ReadyState } from 'react-use-websocket';
import './StatusIndicator.css';

export const StatusIndicator: React.FC = () => {
    const socketContext = useContext(SocketContextProvider);

    if (!socketContext) {
        return <></>;
    }

    const { readyState } = socketContext;

    const getStatusInfo = () => {
        switch (readyState) {
            case ReadyState.CONNECTING:
                return { label: 'Connecting...', className: 'connecting' };
            case ReadyState.OPEN:
                return { label: 'Connected', className: 'connected' };
            case ReadyState.CLOSING:
                return { label: 'Closing...', className: 'connecting' };
            case ReadyState.CLOSED:
                return { label: 'Disconnected', className: 'disconnected' };
            case ReadyState.UNINSTANTIATED:
                return { label: 'Uninitialized', className: 'disconnected' };
            default:
                return { label: 'Unknown', className: 'disconnected' };
        }
    };

    const { label, className } = getStatusInfo();

    return (
        <div className="status-indicator-container">
            <span className={`status-dot ${className}`}></span>
            <span className="status-text">{label}</span>
        </div>
    );
};
