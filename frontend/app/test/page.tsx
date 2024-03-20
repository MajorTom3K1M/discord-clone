"use client"
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/Command";
import { useState, useEffect } from "react";

const Test = () => {
    const [open, setOpen] = useState(false)

    useEffect(() => {
        const down = (e: KeyboardEvent) => {
          if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
            e.preventDefault()
            setOpen((open) => !open)
          }
        }
        document.addEventListener("keydown", down)
        return () => document.removeEventListener("keydown", down)
      }, [])

    return (
        <CommandDialog open={open} onOpenChange={setOpen}>
            <CommandInput placeholder="Type a command or search..." />
            <CommandList asChild>
                <CommandEmpty>No results found.</CommandEmpty>
                {/* <CommandGroup heading="Suggestions"> */}
                    <CommandItem disabled={false}>Calendar</CommandItem>
                    <CommandItem disabled={false}>Search Emoji</CommandItem>
                    <CommandItem disabled={false}>Calculator</CommandItem>
           
                {/* </CommandGroup> */}
            </CommandList>
        </CommandDialog>
    );
}

export default Test;