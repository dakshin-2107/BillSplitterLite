import React from 'react';
import type { BillData } from '../../common/interfaces';
import './BillDebugger.css';

interface BillDebuggerProps {
    billData: BillData;
}

const BillDebugger: React.FC<BillDebuggerProps> = ({ billData }) => {
    return (
        <div className="bill-debugger-container">
            <div className="bill-debugger-header">
                <span className="bill-debugger-title">Bill Data State Debugger</span>
            </div>
            <div className="bill-debugger-content">
                <pre>
                    {JSON.stringify(billData, null, 4)}
                </pre>
            </div>
        </div>
    );
};

export default BillDebugger;
