import React, { useRef, useState, useContext } from 'react';
import './TallyPanel.css';
import type { Tally, SplitData } from '../../../common/interfaces';
import html2canvas from 'html2canvas';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { ActionType } from '../../../common/interfaces';

interface TallyPanelProps {
    tally: Tally | null;
    participants: Record<string, string>;
    splitData: SplitData;
}

const TallyPanel: React.FC<TallyPanelProps> = ({ tally, participants, splitData }) => {
    const SocketContext = useContext(SocketContextProvider);
    const panelRef = useRef<HTMLDivElement>(null)
    const [isEditing, setIsEditing] = useState<boolean>(false);
    const [newTotal, setNewTotal] = useState<number>(splitData.totalAmount);
    const [isCompact, setIsCompact] = useState<boolean>(true);

    const handleCopyImage = async () => {
        if (!panelRef.current) return;

        try {
            const canvas = await html2canvas(panelRef.current, {
                backgroundColor: '#1E1E1E', // Match container background
                scale: 2, // Higher quality
                logging: false,
                useCORS: true
            });

            canvas.toBlob(async (blob) => {
                if (blob) {
                    try {
                        await navigator.clipboard.write([
                            new ClipboardItem({
                                'image/png': blob
                            })
                        ]);
                        console.log("Image copied to clipboard");
                    } catch (err) {
                        console.error("Failed to copy image to clipboard:", err);
                    }
                }
            }, 'image/png');
        } catch (err) {
            console.error("Error generating tally image:", err);
        }
    };

    const handleSaveTotal = () => {
        setIsEditing(false);
        if (SocketContext) {
            SocketContext.publishAction({
                actionId: 0,
                actionType: ActionType.EDIT_BILL_INFO,
                splitId: splitData.splitId,
                itemId: 0,
                total: newTotal
            })
        }
    }

    if (!tally || Object.keys(tally.userShares).length === 0) {
        return (
            <div className="tally-panel empty">
                <h3>Bill Tally</h3>
                <p>Waiting for data...</p>
            </div>
        );
    }

    console.log("Tally", tally);
    console.log("User Shares", Object.entries(tally.userShares));

    return (
        <div className="tally-panel" ref={panelRef}>
            <div className="panel-header">
                <div className="header-actions">
                    <span className="split-label">Per user split</span>
                    <div className="view-toggle">
                        <span className="toggle-label">Compact</span>
                        <label className="switch">
                            <input
                                type="checkbox"
                                checked={isCompact}
                                onChange={(e) => setIsCompact(e.target.checked)}
                            />
                            <span className="slider round"></span>
                        </label>
                    </div>
                    <button className="copy-btn" onClick={handleCopyImage}>
                        Copy
                    </button>
                </div>
            </div>
            <div className="user-shares-list">
                {Object.entries(tally.userShares).map(([userId, userShare]) => (
                    <div key={userId} className="user-share-card">
                        <div className="user-info">
                            <span className="user-name">{participants[userId] || userId}</span>
                            <span className="user-total">{userShare.userShareTotal?.toFixed(2)}</span>
                        </div>
                        <div className="bill-shares-list">
                            {Object.entries(userShare.billShares).map(([billId, billShare]) => (
                                <div key={billId} className="bill-share-group">
                                    <div className="bill-name-header">
                                        <span className="bill-name">{tally.billNameMap[Number(billId)] || `Bill ${billId}`}</span>
                                        <span className="bill-share-total">{billShare.billShareTotal?.toFixed(2)}</span>
                                    </div>
                                    {!isCompact && (
                                        <div className="item-breakdown">
                                            {Object.entries(billShare.itemShares).map(([itemName, amount]) => (
                                                <div key={itemName} className="item-share">
                                                    <span className="item-name">{itemName}</span>
                                                    <span className="item-amount">{amount.toFixed(2)}</span>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>


            <div className="tally-header">
                <h3></h3>
                <div className="total-summary">
                    <div className="summary-item actual">
                        <span>Actual (Y):</span>
                        <div className="edit-total-container">
                            {isEditing ? (
                                <>
                                    <input
                                        type="number"
                                        value={newTotal}
                                        onChange={(e) => setNewTotal(Number(e.target.value))}
                                        className="edit-total-input"
                                        autoFocus
                                    />
                                    <button className="inline-action-btn save" onClick={handleSaveTotal} title="Save">
                                        ✓
                                    </button>
                                    <button className="inline-action-btn cancel" onClick={() => setIsEditing(false)} title="Cancel">
                                        ✕
                                    </button>
                                </>
                            ) : (
                                <>
                                    <button className="edit-trigger" onClick={() => setIsEditing(true)}>
                                        Edit
                                    </button>
                                    <span className="amount">{splitData.totalAmount.toFixed(2)}</span>
                                </>
                            )}
                        </div>
                    </div>
                    <div className="summary-item">
                        <span>Calculated (X):</span>
                        <span className="amount">{tally.calculatedTotal.toFixed(2)}</span>
                    </div>
                    <div className="summary-item difference">
                        <span>Difference (X-Y):</span>
                        <span className={`amount ${tally.totalDifference !== 0 ? 'warning' : 'success'}`}>
                            {tally.totalDifference.toFixed(2)}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TallyPanel;
