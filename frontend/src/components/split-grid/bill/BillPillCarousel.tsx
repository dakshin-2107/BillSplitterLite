import { useRef } from 'react';
import type { BillData } from '@utils/interfaces';
import { formatDate } from '@utils/dateUtils';
import './BillPillCarousel.css';

const SCROLL_AMOUNT_PX = 200;

interface BillPillCarouselProps {
    bills: Record<number, BillData>;
    activeBillId: number;
    onBillSelect: (billId: number) => void;
}

export const BillPillCarousel = ({ bills, activeBillId, onBillSelect }: BillPillCarouselProps) => {
    const scrollContainerRef = useRef<HTMLDivElement>(null);

    const handleScroll = (direction: 'left' | 'right') => {
        if (scrollContainerRef.current) {
            scrollContainerRef.current.scrollBy({
                left: direction === 'left' ? -SCROLL_AMOUNT_PX : SCROLL_AMOUNT_PX,
                behavior: 'smooth',
            });
        }
    };

    return (
        <div className="bill-pill-carousel-container">
            <button className="scroll-arrow left" onClick={() => handleScroll('left')} aria-label="Scroll left">
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
                            <span className="pill-date">{bill.date ? formatDate(bill.date) : ''}</span>
                        </div>
                    );
                })}
            </div>
            <button className="scroll-arrow right" onClick={() => handleScroll('right')} aria-label="Scroll right">
                &gt;
            </button>
        </div>
    );
};
