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
                        <p className="bill-carousel__empty">No bills added yet...</p>
                    )}
                </CarouselContent>
                <CarouselNext className='right-2 bg-input-background' />
            </Carousel>
        </div>
    );
};

export default BillCarousel;
