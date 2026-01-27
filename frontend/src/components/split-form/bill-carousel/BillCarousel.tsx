import React, { useRef, useState, useEffect } from 'react';
import type { BillFormData, SplitFormData } from '../../../common/interfaces';
import './BillCarousel.css';

interface BillCarouselProps {
    splitFormData: SplitFormData;
    onClickDeleteBill: (index: number) => void;
    onClickEditBill: (index: number) => void;
    onStartSession: () => void;
    onResetAll: () => void;
    isSubmitting: boolean;
}

interface BillCarouselItemProps {
    index: number;
    billFormData: BillFormData;
    onClickDeleteBill: (index: number) => void;
    onClickEditBill: (index: number) => void;
}

const BillCarouselItem: React.FC<BillCarouselItemProps> = ({ index, billFormData, onClickDeleteBill, onClickEditBill }) => {
    const { location, date, image } = billFormData;

    // Create preview URL if image exists
    const imageSync = image ? URL.createObjectURL(image) : null;

    return (
        <div className="bill-carousel-item">
            <div className="bill-item-image-wrapper">
                {imageSync ? (
                    <img src={imageSync} alt="Bill" className="bill-item-image" />
                ) : (
                    <span style={{ fontSize: '0.7rem', color: '#666' }}>bill image</span>
                )}
            </div>
            <div className="bill-item-info">
                <div>
                    <p className="bill-location">{location || 'Pizza by alfredo\'s'}</p>
                    <p className="bill-date">{date || '25th Dec, 2025'}</p>
                </div>
                <div className="bill-item-actions">
                    <button className="btn-add-bill" onClick={() => onClickEditBill(index)}>Edit</button>
                    <button className="btn-reset" onClick={() => onClickDeleteBill(index)}>Delete</button>
                </div>
            </div>
        </div>
    )
}

const BillCarousel: React.FC<BillCarouselProps> = ({
    splitFormData,
    onClickDeleteBill,
    onClickEditBill,
    onStartSession,
    onResetAll,
    isSubmitting,
}) => {

    const [showResetButton, setShowResetButton] = useState(false);
    const [showStartSessionButton, setShowStartSessionButton] = useState(false);

    const scrollRef = useRef<HTMLDivElement>(null);

    const scroll = (direction: 'left' | 'right') => {
        if (scrollRef.current) {
            const scrollAmount = 300;
            scrollRef.current.scrollBy({
                left: direction === 'left' ? -scrollAmount : scrollAmount,
                behavior: 'smooth'
            });
        }
    };

    useEffect(() => {
        setShowResetButton(splitFormData.billData.length > 0 ? true : false);
        setShowStartSessionButton((splitFormData.billData.length > 0 && splitFormData.peopleList.length > 1) ? true : false);
    }, [splitFormData]);

    return (
        <div className="bill-carousel-container">
            {/* <h2 className="carousel-title">Bills in this split :</h2> */}
            <div className="carousel-wrapper">
                <button className="nav-button" onClick={() => scroll('left')}>
                    ←
                </button>
                <div className="bill-items-list" ref={scrollRef}>
                    {splitFormData.billData.length > 0 ? (
                        splitFormData.billData.map((bill, index) => (
                            <BillCarouselItem
                                key={index}
                                index={index}
                                billFormData={bill}
                                onClickDeleteBill={onClickDeleteBill}
                                onClickEditBill={onClickEditBill}
                            />
                        ))
                    ) : (
                        <p style={{ color: '#666', padding: '1rem' }}>No bills added yet...</p>
                    )}
                </div>
                <button className="nav-button" onClick={() => scroll('right')}>
                    →
                </button>
            </div>

            <div className="carousel-actions">
                <button
                    className="btn-reset-all btn-reset"
                    onClick={onResetAll}
                    disabled={isSubmitting || !showResetButton}
                >
                    Reset all
                </button>
                <button
                    className="btn-start-session"
                    onClick={onStartSession}
                    disabled={isSubmitting || !showStartSessionButton}
                >
                    {isSubmitting ? 'Starting...' : 'Start split session'}
                </button>
            </div>
        </div>
    )
}

export default BillCarousel;