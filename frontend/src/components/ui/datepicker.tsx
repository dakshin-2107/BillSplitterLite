"use client"

import * as React from "react"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import { format, parse } from "date-fns"
import { CalendarIcon } from "lucide-react"
import { cn } from "@/lib/utils"

interface DatePickerProps {
    value?: string // ISO date string (YYYY-MM-DD)
    onDateChange?: (value: string) => void
    placeholder?: string
    className?: string
}

function DatePicker({ value, onDateChange, placeholder = "Pick a date", className }: DatePickerProps & React.ComponentProps<"input">) {
    const [open, setOpen] = React.useState(false)

    // Parse string value to Date for the Calendar
    const selectedDate = React.useMemo(() => {
        if (!value) return undefined
        try {
            return parse(value, "PPP", new Date())
        } catch {
            return undefined
        }
    }, [value])

    const handleSelect = (date: Date | undefined) => {
        if (date && onDateChange) {
            onDateChange(format(date, "PPP"))
        }
        setOpen(false)
    }

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="defaultEmpty"
                    className={cn(
                        "w-full justify-start text-left font-normal h-input",
                        `rounded-input border border-input-border-color bg-input-background px-4 py-2 text-base text-white 
                        file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-zinc-500 
                        outline-none transition-all focus:border-purple-400`,
                        !selectedDate && "text-muted-foreground",
                        className
                    )}
                >
                    {/* <CalendarIcon className="mr-2 h-4 w-4" /> */}
                    {selectedDate ? format(selectedDate, "PPP") : <span>{placeholder}</span>}
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                    showOutsideDays={false}
                    mode="single"
                    selected={selectedDate}
                    onSelect={handleSelect}
                    defaultMonth={selectedDate}
                />
            </PopoverContent>
        </Popover>
    )
}

export { DatePicker }
