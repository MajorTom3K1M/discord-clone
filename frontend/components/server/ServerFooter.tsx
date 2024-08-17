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
import { Settings, LogOut, ShieldCheck, ShieldAlert } from "lucide-react";
import { useAuth } from "../providers/AuthProvider";

interface ServerFooterProps {
    role?: MemberRole;
}

const roleIconMap = {
    [MemberRole.GUEST]: null,
    [MemberRole.MODERATOR]: <ShieldCheck className='ml-1 h-4 w-4 mr-2 text-indigo-500' />,
    [MemberRole.ADMIN]: <ShieldAlert className='ml-1 h-4 w-4 mr-2 text-rose-500' />,
};

const roleNames = {
    [MemberRole.GUEST]: 'Guest',
    [MemberRole.MODERATOR]: 'Moderator',
    [MemberRole.ADMIN]: 'Admin',
};

export const ServerFooter = ({
    role
}: ServerFooterProps) => {
    const { onOpen } = useModal();
    const { authState, signout } = useAuth();
    const { profile } = authState;

    const imageUrl = profile?.imageUrl;
    const name = profile?.name;
    const isAdmin = role === MemberRole.ADMIN;

    return (
        <div className="flex items-center p-2 dark:bg-[#222327] bg-[#E3E5E8] drop-shadow-md justify-between">
            <div className='flex items-center'>
                <UserAvatar
                    src={imageUrl}
                    className="h-8 w-8 md:h-8 md:w-8 ml-1 mr-1"
                />
                <div className="text-xs">
                    <h3 className='font-bold'>{name}</h3>
                    <p className='flex font-thin'>{roleNames[role ?? MemberRole.GUEST]}{roleIconMap[role ?? MemberRole.GUEST]}</p>
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
                        onClick={() => signout()}
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

