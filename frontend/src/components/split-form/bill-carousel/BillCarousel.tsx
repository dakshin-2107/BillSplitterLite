import React, { useState, useEffect } from 'react';
import type { BillFormData, SplitFormData } from '../../../common/interfaces';
import './BillCarousel.css';
import { CarouselItem, Carousel, CarouselPrevious, CarouselNext, CarouselContent } from '@/components/ui/carousel';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipTrigger, MyToolTip } from '@/components/ui/tooltip';

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
    const imageSync = image ? URL.createObjectURL(image) : null;

    return (
        <>
            <CarouselItem className='carousel-item'>
                {imageSync ? (<img src={imageSync} alt="Bill" className="bill-item-image" />)
                    : (<span style={{ fontSize: '0.7rem', color: '#666' }}>bill image</span>)}
                <div className="bill-item-info">
                    <div>
                        <p className="bill-location">{location || 'Pizza by alfredo\'s'}</p>
                        <p className="bill-date">{date || '25th Dec, 2025'}</p>
                    </div>
                    <div className="bill-item-actions">
                        {/* <Button
                            onClick={() => onClickEditBill(index)}
                            size="lg"
                        >
                            Edit
                        </Button> */}
                        <Button
                            variant="destructive"
                            onClick={() => onClickDeleteBill(index)}
                            size="lg"
                        >
                            Delete
                        </Button>
                    </div>
                </div>
            </CarouselItem>
        </>
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


    useEffect(() => {
        setShowResetButton(splitFormData.billData.length > 0 ? true : false);
        setShowStartSessionButton((splitFormData.billData.length > 0 && splitFormData.peopleList.length > 1) ? true : false);
    }, [splitFormData]);

    return (
        <div className='bill-carousel-container'>
            <Carousel className='bill-carousel-wrapper'>
                <CarouselPrevious className='left-2 bg-input-background' />
                <CarouselContent>
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
                </CarouselContent>
                <CarouselNext className='right-2 bg-input-background' />
            </Carousel>
            <div className="carousel-actions">
                <MyToolTip toolTipText="Deletes all uploaded bills">
                    <Button
                        variant="destructive"
                        onClick={onResetAll}
                        disabled={isSubmitting || !showResetButton}
                    >
                        Delete all
                    </Button>
                </MyToolTip>

                <MyToolTip toolTipText="Start a new session to split all bills">
                    <Button
                        onClick={onStartSession}
                        disabled={isSubmitting || !showStartSessionButton}
                    >
                        {isSubmitting ? 'Starting...' : 'Start split session'}
                    </Button>
                </MyToolTip>
            </div>
        </div>
    )
}

export default BillCarousel;