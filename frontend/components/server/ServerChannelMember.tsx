import { cn } from "@/lib/utils";
import { UserAvatar } from "../UserAvatar";

interface ServerChannelMemberProps {
    src: string;
    name: string;
    className?: string;
}

export const ServerChannelMember = ({
    src,
    name,
    className
}: ServerChannelMemberProps) => {
    return (
        <div className={cn(
            `text-zinc-500 dark:text-zinc-400 hover:text-zinc-600 
            dark:hover:text-zinc-300 transition text-[14px]`,
            className
        )}>
            <div className='relative flex items-center cursor-pointer flex-1 
                hover:bg-zinc-700/10 dark:hover:bg-zinc-700/50 transition 
                rounded-md pt-1 pb-1'>
                    <UserAvatar 
                        className="h-6 w-6 md:h-6 md:w-6 ml-2 mr-2"
                        src={src}
                    />
                    {name}
            </div>
        </div>
    );
};