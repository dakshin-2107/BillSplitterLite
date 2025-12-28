import React from 'react';
import './BillActions.css';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { useContext } from 'react';
import { ActionType } from '../../../common/interfaces';

const BillActions: React.FC = () => {
    const socketContext = useContext(SocketContextProvider);
    const handleShareBill = () => {
        console.log('Share Bill clicked');
    };

    const handleCloseBill = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.BYE_BYE,
                itemId: "0" // dummy itemId for session closing
            });
        }
    };

    return (
        <div className="right-controls">
            <button className="control-btn share-btn" onClick={handleShareBill}>
                Share Bill
            </button>
            <button className="control-btn close-btn" onClick={handleCloseBill}>
                Close Bill
            </button>
        </div>
    );
};

export default BillActions;
