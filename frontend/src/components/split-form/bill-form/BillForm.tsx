import React, { useState, useRef, type ChangeEvent, useEffect } from 'react';
import './BillForm.css';

interface Bill {
    location: string;
    date: string;
    image: File;
}

interface BillFormProps {
    onAddBill: (bill: Bill) => void;
}

const BillForm: React.FC<BillFormProps> = ({ onAddBill }) => {
    const [location, setLocation] = useState('');
    const [date, setDate] = useState('');
    const [image, setImage] = useState<File | null>(null);
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [showResetButton, setShowResetButton] = useState(false);
    const [showAddBillButton, setShowAddBillButton] = useState(false);

    const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setImage(file);
            const reader = new FileReader();
            reader.onloadend = () => {
                setImagePreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleUploadClick = () => {
        fileInputRef.current?.click();
    };

    const handleReset = () => {
        setLocation('');
        setDate('');
        setImage(null);
        setImagePreview(null);
        if (fileInputRef.current) fileInputRef.current.value = '';
    };

    const handleAddBill = () => {
        if (location && date && image) {
            onAddBill({ location, date, image });
            handleReset();
        }
    };

    useEffect(() => {
        setShowResetButton((date || location || image) ? false : true);
        setShowAddBillButton((date && location && image) ? false : true);
    }, [date, location, image]);

    return (
        <div className="bill-form">
            <div className="form-group">
                <label>Location :</label>
                <input
                    type="text"
                    placeholder="Dunder mifflin"
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                />
            </div>

            <div className="form-group">
                <label>Date :</label>
                <input
                    type="date"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                    placeholder="DD/MM/YYYY"
                />
            </div>

            <div className="bill-image-group">
                <label>Bill image :</label>
                <div className="upload-container" onClick={handleUploadClick}>
                    {imagePreview ? (
                        <img src={imagePreview} alt="Bill Preview" className="image-preview" />
                    ) : (
                        <div className="upload-placeholder">
                            <span className="plus-icon">+</span>
                            <p className="upload-text">Click to upload<br />the image of a bill</p>
                        </div>
                    )}
                    <input
                        type="file"
                        ref={fileInputRef}
                        className="file-input"
                        accept="image/*"
                        onChange={handleImageChange}
                    />
                </div>
            </div>

            <div className="form-actions">
                <button
                    className="btn-reset"
                    onClick={handleReset}
                    disabled={showResetButton}
                >
                    Reset
                </button>
                <button
                    className="btn-add-bill"
                    onClick={handleAddBill}
                    disabled={showAddBillButton}
                >
                    Add bill
                </button>
            </div>
        </div>
    );
}

export default BillForm;
