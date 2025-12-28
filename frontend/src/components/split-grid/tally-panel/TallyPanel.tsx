import React, { useRef } from 'react';
import './TallyPanel.css';
import type { Tally, BillData } from '../../../common/interfaces';
import html2canvas from 'html2canvas';

interface TallyPanelProps {
    tally: Tally | null;
    participants: Record<string, string>;
    billData: BillData;
}

const TallyPanel: React.FC<TallyPanelProps> = ({ tally, participants, billData }) => {
    const panelRef = useRef<HTMLDivElement>(null);

    // const handleCopy = () => {
    //     if (!tally) return;

    //     const location = `${billData.location} - ${billData.date}`;

    //     const userBreakdown = Object.entries(tally.userShares)
    //         .map(([userId, userShare]) => {
    //             const name = participants[userId] || userId;
    //             return `${name} - ${userShare.userShareTotal.toFixed(2)}`;
    //         })
    //         .join('\n');

    //     const textToCopy = `Location: ${location}\nDate: ${billData.date}\n\nPer person breakdown:\n${userBreakdown}\n\nTotal - ${tally.calculatedTotal.toFixed(2)}`;
    //     navigator.clipboard.writeText(textToCopy).then(() => {
    //         // Optional: Add a temporary "Copied!" state if needed, but for now just console log
    //         console.log("Tally copied to clipboard");
    //     });
    // };

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

    if (!tally) {
        return (
            <div className="tally-panel empty">
                <h3>Bill Tally</h3>
                <p>Waiting for data...</p>
            </div>
        );
    }

    return (
        <div className="tally-panel" ref={panelRef}>
            <div className="panel-header">
                <h3>Per user split</h3>
                {/* <button className="copy-btn" onClick={handleCopy}>
                    Copy text
                </button> */}
                <button className="copy-btn" onClick={handleCopyImage}>
                    Copy
                </button>
            </div>
            <div className="user-shares-list">
                {Object.entries(tally.userShares).map(([userId, userShare]) => (
                    <div key={userId} className="user-share-card">
                        <div className="user-info">
                            <span className="user-name">{participants[userId] || userId}</span>
                            <span className="user-total">{userShare.userShareTotal.toFixed(2)}</span>
                        </div>
                        <div className="item-breakdown">
                            {Object.entries(userShare.shares).map(([itemName, amount]) => (
                                <div key={itemName} className="item-share">
                                    <span className="item-name">{itemName}</span>
                                    <span className="item-amount">{amount.toFixed(2)}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>

            <div className="tally-header">
                <h3></h3>
                <div className="total-summary">
                    <div className="summary-item">
                        <span>Actual:</span>
                        <span className="amount">{tally.actualTotal.toFixed(2)}</span>
                    </div>
                    <div className="summary-item">
                        <span>Calculated:</span>
                        <span className="amount">{tally.calculatedTotal.toFixed(2)}</span>
                    </div>
                    <div className="summary-item difference">
                        <span>Difference:</span>
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
