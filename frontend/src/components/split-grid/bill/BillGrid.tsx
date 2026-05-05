import type { SplitData } from "../../../common/interfaces";
import React from 'react';
import ItemActionCell from '../item/ItemActionCell';
import TakerCell from '../item/TakerCell';

import './BillGrid.css';

interface BillGridProps {
    splitData: SplitData;
    onAddTakerClick: (itemId: number, itemName: string, billId: number) => void;
    activeBillId: number;
}

export const BillGrid: React.FC<BillGridProps> = ({
    splitData,
    onAddTakerClick,
    activeBillId,
}) => {

    return (
        <>
            <div className="table-section">
                <div className="table-wrapper">
                    <table className="split-table">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>Actions</th>
                                <th>Item Name</th>
                                <th>Price</th>
                                <th>Takers</th>
                            </tr>
                        </thead>
                        <tbody>
                            {Object.entries(splitData.bills)
                                .filter(([billIdStr]) => Number(billIdStr) === activeBillId)
                                .map(([billIdStr, bill]) => {
                                    const billId = Number(billIdStr);
                                    return (
                                        <React.Fragment key={billId}>
                                            {Object.values(bill.items).map((item, itemIndex) => (
                                                <tr key={`${billId}-${item.id || itemIndex}`}>
                                                    <td>{itemIndex + 1}</td>
                                                    <td>
                                                        <ItemActionCell
                                                            item={item}
                                                            billId={billId}
                                                        />
                                                    </td>
                                                    <td>{item.name}</td>
                                                    <td>{item.price.toFixed(2)}</td>
                                                    <td className="takers-cell">
                                                        <TakerCell
                                                            takers={item.takers}
                                                            itemId={item.id}
                                                            itemName={item.name}
                                                            billId={billId}
                                                            onAddTakerClick={onAddTakerClick}
                                                        />
                                                    </td>
                                                </tr>
                                            ))}
                                        </React.Fragment>
                                    )
                                })}
                        </tbody>
                    </table>
                </div>
            </div>
        </>
    );
};
