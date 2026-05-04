import React, { useRef, useState } from 'react';
import './TallyPanel.css';
import type { Tally, SplitData } from '@utils/interfaces';
import { toPng } from 'html-to-image';
import { formatDate } from '@utils/dateUtils';

interface TallyPanelProps {
    tally: Tally | null;
    participants: Record<string, string>;
    splitData: SplitData;
    activeBillId: number;
}

const TallyPanel: React.FC<TallyPanelProps> = ({ tally, participants, splitData, activeBillId }) => {
    const panelRef = useRef<HTMLDivElement>(null)
    const [isCompact, setIsCompact] = useState<boolean>(true);

    const handleCopyImage = async () => {
        if (!panelRef.current) return;
        try {
            const dataUrl = await toPng(panelRef.current, { pixelRatio: 2 });
            const blob = await (await fetch(dataUrl)).blob();
            await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]);
        } catch (err) {
            console.error('Failed to copy tally image:', err);
        }
    };

    const billShare = tally?.billShares[activeBillId];

    if (!tally || !billShare) {
        return (
            <div className="tally-panel empty">
                <h3>Bill Tally</h3>
                <p>Waiting for data...</p>
            </div>
        );
    }

    const billActualTotal = splitData.bills[activeBillId]?.total ?? 0;
    const billCalculatedTotal = billShare.billTotal;
    const billDifference = billCalculatedTotal - billActualTotal;

    const activeBill = splitData.bills[activeBillId];

    return (
        <div className="tally-section-container">
            <div className="tally-controls">
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
                    Copy Image
                </button>
            </div>

            <div className="tally-panel" ref={panelRef}>
                <div className="panel-header">
                    <div className="bill-info">
                        <p className="bill-location">{activeBill?.location || `Bill ${activeBillId}`}</p>
                        <span className="bill-date">{activeBill?.date ? formatDate(activeBill.date) : ''}</span>
                    </div>
                </div>
                <div className="user-shares-list">
                    <div className="bill-share-card">
                        <div className="bill-shares-list">
                            {Object.entries(billShare.userShares).map(([userId, userShare]) => (
                                <div key={userId} className="user-share-group">
                                    <div className="user-info">
                                        <span className="user-name">{participants[userId] || userId}</span>
                                        <span className="user-total">{userShare.userShareTotal?.toFixed(2)}</span>
                                    </div>
                                    {!isCompact && (
                                        <div className="item-breakdown">
                                            {Object.entries(userShare.itemShares).map(([itemName, amount]) => (
                                                <div key={itemName} className="item-share">
                                                    <span className="item-name">{itemName}</span>
                                                    <span className="item-amount">{(amount as number).toFixed(2)}</span>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                <div className="tally-header">
                    <h3></h3>
                    <div className="total-summary">
                        <div className="summary-item actual">
                            <span>Actual (Y):</span>
                            <span className="amount">{billActualTotal.toFixed(2)}</span>
                        </div>
                        <div className="summary-item">
                            <span>Calculated (X):</span>
                            <span className="amount">{billCalculatedTotal.toFixed(2)}</span>
                        </div>
                        <div className="summary-item difference">
                            <span>Difference (X-Y):</span>
                            <span className={`amount ${billDifference !== 0 ? 'warning' : 'success'}`}>
                                {billDifference.toFixed(2)}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TallyPanel;
