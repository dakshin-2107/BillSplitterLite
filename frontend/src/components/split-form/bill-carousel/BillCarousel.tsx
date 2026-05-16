import { Carousel, CarouselPrevious, CarouselNext, CarouselContent } from '@/components/ui/carousel';
import type { BillFormData } from '../../../common/interfaces';
import BillCarouselItem from './BillCarouselItem';
import './BillCarousel.css';

interface BillCarouselProps {
    bills: BillFormData[];
    setBills: React.Dispatch<React.SetStateAction<BillFormData[]>>;
}

const BillCarousel = ({ bills, setBills }: BillCarouselProps) => {
    const handleDelete = (index: number) => {
        setBills(prev => prev.filter((_, i) => i !== index));
    };

    return (
        <div className='bill-carousel-container'>
            <div className="bill-carousel-header">
                <div>
                    <span className="step-label">Step 3</span>
                    <h2 className="bill-carousel-label">
                        Review your bills and start the session
                        {bills.length > 0 && ` · ${bills.length} bill${bills.length !== 1 ? 's' : ''} added`}
                    </h2>
                </div>
            </div>
            <Carousel className='bill-carousel-wrapper'>
                <CarouselPrevious className='left-2 bg-input-background' />
                <CarouselContent className='bill-carousel-content'>
                    {bills.length > 0 ? (
                        bills.map((bill, index) => (
                            <BillCarouselItem
                                key={`${bill.location}-${bill.date}-${index}`}
                                bill={bill}
                                onDelete={() => handleDelete(index)}
                            />
                        ))
                    ) : (
                        <p className="bill-carousel__empty">No bills added yet. Add one below.</p>
                    )}
                </CarouselContent>
                <CarouselNext className='right-2 bg-input-background' />
            </Carousel>
        </div>
    );
};

export default BillCarousel;
