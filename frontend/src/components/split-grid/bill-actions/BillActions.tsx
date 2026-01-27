import React from 'react';
import './BillActions.css';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { useContext, useState } from 'react';
import { ActionType } from '../../../common/interfaces';
import { URLProvider } from '../../../common/urlProvider';

interface ShareBillApiResponse {
    success: boolean;
    message: string;
    sessionId: string;
}

const BillActions: React.FC = () => {
    const socketContext = useContext(SocketContextProvider);
    const [showConfirmModal, setShowConfirmModal] = useState(false);
    const [shareStatus, setShareStatus] = useState<'idle' | 'copying' | 'copied' | 'error'>('idle');

    const handleShareBill = async () => {
        if (shareStatus === 'copying' || shareStatus === 'copied') return;

        setShareStatus('copying');
        try {
            const generateUrl = URLProvider.getGenerateUrl();
            const response = await fetch(generateUrl, {
                method: 'GET',
                credentials: 'include'
            });

            if (!response.ok) throw new Error('Failed to generate session ID');

            const result: ShareBillApiResponse = await response.json();
            if (result.success && result.sessionId) {
                const joinLink = `${URLProvider.getSiteOrigin()}/join/${result.sessionId}`;
                console.log(`Generated join link : ${joinLink}`);
                await navigator.clipboard.writeText(joinLink); // doesn't work in HTTP 
                setShareStatus('copied');

                // Reset status after 3 seconds
                setTimeout(() => setShareStatus('idle'), 3000);
            } else {
                throw new Error(result.message || 'Failed to generate session ID');
            }
        } catch (err) {
            console.error('Error sharing bill:', err);
            setShareStatus('error');
            setTimeout(() => setShareStatus('idle'), 3000);
        }
    };

    const handleCloseBill = () => {
        setShowConfirmModal(true);
    };

    const confirmClose = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.BYE_BYE,
                itemId: 0 // dummy itemId for session closing
            });
        }
        setShowConfirmModal(false);
    };

    const cancelClose = () => {
        setShowConfirmModal(false);
    };

    return (
        <div className="right-controls">
            <button
                className={`control-btn share-btn ${shareStatus}`}
                onClick={handleShareBill}
                disabled={shareStatus === 'copying'}
            >
                {shareStatus === 'idle' && 'Share Bill'}
                {shareStatus === 'copying' && 'Generating...'}
                {shareStatus === 'copied' && 'Link Copied!'}
                {shareStatus === 'error' && 'Retry Share'}
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
