import { Avatar, AvatarImage } from "@/components/ui/Avatar"
import { cn } from "@/lib/utils";

interface UserAvatarProps {
    src?: string;
    className?: string;
}

export const UserAvatar = ({
    src,
    className
}: UserAvatarProps) => {
    return (  
        <Avatar className={cn(
            "h-7 w-7 md:h-10 md:w-10",
            className
        )}>
            <AvatarImage src={src} className="object-cover" />
        </Avatar>
    );
}
 