import React from 'react';
import './BillActions.css';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { useContext, useState } from 'react';
import { ActionType } from '../../../common/interfaces';

const BillActions: React.FC = () => {
    const socketContext = useContext(SocketContextProvider);
    const [showConfirmModal, setShowConfirmModal] = useState(false);

    const handleShareBill = () => {
        console.log('Share Bill clicked');
    };

    const handleCloseBill = () => {
        setShowConfirmModal(true);
    };

    const confirmClose = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.BYE_BYE,
                itemId: "0" // dummy itemId for session closing
            });
        }
        setShowConfirmModal(false);
    };

    const cancelClose = () => {
        setShowConfirmModal(false);
    };

    return (
        <div className="right-controls">
            <button className="control-btn share-btn" onClick={handleShareBill}>
                Share Bill
            </button>
            <button className="control-btn close-btn" onClick={handleCloseBill}>
                Close Bill
            </button>

            {showConfirmModal && (
                <div className="modal-overlay">
                    <div className="modal-content">
                        <h3>Close Bill</h3>
                        <p>Are you sure you want to close this bill session? All data will be deleted.</p>
                        <div className="modal-actions">
                            <button className="modal-btn cancel-btn" onClick={cancelClose}>
                                Cancel
                            </button>
                            <button className="modal-btn confirm-btn" onClick={confirmClose}>
                                Confirm
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default BillActions;
