import { useState, useEffect } from 'react';
import { CarouselItem } from '@/components/ui/carousel';
import { Button } from '@/components/ui/button';
import type { BillFormData } from '../../../common/interfaces';
import './BillCarousel.css';

interface BillCarouselItemProps {
    bill: BillFormData;
    onDelete: () => void;
}

const BillCarouselItem = ({ bill, onDelete }: BillCarouselItemProps) => {
    const [imageUrl, setImageUrl] = useState<string | null>(null);

    useEffect(() => {
        if (!bill.image) { setImageUrl(null); return; }
        const url = URL.createObjectURL(bill.image);
        setImageUrl(url);
        return () => URL.revokeObjectURL(url);
    }, [bill.image]);

    return (
        <CarouselItem className='carousel-item'>
            <div className="bill-item-thumbnail">
                {imageUrl
                    ? <img src={imageUrl} alt="Bill" className="bill-item-image" />
                    : <span className="bill-item-placeholder">?</span>
                }
            </div>
            <div className="bill-item-info">
                <p className="bill-location">{bill.location}</p>
                <p className="bill-date">{bill.date}</p>
                <Button type="button" variant="destructive" size="sm" onClick={onDelete} className="bill-delete-btn">
                    Remove
                </Button>
            </div>
        </CarouselItem>
    );
};

export default BillCarouselItem;
