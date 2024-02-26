import { decodeJwtPayload } from '@/lib/utils';
import { NextRequest } from 'next/server';
import { createUploadthing, type FileRouter } from 'uploadthing/next';

const f = createUploadthing();

const handleAuth = (req: NextRequest) => {
    const token = req.cookies.get('access_token')?.value as string;

    if(!token) throw new Error("Unauthorized");
    const payload = decodeJwtPayload(token);
    return { profileId: payload.profile_id };
};

export const ourFileRouter = {
    serverImage: f({ image: { maxFileSize: "4MB", maxFileCount: 1 } })
        .middleware(({ req }) => handleAuth(req))
        .onUploadComplete(() => {}),
    messageFile: f(["image", "pdf"])
        .middleware(({ req }) => handleAuth(req))
        .onUploadComplete(() => {})
} satisfies FileRouter;

export type OurFileRouter = typeof ourFileRouter;