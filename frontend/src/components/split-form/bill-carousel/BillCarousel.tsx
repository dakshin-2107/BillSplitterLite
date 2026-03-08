import { useState, useEffect } from 'react';
import type { BillFormData } from '../../../common/interfaces';
import './BillCarousel.css';
import { CarouselItem, Carousel, CarouselPrevious, CarouselNext, CarouselContent } from '@/components/ui/carousel';
import { Button } from '@/components/ui/button';

interface BillCarouselProps {
    bills: BillFormData[];
    setBills: React.Dispatch<React.SetStateAction<BillFormData[]>>;
    setEditBill: React.Dispatch<React.SetStateAction<BillFormData | null>>;
}

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
            {imageUrl
                ? <img src={imageUrl} alt="Bill" className="bill-item-image" />
                : <span style={{ fontSize: '0.7rem', color: '#666' }}>bill image</span>
            }
            <div className="bill-item-info">
                <div>
                    <p className="bill-location">{bill.location}</p>
                    <p className="bill-date">{bill.date}</p>
                </div>
                <div className="bill-item-actions">
                    {/* <Button type="button" size="lg" onClick={onEdit}>Edit</Button> */}
                    <Button type="button" variant="destructive" size="lg" onClick={onDelete}>Delete</Button>
                </div>
            </div>
        </CarouselItem>
    );
};

const BillCarousel = ({ bills, setBills, setEditBill }: BillCarouselProps) => {
    const handleDelete = (index: number) => {
        setBills(prev => prev.filter((_, i) => i !== index));
    };

    return (
        <div className='bill-carousel-container'>
            <Carousel className='bill-carousel-wrapper'>
                <CarouselPrevious className='left-2 bg-input-background' />
                <CarouselContent>
                    {bills.length > 0 ? (
                        bills.map((bill, index) => (
                            <BillCarouselItem
                                key={`${bill.location}-${bill.date}-${index}`}
                                bill={bill}
                                onDelete={() => handleDelete(index)}
                            />
                        ))
                    ) : (
                        <p style={{ color: '#666', padding: '1rem' }}>No bills added yet...</p>
                    )}
                </CarouselContent>
                <CarouselNext className='right-2 bg-input-background' />
            </Carousel>
        </div>
    );
};

export default BillCarousel;
