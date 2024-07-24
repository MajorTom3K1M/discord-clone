export interface Profile {
    id: string;
    name: string;
    imageUrl: string;
    email: string;
    servers?: Server[];
    members?: Member[];
    channels?: Channel[];
    created_at: Date;
    updated_at: Date;
}

export enum MemberRole {
    ADMIN = "ADMIN",
    MODERATOR = "MODERATOR",
    GUEST = "GUEST"
}

export interface Member {
    id: string;
    role: MemberRole;
    profileID: string;
    profile: Profile;
    serverID: string;
    server?: Server; // Made optional to prevent circular reference issues
    messages?: Message[]; // Assuming Message is another interface you will define
    directMessages?: DirectMessage[]; // Assuming DirectMessage is another interface you will define
    conversationsInitiated?: Conversation[]; // Assuming Conversation is another interface you will define
    conversationsReceived?: Conversation[]; // Assuming Conversation is another interface you will define
    created_at: Date;
    updated_at: Date;
}

export enum ChannelType {
    TEXT = "TEXT",
    AUDIO = "AUDIO",
    VIDEO = "VIDEO"
}

export interface Channel {
    id: string;
    name: string;
    type: ChannelType;
    profileID: string;
    profile?: Profile;
    serverID: string;
    server?: Server; // Made optional to prevent circular reference issues
    messages?: Message[]; // Assuming Message is another interface you will define
    created_at: Date;
    updated_at: Date;
}

export interface Server {
    id: string; // Assuming uuid.UUID translates to string in TypeScript
    name: string;
    imageUrl: string;
    inviteCode: string;
    profileID: string; // Assuming uuid.UUID translates to string in TypeScript
    profile?: Profile; // Optional, based on "omitempty"
    members?: Member[]; // Optional, based on "omitempty"
    channels?: Channel[]; // Optional, based on "omitempty"
    created_at: Date; // Assuming time.Time translates to Date in TypeScript
    updated_at: Date;
}

export interface Message {
    id: string;
    content: string;
    fileUrl?: string; // Optional because of the pointer type in Go, indicating it can be null
    memberID: string;
    member: Member; // Made optional to avoid deep nesting issues during type checking
    channelID: string;
    channel?: Channel; // Made optional to avoid deep nesting issues during type checking
    deleted: boolean;
    created_at: Date;
    updated_at: Date;
}

export interface DirectMessage {
    id: string;
    content: string;
    fileUrl?: string; // Optional for the same reason as above
    memberID: string;
    member?: Member; // Optional to simplify type structure and usage
    conversationID: string;
    conversation?: Conversation; // Optional to simplify type structure and usage
    deleted: boolean;
    created_at: Date;
    updated_at: Date;
}

export interface Conversation {
    id: string;
    memberOneID: string;
    memberOne: Member; // Optional to simplify type structure and usage
    memberTwoID: string;
    memberTwo: Member; // Optional to simplify type structure and usage
    directMessages?: DirectMessage[]; // Assuming it can be optional
    created_at: Date;
    updated_at: Date;
}

export interface WebRTCMessage {
    candidate?: RTCIceCandidate
    answer?: RTCLocalSessionDescriptionInit
    streamId?: string
    serverId?: string
}