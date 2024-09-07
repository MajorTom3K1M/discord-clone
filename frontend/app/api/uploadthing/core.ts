import { decodeJwtPayload } from '@/lib/utils';
import { NextRequest } from 'next/server';
import { createUploadthing, type FileRouter } from 'uploadthing/next';

const f = createUploadthing();

const handleAuth = (req: NextRequest) => {
    try {
        const token = req.cookies.get('access_token')?.value as string;
    
        if (!token) throw new Error("Unauthorized");
        const payload = decodeJwtPayload(token);
        return { profileId: payload.profile_id };
    } catch (err) {
        console.error(err);
        throw new Error("Unauthorized");
    }
};

export const ourFileRouter = {
    serverImage: f({ image: { maxFileSize: "4MB", maxFileCount: 1 } })
        .middleware(({ req }) => handleAuth(req))
        .onUploadComplete(() => { }),
    messageFile: f(["image", "pdf"])
        .middleware(({ req }) => handleAuth(req))
        .onUploadComplete(() => { }),
    profileImage: f({ image: { maxFileSize: "4MB", maxFileCount: 1 } })
        .middleware(async ({ req }) => {
            console.log("Middleware for courseAttachment");
            return {};
        })
        .onUploadError((err) => {
            console.log("Error uploading profile image", err);
        })
        .onUploadComplete(() => {
            console.log("Profile image uploaded");
        }),
} satisfies FileRouter;

export type OurFileRouter = typeof ourFileRouter;