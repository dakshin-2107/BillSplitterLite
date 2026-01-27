import React, { useState } from 'react';
import './PeoplePanel.css';

interface PeoplePanelProps {
    people: string[];
    setPeople: (people: string[]) => void;
}

const PeoplePanel: React.FC<PeoplePanelProps> = ({ people, setPeople }) => {

    const [name, setName] = useState('');
    const [error, setError] = useState<string | null>(null);

    const validateName = (name: string) => {
        const trimmedName = name.trim();
        if (trimmedName === "") return false;

        // Check uniqueness (case-insensitive)
        if (people.some(p => p.toLowerCase() === trimmedName.toLowerCase())) {
            setError('Name already added');
            return false;
        }

        setError(null);
        return true;
    }

    const handleAdd = () => {
        if (validateName(name)) {
            setPeople([...people, name.trim()]);
            setName('');
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleAdd();
        }
    };

    const handleRemove = (index: number) => {
        setPeople(people.filter((_, i) => i !== index));
    };

    return (
        <div className="people-panel">
            <h2 className="panel-title">People involved :</h2>
            <div className="people-input-group">
                <input
                    type="text"
                    className="people-input"
                    value={name}
                    placeholder="Kevin Malone"
                    onChange={(e) => {
                        setName(e.target.value);
                        if (error) setError(null);
                    }}
                    onKeyPress={handleKeyPress}
                />
                <button className="btn-add-person" onClick={handleAdd}>Add</button>
            </div>
            {error && <p style={{ color: '#ff4444', fontSize: '0.8rem', marginTop: '-1rem', marginBottom: '1rem' }}>{error}</p>}

            <div className="people-list-container">
                {people.length > 0 ? (
                    people.map((person, index) => (
                        <div key={index} className="person-item">
                            <span className="person-name">{person}</span>
                            <button
                                className="btn-remove-person"
                                onClick={() => handleRemove(index)}
                                aria-label="Remove person"
                            >
                                −
                            </button>
                        </div>
                    ))
                ) : (
                    <p style={{ color: '#666', padding: '1rem', textAlign: 'center' }}>No people added yet...</p>
                )}
            </div>
        </div>
    )
}

export default PeoplePanel;