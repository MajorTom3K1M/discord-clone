"use client"

import axios from "@/utils/axios";

import { MemberRole } from '@/types/models';
import { redirect } from 'next/navigation';

import {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator
} from '@/components/ui/DropdownMenu';
import { useModal } from "@/hooks/useModalStore";
import { UserAvatar } from "@/components/UserAvatar";
import { Settings, LogOut } from "lucide-react";

interface ServerFooterProps {
    role?: MemberRole;
}

export const ServerFooter = ({
    role
}: ServerFooterProps) => {
    const { onOpen } = useModal();

    const isAdmin = role === MemberRole.ADMIN;
    const isModerator = isAdmin || role === MemberRole.MODERATOR;

    return (
        <div className="flex items-center p-2 dark:bg-[#222327] bg-[#E3E5E8] drop-shadow-md justify-between">
            <div className='flex items-center'>
                <UserAvatar
                    src={"https://images.pexels.com/photos/220453/pexels-photo-220453.jpeg?auto=compress&cs=tinysrgb&dpr=2&w=80"}
                    className="h-8 w-8 md:h-8 md:w-8 ml-1 mr-1"
                />
                <div className="text-xs">
                    <h3 className='font-bold'>Matt Cooper</h3>
                    <p className='font-thin'>Moderator</p>
                </div>
            </div>
            <DropdownMenu>
                <DropdownMenuTrigger className="ml-2 mr-2 items-center">
                    <Settings className='h-5 w-5' />
                </DropdownMenuTrigger>
                <DropdownMenuContent
                    className="w-56 text-xs font-medium text-black dark:text-neutral-400 space-y-[2px]"
                >
                    <DropdownMenuItem
                        // onClick={() => onOpen("invite", { server })}
                        className="text-rose-500 px-3 py-2 text-sm cursor-pointer"
                    >
                        Sign Out
                        <LogOut className="h-4 w-4 ml-auto" />
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    )
}

