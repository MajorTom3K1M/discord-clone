"use client"

import axios from '@/utils/axios';

import qs from 'query-string';
import { Check, Gavel, Loader2, MoreVertical, Shield, ShieldAlert, ShieldCheck, ShieldQuestion } from 'lucide-react';
import { useEffect, useState } from 'react';
import { MemberRole, Server } from '@/types/models';
import { useRouter } from 'next/navigation';

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle
} from '@/components/ui/Dialog';
import { ScrollArea } from '@/components/ui/ScrollArea';
import { UserAvatar } from '@/components/UserAvatar';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuPortal,
    DropdownMenuSeparator,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuTrigger,
    DropdownMenuSubTrigger
} from '@/components/ui/DropdownMenu'

import { useModal } from '@/hooks/useModalStore';

const roleIconMap = {
    "GUEST": null,
    "MODERATOR": <ShieldCheck className='h-4 w-4 ml-2 text-indigo-500' />,
    "ADMIN": <ShieldAlert className='h-4 w-4 text-rose-500' />
};

export const MembersModal = () => {
    const router = useRouter();
    const { isOpen, onClose, onOpen, type, data } = useModal();
    const [loadingId, setLoadingId] = useState("");

    const isModalOpen = isOpen && type === "members";
    const { server } = data;

    const onKick = async (memberId: string) => {
        try {
            setLoadingId(memberId);

            const response = await axios.delete<{
                message: string,
                server: Server
            }>(`/members/${memberId}/servers/${server?.id}`);

            router.refresh();
            onOpen("members", { server: response.data.server });
        } catch (error) {
            console.error(error);
        } finally {
            setLoadingId("");
        }
    };

    const onRoleChange = async (memberId: string, role: MemberRole) => {
        try {
            setLoadingId(memberId);

            const response = await axios.patch<{
                message: string,
                server: Server
            }>(
                `/members/${memberId}/servers/${server?.id}`,
                { role }
            )

            router.refresh();
            onOpen("members", { server: response.data.server })
        } catch (error) {
            console.log(error);
        } finally {
            setLoadingId("");
        }
    }

    useEffect(() => {
        console.log(server?.members);
        console.log(loadingId);
    }, [loadingId])

    return (
        <Dialog open={isModalOpen} onOpenChange={onClose}>
            <DialogContent className='bg-white text-black overflow-hidden'>
                <DialogHeader className='pt-8 px-6'>
                    <DialogTitle className='text-2xl text-center font-bold'>
                        Manage Members
                    </DialogTitle>
                    <DialogDescription
                        className='text-center text-zinc-500'
                    >
                        {server?.members?.length} Members
                    </DialogDescription>
                </DialogHeader>
                <ScrollArea
                    className='mt-8 max-h-[420px] pr-6'
                >
                    {server?.members?.map((member) =>
                        <div key={member.id} className='flex items-center gap-x-2 mb-6'>
                            <UserAvatar src={member.profile?.imageUrl} />
                            <div className='flex flex-col gap-y-1'>
                                <div className='text-xs font-semibold flex items-center gap-x-1'>
                                    {member.profile?.name}
                                    {roleIconMap[member.role]}
                                </div>
                                <p className='text-xs text-zinc-500'>
                                    {member.profile?.email}
                                </p>
                            </div>
                            {server.profileID !== member.profileID &&
                                loadingId !== member.id && (
                                    <div className='ml-auto'>
                                        <DropdownMenu>
                                            <DropdownMenuTrigger>
                                                <MoreVertical className='h-4 w-4 text-zinc-500' />
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent side='left'>
                                                <DropdownMenuSub>
                                                    <DropdownMenuSubTrigger className='flex items-center'>
                                                        <ShieldQuestion
                                                            className='w-4 h-4 mr-2'
                                                        />
                                                        <span>Role</span>
                                                    </DropdownMenuSubTrigger>
                                                    <DropdownMenuPortal>
                                                        <DropdownMenuSubContent>
                                                            <DropdownMenuItem
                                                                onClick={() => onRoleChange(member.id, MemberRole.GUEST)}
                                                            >
                                                                <Shield className='h-4 w-4 mr-2' />
                                                                Guest
                                                                {member.role === "GUEST" && (
                                                                    <Check
                                                                        className='h-4 w-4 ml-auto'
                                                                    />
                                                                )}
                                                            </DropdownMenuItem>
                                                            <DropdownMenuItem
                                                                onClick={() => onRoleChange(member.id, MemberRole.MODERATOR)}
                                                            >
                                                                <ShieldCheck className='h-4 w-4 mr-2' />
                                                                Moderator
                                                                {member.role === "MODERATOR" && (
                                                                    <Check
                                                                        className='h-4 w-4 ml-auto'
                                                                    />
                                                                )}
                                                            </DropdownMenuItem>
                                                        </DropdownMenuSubContent>
                                                    </DropdownMenuPortal>
                                                </DropdownMenuSub>
                                                <DropdownMenuSeparator />
                                                <DropdownMenuItem onClick={() => onKick(member.id)} >
                                                    <Gavel className='h-4 w-4 mr-2' />
                                                    Kick
                                                </DropdownMenuItem>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    </div>
                                )}
                            {loadingId === member.id && (
                                <Loader2 className='animate-spin text-zinc-500 ml-auto w-4 h-4' />
                            )}
                        </div>
                    )}
                </ScrollArea>
            </DialogContent>
        </Dialog>
    )
};