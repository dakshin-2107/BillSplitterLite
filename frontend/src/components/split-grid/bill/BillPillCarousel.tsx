import React, { useRef } from 'react';
import './BillPillCarousel.css';
import type { BillData } from '../../../common/interfaces';

interface BillPillCarouselProps {
    bills: Record<number, BillData>;
    activeBillId: number;
    onBillSelect: (billId: number) => void;
}

export const BillPillCarousel: React.FC<BillPillCarouselProps> = ({
    bills,
    activeBillId,
    onBillSelect
}) => {
    const scrollContainerRef = useRef<HTMLDivElement>(null);

    const scroll = (direction: 'left' | 'right') => {
        if (scrollContainerRef.current) {
            const scrollAmount = 200;
            scrollContainerRef.current.scrollBy({
                left: direction === 'left' ? -scrollAmount : scrollAmount,
                behavior: 'smooth'
            });
        }
    };

    return (
        <div className="bill-pill-carousel-container">
            <button className="scroll-arrow left" onClick={() => scroll('left')}>
                &lt;
            </button>
            <div className="bill-pills-wrapper" ref={scrollContainerRef}>
                {Object.entries(bills).map(([idStr, bill]) => {
                    const id = Number(idStr);
                    const isActive = id === activeBillId;
                    return (
                        <div
                            key={id}
                            className={`bill-pill ${isActive ? 'active' : ''}`}
                            onClick={() => onBillSelect(id)}
                        >
                            <span className="pill-location">{bill.location || `Bill ${id}`}</span>
                            <span className="pill-date">{bill.date}</span>
                        </div>
                    );
                })}
            </div>
            <button className="scroll-arrow right" onClick={() => scroll('right')}>
                &gt;
            </button>
        </div>
    );
};