import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import './PeoplePanel.css';

interface PeoplePanelProps {
    people: string[];
    setPeople: (people: string[]) => void;
}

const PeoplePanel = ({ people, setPeople }: PeoplePanelProps) => {
    const [name, setName] = useState('');
    const [error, setError] = useState<string | null>(null);

    const validateName = (value: string) => {
        const trimmed = value.trim();
        if (trimmed === '') return false;

        if (people.some(p => p.toLowerCase() === trimmed.toLowerCase())) {
            setError('Name already added');
            return false;
        }

        setError(null);
        return true;
    };

    const handleAdd = () => {
        if (validateName(name)) {
            setPeople([...people, name.trim()]);
            setName('');
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleAdd();
    };

    const handleRemove = (index: number) => {
        setPeople(people.filter((_, i) => i !== index));
    };

    return (
        <div className="people-panel">
            <h2 className="panel-title">People involved :</h2>
            <div className="people-input-group">
                <Input
                    type="text"
                    className="people-input"
                    value={name}
                    placeholder="Kevin Malone"
                    onChange={(e) => {
                        setName(e.target.value);
                        if (error) setError(null);
                    }}
                    onKeyDown={handleKeyDown}
                />
                <Button onClick={handleAdd}>Add</Button>
            </div>
            {error && <p className="people-panel__error">{error}</p>}

            <div className="people-list-container">
                {people.length > 0 ? (
                    people.map((person, index) => (
                        <div key={index} className="person-item">
                            <span className="person-name">{person}</span>
                            <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => handleRemove(index)}
                                aria-label="Remove person"
                            >
                                X
                            </Button>
                        </div>
                    ))
                ) : (
                    <p className="people-panel__empty">No people added yet...</p>
                )}
            </div>
        </div>
    );
};

export default PeoplePanel;
