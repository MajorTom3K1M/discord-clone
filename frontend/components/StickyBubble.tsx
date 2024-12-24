"use client"

import { useState } from "react"
import { Button } from "@/components/ui/Button"
import { Card, CardContent } from "@/components/ui/Card"
import { ChevronUp, ChevronDown } from 'lucide-react'

export default function StickyBubble() {
    const [isExpanded, setIsExpanded] = useState(false)

    const demoAccounts = [
        { username: "demo1@example.com", password: "demopass1" },
        { username: "demo2@example.com", password: "demopass2" },
        { username: "demo3@example.com", password: "demopass3" },
    ]

    return (
        <div className="fixed bottom-4 right-4 z-50">
            <Card className="w-72 overflow-hidden bg-gradient-to-br from-gray-900 to-gray-800 text-white border-none shadow-lg">
                <CardContent className="p-0">
                    <div className="bg-gradient-to-r from-purple-500 to-pink-500 p-4">
                        <div className="flex justify-between items-center">
                            <h3 className="font-bold text-lg">Quick Access</h3>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setIsExpanded(!isExpanded)}
                                className="h-8 w-8 text-white hover:bg-white/20"
                            >
                                {isExpanded ? (
                                    <ChevronDown className="h-5 w-5" />
                                ) : (
                                    <ChevronUp className="h-5 w-5" />
                                )}
                            </Button>
                        </div>
                        <p className="text-sm mt-2">
                            Don&apos;t want to create an account? Use these demo credentials:
                        </p>
                    </div>
                    <div className={`overflow-hidden transition-all duration-300 ease-in-out ${isExpanded ? 'max-h-96' : 'max-h-0'}`}>
                        <div className="p-4 bg-gray-800">
                            {demoAccounts.map((account, index) => (
                                <div key={index} className="mb-4 pb-4 border-b border-gray-700 last:border-b-0 last:mb-0 last:pb-0">
                                    <p className="text-sm mb-1"><strong>Username:</strong> {account.username}</p>
                                    <p className="text-sm"><strong>Password:</strong> {account.password}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}